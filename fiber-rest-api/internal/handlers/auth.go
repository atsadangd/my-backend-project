package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"fiber-rest-api/internal/db"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	_, err = db.DB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", req.Email, string(hash))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "registered"})
}

func Login(c *fiber.Ctx) error {
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
	}

	var id int
	var hashed string
	row := db.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", req.Email)
	switch err := row.Scan(&id, &hashed); err {
	case sql.ErrNoRows:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	case nil:
		// ok
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to query user"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// create JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%d", id),
		"email": req.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to sign token"})
	}

	return c.JSON(fiber.Map{"token": signed})
}

// AuthRequired is middleware that validates JWT and sets the user_id in locals.
func AuthRequired(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header"})
	}
	tokenStr := parts[1]
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}
	subVal, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid subject claim"})
	}
	uid, err := strconv.Atoi(subVal)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id in token"})
	}
	c.Locals("user_id", uid)
	return c.Next()
}

// GetProfile returns the current user's profile information.
func GetProfile(c *fiber.Ctx) error {
	uidRaw := c.Locals("user_id")
	if uidRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	uid, ok := uidRaw.(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user id"})
	}

	var email, firstName, lastName, phone string
	row := db.DB.QueryRow("SELECT email, first_name, last_name, phone FROM users WHERE id = ?", uid)
	switch err := row.Scan(&email, &firstName, &lastName, &phone); err {
	case sql.ErrNoRows:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	case nil:
		return c.JSON(fiber.Map{
			"id":         uid,
			"email":      email,
			"first_name": firstName,
			"last_name":  lastName,
			"phone":      phone,
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch user"})
	}
}

// ProfileUpdate represents allowed profile fields to update.
type ProfileUpdate struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// UpdateProfile updates the current user's profile (first name, last name, phone).
func UpdateProfile(c *fiber.Ctx) error {
	uidRaw := c.Locals("user_id")
	if uidRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	uid, ok := uidRaw.(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user id"})
	}

	var req ProfileUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	_, err := db.DB.Exec("UPDATE users SET first_name = ?, last_name = ?, phone = ? WHERE id = ?", strings.TrimSpace(req.FirstName), strings.TrimSpace(req.LastName), strings.TrimSpace(req.Phone), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update profile"})
	}

	// Return updated profile
	return GetProfile(c)
}

// ProfileUI serves a minimal HTML page that lets a user view and edit their profile using the API.
func ProfileUI(c *fiber.Ctx) error {
	html := `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Edit Profile</title>
    <style>
      body { font-family: Arial, sans-serif; max-width: 600px; margin: 2rem auto; }
      label { display:block; margin-top: 0.5rem; }
      input { width:100%; padding:0.5rem; }
      button { margin-top:1rem; padding:0.6rem 1rem; }
      .token { margin-bottom:1rem; }
    </style>
  </head>
  <body>
    <h1>แก้ไขโปรไฟล์</h1>
    <p>กรอก JWT token (หลังจากเข้าสู่ระบบ) เพื่อเรียกดูและแก้ไขข้อมูล</p>
    <div class="token">
      <label>Token (Bearer)</label>
      <input id="token" placeholder="paste token here" />
      <button id="load">Load profile</button>
    </div>
    <form id="profile" onsubmit="return false;">
      <label>ชื่อ (First name)</label>
      <input id="first_name" />
      <label>นามสกุล (Last name)</label>
      <input id="last_name" />
      <label>เบอร์โทร (Phone)</label>
      <input id="phone" />
      <label>อีเมล (Email - read only)</label>
      <input id="email" readonly />
      <button id="save">Save</button>
    </form>
    <script>
      const loadBtn = document.getElementById('load')
      const saveBtn = document.getElementById('save')
      const tokenInput = document.getElementById('token')

      async function api(path, method='GET', body) {
        const token = tokenInput.value.trim()
        if (!token) { alert('Please provide token'); return null }
        const res = await fetch(path, {
          method,
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
          },
          body: body ? JSON.stringify(body) : undefined
        })
        if (!res.ok) {
          const e = await res.json().catch(()=>({error:'failed'}))
          alert('Error: ' + (e.error || res.status))
          return null
        }
        return await res.json()
      }

      loadBtn.onclick = async () => {
        const data = await api('/profile')
        if (!data) return
        document.getElementById('first_name').value = data.first_name || ''
        document.getElementById('last_name').value = data.last_name || ''
        document.getElementById('phone').value = data.phone || ''
        document.getElementById('email').value = data.email || ''
      }

      saveBtn.onclick = async () => {
        const body = {
          first_name: document.getElementById('first_name').value,
          last_name: document.getElementById('last_name').value,
          phone: document.getElementById('phone').value
        }
        const data = await api('/profile', 'PUT', body)
        if (data) alert('Saved')
      }
    </script>
  </body>
</html>`

	c.Type("html")
	return c.SendString(html)
}
