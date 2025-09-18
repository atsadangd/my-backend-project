package handlers

import (
    "github.com/gofiber/fiber/v2"
)

func GetRoot(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "message": "hello world",
    })
}