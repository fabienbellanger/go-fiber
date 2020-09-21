package main

import (
	"bufio"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func (s *server) handlerStatic(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	})
}

func (s *server) handlerHome(c *fiber.Ctx) error {
	name := c.Params("name")
	return c.SendString("Hello, " + name + " 👋!")
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

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		w.WriteString("[")
		enc := json.NewEncoder(w)
		n := 100000
		for i := 0; i < n; i++ {

			user := User{
				ID:        i + 1,
				Username:  "My Username",
				Lastname:  "My Lastname",
				Firstname: "My Firstname",
			}
			if err := enc.Encode(user); err != nil {
				continue
			}

			if i < n-1 {
				w.WriteString(",")
			}

			// w.Flush()
		}
		w.WriteString("]")
	})

	return nil
}
