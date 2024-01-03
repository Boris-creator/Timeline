package middleware

import (
	"fmt"

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
	messages := map[string]string{
		"required": "is required",
		"exists":   "does not exist",
	}
	message, ok := messages[fieldError.Tag()]
	if !ok {
		message = fieldError.Tag()
	}
	return fmt.Sprintf("The field %s %s", fieldError.Field(), message)
}
