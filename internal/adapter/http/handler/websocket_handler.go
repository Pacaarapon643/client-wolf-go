package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/port"

	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
	GameId   string `json:"game_id"`
	RoomID   string `json:"room_id"`
	UserName string `json:"username"`
	Content  string `json:"content"`
}

type MessageGame struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type Client struct {
	Conn     *websocket.Conn
	RoomID   string
	UserID   string
	GameId   string
	UserName string
	Role     string
	IsDead   bool
	SeerUsed bool // ตรวจแล้วในคืนนี้หรือยัง
	Send     chan []byte
}

type Room struct {
	Id            string
	Client        map[string]*Client
	mu            sync.RWMutex
	Phase         string
	RemainingTime int
	TimerCancel   context.CancelFunc
	IsGameStarted bool
	NightCount    int
	GameId        string
	service       port.Service // เก็บ service ref เพื่อเรียก SummaryVote ตรงจาก server
}

type RoomManager struct {
	Room map[string]*Room
	mu   sync.RWMutex
}

var Manager = &RoomManager{
	Room: make(map[string]*Room),
}

func (rm *RoomManager) JoinRoom(roomId string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if room, ok := rm.Room[roomId]; ok {
		return room
	}

	newRoom := &Room{
		Id:     roomId,
		Client: make(map[string]*Client),
	}

	rm.Room[roomId] = newRoom
	return newRoom
}

func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Client[client.UserID] = client
}

func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Client, client.UserID)
}

func (r *Room) StartTimer(duration int, phase string, roomId string) {
	r.mu.Lock()
	if r.TimerCancel != nil {
		r.TimerCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.TimerCancel = cancel
	r.RemainingTime = duration
	r.Phase = phase
	r.IsGameStarted = true
	r.broadcastTimeSync()
	r.mu.Unlock()

	go func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				r.mu.Lock()
				r.RemainingTime--
				r.broadcastTimeSync()

				if r.RemainingTime <= 0 {
					r.mu.Unlock()
					r.NextPhase(roomId)
					return
				}
				r.mu.Unlock()

			case <-ctx.Done():
				log.Printf("Timer for phase %s stopped.", phase)
				return
			}
		}
	}(ctx)
}

func (r *Room) broadcastTimeSync() {
	msg := map[string]interface{}{
		"type":      "sync_time",
		"phase":     r.Phase,
		"remaining": r.RemainingTime,
	}
	payload, _ := json.Marshal(msg)

	for _, client := range r.Client {
		select {
		case client.Send <- payload:
		default:
		}
	}
}

// NextPhase — ทำ summary ฝั่ง server เลย ไม่ต้อง broadcast กลับไปให้ client
func (r *Room) NextPhase(roomId string) {
	r.mu.RLock()
	currentPhase := r.Phase
	gameId := r.GameId
	r.mu.RUnlock()

	switch currentPhase {
	case "ready_to_start":
		r.mu.Lock()
		r.NightCount++
		for _, c := range r.Client {
			c.SeerUsed = false
		}
		r.mu.Unlock()
		r.BroadcastGameEvent("night", "เข้าสู่คืนที่ "+strconv.Itoa(r.NightCount), r.NightCount)
		r.StartTimer(30, "night", roomId)

	case "night":
		// สรุปผลคืน — เรียกจาก server โดยตรง (ไม่ผ่าน client)
		if r.service != nil && gameId != "" {
			r.doSummary(gameId, "night", roomId)
		}
		r.StartTimer(60, "day", roomId)

	case "day":
		r.StartTimer(30, "vote", roomId)

	case "vote":
		// สรุปผลโหวต — เรียกจาก server โดยตรง
		if r.service != nil && gameId != "" {
			r.doSummary(gameId, "vote", roomId)
		}

		// เริ่มคืนใหม่
		r.mu.Lock()
		r.NightCount++
		for _, c := range r.Client {
			c.SeerUsed = false
		}
		r.mu.Unlock()
		r.BroadcastGameEvent("night", "เข้าสู่คืนที่ "+strconv.Itoa(r.NightCount), r.NightCount)
		r.StartTimer(30, "night", roomId)
	}
}

