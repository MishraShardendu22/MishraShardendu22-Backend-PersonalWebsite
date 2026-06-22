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

func TestVolunteerController(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetVolunteerExperiences - Empty", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/volunteer", GetVolunteerExperiences)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.volunteer_experiences", mtest.FirstBatch))

		req := httptest.NewRequest("GET", "/volunteer", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetVolunteerExperiences - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/volunteer", GetVolunteerExperiences)

		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.volunteer_experiences", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "organisation", Value: "My Org"},
			{Key: "description", Value: "Helping others"},
		}))

		req := httptest.NewRequest("GET", "/volunteer?page=1&limit=5", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetVolunteerExperienceByID - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/volunteer/:id", GetVolunteerExperienceByID)

		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.volunteer_experiences", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "organisation", Value: "My Org"},
		}))

		req := httptest.NewRequest("GET", "/volunteer/"+oid.Hex(), nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetVolunteerExperienceByID - Invalid ID", func(mt *mtest.T) {
		app := fiber.New()
		app.Get("/volunteer/:id", GetVolunteerExperienceByID)

		req := httptest.NewRequest("GET", "/volunteer/invalid-id", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AddVolunteerExperiences - Missing Fields", func(mt *mtest.T) {
		app := fiber.New()
		app.Post("/volunteer", AddVolunteerExperiences)

		body, _ := json.Marshal(models.VolunteerExperience{
			Organisation: "",
		})
		req := httptest.NewRequest("POST", "/volunteer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AddVolunteerExperiences - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Post("/volunteer", AddVolunteerExperiences)

		// Mock insert volunteer experience
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		// Mock find user
		userOid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "email", Value: "admin@example.com"},
		}))
		// Mock update user
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		body, _ := json.Marshal(models.VolunteerExperience{
			Organisation: "My Org",
			VolunteerTimeLine: []models.VolunteerExperienceTimeLine{
				{PositionOfAuthority: "Volunteer Lead", StartDate: "2020-01-01", EndDate: "2021-01-01"},
			},
		})
		req := httptest.NewRequest("POST", "/volunteer", bytes.NewReader(body))
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

	mt.Run("UpdateVolunteerExperiences - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Put("/volunteer/:id", UpdateVolunteerExperiences)

		oid := primitive.NewObjectID()
		// Mock experience find
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.volunteer_experiences", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "organisation", Value: "My Org"},
			{Key: "volunteer_time_line", Value: bson.A{
				bson.D{
					{Key: "position", Value: "Volunteer"},
					{Key: "start_date", Value: "2020-01-01"},
					{Key: "end_date", Value: "2020-06-01"},
				},
			}},
		}))
		// Mock experience update
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		body, _ := json.Marshal(models.VolunteerExperience{
			Organisation: "My Org Updated",
			VolunteerTimeLine: []models.VolunteerExperienceTimeLine{
				{PositionOfAuthority: "Lead Volunteer", StartDate: "2020-06-01", EndDate: "Present"},
			},
		})
		req := httptest.NewRequest("PUT", "/volunteer/"+oid.Hex(), bytes.NewReader(body))
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

	mt.Run("RemoveVolunteerExperiences - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Delete("/volunteer/:id", RemoveVolunteerExperiences)

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
		// Mock delete volunteer experience
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		req := httptest.NewRequest("DELETE", "/volunteer/"+expOid.Hex(), nil)
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
