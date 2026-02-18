package handler

import (
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) Test123(ctx *fiber.Ctx) error {
	_, err := h.s.RoleDistribution(ctx.Context(), "4J5DHE")
	if err != nil {
		return util.HandlerResponse(
			ctx,
			fiber.StatusInternalServerError,
			err,
			"Internal Server Error",
		)
	}
	return util.HandlerResponse(
		ctx,
		fiber.StatusOK,
		nil,
		"Success",
	)
}
