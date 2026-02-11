package repositories

import (
	"context"
	"errors"
	"log"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"
	"werewolf-backend/internal/pkg/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r Repository) CreateRoom(ctx context.Context, obj *models.Room) error {

	err := r.db.WithContext(ctx).Create(obj).Error
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetRoom(ctx context.Context, obj *[]models.Room) error {

	err := r.db.WithContext(ctx).
		Where("is_end = ?", false).
		Find(&obj).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) CountRoom(ctx context.Context, count *int64) error {

	err := r.db.WithContext(ctx).
		Model(&models.Room{}).
		Where("is_end = ?", false).
		Count(count).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) FindRoom(ctx context.Context, obj *models.Room) error {

	err := r.db.WithContext(ctx).
		Where("room_id = ?", obj.RoomId).
		First(&obj).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &util.LocalizedError{
				Code:    404,
				Err:     "Room not found",
				Message: "ไม่พบห้องที่ค้นหา",
			}
		}
		return err
	}

	return nil
}

func (r Repository) JoinRoom(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).
		Model(models.Room{}).
		Where("room_id = ?", id).
		Update("total_player_current", gorm.Expr("total_player_current + ?", 1)).
		Error
	if err != nil {
		return nil
	}

	return nil
}

func (r Repository) JoinRoomMember(ctx context.Context, obj *models.RoomMember, maxRoom int) error {
	log.Println("JoinRoomMember", obj)
	// transaction
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var usedSlots []int

	err := tx.Model(&models.RoomMember{}).
		Where("room_id = ? AND is_left = ?", obj.RoomId, false).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Pluck("slot_index", &usedSlots).
		Error
	if err != nil {
		return err
	}

	slot := util.FindSmallIndexRoom(usedSlots, maxRoom)
	if slot == -1 {
		tx.Rollback()
		return errors.New("room is full")
	}

	obj.SlotIndex = slot

	if len(usedSlots) == 0 {
		obj.IsHost = true
	}

	// 3. Create member
	err = tx.Create(obj).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

func (r Repository) GetRoomMember(ctx context.Context, obj *[]dto.RoomMemberResponse, roomId string) error {
	err := r.db.WithContext(ctx).
		Table("room_members rm").
		Select("rm.is_ready, rm.slot_index, rm.is_host, u.user_name, u.id as user_id").
		Joins("JOIN users u ON u.id = rm.user_id").
		Where("room_id = ? AND is_left = ?", roomId, false).
		Order("slot_index ASC").
		Find(&obj).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) LeaveRoom(ctx context.Context, roomId string, userId uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Model(&models.RoomMember{}).
		Where("room_id = ? AND user_id = ?", roomId, userId).
		Update("is_left", true).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) LeaveRoomJoin(ctx context.Context, roomId string) error {
	err := r.db.WithContext(ctx).
		Model(&models.Room{}).
		Where("room_id = ?", roomId).
		Update("total_player_current", gorm.Expr("total_player_current - ?", 1)).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ReadyRoom(ctx context.Context, roomId string, userId uuid.UUID, action bool) error {
	err := r.db.WithContext(ctx).
		Model(&models.RoomMember{}).
		Where("room_id = ? AND user_id = ? AND is_left = ?", roomId, userId, false).
		Update("is_ready", action).
		Error
	if err != nil {
		return err
	}
	return nil
}
