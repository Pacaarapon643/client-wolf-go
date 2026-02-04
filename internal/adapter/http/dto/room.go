package dto

type CreateRoomRequest struct {
	RoomName    string `json:"room_name"`
	TotalPlayer int    `json:"total_player"`
	CreateBy    string `json:"create_by"`
}
