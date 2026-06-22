package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestStatsController(t *testing.T) {
	// Preserve the original client and restore it after testing
	origClient := httpClient
	defer func() { httpClient = origClient }()

	t.Run("FetchLeetCodeData - Success", func(t *testing.T) {
		httpClient = &http.Client{
			Transport: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					respData := map[string]interface{}{
						"data": map[string]interface{}{
							"matchedUser": map[string]interface{}{
								"profile": map[string]interface{}{
									"realName": "Shardendu Mishra",
								},
							},
						},
					}
					respBytes, _ := json.Marshal(respData)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}

		app := fiber.New()
		app.Post("/leetcode", FetchLeetCodeData)

		req := httptest.NewRequest("POST", "/leetcode", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("FetchGitHubProfile - Success", func(t *testing.T) {
		httpClient = &http.Client{
			Transport: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					respData := map[string]interface{}{
						"login": "MishraShardendu22",
						"name":  "Shardendu Mishra",
					}
					respBytes, _ := json.Marshal(respData)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}

		app := fiber.New()
		app.Get("/github/profile", FetchGitHubProfile)

		req := httptest.NewRequest("GET", "/github/profile", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("FetchGitHubCommits - Success", func(t *testing.T) {
		httpClient = &http.Client{
			Transport: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					var respBytes []byte
					if req.URL.Path == "/users/MishraShardendu22/repos" {
						repos := []RepoInfo{
							{Name: "repo1", Fork: false},
						}
						respBytes, _ = json.Marshal(repos)
					} else {
						commits := []map[string]interface{}{
							{
								"commit": map[string]interface{}{
									"author": map[string]interface{}{
										"date": "2024-07-01T12:00:00Z",
									},
								},
							},
						}
						respBytes, _ = json.Marshal(commits)
					}
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}

		app := fiber.New()
		app.Get("/github/commits", FetchGitHubCommits)

		req := httptest.NewRequest("GET", "/github/commits", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("FetchGitHubLanguages - Success", func(t *testing.T) {
		httpClient = &http.Client{
			Transport: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					var respBytes []byte
					if req.URL.Path == "/users/MishraShardendu22/repos" {
						repos := []RepoInfo{
							{Name: "repo1", Fork: false, LanguagesURL: "https://api.github.com/repos/MishraShardendu22/repo1/languages"},
						}
						respBytes, _ = json.Marshal(repos)
					} else {
						langs := map[string]int{
							"Go":   1000,
							"Rust": 500,
						}
						respBytes, _ = json.Marshal(langs)
					}
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}

		app := fiber.New()
		app.Get("/github/languages", FetchGitHubLanguages)

		req := httptest.NewRequest("GET", "/github/languages", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("FetchGitHubStars - Success", func(t *testing.T) {
		httpClient = &http.Client{
			Transport: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					repos := []RepoInfo{
						{Name: "repo1", StargazersCount: 5},
						{Name: "repo2", StargazersCount: 10},
					}
					respBytes, _ := json.Marshal(repos)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}

		app := fiber.New()
		app.Get("/github/stars", FetchGitHubStars)

		req := httptest.NewRequest("GET", "/github/stars", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("FetchTopStarredRepos - Success", func(t *testing.T) {
		httpClient = &http.Client{
			Transport: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					repos := []RepoInfo{
						{Name: "repo1", StargazersCount: 5, HTMLURL: "url1", Description: "desc1", Language: "Go"},
						{Name: "repo2", StargazersCount: 10, HTMLURL: "url2", Description: "desc2", Language: "Rust"},
					}
					respBytes, _ := json.Marshal(repos)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}

		app := fiber.New()
		app.Get("/github/top-repos", FetchTopStarredRepos)

		req := httptest.NewRequest("GET", "/github/top-repos", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("FetchContributionCalendar - Success", func(t *testing.T) {
		httpClient = &http.Client{
			Transport: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					respData := map[string]interface{}{
						"total": map[string]interface{}{
							"2024": 1500,
						},
					}
					respBytes, _ := json.Marshal(respData)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}

		app := fiber.New()
		app.Get("/github/calendar", FetchContributionCalendar)

		req := httptest.NewRequest("GET", "/github/calendar", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to run app test: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}
