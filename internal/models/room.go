package models

type Room struct {
	BaseModel
	RoomName           string `json:"room_name"`
	CreateBy           string `json:"create_by"`
	TotalPlayer        int    `json:"total_player"`
	TotalPlayerCurrent int    `json:"total_player_current"`
	IsEnd              bool   `json:"is_end"`
}
