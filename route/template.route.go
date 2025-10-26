package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/gofiber/fiber/v2"
)

// SetupTemplateRoutes sets up routes for Templ-based HTML sections
func SetupTemplateRoutes(app *fiber.App) {
	// Root route - serves the full page
	app.Get("/", controller.GetFullPage)

	// Full page demo (alias)
	app.Get("/demo", controller.GetFullPage)

	// Create a route group for template sections
	templates := app.Group("/test")

	// Hero Section
	templates.Get("/hero", controller.GetHeroSection)

	// Skills Section
	templates.Get("/skills", controller.GetSkillsSection)

	// Education Section
	templates.Get("/education", controller.GetEducationSection)

	// Contact/Dashboard Section
	templates.Get("/contact", controller.GetContactSection)

	// Footer Section
	templates.Get("/footer", controller.GetFooterSection)
}
