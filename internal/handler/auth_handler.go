package handler

import (
	"errors"
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

func (h Handler) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return util.HandlerError(
			ctx,
			fiber.StatusBadRequest,
			err.Error(),
			appconst.ErrorBodyParser,
		)
	}

	data, err := h.s.Login(ctx.Context(), req)
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
		data,
		"login success",
	)
}
