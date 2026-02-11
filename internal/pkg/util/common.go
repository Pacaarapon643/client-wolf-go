package util

import (
	"math/rand"
	"time"
)

func GenerateRoomCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func FindSmallIndexRoom(usedSlots []int, maxSlots int) int {
	slotMap := make(map[int]bool)
	for _, slot := range usedSlots {
		slotMap[slot] = true
	}

	for i := 0; i < maxSlots; i++ {
		if !slotMap[i] {
			return i
		}
	}

	return -1
}
