package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/MishraShardendu22/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupProjectRoutes(app *fiber.App, secret string) {
	// Public routes - no authentication required
	app.Get("/api/projects", controller.GetProjects)
	app.Get("/api/projects/:id", controller.GetProjectByID)
	app.Get("/api/projects/kanban", controller.GetProjectsKanban)
	// app.Get("/api/UpdateProjectOrderInitial",controller.UpdateProjectOrder)

	// Admin routes - authentication required
	app.Post("/api/projects", middleware.JWTMiddleware(secret), controller.AddProjects)
	app.Put("/api/projects/:id", middleware.JWTMiddleware(secret), controller.UpdateProjects)
	app.Delete("/api/projects/:id", middleware.JWTMiddleware(secret), controller.RemoveProjects)
	app.Post("/api/projects/updateOrder", middleware.JWTMiddleware(secret), controller.UpdateProjectOrderKanban)
}
