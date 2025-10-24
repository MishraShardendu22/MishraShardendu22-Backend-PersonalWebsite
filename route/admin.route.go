package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/MishraShardendu22/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(router fiber.Router, adminPass string, jwtSecret string) {
	router.Post("/admin/auth", func(c *fiber.Ctx) error {
		return controller.AdminRegisterAndLogin(c, adminPass, jwtSecret)
	})

	router.Get("/admin/auth", middleware.JWTMiddleware(jwtSecret), controller.AdminGet)
}
