package controller

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSearchController(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Search - Missing Query", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/search", Search)

		req := httptest.NewRequest("GET", "/search", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	mt.Run("Search - Empty Results", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/search", Search)

		// buildDocumentIndex makes 4 parallel SimpleFind calls
		for i := 0; i < 4; i++ {
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.coll", mtest.FirstBatch))
		}

		req := httptest.NewRequest("GET", "/search?q=golang", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		}
	})

	mt.Run("GetSearchSuggestions - Empty", func(mt *mtest.T) {
		mgmClient = mt.Client
		mgmDB = mt.DB
		mgmConfig = &mgm.Config{CtxTimeout: 10 * time.Second}

		app := fiber.New()
		app.Get("/search/suggestions", GetSearchSuggestions)

		// GetSearchSuggestions makes 3 parallel SimpleFind calls
		for i := 0; i < 3; i++ {
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.coll", mtest.FirstBatch))
		}

		req := httptest.NewRequest("GET", "/search/suggestions?q=go", nil)
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
