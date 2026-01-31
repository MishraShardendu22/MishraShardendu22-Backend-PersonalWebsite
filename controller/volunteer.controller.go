package controller

import (
	"github.com/MishraShardendu22/models"
	"github.com/MishraShardendu22/util"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetVolunteerExperiences(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 15)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 15
	}

	var exps []models.VolunteerExperience
	if err := mgm.Coll(&models.VolunteerExperience{}).SimpleFind(&exps, bson.M{}); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch volunteer experiences", nil, "")
	}

	if len(exps) == 0 {
		return util.ResponseAPI(c, fiber.StatusOK, "No volunteer experiences found", nil, "")
	}

	exps = ReverseVolunteerExperiences(exps)

	totalExps := len(exps)
	totalPages := (totalExps + limit - 1) / limit
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= totalExps {
		return util.ResponseAPI(c, fiber.StatusOK, "Page out of range", fiber.Map{
			"volunteer_experiences": []models.VolunteerExperience{},
			"page":                  page,
			"limit":                 limit,
			"total":                 totalExps,
			"total_pages":           totalPages,
			"has_next":              false,
			"has_previous":          page > 1,
		}, "")
	}

	if endIndex > totalExps {
		endIndex = totalExps
	}

	paginatedExps := exps[startIndex:endIndex]

	return util.ResponseAPI(c, fiber.StatusOK, "Volunteer experiences retrieved successfully", fiber.Map{
		"volunteer_experiences": paginatedExps,
		"page":                  page,
		"limit":                 limit,
		"total":                 totalExps,
		"total_pages":           totalPages,
		"has_next":              page < totalPages,
		"has_previous":          page > 1,
	}, "")
}

func GetVolunteerExperienceByID(c *fiber.Ctx) error {
	eid := c.Params("id")
	if eid == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Volunteer experience ID is required", nil, "")
	}

	expObjID, err := primitive.ObjectIDFromHex(eid)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid volunteer experience ID", nil, "")
	}

	var e models.VolunteerExperience
	if err := mgm.Coll(&models.VolunteerExperience{}).FindByID(expObjID, &e); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Volunteer experience not found", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Volunteer experience retrieved successfully", e, "")
}

func AddVolunteerExperiences(c *fiber.Ctx) error {
	var e models.VolunteerExperience
	if err := c.BodyParser(&e); err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}

	if e.Organisation == "" || len(e.VolunteerTimeLine) == 0 {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organisation and at least one timeline entry are required", nil, "")
	}

	e.Tokens = util.GenerateTokens([]string{e.Organisation, e.Description}, e.Technologies)

	if err := mgm.Coll(&models.VolunteerExperience{}).Create(&e); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to add volunteer experience", nil, "")
	}

	var user models.User
	if err := mgm.Coll(&models.User{}).First(bson.M{}, &user); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "User not found", nil, "")
	}

	user.Experiences = append(user.Experiences, e.ID)
	if err := mgm.Coll(&models.User{}).Update(&user); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to update user volunteer experiences", nil, "")
	}

	InvalidateSearchCache()
	return util.ResponseAPI(c, fiber.StatusOK, "Volunteer experience added successfully", e, "")
}

func UpdateVolunteerExperiences(c *fiber.Ctx) error {
	eid := c.Params("id")
	if eid == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Volunteer experience ID is required", nil, "")
	}

	expObjID, err := primitive.ObjectIDFromHex(eid)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid volunteer experience ID", nil, "")
	}

	var input models.VolunteerExperience
	if err := c.BodyParser(&input); err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}

	if input.Organisation == "" || len(input.VolunteerTimeLine) == 0 {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organisation and at least one timeline entry are required", nil, "")
	}

	var existing models.VolunteerExperience
	if err := mgm.Coll(&models.VolunteerExperience{}).FindByID(expObjID, &existing); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Volunteer experience not found", nil, "")
	}

	existing.VolunteerTimeLine = append(existing.VolunteerTimeLine, input.VolunteerTimeLine...)

	existing.Organisation = input.Organisation
	existing.Description = input.Description
	existing.Technologies = input.Technologies
	existing.Projects = input.Projects
	existing.OrganisationLogo = input.OrganisationLogo
	existing.Images = input.Images

	existing.Tokens = util.GenerateTokens([]string{existing.Organisation, existing.Description}, existing.Technologies)

	if err := mgm.Coll(&models.VolunteerExperience{}).Update(&existing); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to update volunteer experience", nil, "")
	}

	InvalidateSearchCache()
	return util.ResponseAPI(c, fiber.StatusOK, "Volunteer experience updated successfully", existing, "")
}

func RemoveVolunteerExperiences(c *fiber.Ctx) error {
	var user models.User
	if err := mgm.Coll(&models.User{}).First(bson.M{}, &user); err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "User not found", nil, "")
	}

	eid := c.Params("id")
	if eid == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Volunteer experience ID is required", nil, "")
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
		return util.ResponseAPI(c, fiber.StatusNotFound, "Volunteer experience not found", nil, "")
	}

	user.Experiences = updated
	if err := mgm.Coll(&models.User{}).Update(&user); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to remove volunteer experience", nil, "")
	}

	expObjID, err := primitive.ObjectIDFromHex(eid)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid volunteer experience ID", nil, "")
	}

	proj := &models.VolunteerExperience{}
	proj.SetID(expObjID)
	if err := mgm.Coll(proj).Delete(proj); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to delete volunteer experience", nil, "")
	}

	InvalidateSearchCache()
	return util.ResponseAPI(c, fiber.StatusOK, "Volunteer experience removed successfully", nil, "")
}
