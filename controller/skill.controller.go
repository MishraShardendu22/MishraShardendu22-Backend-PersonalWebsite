package controller

import (
	"sort"
	"strings"

	"github.com/MishraShardendu22/models"
	"github.com/MishraShardendu22/util"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type skillsRequest struct {
	Skills []string `json:"skills"`
}

type skillsResponse struct {
	Skills      []string `json:"skills"`
	Page        int      `json:"page"`
	Limit       int      `json:"limit"`
	Total       int      `json:"total"`
	TotalPages  int      `json:"total_pages"`
	HasNext     bool     `json:"has_next"`
	HasPrevious bool     `json:"has_previous"`
}

func sanitizePagination(page int, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 100
	}
	return page, limit
}

func mergeSkills(existing []string, incoming []string) []string {
	skillMap := make(map[string]string)

	addSkill := func(value string) {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return
		}
		key := strings.ToLower(trimmed)
		if _, exists := skillMap[key]; !exists {
			skillMap[key] = trimmed
		}
	}

	for _, skill := range existing {
		addSkill(skill)
	}
	for _, skill := range incoming {
		addSkill(skill)
	}

	merged := make([]string, 0, len(skillMap))
	for _, value := range skillMap {
		merged = append(merged, value)
	}

	sort.Slice(merged, func(i, j int) bool {
		return strings.ToLower(merged[i]) < strings.ToLower(merged[j])
	})

	return merged
}

func collectSkills() ([]string, error) {
	skillMap := make(map[string]string)
	addSkill := func(value string) {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return
		}
		key := strings.ToLower(trimmed)
		if _, exists := skillMap[key]; !exists {
			skillMap[key] = trimmed
		}
	}

	var projects []models.Project
	if err := mgm.Coll(&models.Project{}).SimpleFind(&projects, bson.M{}); err != nil {
		return nil, err
	}

	for _, project := range projects {
		for _, skill := range project.Skills {
			addSkill(skill)
		}
	}

	var user models.User
	if err := mgm.Coll(&models.User{}).First(bson.M{}, &user); err == nil {
		for _, skill := range user.Skills {
			addSkill(skill)
		}
	}

	skills := make([]string, 0, len(skillMap))
	for _, value := range skillMap {
		skills = append(skills, value)
	}

	sort.Slice(skills, func(i, j int) bool {
		return strings.ToLower(skills[i]) < strings.ToLower(skills[j])
	})

	return skills, nil
}

func paginateSkills(skills []string, page int, limit int) skillsResponse {
	total := len(skills)
	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	startIndex := (page - 1) * limit
	if startIndex < 0 {
		startIndex = 0
	}
	if startIndex > total {
		startIndex = total
	}

	endIndex := startIndex + limit
	if endIndex > total {
		endIndex = total
	}

	pageSkills := []string{}
	if startIndex < endIndex {
		pageSkills = skills[startIndex:endIndex]
	}

	return skillsResponse{
		Skills:      pageSkills,
		Page:        page,
		Limit:       limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}

func GetSkills(c *fiber.Ctx) error {
	page, limit := sanitizePagination(c.QueryInt("page", 1), c.QueryInt("limit", 100))

	skills, err := collectSkills()
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch skills", nil, "")
	}

	payload := paginateSkills(skills, page, limit)
	message := "Skills retrieved successfully"
	if payload.Total == 0 {
		message = "No skills found"
	}

	return util.ResponseAPI(c, fiber.StatusOK, message, payload, "")
}

func AddSkills(c *fiber.Ctx) error {
	var req skillsRequest
	if err := c.BodyParser(&req); err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}
	if len(req.Skills) == 0 {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "At least one skill is required", nil, "")
	}

	var user models.User
	if err := mgm.Coll(&models.User{}).First(bson.M{}, &user); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "User not found", nil, "")
	}

	user.Skills = mergeSkills(user.Skills, req.Skills)
	if err := mgm.Coll(&models.User{}).Update(&user); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to update skills", nil, "")
	}

	InvalidateSearchCache()

	page, limit := sanitizePagination(c.QueryInt("page", 1), c.QueryInt("limit", 100))
	skills, err := collectSkills()
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch skills", nil, "")
	}

	payload := paginateSkills(skills, page, limit)
	return util.ResponseAPI(c, fiber.StatusOK, "Skills updated successfully", payload, "")
}
