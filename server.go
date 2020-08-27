package main

import (
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

	// pprof
	// -----
	s.router.Use(pprof.New())
}
