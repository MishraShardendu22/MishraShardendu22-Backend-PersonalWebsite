package controller

import (
	"github.com/MishraShardendu22/models"
	"github.com/MishraShardendu22/util"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func AddSkills(c *fiber.Ctx) error {
	var payload struct {
		Skills []string `json:"skills"`
	}
	err := c.BodyParser(&payload)

	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}

	if len(payload.Skills) == 0 {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Skills cannot be empty", nil, "")
	}

	// Since there's only one user, get the first user from the database
	user := &models.User{}
	err = mgm.Coll(user).First(bson.M{}, user)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "User not found", nil, "")
	}

	user.Skills = append(user.Skills, payload.Skills...)
	err = mgm.Coll(user).Update(user)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to update skills", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Skills added successfully", user.Skills, "")
}

func GetSkills(c *fiber.Ctx) error {
	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 15)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 15
	}

	user := &models.User{}
	if err := mgm.Coll(user).First(bson.M{}, user); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "User not found", nil, "")
	}

	if len(user.Projects) == 0 {
		return util.ResponseAPI(c, fiber.StatusOK, "No projects found", nil, "")
	}

	var projects []models.Project
	filter := bson.M{"_id": bson.M{"$in": user.Projects}}
	cursor, err := mgm.Coll(&models.Project{}).Find(mgm.Ctx(), filter)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch projects", nil, "")
	}
	if err := cursor.All(mgm.Ctx(), &projects); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to decode projects", nil, "")
	}

	skillSet := make(map[string]struct{}, 0)
	for _, p := range projects {
		for _, s := range p.Skills {
			skillSet[s] = struct{}{}
		}
	}

	if len(skillSet) == 0 {
		return util.ResponseAPI(c, fiber.StatusOK, "No skills found", nil, "")
	}

	skills := make([]string, 0, len(skillSet))
	for s := range skillSet {
		skills = append(skills, s)
	}

	// Calculate pagination
	totalSkills := len(skills)
	totalPages := (totalSkills + limit - 1) / limit
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= totalSkills {
		return util.ResponseAPI(c, fiber.StatusOK, "Page out of range", fiber.Map{
			"skills":       []string{},
			"page":         page,
			"limit":        limit,
			"total":        totalSkills,
			"total_pages":  totalPages,
			"has_next":     false,
			"has_previous": page > 1,
		}, "")
	}

	if endIndex > totalSkills {
		endIndex = totalSkills
	}

	paginatedSkills := skills[startIndex:endIndex]

	return util.ResponseAPI(c, fiber.StatusOK, "Skills retrieved successfully", fiber.Map{
		"skills":       paginatedSkills,
		"page":         page,
		"limit":        limit,
		"total":        totalSkills,
		"total_pages":  totalPages,
		"has_next":     page < totalPages,
		"has_previous": page > 1,
	}, "")
}
