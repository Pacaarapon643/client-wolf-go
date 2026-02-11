package handler

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
	RoomID   string `json:"room_id"`
	UserName string `json:"username"`
	Content  string `json:"content"`
}

type Client struct {
	Conn     *websocket.Conn
	RoomID   string
	UserID   string
	UserName string
	Send     chan []byte
}

type Room struct {
	Id     string
	Client map[string]*Client
	mu     sync.RWMutex
}

type RoomManager struct {
	Room map[string]*Room
	mu   sync.RWMutex
}

type Lobby struct {
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

func (m *RoomManager) BroadcastToLobby(payload []byte) {
	// 1. หาห้องที่ชื่อว่า "lobby" ใน map ของเรา
	lobby, ok := m.Room["lobby"]
	if !ok {
		return // ถ้าไม่มีใครอยู่หน้า lobby เลย ก็ไม่ต้องส่ง
	}

	lobby.mu.RLock() // ล็อคไว้กันคนเข้า/ออกตอนกำลังวนลูป
	defer lobby.mu.RUnlock()

	// 2. วนลูปส่งหาทุกคนที่มีชื่ออยู่ในห้อง lobby
	for _, client := range lobby.Client {
		select {
		case client.Send <- payload:
			// ส่งสำเร็จ ข้อมูลจะไหลเข้า Goroutine Writer ของคนนั้นๆ
		default:
			// ถ้าคนนั้นเน็ตช้าจนท่อเต็ม (256) ให้ข้ามไป ไม่รอ (Non-blocking)
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
	// --- 1. เตรียมล้างข้อมูลตอนจบแน่นอน ---
	// (ย้าย RemoveClient และ NotifyLobby มาไว้ที่นี่)

	// --- 2. รับ Message แรกเพื่อ Join Room ---
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

	// --- 3. การจัดการ Cleanup (ย้ายมาวางตรงนี้หลังจากมีตัวแปร client) ---
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

	// --- 4. ถ้าไม่ใช่ห้อง Lobby ให้จัดการ Logic พิเศษ ---
	if client.RoomID != "lobby" {
		if err := h.s.JoinRoom(context.Background(), client.RoomID); err != nil {
			log.Printf("Failed to join room: %v", err)
		}

		// บอกตัวเองให้โหลดข้อมูลห้องนั้นๆ
		loadMsg, _ := json.Marshal(Message{Type: "load_room"})
		Manager.BroadcastToRoom(loadMsg, client.RoomID)

		// **จุดสำคัญ**: บอก "คนอื่น" ใน Lobby ว่าห้องนี้มีคนเพิ่มแล้ว
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

	// --- 5. ตัวส่งข้อมูล (Writer Loop) ---
	go func() {
		for m := range client.Send {
			if err := ws.WriteMessage(websocket.TextMessage, m); err != nil {
				return
			}
		}
	}()

	// --- 6. ตัวรับข้อมูล (Reader Loop) ค้างไว้จนกว่าจะปิดการเชื่อมต่อ ---
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break // หลุด Loop นี้จะไปทำ defer ข้างบน
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

			loadMsg, _ := json.Marshal(Message{Type: "load_room"})
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
		UserName: joinMsg.UserName,
		Send:     make(chan []byte, 256),
	}

	defer func() {
		room.RemoveClient(client)
		ws.Close()
	}()

	room.AddClient(client)

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
		json.Unmarshal(msg, &msgData)

		log.Println("msgData: ", msgData)
		if msgData.Type == "ready" {
			err := h.s.ReadyRoom(context.Background(), client.RoomID, client.UserID, msgData.Content)
			if err != nil {
				log.Println("Ready Room Error:", err)
			}

			loadMsg, _ := json.Marshal(Message{Type: "load_room"})
			Manager.BroadcastToRoom(loadMsg, client.RoomID)

		}
	}
}
