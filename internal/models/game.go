package models

type Game struct {
	BaseModel
	GameId    string `json:"game_id"`
	RoomId    string `json:"room_id"`
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	Role      string `json:"role"`
	SlotIndex int    `json:"slot_index"`
	IsDead    bool   `json:"is_dead"`
	Img       string `json:"img"`
	IsJoin    bool   `json:"is_join"`
}

func (Game) TableName() string {
	return "games"
}
