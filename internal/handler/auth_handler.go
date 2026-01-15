package handler

import (
	"werewolf-backend/internal/appconst"
	"werewolf-backend/internal/dto"
	"werewolf-backend/internal/util"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) Register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			err.Error(),
			appconst.ErrorBodyParser,
		)
	}

	if errors := util.ValidateStruct(req); errors != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			"",
			errors[0]["message"],
		)
	}

	err := h.s.Register(ctx.Context(), req)
	if err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusInternalServerError,
			err.Error(),
			"ไม่สามารถสมัครสมาชิกได้",
		)
	}

	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		nil,
		"register success",
	)
}
