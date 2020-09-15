package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func fiberProcessTimer() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		now := time.Now()

		c.Next()

		duration := time.Since(now)

		c.Set("X-Process-Time", fmt.Sprintf("%v", duration))

		return nil
	}
}
