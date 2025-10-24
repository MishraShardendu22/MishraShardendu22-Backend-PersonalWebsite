package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/MishraShardendu22/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupExpRoutes(router fiber.Router, secret string) {
	// Public routes - no authentication required
	router.Get("/experiences", controller.GetExperiences)
	router.Get("/experiences/:id", controller.GetExperienceByID)

	// Admin routes - authentication required
	router.Post("/experiences", middleware.JWTMiddleware(secret), controller.AddExperiences)
	router.Put("/experiences/:id", middleware.JWTMiddleware(secret), controller.UpdateExperiences)
	router.Delete("/experiences/:id", middleware.JWTMiddleware(secret), controller.RemoveExperiences)
}
