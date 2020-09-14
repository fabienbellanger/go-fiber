package main

import (
	"github.com/gofiber/fiber/v2"
)

func (s *server) handlerGetUser(c *fiber.Ctx) error {
	u := s.store.user().getUser()
	return c.JSON(u)
}
