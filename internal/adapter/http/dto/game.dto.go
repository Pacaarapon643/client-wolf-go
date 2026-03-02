package dto

type QueryRoleGame struct {
	RoomId string `query:"room_id" validate:"required"`
	UserId string `query:"user_id" validate:"required"`
}

type Game struct {
	GameId    string `json:"game_id"`
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	SlotIndex int    `json:"slot_index"`
	IsDead    bool   `json:"is_dead"`
	Img       string `json:"img"`
	IsJoin    bool   `json:"is_join"`
	Role      string `json:"role"`
}

type GameResponse struct {
	GameId     string `json:"game_id"`
	UserId     string `json:"user_id"`
	UserName   string `json:"user_name"`
	SlotIndex  int    `json:"slot_index"`
	IsDead     bool   `json:"is_dead"`
	Img        string `json:"img"`
	IsJoin     bool   `json:"is_join"`
	IsWerewolf bool   `json:"is_werewolf"`
}

type Vote struct {
	UserID string `json:"user_id"`
	Vote   *int   `json:"vote"`
	Role   string `json:"role"`
}

type SummaryVote struct {
	Index int `json:"index"`
	Count int `json:"count"`
}

type GameEvent struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	EventType  string `json:"event_type"` // "death", "saved", "night", "execute", "info"
	NightCount int    `json:"night_count"`
}

type WinResult struct {
	Winner  string `json:"winner"` // "werewolf", "villager", ""
	Message string `json:"message"`
}

type AliveCount struct {
	Werewolf    int `json:"werewolf"`
	NonWerewolf int `json:"non_werewolf"`
}

type SummaryResult struct {
	Result   string `json:"result"` // "dead", "saved", "tie"
	Message  string `json:"message"`
	DeadName string `json:"dead_name"`
	DeadSlot int    `json:"dead_slot"`
}
