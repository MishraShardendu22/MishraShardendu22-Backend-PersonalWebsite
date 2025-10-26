package controller

// Template Controller - Renders Templ templates with real data from database and APIs
// Uses existing controller functions and database queries for optimal performance
// Implements caching for API calls to reduce external requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	static "github.com/MishraShardendu22/Static"
	"github.com/MishraShardendu22/models"
	"github.com/MishraShardendu22/templates"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

var templHttpClient = &http.Client{Timeout: 10 * time.Second}

// Cache structure for API responses
type apiCache struct {
	data      interface{}
	timestamp time.Time
	mu        sync.RWMutex
}

var (
	githubCache   = &apiCache{}
	leetcodeCache = &apiCache{}
	cacheDuration = 5 * time.Minute // Cache API responses for 5 minutes
)

// fetchGitHubProfile fetches real GitHub profile data with caching
func fetchGitHubProfile() (map[string]interface{}, error) {
	// Check cache first
	githubCache.mu.RLock()
	if githubCache.data != nil && time.Since(githubCache.timestamp) < cacheDuration {
		data := githubCache.data.(map[string]interface{})
		githubCache.mu.RUnlock()
		return data, nil
	}
	githubCache.mu.RUnlock()

	// Fetch from API
	token := os.Getenv("GITHUB_TOKEN")
	username := static.StaticData.SocialLinks.GitHub.Username

	req, _ := http.NewRequest("GET", "https://api.github.com/users/"+username, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("User-Agent", "portfolio-backend")

	resp, err := templHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	// Update cache
	githubCache.mu.Lock()
	githubCache.data = data
	githubCache.timestamp = time.Now()
	githubCache.mu.Unlock()

	return data, nil
}

// fetchLeetCodeStats fetches real LeetCode statistics with caching
func fetchLeetCodeStats() (map[string]interface{}, error) {
	// Check cache first
	leetcodeCache.mu.RLock()
	if leetcodeCache.data != nil && time.Since(leetcodeCache.timestamp) < cacheDuration {
		data := leetcodeCache.data.(map[string]interface{})
		leetcodeCache.mu.RUnlock()
		return data, nil
	}
	leetcodeCache.mu.RUnlock()

	// Fetch from API
	query := `{
		matchedUser(username: "` + static.StaticData.SocialLinks.LeetCode.Username + `") {
			submitStats {
				acSubmissionNum {
					difficulty
					count
				}
			}
		}
	}`

	payload := map[string]string{"query": query}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com")

	resp, err := templHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jsonResponse map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &jsonResponse)

	// Update cache
	leetcodeCache.mu.Lock()
	leetcodeCache.data = jsonResponse
	leetcodeCache.timestamp = time.Now()
	leetcodeCache.mu.Unlock()

	return jsonResponse, nil
}

// getSkillsFromDB fetches skills from database using existing controller logic
// This reuses the same logic as GetSkills controller for consistency
func getSkillsFromDB() ([]string, error) {
	user := &models.User{}
	if err := mgm.Coll(user).First(bson.M{}, user); err != nil {
		return nil, err
	}

	// If no projects, return user skills directly
	if len(user.Projects) == 0 {
		return user.Skills, nil
	}

	// Fetch all projects to extract skills
	var projects []models.Project
	filter := bson.M{"_id": bson.M{"$in": user.Projects}}
	cursor, err := mgm.Coll(&models.Project{}).Find(mgm.Ctx(), filter)
	if err != nil {
		return user.Skills, nil
	}
	defer cursor.Close(mgm.Ctx())

	if err := cursor.All(mgm.Ctx(), &projects); err != nil {
		return user.Skills, nil
	}

	// Use map to deduplicate skills from all projects
	skillSet := make(map[string]struct{})
	for _, p := range projects {
		for _, s := range p.Skills {
			skillSet[s] = struct{}{}
		}
	}

	// Convert map to slice
	skills := make([]string, 0, len(skillSet))
	for s := range skillSet {
		skills = append(skills, s)
	}

	return skills, nil
}

