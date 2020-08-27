package main

import (
	"github.com/gofiber/fiber"
)

func (s *server) handlerHome(c *fiber.Ctx) {
	c.Send("Hello, World 👋!")
}
