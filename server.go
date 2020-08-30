package main

import (
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
