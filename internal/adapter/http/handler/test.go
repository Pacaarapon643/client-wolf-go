package handler

import (
	"log"
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) Test123(ctx *fiber.Ctx) error {
	err := h.s.RoleDistribution(ctx.Context(), "34F8XR")
	log.Println("err", err)
	return util.HandlerResponse(
		ctx,
		fiber.StatusInternalServerError,
		nil,
		"Internal Server Error",
	)
}
