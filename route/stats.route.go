package route

import (
	"github.com/MishraShardendu22/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupStatsRoutes(router fiber.Router) {
	// GitHub Stats Routes - All public, no authentication required
	router.Get("/github", controller.FetchGitHubProfile)
	router.Get("/github/stars", controller.FetchGitHubStars)
	router.Get("/github/commits", controller.FetchGitHubCommits)
	router.Get("/github/languages", controller.FetchGitHubLanguages)
	router.Get("/github/top-repos", controller.FetchTopStarredRepos)
	router.Get("/github/calendar", controller.FetchContributionCalendar)

	// LeetCode Stats Routes - All public, no authentication required
	router.Get("/leetcode", controller.FetchLeetCodeData)
}
