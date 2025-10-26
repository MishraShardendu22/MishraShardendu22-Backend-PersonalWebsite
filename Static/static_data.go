package static

type PersonalInfo struct {
	Name        string
	GivenName   string
	FamilyName  string
	Email       string
	Telephone   string
	JobTitle    string
	Description string
	Bio         Bio
	Image       Image
}

type Bio struct {
	Short string
	Long  string
}

type Image struct {
	Profile       string
	OgImage       string
	FallbackImage string
}

// SocialLinks contains all social media links
type SocialLinks struct {
	Twitter  TwitterLink
	GitHub   GitHubLinks
	LinkedIn LinkedInLink
	YouTube  YouTubeLink
	LeetCode LeetCodeLink
	Resume   string
}

type TwitterLink struct {
	URL      string
	Handle   string
	Username string
}

type GitHubLinks struct {
	Personal     string
	Learning     string
	Organization string
	Username     string
}

type LinkedInLink struct {
	URL      string
	Username string
}

type YouTubeLink struct {
	URL     string
	Channel string
}

type LeetCodeLink struct {
	URL      string
	Username string
}

// EducationData contains static education information
type EducationData struct {
	College College
	School  School
}

type College struct {
	Name     string
	Batch    string
	Website  string
	Location string
	Degree   string
	Field    string
}

type School struct {
	Name     string
	Batch    string
	Location string
	Class12  Class12
	Class10  Class10
}

type Class12 struct {
	Percentage string
	Stream     string
	Year       string
}

type Class10 struct {
	Percentage string
	Year       string
}

// StaticData holds all static configuration
var StaticData = struct {
	PersonalInfo  PersonalInfo
	SocialLinks   SocialLinks
	EducationData EducationData
	Languages     []string
}{
	PersonalInfo: PersonalInfo{
		Name:        "Shardendu Mishra",
		GivenName:   "Shardendu Mishra",
		FamilyName:  "Mishra",
		Email:       "shardendumishra01@gmail.com",
		Telephone:   "+91-8707359576",
		JobTitle:    "Software Developer",
		Description: "Software Developer and Engineer passionate about building impactful applications with modern technologies. Specializing in Go, React, and cloud-native solutions.",
		Bio: Bio{
			Short: "Software Developer and Engineer passionate about building impactful applications with modern technologies.",
			Long:  "Software Developer and Engineer passionate about building impactful applications with modern technologies. Specializing in Go, React, and cloud-native solutions.",
		},
		Image: Image{
			Profile:       "https://github.com/MishraShardendu22.png",
			OgImage:       "https://github.com/MishraShardendu22.png",
			FallbackImage: "https://github.com/MishraShardendu22.png",
		},
	},
	SocialLinks: SocialLinks{
		Twitter: TwitterLink{
			URL:      "https://x.com/Shardendu_M",
			Handle:   "@Shardendu_M",
			Username: "Shardendu_M",
		},
		GitHub: GitHubLinks{
			Personal:     "https://github.com/MishraShardendu22",
			Learning:     "https://github.com/ShardenduMishra22",
			Organization: "https://github.com/Team-Parashuram",
			Username:     "MishraShardendu22",
		},
		LinkedIn: LinkedInLink{
			URL:      "https://www.linkedin.com/in/shardendumishra22",
			Username: "shardendumishra22",
		},
		YouTube: YouTubeLink{
			URL:     "https://www.youtube.com/@Shardendu_Mishra",
			Channel: "@Shardendu_Mishra",
		},
		LeetCode: LeetCodeLink{
			URL:      "https://leetcode.com/u/ShardenduMishra22",
			Username: "ShardenduMishra22",
		},
		Resume: "https://drive.google.com/file/d/1F-ORaZyX8iMmBFhX2i-rtn21rdDMnsew/view?usp=sharing",
	},
	EducationData: EducationData{
		College: College{
			Name:     "Indian Institute of Information Technology, Dharwad",
			Batch:    "2023-2027",
			Website:  "https://www.iiitdwd.ac.in/",
			Location: "Dharwad, Karnataka, India",
			Degree:   "Bachelor of Technology",
			Field:    "Computer Science and Engineering",
		},
		School: School{
			Name:     "Delhi Public School, Kalyanpur",
			Batch:    "2008-2022",
			Location: "Kanpur, Uttar Pradesh, India",
			Class12: Class12{
				Percentage: "96.4%",
				Stream:     "PCM and Computer Science",
				Year:       "2022",
			},
			Class10: Class10{
				Percentage: "84%",
				Year:       "2020",
			},
		},
	},
	Languages: []string{"Hindi", "French", "English"},
}

// FooterLinks contains static footer navigation links
var FooterLinks = struct {
	QuickLinks     []Link
	SocialLinks    []Link
	CodingProfiles []Link
	MoreProfiles   []Link
	TechStack      []string
}{
	QuickLinks: []Link{
		{Name: "Projects", URL: "https://mishrashardendu22.is-a.dev/#projects"},
		{Name: "Experience", URL: "https://mishrashardendu22.is-a.dev/#experiences"},
		{Name: "Certifications", URL: "https://mishrashardendu22.is-a.dev/#certifications"},
		{Name: "Admin Dashboard", URL: "https://mishrashardendu22.is-a.dev/admin/dashboard"},
	},
	SocialLinks: []Link{
		{Name: "Twitter / X", URL: StaticData.SocialLinks.Twitter.URL},
		{Name: "LinkedIn", URL: StaticData.SocialLinks.LinkedIn.URL},
		{Name: "YouTube", URL: StaticData.SocialLinks.YouTube.URL},
		{Name: "Email", URL: "mailto:" + StaticData.PersonalInfo.Email},
	},
	CodingProfiles: []Link{
		{Name: "GitHub Main", URL: StaticData.SocialLinks.GitHub.Personal},
		{Name: "GitHub Learning", URL: StaticData.SocialLinks.GitHub.Learning},
		{Name: "Team Parashuram", URL: StaticData.SocialLinks.GitHub.Organization},
		{Name: "LeetCode", URL: StaticData.SocialLinks.LeetCode.URL},
	},
	MoreProfiles: []Link{
		{Name: "CodeChef", URL: "https://mishrashardendu22.is-a.dev/coming_soon"},
		{Name: "Codeforces", URL: "https://mishrashardendu22.is-a.dev/coming_soon"},
		{Name: "Education", URL: "https://mishrashardendu22.is-a.dev/#education"},
		{Name: "Contact Me", URL: "https://mishrashardendu22.is-a.dev/#contact"},
	},
	TechStack: []string{
		"Go", "Next.js", "TypeScript", "TailwindCSS",
		"MongoDB", "PostgreSQL", "Docker", "Kubernetes",
	},
}

type Link struct {
	Name string
	URL  string
}
