package controller

import (
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/MishraShardendu22/models"
	"github.com/MishraShardendu22/util"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	bm25K1        = 1.2
	bm25B         = 0.75
	indexCacheTTL = 5 * time.Minute
)

// In-memory cache for document index
var (
	cachedIndex     []models.SearchDocument
	cacheTimestamp  time.Time
	cacheMutex      sync.RWMutex
)

type scoredDocument struct {
	score float64
	doc   *models.SearchDocument
}

func calculateBM25Score(
	doc *models.SearchDocument,
	queryTokens []string,
	avgDocLength float64,
	docCount int,
	docFrequencies map[string]int,
	termFreq map[string]int,
) float64 {
	docLength := float64(len(doc.Tokens))
	var score float64

	for _, queryToken := range queryTokens {
		tf := float64(termFreq[queryToken])
		if tf == 0 {
			continue
		}

		df := float64(docFrequencies[queryToken])
		if df == 0 {
			df = 1
		}

		idf := math.Log((float64(docCount)-df+0.5)/(df+0.5) + 1)
		tfComponent := (tf * (bm25K1 + 1)) / (tf + bm25K1*(1-bm25B+bm25B*(docLength/avgDocLength)))
		score += idf * tfComponent
	}

	titleTokens := util.Tokenize(doc.Title)
	for _, qt := range queryTokens {
		for _, tt := range titleTokens {
			if tt == qt {
				score *= 1.5
				break
			}
		}
	}

	for _, qt := range queryTokens {
		for _, skill := range doc.Skills {
			if strings.Contains(strings.ToLower(skill), qt) {
				score *= 1.3
				break
			}
		}
	}

	return score
}

func InvalidateSearchCache() {
	cacheMutex.Lock()
	cachedIndex = nil
	cacheTimestamp = time.Time{}
	cacheMutex.Unlock()
}

func getDocumentIndex() ([]models.SearchDocument, error) {
	cacheMutex.RLock()
	if cachedIndex != nil && time.Since(cacheTimestamp) < indexCacheTTL {
		result := make([]models.SearchDocument, len(cachedIndex))
		copy(result, cachedIndex)
		cacheMutex.RUnlock()
		return result, nil
	}
	cacheMutex.RUnlock()

	documents, err := buildDocumentIndex()
	if err != nil {
		return nil, err
	}

	cacheMutex.Lock()
	cachedIndex = make([]models.SearchDocument, len(documents))
	copy(cachedIndex, documents)
	cacheTimestamp = time.Now()
	cacheMutex.Unlock()

	return documents, nil
}

