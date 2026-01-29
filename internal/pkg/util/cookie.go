package util

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateCookie(ctx *fiber.Ctx, name string, accessToken *string, times time.Time) {
	ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    *accessToken,
		Expires:  times,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func DeleteCookie(ctx *fiber.Ctx, name string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}