// parseGitHubStats extracts repos and followers from GitHub API response
func parseGitHubStats(data map[string]interface{}) (repos int, followers int) {
	if data == nil {
		return 0, 0
	}
	if r, ok := data["public_repos"].(float64); ok {
		repos = int(r)
	}
	if f, ok := data["followers"].(float64); ok {
		followers = int(f)
	}
	return repos, followers
}

// parseLeetCodeStats extracts problem counts from LeetCode API response
func parseLeetCodeStats(data map[string]interface{}) (easy, medium, hard int) {
	if data == nil {
		return 0, 0, 0
	}

	respData, ok := data["data"].(map[string]interface{})
	if !ok {
		return 0, 0, 0
	}

	user, ok := respData["matchedUser"].(map[string]interface{})
	if !ok {
		return 0, 0, 0
	}

	stats, ok := user["submitStats"].(map[string]interface{})
	if !ok {
		return 0, 0, 0
	}

	submissions, ok := stats["acSubmissionNum"].([]interface{})
	if !ok {
		return 0, 0, 0
	}

	for _, sub := range submissions {
		s, ok := sub.(map[string]interface{})
		if !ok {
			continue
		}

		difficulty, _ := s["difficulty"].(string)
		count := int(s["count"].(float64))

		switch difficulty {
		case "Easy":
			easy = count
		case "Medium":
			medium = count
		case "Hard":
			hard = count
		}
	}

	return easy, medium, hard
}

