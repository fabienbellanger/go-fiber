package main

import (
	"github.com/gofiber/fiber"
)

func (s *server) handlerGetUser(c *fiber.Ctx) {
	u := s.store.user().getUser()
	c.JSON(u)
}
