package controller

import (
	"testing"

	"github.com/MishraShardendu22/models"
)

func TestReverseExperiences(t *testing.T) {
	exp1 := models.Experience{CompanyName: "Company 1"}
	exp2 := models.Experience{CompanyName: "Company 2"}

	exps := []models.Experience{exp1, exp2}
	reversed := ReverseExperiences(exps)

	if len(reversed) != 2 {
		t.Fatalf("expected length 2, got %d", len(reversed))
	}
	if reversed[0].CompanyName != "Company 2" || reversed[1].CompanyName != "Company 1" {
		t.Errorf("reverse failed: got order %s, %s", reversed[0].CompanyName, reversed[1].CompanyName)
	}
}

func TestReverseVolunteerExperiences(t *testing.T) {
	exp1 := models.VolunteerExperience{Organisation: "Org 1"}
	exp2 := models.VolunteerExperience{Organisation: "Org 2"}

	exps := []models.VolunteerExperience{exp1, exp2}
	reversed := ReverseVolunteerExperiences(exps)

	if len(reversed) != 2 {
		t.Fatalf("expected length 2, got %d", len(reversed))
	}
	if reversed[0].Organisation != "Org 2" || reversed[1].Organisation != "Org 1" {
		t.Errorf("reverse failed: got order %s, %s", reversed[0].Organisation, reversed[1].Organisation)
	}
}
