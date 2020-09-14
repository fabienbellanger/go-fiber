package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

type server struct {
	store  store
	router *fiber.App
	mode   string
}

func newServer() *server {
	s := &server{
		router: fiber.New(errorHandling()),
		mode:   "production",
	}

	s.initHTTPServer()
	s.routes()
	s.initPprof()

	// Custom 404 (after all routes)
	// -----------------------------
	s.router.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "Resource Not Found",
		})
	})

	return s
}

func (s *server) initHTTPServer() {
	// Mode
	// ----
	s.mode = viper.GetString("environment")

	// CORS
	// ----
	s.router.Use(cors.New())

	// Logger
	// ------
	if s.mode != "production" {
		s.router.Use(logger.New(logger.Config{
			Next:       nil,
			Format:     "[${time}] ${status} - ${latency} - ${method} ${path}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "Local",
			Output:     os.Stderr,
		}))
	}

	// Recover
	// -------
	s.router.Use(recover.New())
}

func (s *server) initPprof() {
	if viper.GetBool("debug.pprof") {
		private := s.router.Group("/debug/pprof")

		// Basic Auth
		// ----------
		cfg := basicauth.Config{
			Users: map[string]string{
				viper.GetString("debug.basicAuthUsername"): viper.GetString("debug.basicAuthPassword"),
			},
		}
		private.Use(basicauth.New(cfg))

		// pprof (The handled paths all begin with /debug/pprof/)
		// ------------------------------------------------------
		s.router.Use(pprof.New())
	}
}

func errorHandling() fiber.Config {
	return fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// Retreive the custom statuscode if it's an fiber.*Error
			e, ok := err.(*fiber.Error)

			if ok {
				code = e.Code
			}

			log.Printf("Error:%v - Code:%v", e, code)

			if e != nil {
				return c.JSON(e)
			}

			if code == 500 {
				return c.JSON(fiber.Map{
					"code":    code,
					"message": "Internal Server Error",
				})
			}

			return nil
		},
	}
}
