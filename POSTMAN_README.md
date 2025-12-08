# Postman Collection - Personal Website Backend API

Complete API collection for the Personal Website Backend with proper variable management and automation.

## 📦 Files Included

- `postman_collection.json` - Main API collection with all endpoints
- `postman_environment.json` - Environment variables configuration

## 🚀 Quick Start

### 1. Import into Postman

1. Open Postman
2. Click **Import** button (top left)
3. Import both files:
   - `postman_collection.json`
   - `postman_environment.json`

### 2. Configure Environment

1. Select the imported environment from the dropdown (top right)
2. Click the eye icon to view/edit variables
3. Update these values:
   - `base_url`: Your API base URL (default: `http://localhost:5000`)
   - `admin_email`: Your admin email
   - `admin_password`: Your admin password

### 3. Start Testing

1. Run **Admin Login** request first
2. JWT token will be automatically saved
3. All protected endpoints will use the saved token automatically

## 🔑 Authentication

### Automatic Token Management

The collection includes automatic JWT token handling:

1. **Login Request** has a test script that automatically saves the JWT token
2. Protected endpoints use `{{jwt_token}}` variable automatically
3. No manual copy-paste needed!

### How it Works

When you run the **Admin Login** request:
```javascript
// Automatic test script extracts and saves token
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.collectionVariables.set("jwt_token", jsonData.data.token);
}
```

## 📋 Collection Variables

Variables are automatically managed but can be manually set if needed:

| Variable | Description | Auto-saved | Usage |
|----------|-------------|------------|-------|
| `base_url` | API base URL | ❌ | All requests |
| `admin_email` | Admin email | ❌ | Login |
| `admin_password` | Admin password | ❌ | Login |
| `jwt_token` | JWT authentication token | ✅ | Protected endpoints |
| `project_id` | Last created/used project ID | ✅ | Project operations |
| `experience_id` | Last created/used experience ID | ✅ | Experience operations |
| `certification_id` | Last created/used cert ID | ✅ | Certification operations |
| `volunteer_id` | Last created/used volunteer ID | ✅ | Volunteer operations |

## 📂 Endpoint Structure

### Public Endpoints (No Auth Required)
- ✅ GET All Resources
- ✅ GET Resource by ID
- ✅ GitHub Stats
- ✅ LeetCode Stats
- ✅ Timeline

### Protected Endpoints (JWT Required)
- 🔒 POST (Create)
- 🔒 PUT (Update)
- 🔒 DELETE (Remove)

## 🔗 API Endpoints Overview

### Admin
- `POST /api/v1/admin/auth` - Login (auto-saves JWT)
- `GET /api/v1/admin/auth` - Get admin info

### Projects
- `GET /api/v1/projects` - Get all projects
- `GET /api/v1/projects/kanban` - Get projects with order
- `GET /api/v1/projects/:id` - Get project by ID
- `POST /api/v1/projects` - Add project (protected)
- `PUT /api/v1/projects/:id` - Update project (protected)
- `POST /api/v1/projects/updateOrder` - Update order (protected)
- `DELETE /api/v1/projects/:id` - Delete project (protected)

### Experiences
- `GET /api/v1/experiences` - Get all experiences
- `GET /api/v1/experiences/:id` - Get experience by ID
- `POST /api/v1/experiences` - Add experience (protected)
- `PUT /api/v1/experiences/:id` - Update experience (protected)
- `DELETE /api/v1/experiences/:id` - Delete experience (protected)

### Certifications
- `GET /api/v1/certifications` - Get all certifications
- `GET /api/v1/certifications/:id` - Get certification by ID
- `POST /api/v1/certifications` - Add certification (protected)
- `PUT /api/v1/certifications/:id` - Update certification (protected)
- `DELETE /api/v1/certifications/:id` - Delete certification (protected)

### Volunteer Experiences
- `GET /api/v1/volunteer/experiences` - Get all volunteer experiences
- `GET /api/v1/volunteer/experiences/:id` - Get by ID
- `POST /api/v1/volunteer/experiences` - Add (protected)
- `PUT /api/v1/volunteer/experiences/:id` - Update (protected)
- `DELETE /api/v1/volunteer/experiences/:id` - Delete (protected)

### Skills
- `GET /api/v1/skills` - Get all skills
- `POST /api/v1/skills` - Add skills (protected)

