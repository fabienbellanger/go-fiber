package main

import (
	"github.com/gofiber/basicauth"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/pprof"
)

type server struct {
	store  store
	router *fiber.App
}

func newServer() *server {
	s := &server{
		router: fiber.New(),
	}

	s.initHTTPServer()
	s.routes()
	s.initPprof()

	return s
}

func (s *server) initHTTPServer() {
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
	// Basic Auth
	// ----------
	// TODO: Put in config.toml
	cfg := basicauth.Config{
		Users: map[string]string{
			"john":  "doe",
			"admin": "123456",
		},
	}
	s.router.Use(basicauth.New(cfg))

	// pprof
	// -----
	s.router.Use(pprof.New())
}
