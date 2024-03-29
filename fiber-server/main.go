package main

import (
	"fiber-server/db"
	"fiber-server/middleware"
	"fiber-server/validators"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	middleware.ValidatorInstance = validator.New()
	middleware.ValidatorInstance.RegisterValidation("exists", validators.Exists)
	middleware.ValidatorInstance.RegisterValidation("unique", validators.Unique)

	db.ConnectDB()
	engine := html.New("./public/views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New())

	RegisterRoutes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