### GitHub Stats
- `GET /api/v1/github` - GitHub profile
- `GET /api/v1/github/stars` - Total stars
- `GET /api/v1/github/commits` - Total commits
- `GET /api/v1/github/languages` - Language stats
- `GET /api/v1/github/top-repos` - Top repositories
- `GET /api/v1/github/calendar` - Contribution calendar

### LeetCode Stats
- `GET /api/v1/leetcode` - LeetCode statistics

### Timeline
- `GET /api/v1/timeline` - Experience timeline

## 🎯 Usage Tips

### For Frontend Developers

1. **Use Environment Variables**: All URLs and IDs use variables
2. **No Hardcoding**: Change `base_url` in environment to switch between dev/staging/prod
3. **Sample Data Included**: All POST/PUT requests have example payloads
4. **ID Auto-capture**: Resource IDs are automatically saved after creation

### Testing Workflow

1. **Initial Setup**:
   ```
   1. Import collection + environment
   2. Set base_url, admin_email, admin_password
   3. Run Admin Login
   ```

2. **Create Resource**:
   ```
   1. Run POST request (e.g., Add Project)
   2. ID automatically saved to {{project_id}}
   3. Use in subsequent requests
   ```

3. **Update/Delete**:
   ```
   1. Use saved {{project_id}} variable
   2. No manual ID entry needed
   ```

## 🔄 Environment Switching

Create multiple environments for different stages:

### Development
```json
{
  "base_url": "http://localhost:5000"
}
```

### Staging
```json
{
  "base_url": "https://staging-api.example.com"
}
```

### Production
```json
{
  "base_url": "https://api.example.com"
}
```

Switch between them using the environment dropdown in Postman.

## 📝 Sample Request Bodies

All requests include realistic sample data. Key examples:

### Project
```json
{
  "order": 1,
  "skills": ["Go", "MongoDB", "Fiber"],
  "description": "Detailed project description",
  "project_name": "Project Name",
  "project_video": "https://example.com/video.mp4",
  "project_live_link": "https://example.com",
  "small_description": "Short project summary",
  "project_repository": "https://github.com/user/repo"
}
```

### Experience
```json
{
  "images": ["https://example.com/image1.jpg"],
  "technologies": ["Go", "React", "MongoDB"],
  "created_by": "Admin User",
  "description": "Detailed experience description",
  "company_name": "Company Name",
  "company_logo": "https://example.com/logo.png",
  "certificate_url": "https://example.com/certificate.pdf",
  "projects": [],
  "experience_time_line": [
    {
      "position": "Software Engineer",
      "start_date": "2023-01-01",
      "end_date": "2024-01-01"
    }
  ]
}
```

## 🛡️ Security Notes

- JWT tokens are stored in collection variables (not in environment for security)
- Use Postman's secret variable type for sensitive data
- Never commit real credentials in exported collections
- The environment file marks `admin_password` and `jwt_token` as secret

## 🤝 Sharing with Frontend Team

### Option 1: Export and Share Files
```bash
# Share these files with your team
postman_collection.json
postman_environment.json
POSTMAN_README.md
```

### Option 2: Postman Workspace
1. Create a team workspace in Postman
2. Share the collection directly
3. Team members can sync automatically

### What to Update Before Sharing
1. Set `base_url` to your deployed API URL
2. Remove any real credentials
3. Add placeholder values in environment

## 📚 Additional Resources

- [Postman Documentation](https://learning.postman.com/)
- [API Documentation](./API_DOCS.md)
- [Backend Repository](https://github.com/MishraShardendu22/Backend-PersonalWebsite)

## 🐛 Troubleshooting

### Token Not Saving
- Check the "Tests" tab in Admin Login request
- Verify response structure matches expected format

### 401 Unauthorized
- Re-run Admin Login to refresh token
- Check if token is set in collection variables

### Wrong Base URL
- Verify selected environment
- Check `base_url` value in environment variables

### Variables Not Working
- Ensure environment is selected (top right dropdown)
- Check variable names match exactly (case-sensitive)

## 📞 Support

For issues or questions:
- Check API_DOCS.md for detailed endpoint documentation
- Review the backend repository README
- Contact the backend team

---

**Ready to use!** Import the collection and start testing your API. 🚀
