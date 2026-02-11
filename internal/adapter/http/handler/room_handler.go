package handler

import (
	"errors"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"
	"werewolf-backend/internal/pkg/appconst"
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h Handler) CreateRoom(ctx *fiber.Ctx) error {
	var req dto.CreateRoomRequest
	if err := ctx.BodyParser(&req); err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			err.Error(),
			"Bad Request",
		)
	}

	obj := models.Room{
		RoomName:    req.RoomName,
		RoomId:      util.GenerateRoomCode(),
		RoomStatus:  appconst.RoomStatusWaiting,
		TotalPlayer: req.TotalPlayer,
		CreateBy:    req.CreateBy,
	}

	err := h.s.CreateRoom(ctx.Context(), &obj)
	if err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusInternalServerError,
			err.Error(),
			"Internal Server Error",
		)
	}

	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		obj,
		"Success",
	)
}

func (h Handler) GetRoom(ctx *fiber.Ctx) error {

	var obj []models.Room

	err := h.s.GetRoom(ctx.Context(), &obj)
	if err != nil {
		var detailedError *util.LocalizedError
		if errors.As(err, &detailedError) {
			return util.HandlerError(
				ctx,
				detailedError.Code,
				detailedError.Err,
				detailedError.Message,
			)
		} else {
			return util.HandlerError(
				ctx,
				fiber.StatusInternalServerError,
				detailedError.Err,
				appconst.InternalServer,
			)
		}
	}

	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		obj,
		"Success",
	)
}

func (h Handler) CountRoom(ctx *fiber.Ctx) error {

	var count int64
	err := h.s.CountRoom(ctx.Context(), &count)
	if err != nil {
		var detailedError *util.LocalizedError
		if errors.As(err, &detailedError) {
			return util.HandlerError(
				ctx,
				detailedError.Code,
				detailedError.Err,
				detailedError.Message,
			)
		} else {
			return util.HandlerError(
				ctx,
				fiber.StatusInternalServerError,
				detailedError.Err,
				appconst.InternalServer,
			)
		}
	}

	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		count,
		"Success",
	)
}

func (h Handler) FindRoom(ctx *fiber.Ctx) error {

	id := ctx.Query("room_id")
	if id == "" {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			"Bad Request",
			"Room ID is required",
		)
	}

	var obj models.Room
	obj.RoomId = id
	err := h.s.FindRoom(ctx.Context(), &obj)
	if err != nil {
		var detailedError *util.LocalizedError
		if errors.As(err, &detailedError) {
			return util.HandlerError(
				ctx,
				detailedError.Code,
				detailedError.Err,
				detailedError.Message,
			)
		} else {
			return util.HandlerError(
				ctx,
				fiber.StatusInternalServerError,
				detailedError.Err,
				appconst.InternalServer,
			)
		}
	}

	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		obj,
		"Success",
	)
}

func (h Handler) JoinRoom(ctx *fiber.Ctx) error {
	id := ctx.Query("room_id")
	if id == "" {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			"Bad Request",
			"Room ID is required",
		)
	}

	err := h.s.JoinRoom(ctx.Context(), id)
	if err != nil {
		var detailedError *util.LocalizedError
		if errors.As(err, &detailedError) {
			return util.HandlerError(
				ctx,
				detailedError.Code,
				detailedError.Err,
				detailedError.Message,
			)
		} else {
			return util.HandlerError(
				ctx,
				fiber.StatusInternalServerError,
				detailedError.Err,
				appconst.InternalServer,
			)
		}
	}

	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		nil,
		"Success",
	)
}

func (h Handler) JoinRoomMember(ctx *fiber.Ctx) error {

	var req dto.JoinRoomMemberRequest
	if err := ctx.BodyParser(&req); err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			err.Error(),
			"Bad Request",
		)
	}

	// แปลง UserId จาก string เป็น UUID เพราะ User.ID เป็น UUID
	// แต่ RoomId ยังเป็น string เพราะ Room.RoomId เป็น string (room code)
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			err.Error(),
			"Invalid user ID format",
		)
	}

	obj := models.RoomMember{
		RoomId: req.RoomId, // string - room code
		UserId: userUUID,   // UUID - user ID
	}

	err = h.s.JoinRoomMember(ctx.Context(), &obj, req.MaxRoom)
	if err != nil {
		var detailedError *util.LocalizedError
		if errors.As(err, &detailedError) {
			return util.HandlerError(
				ctx,
				detailedError.Code,
				detailedError.Err,
				detailedError.Message,
			)
		} else {
			return util.HandlerError(
				ctx,
				fiber.StatusInternalServerError,
				err.Error(),
				appconst.InternalServer,
			)
		}
	}

	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		nil,
		"Success",
	)
}

func (h Handler) GetRoomMember(ctx *fiber.Ctx) error {

	var query dto.QueryRoomMember
	if err := ctx.QueryParser(&query); err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			err.Error(),
			"Bad Request",
		)
	}

	var obj []dto.RoomMemberResponse
	err := h.s.GetRoomMember(ctx.Context(), &obj, query)
	if err != nil {
		var detailedError *util.LocalizedError
		if errors.As(err, &detailedError) {
			return util.HandlerError(
				ctx,
				detailedError.Code,
				detailedError.Err,
				detailedError.Message,
			)
		} else {
			return util.HandlerError(
				ctx,
				fiber.StatusInternalServerError,
				detailedError.Err,
				appconst.InternalServer,
			)
		}
	}

	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		obj,
		"Success",
	)
}

func (h Handler) LeaveRoom(ctx *fiber.Ctx) error {

	var req dto.LeaveRoomRequest
	if err := ctx.BodyParser(&req); err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			err.Error(),
			"Bad Request",
		)
	}

	err := h.s.LeaveRoom(ctx.Context(), req.RoomId, req.UserId)
	if err != nil {
		var detailedError *util.LocalizedError
		if errors.As(err, &detailedError) {
			return util.HandlerError(
				ctx,
				detailedError.Code,
				detailedError.Err,
				detailedError.Message,
			)
		} else {
			return util.HandlerError(
				ctx,
				fiber.StatusInternalServerError,
				err.Error(),
				appconst.InternalServer,
			)
		}
	}

	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		nil,
		"Success",
	)
}
