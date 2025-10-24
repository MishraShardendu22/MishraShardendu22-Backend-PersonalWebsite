package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupTimeline(router fiber.Router, secret string) {

	router.Get("/timeline", controller.ExperienceTimeline)
}
