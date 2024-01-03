package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func IsAuthorised(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userData, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("user", userData)

	return c.Next()
}
