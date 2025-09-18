package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func SwaggerJSON(c *fiber.Ctx) error {
	// OpenAPI 3.0 JSON describing auth and profile endpoints
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
              "schema": { "$ref": "#/components/schemas/AuthRequest" }
            }
          }
        },
        "responses": { "201": { "description": "registered" }, "409": { "description": "email already registered" } }
      }
    },
    "/auth/login": {
      "post": {
        "summary": "Login and receive JWT",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/AuthRequest" }
            }
          }
        },
        "responses": { "200": { "description": "token returned", "content": { "application/json": { "schema": { "$ref": "#/components/schemas/TokenResponse" } } } }, "401": { "description": "invalid credentials" } }
      }
    },
    "/profile": {
      "get": {
        "summary": "Get current user's profile",
        "security": [ { "bearerAuth": [] } ],
        "responses": { "200": { "description": "profile returned", "content": { "application/json": { "schema": { "$ref": "#/components/schemas/Profile" } } } }, "401": { "description": "unauthorized" } }
      },
      "put": {
        "summary": "Update current user's profile",
        "security": [ { "bearerAuth": [] } ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": { "schema": { "$ref": "#/components/schemas/ProfileUpdate" } }
          }
        },
        "responses": { "200": { "description": "updated profile", "content": { "application/json": { "schema": { "$ref": "#/components/schemas/Profile" } } } }, "400": { "description": "validation error" }, "401": { "description": "unauthorized" } }
      }
    },
    "/profile/avatar": {
      "post": {
        "summary": "Upload avatar for current user",
        "security": [ { "bearerAuth": [] } ],
        "requestBody": {
          "required": true,
          "content": {
            "multipart/form-data": {
              "schema": {
                "type": "object",
                "properties": {
                  "avatar": { "type": "string", "format": "binary" }
                },
                "required": ["avatar"]
              }
            }
          }
        },
        "responses": { "200": { "description": "avatar uploaded", "content": { "application/json": { "schema": { "$ref": "#/components/schemas/AvatarResponse" } } } }, "400": { "description": "invalid file" }, "401": { "description": "unauthorized" } }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "bearerAuth": { "type": "http", "scheme": "bearer", "bearerFormat": "JWT" }
    },
    "schemas": {
      "AuthRequest": {
        "type": "object",
        "properties": { "email": { "type": "string" }, "password": { "type": "string" } },
        "required": ["email", "password"]
      },
      "TokenResponse": {
        "type": "object",
        "properties": { "token": { "type": "string" } }
      },
      "Profile": {
        "type": "object",
        "properties": {
          "id": { "type": "integer" },
          "email": { "type": "string" },
          "first_name": { "type": "string" },
          "last_name": { "type": "string" },
          "phone": { "type": "string" },
          "avatar": { "type": "string" }
        }
      },
      "ProfileUpdate": {
        "type": "object",
        "properties": {
          "first_name": { "type": "string" },
          "last_name": { "type": "string" },
          "phone": { "type": "string" }
        }
      },
      "AvatarResponse": {
        "type": "object",
        "properties": { "avatar": { "type": "string" } }
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
