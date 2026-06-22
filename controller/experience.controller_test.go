package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MishraShardendu22/models"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestExperienceController(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetExperiences - Empty", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/experiences", GetExperiences)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.experiences", mtest.FirstBatch))

		req := httptest.NewRequest("GET", "/experiences", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetExperiences - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/experiences", GetExperiences)

		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.experiences", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "company_name", Value: "Acme Corp"},
			{Key: "description", Value: "Software Engineer role"},
		}))

		req := httptest.NewRequest("GET", "/experiences?page=1&limit=5", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetExperienceByID - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/experiences/:id", GetExperienceByID)

		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.experiences", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "company_name", Value: "Acme Corp"},
		}))

		req := httptest.NewRequest("GET", "/experiences/"+oid.Hex(), nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetExperienceByID - Invalid ID", func(mt *mtest.T) {
		app := fiber.New()
		app.Get("/experiences/:id", GetExperienceByID)

		req := httptest.NewRequest("GET", "/experiences/invalid-id", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AddExperiences - Missing Fields", func(mt *mtest.T) {
		app := fiber.New()
		app.Post("/experiences", AddExperiences)

		body, _ := json.Marshal(models.Experience{
			CompanyName: "",
		})
		req := httptest.NewRequest("POST", "/experiences", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AddExperiences - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Post("/experiences", AddExperiences)

		// Mock insert experience
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		// Mock find user
		userOid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "email", Value: "admin@example.com"},
		}))
		// Mock update user
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		body, _ := json.Marshal(models.Experience{
			CompanyName: "Acme Corp",
			ExperienceTimeline: []models.ExperienceTimeLine{
				{Position: "Software Engineer", StartDate: "2020-01-01", EndDate: "2021-01-01"},
			},
		})
		req := httptest.NewRequest("POST", "/experiences", bytes.NewReader(body))
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

	mt.Run("UpdateExperiences - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Put("/experiences/:id", UpdateExperiences)

		oid := primitive.NewObjectID()
		// Mock experience find
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.experiences", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "company_name", Value: "Acme Corp"},
			{Key: "experience_time_line", Value: bson.A{
				bson.D{
					{Key: "position", Value: "Junior dev"},
					{Key: "start_date", Value: "2020-01-01"},
					{Key: "end_date", Value: "2020-06-01"},
				},
			}},
		}))
		// Mock experience update
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		body, _ := json.Marshal(models.Experience{
			CompanyName: "Acme Corp Updated",
			ExperienceTimeline: []models.ExperienceTimeLine{
				{Position: "Senior Engineer", StartDate: "2020-06-01", EndDate: "Present"},
			},
		})
		req := httptest.NewRequest("PUT", "/experiences/"+oid.Hex(), bytes.NewReader(body))
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

	mt.Run("RemoveExperiences - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Delete("/experiences/:id", RemoveExperiences)

		expOid := primitive.NewObjectID()
		userOid := primitive.NewObjectID()

		// Mock find user
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "email", Value: "admin@example.com"},
			{Key: "experiences", Value: bson.A{expOid}},
		}))
		// Mock update user
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		// Mock delete experience
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		req := httptest.NewRequest("DELETE", "/experiences/"+expOid.Hex(), nil)
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
