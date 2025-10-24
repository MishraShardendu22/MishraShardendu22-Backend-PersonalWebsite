package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/MishraShardendu22/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupCertificationRoutes(router fiber.Router, secret string) {
	// Public routes - no authentication required
	router.Get("/certifications", controller.GetCertifications)
	router.Get("/certifications/:id", controller.GetCertificationByID)

	// Admin routes - authentication required
	router.Post("/certifications", middleware.JWTMiddleware(secret), controller.AddCertification)
	router.Put("/certifications/:id", middleware.JWTMiddleware(secret), controller.UpdateCertification)
	router.Delete("/certifications/:id", middleware.JWTMiddleware(secret), controller.RemoveCertification)
}
