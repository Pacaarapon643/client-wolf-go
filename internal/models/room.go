package models

type Room struct {
	BaseModel
	RoomName           string `json:"room_name"`
	RoomId             string `json:"room_id"`
	RoomStatus         string `json:"room_status"`
	IsEnd              bool   `json:"is_end"`
	CreateBy           string `json:"create_by"`
	TotalPlayer        int    `json:"total_player"`
	TotalPlayerCurrent int    `json:"total_player_current"`
}

func (Room) TableName() string {
	return "rooms"
}
