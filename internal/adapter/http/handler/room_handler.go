package handler

import (
	"errors"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"
	"werewolf-backend/internal/pkg/appconst"
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
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
