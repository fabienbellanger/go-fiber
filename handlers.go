package main

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func (s *server) handlerHome(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func (s *server) handlerBigJSON(c *fiber.Ctx) error {
	type User struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Password  string `json:"-"`
		Lastname  string `json:"lastname"`
		Firstname string `json:"firstname"`
	}
	var users []User
	for i := 0; i < 100000; i++ {
		users = append(users, User{
			ID:        i + 1,
			Username:  "My Username",
			Lastname:  "My Lastname",
			Firstname: "My Firstname",
		})
	}
	return c.JSON(&users)
}

func (s *server) handlerBigJSONStream(c *fiber.Ctx) error {
	type User struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Password  string `json:"-"`
		Lastname  string `json:"lastname"`
		Firstname string `json:"firstname"`
	}

	c.Set("Content-Type", "application/json")

	c.Write([]byte("["))

	n := 100000
	for i := 0; i < n; i++ {
		user, err := json.Marshal(User{
			ID:        i + 1,
			Username:  "My Username",
			Lastname:  "My Lastname",
			Firstname: "My Firstname",
		})
		if err != nil {
			continue
		}

		c.Write(user)

		if i < n-1 {
			c.Write([]byte(","))
		}
	}
	c.Write([]byte("]"))

	return c.SendStatus(200)
}
