package main

import (
	"log"

	"github.com/gofiber/basicauth"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/pprof"
	"github.com/spf13/viper"
)

type server struct {
	store  store
	router *fiber.App
	mode   string
}

func newServer() *server {
	s := &server{
		router: fiber.New(),
		mode:   "production",
	}

	s.initHTTPServer()
	s.routes()
	s.initPprof()
	s.errorHandling()

	return s
}

func (s *server) initHTTPServer() {
	// Mode
	// ----
	s.mode = viper.GetString("environment")

	// Default RequestID
	// -----------------
	s.router.Use(middleware.RequestID())

	// CORS
	// ----
	s.router.Use(cors.New())

	// Logger
	// ------
	s.router.Use(middleware.Logger())

	// Recover
	// -------
	s.router.Use(middleware.Recover())

	// Custom 404 (after all routes)
	// -----------------------------
	s.router.Use(func(ctx *fiber.Ctx) {
		ctx.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "Resource Not Found",
		})
	})
}

func (s *server) initPprof() {
	if viper.GetBool("debug.pprof") {
		private := s.router.Group("private", func(c *fiber.Ctx) {
			c.Next()
		})

		// Basic Auth
		// ----------
		cfg := basicauth.Config{
			Users: map[string]string{
				viper.GetString("debug.basicAuthUsername"): viper.GetString("debug.basicAuthPassword"),
			},
		}
		private.Use(basicauth.New(cfg))

		// pprof
		// -----
		private.Use(pprof.New())
	}
}

func (s *server) errorHandling() {
	s.router.Settings.ErrorHandler = func(c *fiber.Ctx, err error) {
		code := fiber.StatusInternalServerError

		// Retreive the custom statuscode if it's an fiber.*Error
		e, ok := err.(*fiber.Error)

		if ok {
			code = e.Code
		}

		log.Printf("Error:%v - Code:%v", e, code)

		if e != nil {
			c.JSON(e)
		}

		if code == 500 {
			c.JSON(fiber.Map{
				"code":    code,
				"message": "Internal Server Error",
			})
		}
	}
}
