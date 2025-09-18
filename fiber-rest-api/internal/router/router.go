package router

import (
    "github.com/gofiber/fiber/v2"
    "fiber-rest-api/internal/handlers"
)

func SetupRoutes(app *fiber.App) {
    app.Get("/", handlers.GetRoot)
}