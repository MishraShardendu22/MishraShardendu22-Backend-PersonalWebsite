package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func TestJWTMiddleware(t *testing.T) {
	app := fiber.New()
	secret := "secret123"

	app.Get("/test", JWTMiddleware(secret), func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	t.Run("No Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("No Bearer Prefix", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidToken")
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("Invalid Token Signature", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    "123",
			"email": "test@example.com",
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
		}
	})
}
