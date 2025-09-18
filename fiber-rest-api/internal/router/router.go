package router

import (
	"fiber-rest-api/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", handlers.GetRoot)

	// auth endpoints
	app.Post("/auth/register", handlers.Register)
	app.Post("/auth/login", handlers.Login)

	// profile endpoints (protected)
	app.Get("/profile", handlers.AuthRequired, handlers.GetProfile)
	app.Put("/profile", handlers.AuthRequired, handlers.UpdateProfile)
	// minimal UI to edit profile
	app.Get("/profile/ui", handlers.ProfileUI)

	// swagger
	app.Get("/docs/swagger.json", handlers.SwaggerJSON)
	app.Get("/docs", handlers.SwaggerUI)
}
