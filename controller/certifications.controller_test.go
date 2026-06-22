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

func TestCertificationsController(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetCertifications - Empty", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/certifications", GetCertifications)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.certification_or_achievements", mtest.FirstBatch))

		req := httptest.NewRequest("GET", "/certifications", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetCertifications - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/certifications", GetCertifications)

		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.certification_or_achievements", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "title", Value: "AWS Certified Solutions Architect"},
			{Key: "description", Value: "Cloud architect certification"},
			{Key: "issuer", Value: "AWS"},
		}))

		req := httptest.NewRequest("GET", "/certifications?page=1&limit=5", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetCertificationByID - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/certifications/:id", GetCertificationByID)

		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.certification_or_achievements", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: oid},
			{Key: "title", Value: "AWS Certified Solutions Architect"},
			{Key: "description", Value: "Cloud architect certification"},
			{Key: "issuer", Value: "AWS"},
		}))

		req := httptest.NewRequest("GET", "/certifications/"+oid.Hex(), nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetCertificationByID - Invalid ID", func(mt *mtest.T) {
		app := fiber.New()
		app.Get("/certifications/:id", GetCertificationByID)

		req := httptest.NewRequest("GET", "/certifications/invalid-id", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AddCertification - Missing Fields", func(mt *mtest.T) {
		app := fiber.New()
		app.Post("/certifications", AddCertification)

		body, _ := json.Marshal(models.CertificationOrAchievements{
			Title: "",
		})
		req := httptest.NewRequest("POST", "/certifications", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("AddCertification - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Post("/certifications", AddCertification)

		// Mock the insertion of the certification
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		// Mock the first user find
		userOid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "email", Value: "admin@example.com"},
		}))
		// Mock the user update
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		body, _ := json.Marshal(models.CertificationOrAchievements{
			Title:       "AWS Certified Solutions Architect",
			Description: "Cloud architect certification",
			Issuer:      "AWS",
		})
		req := httptest.NewRequest("POST", "/certifications", bytes.NewReader(body))
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

	mt.Run("UpdateCertification - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Put("/certifications/:id", UpdateCertification)

		// Mock the update operation
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		oid := primitive.NewObjectID()
		body, _ := json.Marshal(models.CertificationOrAchievements{
			Title:       "Updated Title",
			Description: "Updated Description",
			Issuer:      "Updated Issuer",
		})
		req := httptest.NewRequest("PUT", "/certifications/"+oid.Hex(), bytes.NewReader(body))
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

	mt.Run("RemoveCertification - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Delete("/certifications/:id", RemoveCertification)

		certOid := primitive.NewObjectID()
		userOid := primitive.NewObjectID()

		// Mock User Find
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: userOid},
			{Key: "email", Value: "admin@example.com"},
			{Key: "certifications", Value: bson.A{certOid}},
		}))
		// Mock User Update
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		// Mock Cert Delete
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		req := httptest.NewRequest("DELETE", "/certifications/"+certOid.Hex(), nil)
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
