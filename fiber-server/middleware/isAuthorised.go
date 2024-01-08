package middleware

import (
	jwtutils "fiber-server/jwtUtils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func IsAuthorised(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")
	token, err := jwtutils.ParseTokenString(tokenString)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userData, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("user", userData)

	return c.Next()
}