// doSummary — ทำ summary ครั้งเดียวจากฝั่ง server
func (r *Room) doSummary(gameId string, phase string, roomId string) {
	ctx := context.Background()

	summaryResult, err := r.service.SummaryVote(ctx, gameId, phase)
	if err != nil {
		log.Println("Summary Error:", err)
		return
	}

	// Broadcast event ผลลัพธ์
	if summaryResult != nil {
		r.BroadcastGameEvent(summaryResult.Result, summaryResult.Message, r.NightCount)
	}

	// ถ้ามีคนตาย mark client ว่าตาย
	if summaryResult != nil && summaryResult.Result == "dead" {
		r.markClientDead(summaryResult.DeadSlot, gameId)
	}

	// Reload game player data
	loadMsg, _ := json.Marshal(Message{Type: "load_game"})
	Manager.BroadcastToRoom(loadMsg, roomId)

	// Reset votes
	if err := r.service.ResetVotes(ctx, gameId); err != nil {
		log.Println("Reset Votes Error:", err)
	}

	// Check win condition
	winResult, err := r.service.CheckWinCondition(ctx, gameId)
	if err != nil {
		log.Println("Check Win Error:", err)
		return
	}
	if winResult != nil {
		winMsg := map[string]interface{}{
			"type":    "game_over",
			"winner":  winResult.Winner,
			"message": winResult.Message,
		}
		payload, _ := json.Marshal(winMsg)
		Manager.BroadcastToRoom(payload, roomId)

		// หยุด timer + reset Room state ทั้งหมด
		r.mu.Lock()
		if r.TimerCancel != nil {
			r.TimerCancel()
		}
		r.IsGameStarted = false
		r.NightCount = 0
		r.Phase = ""
		r.GameId = ""
		// reset client state
		for _, c := range r.Client {
			c.IsDead = false
			c.SeerUsed = false
		}
		r.mu.Unlock()
	}
}

