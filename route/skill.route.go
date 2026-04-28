package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/MishraShardendu22/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupSkillRoutes(router fiber.Router, secret string) {
	// Public routes - no authentication required
	router.Get("/skills", controller.GetSkills)

	// Admin routes - authentication required
	router.Post("/skills", middleware.JWTMiddleware(secret), controller.AddSkills)
}