func buildDocumentIndex() ([]models.SearchDocument, error) {
	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		documents = make([]models.SearchDocument, 0, 100)
	)

	wg.Add(4)

	go func() {
		defer wg.Done()
		var projects []models.Project
		if err := mgm.Coll(&models.Project{}).SimpleFind(&projects, bson.M{}); err != nil {
			return
		}

		localDocs := make([]models.SearchDocument, 0, len(projects))

		for i := range projects {
			p := &projects[i]
			localDocs = append(localDocs, models.SearchDocument{
				ID:          p.ID.Hex(),
				Type:        "project",
				Title:       p.ProjectName,
				Subtitle:    p.SmallDescription,
				Description: p.Description,
				Skills:      p.Skills,
				Tokens:      p.Tokens,
				URL:         "/projects/" + p.ID.Hex(),
			})
		}

		mu.Lock()
		documents = append(documents, localDocs...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		var experiences []models.Experience
		if err := mgm.Coll(&models.Experience{}).SimpleFind(&experiences, bson.M{}); err != nil {
			return
		}

		localDocs := make([]models.SearchDocument, 0, len(experiences))

		for i := range experiences {
			e := &experiences[i]
			subtitle := ""
			if len(e.ExperienceTimeline) > 0 {
				subtitle = e.ExperienceTimeline[0].Position
			}

			localDocs = append(localDocs, models.SearchDocument{
				ID:          e.ID.Hex(),
				Type:        "experience",
				Title:       e.CompanyName,
				Subtitle:    subtitle,
				Description: e.Description,
				Skills:      e.Technologies,
				Tokens:      e.Tokens,
				URL:         "/experiences/" + e.ID.Hex(),
			})
		}

		mu.Lock()
		documents = append(documents, localDocs...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		var certifications []models.CertificationOrAchievements
		if err := mgm.Coll(&models.CertificationOrAchievements{}).SimpleFind(&certifications, bson.M{}); err != nil {
			return
		}

		localDocs := make([]models.SearchDocument, 0, len(certifications))

		for i := range certifications {
			c := &certifications[i]
			localDocs = append(localDocs, models.SearchDocument{
				ID:          c.ID.Hex(),
				Type:        "certificate",
				Title:       c.Title,
				Subtitle:    c.Issuer,
				Description: c.Description,
				Skills:      c.Skills,
				Tokens:      c.Tokens,
				URL:         "/certificates/" + c.ID.Hex(),
			})
		}

		mu.Lock()
		documents = append(documents, localDocs...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		var volunteers []models.VolunteerExperience
		if err := mgm.Coll(&models.VolunteerExperience{}).SimpleFind(&volunteers, bson.M{}); err != nil {
			return
		}

		localDocs := make([]models.SearchDocument, 0, len(volunteers))

		for i := range volunteers {
			v := &volunteers[i]
			subtitle := ""
			if len(v.VolunteerTimeLine) > 0 {
				subtitle = v.VolunteerTimeLine[0].PositionOfAuthority
			}

			localDocs = append(localDocs, models.SearchDocument{
				ID:          v.ID.Hex(),
				Type:        "volunteer",
				Title:       v.Organisation,
				Subtitle:    subtitle,
				Description: v.Description,
				Skills:      v.Technologies,
				Tokens:      v.Tokens,
				URL:         "/volunteer/" + v.ID.Hex(),
			})
		}

		mu.Lock()
		documents = append(documents, localDocs...)
		mu.Unlock()
	}()

	wg.Wait()
	return documents, nil
}

func Search(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if query == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Query parameter 'q' is required", nil, "")
	}

	typeFilter := c.Query("type", "")
	limit := c.QueryInt("limit", 10)
	if limit < 1 || limit > 50 {
		limit = 10
	}

	documents, err := getDocumentIndex()
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to get search index", nil, "")
	}

	if typeFilter != "" {
		filtered := make([]models.SearchDocument, 0, len(documents)/4)
		for i := range documents {
			if documents[i].Type == typeFilter {
				filtered = append(filtered, documents[i])
			}
		}
		documents = filtered
	}

	if len(documents) == 0 {
		return util.ResponseAPI(c, fiber.StatusOK, "Search completed", models.SearchResponse{
			Results:    []models.SearchResult{},
			Query:      query,
			TotalCount: 0,
		}, "")
	}

	// Tokenize query at runtime (document tokens are pre-stored in DB)
	queryTokens := util.Tokenize(query)
	if len(queryTokens) == 0 {
		return util.ResponseAPI(c, fiber.StatusOK, "Search completed", models.SearchResponse{
			Results:    []models.SearchResult{},
			Query:      query,
			TotalCount: 0,
		}, "")
	}

	totalTokens := 0
	for i := range documents {
		totalTokens += len(documents[i].Tokens)
	}
	avgDocLength := float64(totalTokens) / float64(len(documents))

	docFrequencies := make(map[string]int, len(queryTokens))
	for _, token := range queryTokens {
		for i := range documents {
			for _, docToken := range documents[i].Tokens {
				if docToken == token {
					docFrequencies[token]++
					break
				}
			}
		}
	}

	scoredDocs := make([]scoredDocument, 0, len(documents)/2)
	termFreq := make(map[string]int, 32)

	for i := range documents {
		doc := &documents[i]

		for k := range termFreq {
			delete(termFreq, k)
		}
		for _, token := range doc.Tokens {
			termFreq[token]++
		}

		score := calculateBM25Score(doc, queryTokens, avgDocLength, len(documents), docFrequencies, termFreq)
		if score > 0 {
			scoredDocs = append(scoredDocs, scoredDocument{score: score, doc: doc})
		}
	}

	sort.Slice(scoredDocs, func(i, j int) bool {
		return scoredDocs[i].score > scoredDocs[j].score
	})

	resultCount := limit
	if len(scoredDocs) < limit {
		resultCount = len(scoredDocs)
	}
	results := make([]models.SearchResult, 0, resultCount)

	seenIDs := make(map[string]struct{}, resultCount)

	for i := 0; i < len(scoredDocs) && len(results) < resultCount; i++ {
		sd := &scoredDocs[i]

		if _, exists := seenIDs[sd.doc.ID]; exists {
			continue
		}
		seenIDs[sd.doc.ID] = struct{}{}

		description := sd.doc.Description
		if len(description) > 150 {
			description = description[:150] + "..."
		}

		results = append(results, models.SearchResult{
			ID:          sd.doc.ID,
			Type:        sd.doc.Type,
			Title:       sd.doc.Title,
			Subtitle:    sd.doc.Subtitle,
			Description: description,
			Skills:      sd.doc.Skills,
			Score:       math.Round(sd.score*100) / 100,
			URL:         sd.doc.URL,
		})
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Search completed", models.SearchResponse{
		Results:    results,
		Query:      query,
		TotalCount: len(results),
	}, "")
}

