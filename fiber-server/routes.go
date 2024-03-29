package main

import (
	eventHandlers "fiber-server/handlers/event"
	"fiber-server/handlers/oauth"
	userHandlers "fiber-server/handlers/user"
	"fiber-server/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Static("/", "./public/static", fiber.Static{
		CacheDuration: 0,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})
	eventsApi := app.Group("/api/events")
	eventsApi.Post(
		"/", middleware.Validate(eventHandlers.SaveEvent{}), middleware.IsAuthorised, eventHandlers.StoreEvent,
	)
	eventsApi.Put(
		"/", middleware.Validate(eventHandlers.SaveEvent{}), middleware.IsAuthorised, eventHandlers.UpdateEvent,
	)
	eventsApi.Post("/search", eventHandlers.SearchEvents)
	authApi := app.Group("/api/auth")
	authApi.Post(
		"/register", middleware.Validate(userHandlers.RegisterUserRequest{}), userHandlers.Register,
	)
	authApi.Post(
		"/login", middleware.Validate(userHandlers.LoginUserRequest{}), userHandlers.Login,
	)

	app.Get("/github/login", oauth.GithubLogin)
	app.Get("/github/callback", oauth.GithubCallback)
}
