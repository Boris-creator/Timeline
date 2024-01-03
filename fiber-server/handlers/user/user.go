package user

import (
	"fiber-server/auth"

	"github.com/gofiber/fiber/v2"
)

type RegisterUserRequest struct {
	Login    string `json:"login" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=70,min=4"`
}
type LoginUserRequest struct {
	Login    string `json:"login" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=70,min=4"`
}
type LoginUserResponse struct {
	Token string `json:"token"`
}

func Register(c *fiber.Ctx) error {
	var req = auth.Credentials{}
	c.BodyParser(&req)
	user, err := auth.Register(req)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).Send([]byte(err.Error()))
	}
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var req = auth.Credentials{}
	c.BodyParser(&req)
	user, err := auth.FindUserByCredentials(req)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).Send([]byte(err.Error()))
	}
	token := auth.GenerateToken(user)
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
	})
	return c.JSON(LoginUserResponse{
		Token: token,
	})
}