// markClientDead — หา user_id ของ slot ที่ตาย แล้ว set Client.IsDead = true
func (r *Room) markClientDead(deadSlot int, gameId string) {
	// หา player data ทั้งหมดจาก GetGame
	result, err := r.service.GetGame(context.Background(), gameId, "")
	if err != nil {
		log.Println("markClientDead: GetGame error:", err)
		return
	}

	// หา user_id ของ slot ที่ตาย
	deadUserId := ""
	if players, ok := result.([]dto.GameResponse); ok {
		for _, p := range players {
			if p.SlotIndex == deadSlot {
				deadUserId = p.UserId
				break
			}
		}
	}

	if deadUserId == "" {
		log.Printf("markClientDead: could not find userId for slot %d", deadSlot)
		return
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// set IsDead = true สำหรับ client ที่ตาย
	for _, client := range r.Client {
		if client.UserID == deadUserId {
			client.IsDead = true
			log.Printf("[GAME] Marked client dead: User=%s Slot=%d", client.UserID, deadSlot)
			break
		}
	}

	// ส่ง death notification ให้ทุกคนรู้ว่า slot ไหนตาย
	deathMsg := map[string]interface{}{
		"type":      "player_died",
		"dead_slot": deadSlot,
	}
	payload, _ := json.Marshal(deathMsg)
	for _, client := range r.Client {
		select {
		case client.Send <- payload:
		default:
		}
	}
}

// BroadcastGameEvent ส่ง event log ไปให้ทุกคนในห้อง
func (r *Room) BroadcastGameEvent(eventType string, message string, nightCount int) {
	msg := map[string]interface{}{
		"type":        "game_event",
		"event_type":  eventType,
		"message":     message,
		"night_count": nightCount,
	}
	payload, _ := json.Marshal(msg)

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, client := range r.Client {
		select {
		case client.Send <- payload:
		default:
		}
	}
}

// BroadcastToWerewolves ส่ง message ให้เฉพาะหมาป่า
func (r *Room) BroadcastToWerewolves(payload []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, client := range r.Client {
		if client.Role == "werewolf" {
			select {
			case client.Send <- payload:
			default:
			}
		}
	}
}

// BroadcastToDeadPlayers ส่ง message ให้เฉพาะคนตาย
func (r *Room) BroadcastToDeadPlayers(payload []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, client := range r.Client {
		if client.IsDead {
			select {
			case client.Send <- payload:
			default:
			}
		}
	}
}

func sendToClient(client *Client, payload []byte) {
	select {
	case client.Send <- payload:
	default:
	}
}

func (m *RoomManager) BroadcastToLobby(payload []byte) {
	lobby, ok := m.Room["lobby"]
	if !ok {
		return
	}

	lobby.mu.RLock()
	defer lobby.mu.RUnlock()

	for _, client := range lobby.Client {
		select {
		case client.Send <- payload:
		default:
		}
	}
}

func (m *RoomManager) BroadcastToRoom(payload []byte, roomId string) {
	lobby, ok := m.Room[roomId]
	if !ok {
		return
	}

	lobby.mu.RLock()
	defer lobby.mu.RUnlock()

	for _, client := range lobby.Client {
		select {
		case client.Send <- payload:
		default:
		}
	}
}

func (h Handler) RoomWebSocket(ws *websocket.Conn) {
	_, msgBytes, err := ws.ReadMessage()
	if err != nil {
		return
	}

	var joinMsg Message
	if err := json.Unmarshal(msgBytes, &joinMsg); err != nil {
		log.Printf("Failed to unmarshal join message: %v", err)
		return
	}

	room := Manager.JoinRoom(joinMsg.RoomID)
	client := &Client{
		Conn: ws, RoomID: joinMsg.RoomID, UserID: joinMsg.UserID,
		UserName: joinMsg.UserName, Send: make(chan []byte, 256),
	}

	defer func() {
		room.RemoveClient(client)
		if client.RoomID != "lobby" {

			err = h.s.LeaveRoom(context.Background(), client.RoomID, client.UserID)
			if err != nil {
				log.Println("Leave Room Error:", err)
			}

			err = h.s.LeaveRoomJoin(context.Background(), client.RoomID)
			if err != nil {
				log.Println("Leave Room Join Error:", err)
			}

			msgLoadLobby := Message{
				Type:     "lobby",
				UserID:   client.UserID,
				RoomID:   "lobby",
				UserName: client.UserName,
			}

			payload, err := json.Marshal(msgLoadLobby)
			if err != nil {
				log.Println("Marshal Error:", err)
				return
			}
			log.Println("แจ้งเตือนๆ lobby")
			Manager.BroadcastToLobby(payload)

			log.Println("โหลด room")
			loadMsg, _ := json.Marshal(Message{Type: "load_room"})
			Manager.BroadcastToRoom(loadMsg, client.RoomID)
		}
		ws.Close()
	}()

	room.AddClient(client)

	if client.RoomID != "lobby" {
		if err := h.s.JoinRoom(context.Background(), client.RoomID); err != nil {
			log.Printf("Failed to join room: %v", err)
		}

		loadMsg, _ := json.Marshal(Message{Type: "load_room"})
		Manager.BroadcastToRoom(loadMsg, client.RoomID)

		msgLoadLobby := Message{
			Type:     "lobby",
			UserID:   client.UserID,
			RoomID:   "lobby",
			UserName: client.UserName,
		}

		payload, err := json.Marshal(msgLoadLobby)
		if err != nil {
			log.Println("Marshal Error:", err)
			return
		}

		Manager.BroadcastToLobby(payload)
	}

	go func() {
		for m := range client.Send {
			if err := ws.WriteMessage(websocket.TextMessage, m); err != nil {
				return
			}
		}
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		var msgData Message
		if err := json.Unmarshal(msg, &msgData); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		log.Println("msgData: ", msgData)
		if msgData.Type == "ready" {
			err := h.s.ReadyRoom(context.Background(), client.RoomID, client.UserID, msgData.Content)
			if err != nil {
				log.Println("Ready Room Error:", err)
			}

			loadMsg, err := json.Marshal(Message{Type: "load_room"})
			if err != nil {
				log.Println("Marshal Error:", err)
				return
			}
			Manager.BroadcastToRoom(loadMsg, client.RoomID)

		} else if msgData.Type == "start_game" {
			game_id, err := h.s.RoleDistribution(context.Background(), client.RoomID)
			if err != nil {
				log.Println("Role Distribution Error:", err)
				return
			}

			loadMsg, err := json.Marshal(Message{Type: "start_game", Content: *game_id})
			if err != nil {
				log.Println("Marshal Error:", err)
				return
			}
			Manager.BroadcastToRoom(loadMsg, client.RoomID)
		}
	}
}

func (h Handler) GameWebSocket(ws *websocket.Conn) {

	_, msgBytes, err := ws.ReadMessage()
	if err != nil {
		return
	}

	var joinMsg Message
	if err := json.Unmarshal(msgBytes, &joinMsg); err != nil {
		log.Printf("Failed to unmarshal game join message: %v", err)
		return
	}

	room := Manager.JoinRoom(joinMsg.RoomID)
	client := &Client{
		Conn:     ws,
		RoomID:   joinMsg.RoomID,
		UserID:   joinMsg.UserID,
		GameId:   joinMsg.GameId,
		UserName: joinMsg.UserName,
		Send:     make(chan []byte, 256),
	}

	defer func() {
		room.RemoveClient(client)
		err = h.s.LeaveGame(context.Background(), client.GameId, client.UserID)
		if err != nil {
			log.Println("Leave Game Error:", err)
		}
		loadMsg, _ := json.Marshal(Message{Type: "load_game"})
		Manager.BroadcastToRoom(loadMsg, client.RoomID)
		ws.Close()
	}()

	room.AddClient(client)
	err = h.s.JoinGame(context.Background(), client.GameId, client.UserID)
	if err != nil {
		log.Println("Join Game Error:", err)
		return
	}

	// ดึง role ของผู้เล่นนี้มาเก็บไว้ใน Client
	roleResult, err := h.s.GetRoleGame(context.Background(), client.RoomID, client.UserID)
	if err == nil {
		if roleStr, ok := roleResult.(string); ok {
			client.Role = roleStr
		}
	}

	// เก็บ gameId + service ไว้ใน Room เพื่อให้ NextPhase เรียก SummaryVote ได้
	room.mu.Lock()
	room.GameId = client.GameId
	room.service = h.s
	room.mu.Unlock()

	loadMsg, _ := json.Marshal(Message{Type: "load_game"})
	Manager.BroadcastToRoom(loadMsg, client.RoomID)

	go func() {
		for m := range client.Send {
			if err := ws.WriteMessage(websocket.TextMessage, m); err != nil {
				return
			}
		}
	}()

	room.mu.RLock()
	if room.IsGameStarted {
		syncMsg := map[string]interface{}{
			"type":      "sync_time",
			"phase":     room.Phase,
			"remaining": room.RemainingTime,
		}
		payload, _ := json.Marshal(syncMsg)
		client.Send <- payload
	} else {
		stargame, err := h.s.StartGame(context.Background(), client.RoomID)
		if err != nil {
			log.Println("Start Game Error:", err)
			room.mu.RUnlock()
			return
		}
		if stargame {
			err := h.s.UpdateStartGame(context.Background(), client.RoomID)
			if err != nil {
				log.Println("Update Start Game Error:", err)
				room.mu.RUnlock()
				return
			}
			loadMsg, _ := json.Marshal(Message{Type: "start_game"})
			Manager.BroadcastToRoom(loadMsg, client.RoomID)

			loadMsg, _ = json.Marshal(Message{Type: "status_room"})
			Manager.BroadcastToRoom(loadMsg, client.RoomID)
		}
	}
	room.mu.RUnlock()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		var msgData Message
		if err := json.Unmarshal(msg, &msgData); err != nil {
			log.Printf("Failed to unmarshal game message: %v", err)
			continue
		}

		log.Println("msgData: ", msgData)

		// แชทปกติ (เฉพาะคนที่ยังมีชีวิต)
		if msgData.Type == "chat" {
			chatMsg, _ := json.Marshal(MessageGame{Type: "chat", Content: msgData.Content, Sender: client.UserName, Timestamp: time.Now().Format("15:04:05")})
			Manager.BroadcastToRoom(chatMsg, client.RoomID)
		}

		// แชทคนตาย (เห็นเฉพาะคนตาย)
		if msgData.Type == "dead_chat" {
			chatMsg, _ := json.Marshal(MessageGame{Type: "dead_chat", Content: msgData.Content, Sender: client.UserName, Timestamp: time.Now().Format("15:04:05")})
			room.BroadcastToDeadPlayers(chatMsg)
		}

		// แชทหมาป่า (เห็นเฉพาะหมาป่า)
		if msgData.Type == "wolf_chat" {
			chatMsg, _ := json.Marshal(MessageGame{Type: "wolf_chat", Content: msgData.Content, Sender: client.UserName, Timestamp: time.Now().Format("15:04:05")})
			room.BroadcastToWerewolves(chatMsg)
		}

		// ใช้สำหรับเริ่มเกม (guard: ให้เริ่มได้ครั้งเดียวเท่านั้น)
		if msgData.Type == "start_game" {
			room.mu.Lock()
			alreadyStarted := room.IsGameStarted
			room.mu.Unlock()
			if !alreadyStarted {
				log.Println("[GAME] Starting timer for room:", client.RoomID)
				room.StartTimer(5, "ready_to_start", client.RoomID)
			} else {
				log.Println("[GAME] Timer already started, ignoring duplicate start_game")
			}
		}

		// โหวต — คลิกคนเลยจะส่ง vote ทันที
		if msgData.Type == "vote" {
			log.Printf("[VOTE] User=%s, GameId=%s, Target=%s", client.UserID, client.GameId, msgData.Content)
			_, err := h.s.Vote(context.Background(), client.GameId, client.UserID, msgData.Content, "")
			if err != nil {
				log.Println("Vote Error:", err)
			}
		}

		// ยกเลิก/ข้ามโหวต
		if msgData.Type == "cancel_vote" {
			err := h.s.CancelVote(context.Background(), client.GameId, client.UserID)
			if err != nil {
				log.Println("Cancel Vote Error:", err)
			}
		}

		// Seer ตรวจสอบ role (ใช้ได้ครั้งเดียวต่อคืน)
		if msgData.Type == "seer_check" {
			if client.SeerUsed {
				log.Println("[GAME] Seer already used this night, ignoring")
				continue
			}
			result, err := h.s.SeerCheck(context.Background(), client.GameId, msgData.Content)
			if err != nil {
				log.Println("Seer Check Error:", err)
				continue
			}
			client.SeerUsed = true
			seerMsg := map[string]interface{}{
				"type":    "seer_result",
				"content": result,
				"target":  msgData.Content,
			}
			payload, _ := json.Marshal(seerMsg)
			sendToClient(client, payload)
		}

		// Guard ปกป้อง
		if msgData.Type == "guard_protect" {
			_, err := h.s.Vote(context.Background(), client.GameId, client.UserID, msgData.Content, "guard")
			if err != nil {
				log.Println("Guard Protect Error:", err)
			}
		}

		// NOTE: "summary" type is no longer handled here.
		// Summary is called directly from NextPhase on the server side to prevent duplicates.
	}
}
