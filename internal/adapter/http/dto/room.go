package dto

type CreateRoomRequest struct {
	RoomName    string `json:"room_name"`
	TotalPlayer int    `json:"total_player"`
	CreateBy    string `json:"create_by"`
}

type JoinRoomMemberRequest struct {
	RoomId  string `json:"room_id"`
	UserId  string `json:"user_id"`
	MaxRoom int    `json:"max_room"`
}
