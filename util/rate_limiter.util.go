package util

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupCRUDAPILimiter(logger *slog.Logger) fiber.Handler {
	CrudAPILimiter := limiter.New(limiter.Config{
		Max:        120,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warn("CRUD API rate limit hit",
				"ip", c.IP(),
				"path", c.Path(),
			)
			return ResponseAPI(c, fiber.StatusTooManyRequests,
				"Too many requests. Please slow down.",
				fiber.Map{"retry_after": 60, "endpoint": c.Path()},
				"")
		},
	})

	return CrudAPILimiter
}

func SetupExternalAPILimiter(logger *slog.Logger) fiber.Handler {
	ExternalAPILimiter := limiter.New(limiter.Config{
		Max:        120,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warn("External API rate limit hit",
				"ip", c.IP(),
				"path", c.Path(),
			)

			return ResponseAPI(c, fiber.StatusTooManyRequests,
				"Too many requests to external APIs. Please wait before retrying.",
				fiber.Map{"retry_after": 60, "endpoint": c.Path()},
				"")
		},
	})

	return ExternalAPILimiter
}
