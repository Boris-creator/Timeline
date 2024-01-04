package middleware

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type VError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

var ValidatorInstance *validator.Validate

func Validate[T any](structure T) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var errorsBag []*VError
		parseError := ctx.BodyParser(&structure)
		if parseError != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).SendString(fmt.Sprintf("%s", parseError))
		}
		err := ValidatorInstance.Struct(structure)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				var el VError
				el.Field = err.Field()
				el.Error = getErrorMessage(err)
				errorsBag = append(errorsBag, &el)
			}
			return ctx.Status(fiber.ErrUnprocessableEntity.Code).JSON(errorsBag)
		}
		return ctx.Next()
	}
}

func getErrorMessage(fieldError validator.FieldError) string {
	messages := map[string]func(...string) string{
		"required": func(s ...string) string {
			return fmt.Sprintf("The field %s is required", s[0])
		},
		"exists": func(s ...string) string {
			return fmt.Sprintf("Such %s does not exist in %s", s[0], s[1])
		},
		"oneof": func(s ...string) string {
			return fmt.Sprintf("The field %s must be one of %s", s[0], strings.Join(strings.Split(s[1], " "), ", "))
		},
	}
	getMessage, ok := messages[fieldError.Tag()]
	if !ok {
		getMessage = func(s ...string) string {
			return fmt.Sprintf("The field %s %s", s[0], fieldError.Tag())
		}
	}
	return getMessage(fieldError.Field(), fieldError.Param())
}
