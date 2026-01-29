package middleware

import (
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

type JWTMiddleware struct {
	jwtService *util.JWTService
}

func NewJWTMiddleware(jwtService *util.JWTService) *JWTMiddleware {
	return &JWTMiddleware{
		jwtService: jwtService,
	}
}

// Protected middleware สำหรับ routes ที่ต้อง authentication
func (m *JWTMiddleware) MiddlewareAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {

		tokenString := c.Cookies("auth")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			return err
		}

		c.Locals("user_id", claims.UserID)
		return c.Next()
	}
}

// RequireRole middleware สำหรับตรวจสอบสิทธิ์ตาม role
func (m *JWTMiddleware) RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Missing role information",
			})
		}
		userRole := role.(string)

		// ตรวจสอบว่า role ตรงกับที่อนุญาตหรือไม่
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}