// Helper function to render with HTML wrapper
func renderWithWrapper(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	c.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Portfolio Section</title>
	<script src="https://cdn.tailwindcss.com"></script>
	<style>
		@keyframes pulse {
			0%, 100% { opacity: 1; }
			50% { opacity: 0.5; }
		}
		.animate-pulse {
			animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
		}
	</style>
</head>
<body class="bg-gray-950 text-white antialiased">
`))
	component.Render(c.Context(), c.Response().BodyWriter())
	c.Write([]byte(`
</body>
</html>`))
	return nil
}

// GetHeroSection renders the hero section with real data from static config
func GetHeroSection(c *fiber.Ctx) error {
	data := templates.HeroData{
		ProfileImage: static.StaticData.PersonalInfo.Image.Profile,
		Name:         static.StaticData.PersonalInfo.Name,
		FirstName:    static.StaticData.PersonalInfo.GivenName,
		LastName:     static.StaticData.PersonalInfo.FamilyName,
		Email:        static.StaticData.PersonalInfo.Email,
		GithubLink:   static.StaticData.SocialLinks.GitHub.Personal,
		LinkedInLink: static.StaticData.SocialLinks.LinkedIn.URL,
	}

	component := templates.HeroSection(data)
	return renderWithWrapper(c, component)
}

// GetSkillsSection renders the skills section with real data from database
func GetSkillsSection(c *fiber.Ctx) error {
	// Fetch skills from database
	skills, err := getSkillsFromDB()
	if err != nil || len(skills) == 0 {
		// Fallback to empty if database fetch fails
		skills = []string{}
	}

	// Pagination settings
	itemsPerPage := 24 // Show 24 skills per page (4 rows x 6 columns on desktop)
	totalSkills := len(skills)
	totalPages := (totalSkills + itemsPerPage - 1) / itemsPerPage

	// Get current page from query parameter
	page := c.QueryInt("skills_page", 1)
	if page < 1 {
		page = 1
	}
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// Calculate pagination bounds
	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage
	if endIndex > totalSkills {
		endIndex = totalSkills
	}

	// Get skills for current page
	var pageSkills []string
	if totalSkills > 0 {
		pageSkills = skills[startIndex:endIndex]
	} else {
		pageSkills = []string{}
	}

	data := templates.SkillsData{
		StartIndex:  startIndex + 1,
		EndIndex:    endIndex,
		TotalSkills: totalSkills,
		Skills:      pageSkills,
		CurrentPage: page,
		TotalPages:  totalPages,
	}

	component := templates.SkillsSection(data)
	return renderWithWrapper(c, component)
}

// GetEducationSection renders the education section with real data from static config
func GetEducationSection(c *fiber.Ctx) error {
	data := templates.EducationData{
		CollegeName:        static.StaticData.EducationData.College.Name,
		CollegeBatch:       static.StaticData.EducationData.College.Batch,
		CollegeWebsite:     static.StaticData.EducationData.College.Website,
		CollegeLocation:    static.StaticData.EducationData.College.Location,
		SchoolName:         static.StaticData.EducationData.School.Name,
		SchoolBatch:        static.StaticData.EducationData.School.Batch,
		SchoolLocation:     static.StaticData.EducationData.School.Location,
		Class12Percentage:  static.StaticData.EducationData.School.Class12.Percentage,
		Class12Course:      static.StaticData.EducationData.School.Class12.Stream,
		Class10Percentage:  static.StaticData.EducationData.School.Class10.Percentage,
		Languages:          static.StaticData.Languages,
		ResumeViewLink:     static.StaticData.SocialLinks.Resume,
		ResumeDownloadLink: static.StaticData.SocialLinks.Resume,
	}

	component := templates.EducationSection(data)
	return renderWithWrapper(c, component)
}

// GetContactSection renders the contact section with real data from GitHub and LeetCode APIs
func GetContactSection(c *fiber.Ctx) error {
	// Fetch real data from APIs (uses caching)
	githubData, _ := fetchGitHubProfile()
	leetcodeData, _ := fetchLeetCodeStats()

	// Parse stats using helper functions
	githubRepos, githubFollowers := parseGitHubStats(githubData)
	leetcodeEasy, leetcodeMedium, leetcodeHard := parseLeetCodeStats(leetcodeData)

	// Fetch skills from database
	skills, _ := getSkillsFromDB()
	if len(skills) == 0 {
		skills = []string{"Go", "React", "TypeScript", "Next.js", "Node.js", "MongoDB", "PostgreSQL", "Docker", "Kubernetes"}
	}

	data := templates.ContactData{
		GitHubAvatar:    "https://github.com/" + static.StaticData.SocialLinks.GitHub.Username + ".png",
		GitHubUsername:  static.StaticData.SocialLinks.GitHub.Username,
		GitHubBio:       static.StaticData.PersonalInfo.Description,
		GitHubRepos:     githubRepos,
		GitHubFollowers: githubFollowers,
		TotalStars:      0, // Could fetch from another API call if needed
		LeetCodeEasy:    leetcodeEasy,
		LeetCodeMedium:  leetcodeMedium,
		LeetCodeHard:    leetcodeHard,
		LeetCodeTotal:   leetcodeEasy + leetcodeMedium + leetcodeHard,
		RecentCommits:   []templates.ContactCommit{},
		Technologies:    skills,
	}

	component := templates.ContactSection(data)
	return renderWithWrapper(c, component)
}

// GetFooterSection renders the footer section with real data from static config
func GetFooterSection(c *fiber.Ctx) error {
	data := templates.FooterData{
		TechStack:      static.FooterLinks.TechStack,
		CurrentYear:    time.Now().Year(),
		FormAction:     "/api/contact",
		QuickLinks:     convertToTemplateLinks(static.FooterLinks.QuickLinks),
		SocialLinks:    convertToTemplateLinks(static.FooterLinks.SocialLinks),
		CodingProfiles: convertToTemplateLinks(static.FooterLinks.CodingProfiles),
		MoreProfiles:   convertToTemplateLinks(static.FooterLinks.MoreProfiles),
	}

	component := templates.FooterSection(data)
	return renderWithWrapper(c, component)
}

// convertToTemplateLinks converts static.Link to templates.Link
func convertToTemplateLinks(links []static.Link) []templates.Link {
	result := make([]templates.Link, len(links))
	for i, link := range links {
		result[i] = templates.Link{
			Name: link.Name,
			URL:  link.URL,
		}
	}
	return result
}

// GetFullPage renders all sections together with real data
// Uses helper functions and caching for optimal performance
func GetFullPage(c *fiber.Ctx) error {
	// Fetch dynamic data from database
	skills, err := getSkillsFromDB()
	if err != nil || len(skills) == 0 {
		skills = []string{}
	}

	// Fetch API data (uses caching internally)
	githubData, _ := fetchGitHubProfile()
	leetcodeData, _ := fetchLeetCodeStats()

	// Parse API responses using helper functions
	githubRepos, githubFollowers := parseGitHubStats(githubData)
	leetcodeEasy, leetcodeMedium, leetcodeHard := parseLeetCodeStats(leetcodeData)

	// Build hero data from static config
	heroData := templates.HeroData{
		ProfileImage: static.StaticData.PersonalInfo.Image.Profile,
		Name:         static.StaticData.PersonalInfo.Name,
		FirstName:    static.StaticData.PersonalInfo.GivenName,
		LastName:     static.StaticData.PersonalInfo.FamilyName,
		Email:        static.StaticData.PersonalInfo.Email,
		GithubLink:   static.StaticData.SocialLinks.GitHub.Personal,
		LinkedInLink: static.StaticData.SocialLinks.LinkedIn.URL,
	}

	// Build skills data from database with pagination
	itemsPerPage := 24
	totalSkills := len(skills)
	totalPages := (totalSkills + itemsPerPage - 1) / itemsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	page := c.QueryInt("skills_page", 1)
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage
	if endIndex > totalSkills {
		endIndex = totalSkills
	}

	var pageSkills []string
	if totalSkills > 0 {
		pageSkills = skills[startIndex:endIndex]
	} else {
		pageSkills = []string{}
	}

	skillsData := templates.SkillsData{
		StartIndex:  startIndex + 1,
		EndIndex:    endIndex,
		TotalSkills: totalSkills,
		Skills:      pageSkills,
		CurrentPage: page,
		TotalPages:  totalPages,
	}

	// Build education data from static config
	educationData := templates.EducationData{
		CollegeName:        static.StaticData.EducationData.College.Name,
		CollegeBatch:       static.StaticData.EducationData.College.Batch,
		CollegeWebsite:     static.StaticData.EducationData.College.Website,
		CollegeLocation:    static.StaticData.EducationData.College.Location,
		SchoolName:         static.StaticData.EducationData.School.Name,
		SchoolBatch:        static.StaticData.EducationData.School.Batch,
		SchoolLocation:     static.StaticData.EducationData.School.Location,
		Class12Percentage:  static.StaticData.EducationData.School.Class12.Percentage,
		Class12Course:      static.StaticData.EducationData.School.Class12.Stream,
		Class10Percentage:  static.StaticData.EducationData.School.Class10.Percentage,
		Languages:          static.StaticData.Languages,
		ResumeViewLink:     static.StaticData.SocialLinks.Resume,
		ResumeDownloadLink: static.StaticData.SocialLinks.Resume,
	}

	// Build contact data from APIs and database
	contactData := templates.ContactData{
		GitHubAvatar:    "https://github.com/" + static.StaticData.SocialLinks.GitHub.Username + ".png",
		GitHubUsername:  static.StaticData.SocialLinks.GitHub.Username,
		GitHubBio:       static.StaticData.PersonalInfo.Description,
		GitHubRepos:     githubRepos,
		GitHubFollowers: githubFollowers,
		TotalStars:      0,
		LeetCodeEasy:    leetcodeEasy,
		LeetCodeMedium:  leetcodeMedium,
		LeetCodeHard:    leetcodeHard,
		LeetCodeTotal:   leetcodeEasy + leetcodeMedium + leetcodeHard,
		RecentCommits:   []templates.ContactCommit{},
		Technologies:    skills,
	}

	// Build footer data from static config
	footerData := templates.FooterData{
		TechStack:      static.FooterLinks.TechStack,
		CurrentYear:    time.Now().Year(),
		FormAction:     "/api/contact",
		QuickLinks:     convertToTemplateLinks(static.FooterLinks.QuickLinks),
		SocialLinks:    convertToTemplateLinks(static.FooterLinks.SocialLinks),
		CodingProfiles: convertToTemplateLinks(static.FooterLinks.CodingProfiles),
		MoreProfiles:   convertToTemplateLinks(static.FooterLinks.MoreProfiles),
	}

	component := templates.FullPageLayout(heroData, skillsData, educationData, contactData, footerData)
	return adaptor.HTTPHandler(templ.Handler(component))(c)
}
