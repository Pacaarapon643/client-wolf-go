package handler

import (
	"errors"
	"werewolf-backend/internal/pkg/appconst"
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CountUser(ctx *fiber.Ctx) error {
	var count int64
	err := h.s.CountUser(ctx.UserContext(), &count)
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
