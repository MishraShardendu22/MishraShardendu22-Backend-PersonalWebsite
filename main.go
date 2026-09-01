package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/MishraShardendu22/database"
	"github.com/MishraShardendu22/models"
	"github.com/MishraShardendu22/route"
	"github.com/MishraShardendu22/util"
	"github.com/MishraShardendu22/web"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func loadConfig() *models.Config {
	config := &models.Config{
		Port:             util.GetEnv("PORT", "5000"),
		Environment:      util.GetEnv("ENVIRONMENT", "development"),
		CorsAllowOrigins: util.GetEnv("CORS_ALLOW_ORIGINS", "*"),
		LogLevel:         util.GetEnv("LOG_LEVEL", "info"),
		MONGODB_URI:      util.GetEnv("MONGODB_URI", "some_default_mongo_uri"),
		DbName:           util.GetEnv("DB_NAME", "test"),
		AdminPass:        util.GetEnv("ADMIN_PASS", ""),
		JWT_SECRET:       util.GetEnv("JWT_SECRET", ""),
	}
	return config
}

func setupLogger(config *models.Config) {
	var level slog.Level
	switch config.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handler slog.Handler
	if config.Environment == "development" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func setupMiddleware(app *fiber.App, config *models.Config) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: config.Environment == "development",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:  config.CorsAllowOrigins,
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length",
		MaxAge:        86400,
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Global rate limiter - applies to all routes
	// Limits requests per IP address to prevent abuse
	app.Use(limiter.New(limiter.Config{
		Max:        100,             // Maximum number of requests
		Expiration: 1 * time.Minute, // Time window
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP address
		},
		LimitReached: func(c *fiber.Ctx) error {
			return util.ResponseAPI(c, fiber.StatusTooManyRequests,
				"Too many requests from this IP. Please try again later.",
				fiber.Map{"retry_after": 60},
				"")
		},
		SkipFailedRequests:     false, // Count failed requests
		SkipSuccessfulRequests: false, // Count all requests
		Storage:                nil,   // Use in-memory storage (consider Redis for production)
	}))
}

func gracefulShutdown(app *fiber.App, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}

func init() {
	currEnv := "development"

	if currEnv == "development" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: error loading .env file: %v", err)
		}
	}
}

func main() {
	config := loadConfig()
	if err := database.ConnectDatabase(config.DbName, config.MONGODB_URI); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	setupLogger(config)
	logger := slog.Default()

	logger.Info("Starting Portfolio Backend",
		"environment", config.Environment,
		"port", config.Port,
		"log_level", config.LogLevel,
	)

	app := fiber.New(fiber.Config{
		AppName:      "Portfolio Backend",
		ServerHeader: "Portfolio-Backend",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error("request error", slog.Group("req",
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.String("error", err.Error()),
			))

			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			// Content negotiation only. Programmatic clients keep the exact
			// JSON envelope they have always received; a browser that lands on
			// an unknown route gets a readable page instead of a raw payload.
			if code == fiber.StatusNotFound && wantsHTML(c) {
				c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
				return c.Status(code).SendString(web.NotFoundPage(c.Path()))
			}

			return util.ResponseAPI(c, code, message, nil, "")
		},
	})

	setupMiddleware(app, config)

	SetUpRoutes(app, logger, config)

	go func() {
		logger.Info("Server starting", "port", config.Port)
		if err := app.Listen(":" + config.Port); err != nil {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	gracefulShutdown(app, logger)
}

// wantsHTML reports whether the caller is a browser navigating to the service
// rather than a program calling the API. It is deliberately conservative: only
// a request that explicitly prefers HTML over JSON gets an HTML response.
func wantsHTML(c *fiber.Ctx) bool {
	accept := c.Get(fiber.HeaderAccept)
	if !strings.Contains(accept, "text/html") {
		return false
	}
	return !strings.Contains(accept, "application/json")
}

func SetUpRoutes(app *fiber.App, logger *slog.Logger, config *models.Config) {
	crudGroup := app.Group("/api", util.SetupCRUDAPILimiter(logger))
	route.SetupSearchRoutes(crudGroup)

	statsGroup := app.Group("/api", util.SetupExternalAPILimiter(logger))
	route.SetupStatsRoutes(statsGroup)

	route.SetupTimeline(crudGroup, config.JWT_SECRET)
	route.SetupExpRoutes(crudGroup, config.JWT_SECRET)
	route.SetupProjectRoutes(crudGroup, config.JWT_SECRET)
	route.SetupSkillRoutes(crudGroup, config.JWT_SECRET)
	route.SetupVolunteerExpRoutes(crudGroup, config.JWT_SECRET)
	route.SetupCertificationRoutes(crudGroup, config.JWT_SECRET)
	route.SetupAdminRoutes(crudGroup, config.AdminPass, config.JWT_SECRET)

	app.Get("/api/test123", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Working fine",
		})
	})

	// Human-facing status page. Programmatic callers that ask for JSON still
	// get JSON, so nothing that consumes this service is affected.
	app.Get("/", func(c *fiber.Ctx) error {
		if !wantsHTML(c) {
			return util.ResponseAPI(c, fiber.StatusOK, "Portfolio API is running",
				fiber.Map{
					"service":     "Portfolio API",
					"environment": config.Environment,
					"docs":        "/",
				}, "")
		}
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.Status(fiber.StatusOK).SendString(
			web.StatusPage(config.Environment, runtime.Version()),
		)
	})
}
