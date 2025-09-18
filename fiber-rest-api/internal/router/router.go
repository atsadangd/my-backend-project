package router

import (
    "github.com/gofiber/fiber/v2"
    "fiber-rest-api/internal/handlers"
)

func SetupRoutes(app *fiber.App) {
    app.Get("/", handlers.GetRoot)

    // auth endpoints
    app.Post("/auth/register", handlers.Register)
    app.Post("/auth/login", handlers.Login)

    // swagger
    app.Get("/docs/swagger.json", handlers.SwaggerJSON)
    app.Get("/docs", handlers.SwaggerUI)
}