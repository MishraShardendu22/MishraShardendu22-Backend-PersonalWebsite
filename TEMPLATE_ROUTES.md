# Template Routes Documentation

This document provides comprehensive documentation for all Templ-based HTML template routes in the portfolio backend.

## Table of Contents
- [Overview](#overview)
- [Base Configuration](#base-configuration)
- [Routes](#routes)
  - [Full Page Routes](#full-page-routes)
  - [Component Section Routes](#component-section-routes)
- [Theme Support](#theme-support)
- [Response Format](#response-format)
- [Usage Examples](#usage-examples)

## Overview

The template routes serve server-side rendered HTML using [Templ](https://templ.guide/), a type-safe templating language for Go. All routes return fully styled, responsive HTML with built-in dark/light theme support.

## Base Configuration

**Base URL**: `/`  
**Response Type**: `text/html`  
**Framework**: Fiber v2 + Templ  
**Theme System**: CSS Custom Properties with Tailwind CSS

## Routes

### Full Page Routes

#### Get Full Portfolio Page
```
GET /
GET /demo
```

**Description**: Renders the complete portfolio page with all sections (Hero, Skills, Education, Contact, Footer).

**Response**: Full HTML page with:
- Responsive layout
- Dark/light theme toggle
- All portfolio sections combined
- Tailwind CSS styling
- Theme persistence via localStorage

**Features**:
- **Auto Theme Detection**: Respects system preferences
- **Theme Toggle**: JavaScript-based theme switcher included
- **Responsive Design**: Mobile-first approach
- **SEO Optimized**: Proper meta tags and semantic HTML

**Example**:
```bash
curl http://localhost:8080/
# or
curl http://localhost:8080/demo
```

---

### Component Section Routes

All component routes are under the `/test` prefix for testing individual sections.

#### 1. Hero Section
```
GET /test/hero
```

**Description**: Renders the hero/landing section of the portfolio.

**Content Includes**:
- Name and professional title
- Tagline/bio
- Profile image/avatar
- Social media links (GitHub, LinkedIn, Twitter, etc.)
- Call-to-action buttons
- Animated background effects

**Data Structure**:
```go
type HeroData struct {
    Name        string
    Title       string
    Tagline     string
    Avatar      string
    SocialLinks []SocialLink
}

type SocialLink struct {
    Platform string
    URL      string
    Icon     string
}
```

**Theme Variables Used**:
- `--background`, `--foreground`
- `--primary`, `--primary-foreground`
- `--card`, `--card-foreground`
- `--accent`, `--accent-foreground`

---

#### 2. Skills Section
```
GET /test/skills
```

**Description**: Displays technical skills organized by categories.

**Content Includes**:
- Skill categories (Frontend, Backend, DevOps, etc.)
- Individual skills with proficiency levels
- Skill icons/logos
- Interactive hover effects
- Progress indicators

**Data Structure**:
```go
type SkillsData struct {
    Categories []SkillCategory
}

type SkillCategory struct {
    Name   string
    Skills []Skill
}

type Skill struct {
    Name        string
    Icon        string
    Proficiency int // 0-100
}
```

**Theme Variables Used**:
- `--card`, `--card-foreground`
- `--border`
- `--muted`, `--muted-foreground`
- `--accent`

---

#### 3. Education Section
```
GET /test/education
```

**Description**: Shows educational background and academic achievements.

**Content Includes**:
- Educational institutions
- Degrees and certifications
- Study periods (start/end dates)
- Grades/achievements
- Relevant coursework
- Timeline visualization

**Data Structure**:
```go
type EducationData struct {
    Institutions []Education
}

type Education struct {
    Institution string
    Degree      string
    Field       string
    StartDate   string
    EndDate     string
    Grade       string
    Achievements []string
}
```

**Theme Variables Used**:
- `--card`, `--card-foreground`
- `--border`
- `--primary`, `--primary-foreground`
- `--muted-foreground`

---

#### 4. Contact/Dashboard Section
```
GET /test/contact
```

**Description**: Developer dashboard with real-time statistics and contact form.

**Content Includes**:
- **GitHub Profile Stats**:
  - Profile avatar and bio
  - Repository count
  - Followers count
  - Total stars
  
- **LeetCode Statistics**:
  - Easy problems solved
  - Medium problems solved
  - Hard problems solved
  - Total problems solved
  
- **Contact Form**:
  - Name input
  - Email input
  - Message textarea
  - Email submission (mailto link)
  
- **Technology Stack**:
  - List of technologies used
  - Tech badges

**Data Structure**:
```go
type ContactData struct {
    GitHubAvatar     string
    GitHubUsername   string
    GitHubBio        string
    GitHubRepos      int
    GitHubFollowers  int
    TotalStars       int
    LeetCodeEasy     int
    LeetCodeMedium   int
    LeetCodeHard     int
    LeetCodeTotal    int
    RecentCommits    []ContactCommit
    Technologies     []string
}

type ContactCommit struct {
    Message string
    Date    string
    Repo    string
}
```

**Theme Variables Used**:
- `--background`, `--foreground`
- `--card`, `--card-foreground`
- `--border`, `--input`
- `--muted`, `--muted-foreground`
- `--ring` (focus states)

**Form Behavior**:
- Client-side form handling
- Opens default email client with pre-filled subject and body
- Email sent to: `shardendumishra01@gmail.com`

---

#### 5. Footer Section
```
GET /test/footer
```

**Description**: Portfolio footer with additional information and links.

**Content Includes**:
- Copyright information
- Additional navigation links
- Social media links (duplicate/alternative placement)
- Site credits
- Footer navigation

**Data Structure**:
```go
type FooterData struct {
    Copyright    string
    Links        []FooterLink
    SocialLinks  []SocialLink
    Year         int
}

type FooterLink struct {
    Text string
    URL  string
}
```

**Theme Variables Used**:
- `--background`, `--foreground`
- `--muted`, `--muted-foreground`
- `--border`

---

## Theme Support

### CSS Custom Properties

All templates use CSS custom properties that automatically adapt to light/dark themes:

**Light Theme**:
```css
--background: 0 0% 100%;
--foreground: 222.2 84% 4.9%;
--primary: 142.1 76.2% 36.3%;
--card: 0 0% 100%;
--border: 214.3 31.8% 91.4%;
```

**Dark Theme**:
```css
--background: 222.2 84% 4.9%;
--foreground: 210 40% 98%;
--primary: 142.1 70.6% 45.3%;
--card: 222.2 84% 10%;
--border: 217.2 32.6% 17.5%;
```

### Theme Toggle

The theme can be toggled using JavaScript (included in full page):

```javascript
function toggleTheme() {
    document.documentElement.classList.toggle('dark');
    localStorage.setItem('theme', 
        document.documentElement.classList.contains('dark') ? 'dark' : 'light'
    );
}
```

### Automatic Theme Detection

On page load, the theme is automatically set based on:
1. User's previous preference (localStorage)
2. System preference (prefers-color-scheme media query)

---

## Response Format

All routes return:
- **Content-Type**: `text/html; charset=utf-8`
- **Status Code**: `200 OK` (on success)
- Fully rendered HTML with inline Tailwind CSS
- Responsive design (mobile-first)
- Accessibility features (semantic HTML, ARIA labels)

---

## Usage Examples

### Full Page Access
```bash
# Access the complete portfolio
curl http://localhost:8080/

# Browser access
http://localhost:8080/
```

### Individual Section Testing
```bash
# Test hero section
curl http://localhost:8080/test/hero

# Test skills section
curl http://localhost:8080/test/skills

# Test education section
curl http://localhost:8080/test/education

# Test contact section
curl http://localhost:8080/test/contact

# Test footer section
curl http://localhost:8080/test/footer
```

### Using in Frontend Applications

```javascript
// Fetch individual sections for SPA integration
async function loadSection(section) {
    const response = await fetch(`http://localhost:8080/test/${section}`);
    const html = await response.text();
    document.getElementById('content').innerHTML = html;
}

// Example: Load skills section
loadSection('skills');
```

### iframe Embedding

```html
<!-- Embed individual sections in iframes -->
<iframe src="http://localhost:8080/test/hero" 
        style="width: 100%; height: 600px; border: none;">
</iframe>
```

---

## Implementation Details

### Route Setup
Routes are configured in `/route/template.route.go`:

```go
func SetupTemplateRoutes(app *fiber.App) {
    // Full page routes
    app.Get("/", controller.GetFullPage)
    app.Get("/demo", controller.GetFullPage)
    
    // Component routes under /test
    templates := app.Group("/test")
    templates.Get("/hero", controller.GetHeroSection)
    templates.Get("/skills", controller.GetSkillsSection)
    templates.Get("/education", controller.GetEducationSection)
    templates.Get("/contact", controller.GetContactSection)
    templates.Get("/footer", controller.GetFooterSection)
}
```

### Controller Location
All template controllers are in `/controller/template.controller.go`

### Template Files
Templ source files are located in `/templates/`:
- `layout.templ` - Full page layout
- `hero.templ` - Hero section
- `skills.templ` - Skills section
- `education.templ` - Education section
- `contact.templ` - Contact/dashboard section
- `footer.templ` - Footer section

Generated Go files (`*_templ.go`) are auto-generated by the Templ compiler.

---

## Best Practices

1. **Testing Individual Components**: Use `/test/*` routes to test and debug individual sections
2. **Theme Testing**: Toggle between light/dark themes to ensure proper styling
3. **Responsive Testing**: Test on different screen sizes (mobile, tablet, desktop)
4. **Performance**: Full page includes all sections, use component routes for partial updates
5. **Caching**: Consider implementing HTTP caching headers for production

---

## Technical Stack

- **Backend**: Go + Fiber v2
- **Templates**: Templ (type-safe Go templates)
- **Styling**: Tailwind CSS (via CDN)
- **Theme System**: CSS Custom Properties
- **JavaScript**: Vanilla JS for theme toggle

---

## Development Commands

```bash
# Generate Templ templates
templ generate

# Run the server
go run main.go

# Build for production
go build -o portfolio-backend
```

---

## Notes

- All sections are responsive and mobile-friendly
- Dark/light theme support is built into all components
- Contact form uses mailto: protocol (no backend submission)
- Static data is currently hardcoded (can be made dynamic with database integration)
- All routes are GET requests (no authentication required)

---

## Future Enhancements

- [ ] Add API endpoints for dynamic data updates
- [ ] Implement admin panel for content management
- [ ] Add analytics tracking
- [ ] Server-side theme preference detection
- [ ] Form submission backend integration
- [ ] Add more sections (Projects, Experience, Certifications, Volunteering)
- [ ] Implement rate limiting for contact form
- [ ] Add CAPTCHA for spam prevention

---

## Support

For issues or questions, please refer to:
- Main README: `/README.md`
- API Documentation: `/API_DOCS.md`
- Contributing Guidelines: `/CONTRIBUTING.md`

---

**Last Updated**: October 26, 2025  
**Version**: 1.0.0  
**Maintainer**: Shardendu Mishra
