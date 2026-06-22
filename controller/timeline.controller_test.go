package controller

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestTimelineController(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("ExperienceTimeline - Empty Experiences", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/timeline", ExperienceTimeline)

		// Mock experiences SimpleFind returning empty
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.experiences", mtest.FirstBatch))

		req := httptest.NewRequest("GET", "/timeline", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("ExperienceTimeline - Success", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/timeline", ExperienceTimeline)

		// 1. Mock experiences SimpleFind
		expOid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.experiences", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expOid},
			{Key: "company_name", Value: "Acme Corp"},
			{Key: "experience_time_line", Value: bson.A{
				bson.D{
					{Key: "position", Value: "SWE"},
					{Key: "start_date", Value: "2020-01-01"},
					{Key: "end_date", Value: "2021-01-01"},
				},
			}},
		}))

		// 2. Mock volunteer experiences SimpleFind
		volOid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.volunteer_experiences", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: volOid},
			{Key: "organisation", Value: "NGO"},
			{Key: "volunteer_time_line", Value: bson.A{
				bson.D{
					{Key: "position", Value: "Volunteer"},
					{Key: "start_date", Value: "2019-01-01"},
					{Key: "end_date", Value: "2019-12-31"},
				},
			}},
		}))

		req := httptest.NewRequest("GET", "/timeline", nil)
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
