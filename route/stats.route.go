package route

import (
	"log/slog"

	"github.com/MishraShardendu22/controller"
	"github.com/MishraShardendu22/util"
	"github.com/gofiber/fiber/v2"
)

func SetupStatsRoutes(router fiber.Router, logger *slog.Logger) {
	router.Get("/github", util.SetupExternalAPILimiter(logger), controller.FetchGitHubProfile)
	router.Get("/leetcode", util.SetupExternalAPILimiter(logger), controller.FetchLeetCodeData)
	router.Get("/github/stars", util.SetupExternalAPILimiter(logger), controller.FetchGitHubStars)
	router.Get("/github/commits", util.SetupExternalAPILimiter(logger), controller.FetchGitHubCommits)
	router.Get("/github/languages", util.SetupExternalAPILimiter(logger), controller.FetchGitHubLanguages)
	router.Get("/github/top-repos", util.SetupExternalAPILimiter(logger), controller.FetchTopStarredRepos)
	router.Get("/github/calendar", util.SetupExternalAPILimiter(logger), controller.FetchContributionCalendar)
}
