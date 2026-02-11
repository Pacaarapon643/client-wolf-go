package models

import "github.com/google/uuid"

type RoomMember struct {
	BaseModel
	RoomId    string    `json:"room_id"`
	UserId    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	IsReady   bool      `json:"is_ready"`
	SlotIndex int       `json:"slot_index"`
	IsHost    bool      `json:"is_host"`
	IsLeft    bool      `json:"is_left"`
	LeftAt    string    `json:"left_at"`
}

func (RoomMember) TableName() string {
	return "room_members"
}
