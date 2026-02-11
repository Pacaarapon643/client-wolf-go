package dto

type RoomMemberResponse struct {
	IsReady   bool   `json:"is_ready"`
	SlotIndex int    `json:"slot_index"`
	IsHost    bool   `json:"is_host"`
	UserName  string `json:"user_name"`
	UserId    string `json:"user_id"`
	EmptySlot bool   `json:"empty_slot"`
}

type LeaveRoomRequest struct {
	RoomId string `json:"room_id"`
	UserId string `json:"user_id"`
}

type QueryRoomMember struct {
	RoomId  string `query:"room_id"`
	MaxRoom int    `query:"max_room"`
}