func GetSearchSuggestions(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if len(query) < 2 {
		return util.ResponseAPI(c, fiber.StatusOK, "Suggestions", fiber.Map{
			"suggestions": []string{},
		}, "")
	}

	query = strings.ToLower(query)

	var (
		wg            sync.WaitGroup
		mu            sync.Mutex
		suggestionSet = make(map[string]struct{}, 32)
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		var projects []models.Project
		if err := mgm.Coll(&models.Project{}).SimpleFind(&projects, bson.M{}); err != nil {
			return
		}

		localSuggestions := make([]string, 0, 16)
		for i := range projects {
			p := &projects[i]
			for _, skill := range p.Skills {
				if strings.Contains(strings.ToLower(skill), query) {
					localSuggestions = append(localSuggestions, skill)
				}
			}
			if strings.Contains(strings.ToLower(p.ProjectName), query) {
				localSuggestions = append(localSuggestions, p.ProjectName)
			}
		}

		mu.Lock()
		for _, s := range localSuggestions {
			suggestionSet[s] = struct{}{}
		}
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		var experiences []models.Experience
		if err := mgm.Coll(&models.Experience{}).SimpleFind(&experiences, bson.M{}); err != nil {
			return
		}

		localSuggestions := make([]string, 0, 16)
		for i := range experiences {
			e := &experiences[i]
			for _, tech := range e.Technologies {
				if strings.Contains(strings.ToLower(tech), query) {
					localSuggestions = append(localSuggestions, tech)
				}
			}
			if strings.Contains(strings.ToLower(e.CompanyName), query) {
				localSuggestions = append(localSuggestions, e.CompanyName)
			}
		}

		mu.Lock()
		for _, s := range localSuggestions {
			suggestionSet[s] = struct{}{}
		}
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		var certifications []models.CertificationOrAchievements
		if err := mgm.Coll(&models.CertificationOrAchievements{}).SimpleFind(&certifications, bson.M{}); err != nil {
			return
		}

		localSuggestions := make([]string, 0, 16)
		for i := range certifications {
			cert := &certifications[i]
			for _, skill := range cert.Skills {
				if strings.Contains(strings.ToLower(skill), query) {
					localSuggestions = append(localSuggestions, skill)
				}
			}
			if strings.Contains(strings.ToLower(cert.Title), query) {
				localSuggestions = append(localSuggestions, cert.Title)
			}
		}

		mu.Lock()
		for _, s := range localSuggestions {
			suggestionSet[s] = struct{}{}
		}
		mu.Unlock()
	}()

	wg.Wait()

	suggestions := make([]string, 0, len(suggestionSet))
	for s := range suggestionSet {
		suggestions = append(suggestions, s)
	}

	sort.Strings(suggestions)
	if len(suggestions) > 8 {
		suggestions = suggestions[:8]
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Suggestions", fiber.Map{
		"suggestions": suggestions,
	}, "")
}
