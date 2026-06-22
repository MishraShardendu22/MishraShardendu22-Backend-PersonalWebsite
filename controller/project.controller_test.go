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

func TestProjectController(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetProjects - Empty", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/projects", GetProjects)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.projects", mtest.FirstBatch))

		req := httptest.NewRequest("GET", "/projects", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetProjects - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/projects", GetProjects)

		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.projects", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "project_name", Value: "My Project"},
			{Key: "small_description", Value: "Small description"},
			{Key: "description", Value: "Description"},
			{Key: "order", Value: 1},
		}))

		req := httptest.NewRequest("GET", "/projects?page=1&limit=5", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetProjectByID - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/projects/:id", GetProjectByID)

		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.projects", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "project_name", Value: "My Project"},
		}))

		req := httptest.NewRequest("GET", "/projects/"+oid.Hex(), nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetProjectByID - Invalid ID", func(mt *mtest.T) {
		app := fiber.New()
		app.Get("/projects/:id", GetProjectByID)

		req := httptest.NewRequest("GET", "/projects/invalid-id", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AddProjects - Missing Fields", func(mt *mtest.T) {
		app := fiber.New()
		app.Post("/projects", AddProjects)

		body, _ := json.Marshal(models.Project{
			ProjectName: "",
		})
		req := httptest.NewRequest("POST", "/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AddProjects - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Post("/projects", AddProjects)

		// Mock insert project
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		// Mock find user
		userOid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "email", Value: "admin@example.com"},
		}))
		// Mock update user
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		body, _ := json.Marshal(models.Project{
			ProjectName:      "My Project",
			SmallDescription: "Small description",
			Description:      "Full description",
		})
		req := httptest.NewRequest("POST", "/projects", bytes.NewReader(body))
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

	mt.Run("UpdateProjects - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Put("/projects/:id", UpdateProjects)

		oid := primitive.NewObjectID()
		// Mock project updateByID
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		body, _ := json.Marshal(models.Project{
			ProjectName:      "My Project Updated",
			SmallDescription: "Small description updated",
			Description:      "Full description updated",
		})
		req := httptest.NewRequest("PUT", "/projects/"+oid.Hex(), bytes.NewReader(body))
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

	mt.Run("RemoveProjects - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Delete("/projects/:id", RemoveProjects)

		projOid := primitive.NewObjectID()
		userOid := primitive.NewObjectID()

		// Mock find user
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "email", Value: "admin@example.com"},
			{Key: "projects", Value: bson.A{projOid}},
		}))
		// Mock update user
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		// Mock delete project
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		req := httptest.NewRequest("DELETE", "/projects/"+projOid.Hex(), nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("UpdateProjectOrderKanban - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Put("/projects/order/kanban", UpdateProjectOrderKanban)

		// Mock UpdateOne
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		projID := primitive.NewObjectID()
		body, _ := json.Marshal([]models.UpdatedProject{
			{ProjectID: projID, Order: 2},
		})
		req := httptest.NewRequest("PUT", "/projects/order/kanban", bytes.NewReader(body))
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

	mt.Run("GetProjectsKanban - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/projects/kanban", GetProjectsKanban)

		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.projects", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "project_name", Value: "My Project"},
			{Key: "order", Value: 1},
		}))

		req := httptest.NewRequest("GET", "/projects/kanban", nil)
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
