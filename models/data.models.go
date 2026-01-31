// models/models.go
package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Memory alignment order (64-bit): []T (24) > string (16) > ObjectID (12) > int/float64 (8) > int32 (4) > bool (1)

type User struct {
	mgm.DefaultModel `bson:",inline" json:"inline"`
	Projects         []primitive.ObjectID `bson:"projects" json:"projects"`
	Experiences      []primitive.ObjectID `bson:"experiences" json:"experiences"`
	Certifications   []primitive.ObjectID `bson:"certifications" json:"certifications"`
	Skills           []string             `bson:"skills" json:"skills"`
	Email            string               `bson:"email" json:"email"`
	Password         string               `bson:"password" json:"password"`
	AdminPass        string               `bson:"admin_pass" json:"admin_pass"`
}

type Project struct {
	mgm.DefaultModel  `bson:",inline" json:"inline"`
	Skills            []string `bson:"skills" json:"skills"`
	Tokens            []string `bson:"tokens" json:"-"`
	Description       string   `bson:"description" json:"description"`
	ProjectName       string   `bson:"project_name" json:"project_name"`
	ProjectVideo      string   `bson:"project_video" json:"project_video"`
	ProjectLiveLink   string   `bson:"project_live_link" json:"project_live_link"`
	SmallDescription  string   `bson:"small_description" json:"small_description"`
	ProjectRepository string   `bson:"project_repository" json:"project_repository"`
	Order             int      `bson:"order" json:"order"`
}

type Experience struct {
	mgm.DefaultModel   `bson:",inline" json:"inline"`
	Images             []string             `bson:"images" json:"images"`
	Technologies       []string             `bson:"technologies" json:"technologies"`
	Tokens             []string             `bson:"tokens" json:"-"`
	Projects           []primitive.ObjectID `bson:"projects" json:"projects"`
	ExperienceTimeline []ExperienceTimeLine `bson:"experience_time_line" json:"experience_time_line"`
	CreatedBy          string               `bson:"created_by" json:"created_by"`
	Description        string               `bson:"description" json:"description"`
	CompanyName        string               `bson:"company_name" json:"company_name"`
	CompanyLogo        string               `bson:"company_logo" json:"company_logo"`
	CertificateURL     string               `bson:"certificate_url" json:"certificate_url"`
}

type CertificationOrAchievements struct {
	mgm.DefaultModel `bson:",inline" json:"inline"`
	Projects         []primitive.ObjectID `bson:"projects" json:"projects"`
	Skills           []string             `bson:"skills" json:"skills"`
	Images           []string             `bson:"images" json:"images"`
	Tokens           []string             `bson:"tokens" json:"-"`
	Title            string               `bson:"title" json:"title"`
	Issuer           string               `bson:"issuer" json:"issuer"`
	IssueDate        string               `bson:"issue_date" json:"issue_date"`
	ExpiryDate       string               `bson:"expiry_date" json:"expiry_date"`
	Description      string               `bson:"description" json:"description"`
	CertificateURL   string               `bson:"certificate_url" json:"certificate_url"`
}

type VolunteerExperience struct {
	mgm.DefaultModel  `bson:",inline" json:"inline"`
	Images            []string                      `bson:"images" json:"images"`
	Technologies      []string                      `bson:"technologies" json:"technologies"`
	Tokens            []string                      `bson:"tokens" json:"-"`
	Projects          []primitive.ObjectID          `bson:"projects" json:"projects"`
	VolunteerTimeLine []VolunteerExperienceTimeLine `bson:"volunteer_time_line" json:"volunteer_time_line"`
	CreatedBy         string                        `bson:"created_by" json:"created_by"`
	Description       string                        `bson:"description" json:"description"`
	Organisation      string                        `bson:"organisation" json:"organisation"`
	OrganisationLogo  string                        `bson:"organisation_logo" json:"organisation_logo"`
}

type ExperienceTimeLine struct {
	Position  string `bson:"position" json:"position"`
	EndDate   string `bson:"end_date" json:"end_date"`
	StartDate string `bson:"start_date" json:"start_date"`
}

type VolunteerExperienceTimeLine struct {
	PositionOfAuthority string `bson:"position" json:"position"`
	EndDate             string `bson:"end_date" json:"end_date"`
	StartDate           string `bson:"start_date" json:"start_date"`
}

type UpdatedProject struct {
	ProjectID primitive.ObjectID `bson:"project_id" json:"project_id"`
	Order     int                `bson:"order" json:"order"`
}

type ProjectKanban struct {
	ProjectID    primitive.ObjectID `bson:"project_id" json:"project_id"`
	ProjectTitle string             `bson:"project_title" json:"project_title"`
	Order        int                `bson:"order" json:"order"`
}

type SearchResult struct {
	Skills      []string `json:"skills,omitempty"`
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle,omitempty"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Score       float64  `json:"score"`
}

type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	Query      string         `json:"query"`
	TotalCount int            `json:"total_count"`
}

type SearchDocument struct {
	Skills      []string `json:"-"`
	Tokens      []string `json:"-"`
	ID          string   `json:"-"`
	Type        string   `json:"-"`
	Title       string   `json:"-"`
	Subtitle    string   `json:"-"`
	Description string   `json:"-"`
	URL         string   `json:"-"`
}
