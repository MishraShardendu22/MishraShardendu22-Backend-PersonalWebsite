package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/MishraShardendu22/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupVolunteerExpRoutes(router fiber.Router, secret string) {
	// Public routes - no authentication required
	router.Get("/volunteer/experiences", controller.GetVolunteerExperiences)
	router.Get("/volunteer/experiences/:id", controller.GetVolunteerExperienceByID)

	// Admin routes - authentication required
	router.Post("/volunteer/experiences", middleware.JWTMiddleware(secret), controller.AddVolunteerExperiences)
	router.Put("/volunteer/experiences/:id", middleware.JWTMiddleware(secret), controller.UpdateVolunteerExperiences)
	router.Delete("/volunteer/experiences/:id", middleware.JWTMiddleware(secret), controller.RemoveVolunteerExperiences)
}
