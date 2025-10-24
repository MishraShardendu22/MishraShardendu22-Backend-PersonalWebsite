package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/MishraShardendu22/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupProjectRoutes(router fiber.Router, secret string) {
	// Public routes - no authentication required
	router.Get("/projects", controller.GetProjects)
	router.Get("/projects/kanban", controller.GetProjectsKanban)
	// router.Get("/UpdateProjectOrderInitial",controller.UpdateProjectOrder)

	// Admin routes - authentication required
	router.Post("/projects", middleware.JWTMiddleware(secret), controller.AddProjects)
	router.Post("/projects/updateOrder", middleware.JWTMiddleware(secret), controller.UpdateProjectOrderKanban)

	router.Get("/projects/:id", controller.GetProjectByID)
	router.Put("/projects/:id", middleware.JWTMiddleware(secret), controller.UpdateProjects)
	router.Delete("/projects/:id", middleware.JWTMiddleware(secret), controller.RemoveProjects)
}

/*
Route order is important in Fiber routing:
- Specific routes like "/projects/kanban" must be defined BEFORE parameterized routes like "/projects/:id"
- If "/projects/:id" comes first, the router would treat "kanban" as a value for the :id parameter
- Current implementation is correct: /kanban route is defined before /:id route
*/
