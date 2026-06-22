package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"
	_ "unsafe"

	"github.com/MishraShardendu22/models"
	"github.com/MishraShardendu22/util"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

//go:linkname mgmClient github.com/kamva/mgm/v3.client
var mgmClient *mongo.Client

//go:linkname mgmDB github.com/kamva/mgm/v3.db
var mgmDB *mongo.Database

//go:linkname mgmConfig github.com/kamva/mgm/v3.config
var mgmConfig *mgm.Config

func TestAdminController(t *testing.T) {
	// Initialize mtest
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	// Link mgm client/db/config to the mtest mock client/db
	mgmClient = mt.Client
	mgmDB = mt.DB
	mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

	mt.Run("AdminRegisterAndLogin - Invalid Body", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Post("/admin", func(c *fiber.Ctx) error {
			return AdminRegisterAndLogin(c, "admin123", "jwtsecret")
		})

		req := httptest.NewRequest("POST", "/admin", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AdminRegisterAndLogin - Unauthorized Pass", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Post("/admin", func(c *fiber.Ctx) error {
			return AdminRegisterAndLogin(c, "admin123", "jwtsecret")
		})

		body, _ := json.Marshal(models.User{
			AdminPass: "wrong_pass",
		})
		req := httptest.NewRequest("POST", "/admin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", resp.StatusCode)
		}
	})

	mt.Run("AdminRegisterAndLogin - Missing Fields", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Post("/admin", func(c *fiber.Ctx) error {
			return AdminRegisterAndLogin(c, "admin123", "jwtsecret")
		})

		body, _ := json.Marshal(models.User{
			AdminPass: "admin123",
			Email:     "",
		})
		req := httptest.NewRequest("POST", "/admin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AdminRegisterAndLogin - User Exists Valid Password", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Post("/admin", func(c *fiber.Ctx) error {
			return AdminRegisterAndLogin(c, "admin123", "jwtsecret")
		})

		oid, _ := primitive.ObjectIDFromHex("60f78b1234567890abcdef01")
		hash := util.HashPassword("testpass")
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test_db.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "email", Value: "admin@example.com"},
			{Key: "password", Value: hash},
		}))

		body, _ := json.Marshal(models.User{
			AdminPass: "admin123",
			Email:     "admin@example.com",
			Password:  "testpass",
		})
		req := httptest.NewRequest("POST", "/admin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusAccepted {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 202, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("AdminGet - Unauthorized", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/admin", func(c *fiber.Ctx) error {
			c.Locals("user_id", "") // Avoid nil pointer panic by setting empty string
			return AdminGet(c)
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", resp.StatusCode)
		}
	})

	mt.Run("AdminGet - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/admin", func(c *fiber.Ctx) error {
			c.Locals("user_id", "60f78b1234567890abcdef01")
			return AdminGet(c)
		})

		oid, _ := primitive.ObjectIDFromHex("60f78b1234567890abcdef01")
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test_db.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "email", Value: "admin@example.com"},
		}))

		req := httptest.NewRequest("GET", "/admin", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		var response map[string]interface{}
		json.Unmarshal(bodyBytes, &response)
		data := response["data"].(map[string]interface{})

		if data["email"] != "admin@example.com" {
			t.Errorf("expected email 'admin@example.com', got '%s'", data["email"])
		}
	})
}
