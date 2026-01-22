package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupSearchRoutes(router fiber.Router) {
	searchGroup := router.Group("/search")

	searchGroup.Get("/", controller.Search)
	searchGroup.Get("/suggestions", controller.GetSearchSuggestions)
}
