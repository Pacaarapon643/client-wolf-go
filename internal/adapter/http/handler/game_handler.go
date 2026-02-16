package handler

import (
	"errors"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/pkg/appconst"
	"werewolf-backend/internal/pkg/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (h Handler) GetRoleGame(ctx *fiber.Ctx) error {

	var query dto.QueryRoleGame
	if err := ctx.QueryParser(&query); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(query); err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			err.Error(),
			"ข้อมูลไม่ครบตามความต้องการ",
		)
	}

	role, err := h.s.GetRoleGame(ctx.Context(), query.RoomId, query.UserId)
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
		role,
		"success",
	)

}

func (h Handler) GetGame(ctx *fiber.Ctx) error {
	gameId := ctx.Query("game_id")
	if gameId == "" {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			"",
			"game_id is required",
		)
	}

	game, err := h.s.GetGame(ctx.Context(), gameId)
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
		game,
		"success",
	)
}