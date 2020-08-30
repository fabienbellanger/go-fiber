package main

import (
	"github.com/gofiber/fiber"
)

func (s *server) handlerHome(c *fiber.Ctx) {
	c.Send("Hello, World 👋!")
}

func (s *server) handlerBigJSON(c *fiber.Ctx) {
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
	c.JSON(&users)
}
