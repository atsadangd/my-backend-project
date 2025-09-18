package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func SwaggerJSON(c *fiber.Ctx) error {
	// Minimal OpenAPI 3.0 JSON describing the auth endpoints
	spec := `{
  "openapi": "3.0.0",
  "info": { "title": "Fiber REST API", "version": "1.0.0" },
  "paths": {
    "/auth/register": {
      "post": {
        "summary": "Register a new user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "email": { "type": "string" },
                  "password": { "type": "string" }
                },
                "required": ["email", "password"]
              }
            }
          }
        },
        "responses": {
          "201": { "description": "registered" },
          "409": { "description": "email already registered" }
        }
      }
    },
    "/auth/login": {
      "post": {
        "summary": "Login and receive JWT",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "email": { "type": "string" },
                  "password": { "type": "string" }
                },
                "required": ["email", "password"]
              }
            }
          }
        },
        "responses": {
          "200": { "description": "token returned" },
          "401": { "description": "invalid credentials" }
        }
      }
    }
  }
}`

	c.Set("Content-Type", "application/json")
	return c.SendString(spec)
}

func SwaggerUI(c *fiber.Ctx) error {
	html := `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
    <script>
      const ui = SwaggerUIBundle({
        url: '/docs/swagger.json',
        dom_id: '#swagger-ui',
      });
    </script>
  </body>
</html>`

	c.Type("html")
	return c.SendString(html)
}
