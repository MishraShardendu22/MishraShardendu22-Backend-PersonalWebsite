package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"io"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSkillController(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetSkills - Empty", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/skills", GetSkills)

		// collectSkills calls:
		// 1. projects SimpleFind
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.projects", mtest.FirstBatch))
		// 2. users First
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch))

		req := httptest.NewRequest("GET", "/skills", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetSkills - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/skills", GetSkills)

		// 1. projects SimpleFind
		projOid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.projects", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: projOid},
			{Key: "skills", Value: bson.A{"Go", "Rust"}},
		}))
		// 2. users First
		userOid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "skills", Value: bson.A{"TypeScript", "Go"}},
		}))

		req := httptest.NewRequest("GET", "/skills?page=1&limit=5", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("AddSkills - Missing fields", func(mt *mtest.T) {
		app := fiber.New()
		app.Post("/skills", AddSkills)

		body, _ := json.Marshal(skillsRequest{
			Skills: []string{},
		})
		req := httptest.NewRequest("POST", "/skills", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AddSkills - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Post("/skills", AddSkills)

		userOid := primitive.NewObjectID()
		// 1. users First (in AddSkills)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "skills", Value: bson.A{"Go"}},
		}))
		// 2. users Update (in AddSkills)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		// 3. projects SimpleFind (in collectSkills)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.projects", mtest.FirstBatch))
		// 4. users First (in collectSkills)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "skills", Value: bson.A{"Go", "Rust"}},
		}))

		body, _ := json.Marshal(skillsRequest{
			Skills: []string{"Rust"},
		})
		req := httptest.NewRequest("POST", "/skills", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})
}
