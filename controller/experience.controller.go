package controller

import (
	"github.com/MishraShardendu22/models"
	"github.com/MishraShardendu22/util"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetExperiences(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 15)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 15
	}

	var exps []models.Experience
	if err := mgm.Coll(&models.Experience{}).SimpleFind(&exps, bson.M{}); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch experiences", nil, "")
	}

	if len(exps) == 0 {
		return util.ResponseAPI(c, fiber.StatusOK, "No experiences found", nil, "")
	}

	exps = ReverseExperiences(exps)

	totalExps := len(exps)
	totalPages := (totalExps + limit - 1) / limit
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= totalExps {
		return util.ResponseAPI(c, fiber.StatusOK, "Page out of range", fiber.Map{
			"experiences":  []models.Experience{},
			"page":         page,
			"limit":        limit,
			"total":        totalExps,
			"total_pages":  totalPages,
			"has_next":     false,
			"has_previous": page > 1,
		}, "")
	}

	if endIndex > totalExps {
		endIndex = totalExps
	}

	paginatedExps := exps[startIndex:endIndex]

	return util.ResponseAPI(c, fiber.StatusOK, "Experiences retrieved successfully", fiber.Map{
		"experiences":  paginatedExps,
		"page":         page,
		"limit":        limit,
		"total":        totalExps,
		"total_pages":  totalPages,
		"has_next":     page < totalPages,
		"has_previous": page > 1,
	}, "")
}

func GetExperienceByID(c *fiber.Ctx) error {
	eid := c.Params("id")
	if eid == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Experience ID is required", nil, "")
	}

	expObjID, err := primitive.ObjectIDFromHex(eid)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid experience ID", nil, "")
	}

	var e models.Experience
	if err := mgm.Coll(&models.Experience{}).FindByID(expObjID, &e); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Experience not found", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Experience retrieved successfully", e, "")
}

func AddExperiences(c *fiber.Ctx) error {
	var e models.Experience
	if err := c.BodyParser(&e); err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}

	if e.CompanyName == "" || len(e.ExperienceTimeline) == 0 {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Company name, position and start date are required", nil, "")
	}

	e.Tokens = util.GenerateTokens([]string{e.CompanyName, e.Description}, e.Technologies)

	if err := mgm.Coll(&models.Experience{}).Create(&e); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to add experience", nil, "")
	}

	var user models.User
	if err := mgm.Coll(&models.User{}).First(bson.M{}, &user); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "User not found", nil, "")
	}

	user.Experiences = append(user.Experiences, e.ID)
	if err := mgm.Coll(&models.User{}).Update(&user); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to update user experiences", nil, "")
	}

	InvalidateSearchCache()
	return util.ResponseAPI(c, fiber.StatusOK, "Experience added successfully", e, "")
}

func UpdateExperiences(c *fiber.Ctx) error {
	eid := c.Params("id")
	if eid == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Experience ID is required", nil, "")
	}

	expObjID, err := primitive.ObjectIDFromHex(eid)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid experience ID", nil, "")
	}

	var input models.Experience
	if err := c.BodyParser(&input); err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}

	if input.CompanyName == "" || len(input.ExperienceTimeline) == 0 {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Company name and at least one timeline entry are required", nil, "")
	}

	var existing models.Experience
	if err := mgm.Coll(&models.Experience{}).FindByID(expObjID, &existing); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Experience not found", nil, "")
	}

	existing.ExperienceTimeline = append(existing.ExperienceTimeline, input.ExperienceTimeline...)

	existing.CompanyName = input.CompanyName
	existing.Description = input.Description
	existing.Technologies = input.Technologies
	existing.Projects = input.Projects
	existing.CompanyLogo = input.CompanyLogo
	existing.CertificateURL = input.CertificateURL
	existing.Images = input.Images

	existing.Tokens = util.GenerateTokens([]string{existing.CompanyName, existing.Description}, existing.Technologies)

	if err := mgm.Coll(&models.Experience{}).Update(&existing); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to update experience", nil, "")
	}

	InvalidateSearchCache()
	return util.ResponseAPI(c, fiber.StatusOK, "Experience updated successfully", existing, "")
}

func RemoveExperiences(c *fiber.Ctx) error {
	var user models.User
	if err := mgm.Coll(&models.User{}).First(bson.M{}, &user); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "User not found", nil, "")
	}

	eid := c.Params("id")
	if eid == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Experience ID is required", nil, "")
	}

	var updated []primitive.ObjectID
	found := false
	for _, expID := range user.Experiences {
		if expID.Hex() == eid {
			found = true
			continue
		}
		updated = append(updated, expID)
	}

	if !found {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Experience not found", nil, "")
	}

	user.Experiences = updated
	if err := mgm.Coll(&models.User{}).Update(&user); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to remove experience", nil, "")
	}

	expObjID, err := primitive.ObjectIDFromHex(eid)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid experience ID", nil, "")
	}

	proj := &models.Experience{}
	proj.SetID(expObjID)
	if err := mgm.Coll(proj).Delete(proj); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to delete experience", nil, "")
	}

	InvalidateSearchCache()
	return util.ResponseAPI(c, fiber.StatusOK, "Experience removed successfully", nil, "")
}
