package util

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  int         `json:"status"`
	Result  string      `json:"result"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type Error struct {
	Status  int    `json:"status"`
	Result  string `json:"result"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func HandlerResponse(ctx *fiber.Ctx, statusCode int, data interface{}, msg string) error {
	response := Response{
		Status:  statusCode,
		Result:  http.StatusText(statusCode),
		Data:    data,
		Message: msg,
	}

	return ctx.Status(statusCode).JSON(response)
}

func HandlerError(ctx *fiber.Ctx, statusCode int, err string, msg string) error {
	log.Println("err:", err)
	response := Error{
		Status:  statusCode,
		Result:  http.StatusText(statusCode),
		Error:   err,
		Message: msg,
	}

	return ctx.Status(statusCode).JSON(response)
}

func HandlerErrorBodyParser(ctx *fiber.Ctx, statusCode int, data interface{}, msg string) error {
	response := Response{
		Status:  statusCode,
		Result:  http.StatusText(statusCode),
		Data:    data,
		Message: msg,
	}

	return ctx.Status(statusCode).JSON(response)
}
