package util

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestResponseAPI(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return ResponseAPI(c, 200, "Success message", map[string]string{"key": "value"}, "dummy_token")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test fiber app: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if body["status"] != float64(200) {
		t.Errorf("Expected body status 200, got %v", body["status"])
	}

	if body["message"] != "Success message" {
		t.Errorf("Expected body message 'Success message', got %v", body["message"])
	}

	if body["token"] != "dummy_token" {
		t.Errorf("Expected body token 'dummy_token', got %v", body["token"])
	}
}
