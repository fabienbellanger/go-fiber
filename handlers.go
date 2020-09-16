package main

import (
	"bufio"

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

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		w.WriteString("[")
		// enc := json.NewEncoder(w)
		n := 100000
		for i := 0; i < n; i++ {
			if i > 0 {
				w.WriteString(",")
			}

			w.WriteString(`{"id": 1, "username": "My Username", "lastname": "My Lastname", "firstname": "My Firstname"}`)
			w.Flush()

			// user := User{
			// 	ID:        i + 1,
			// 	Username:  "My Username",
			// 	Lastname:  "My Lastname",
			// 	Firstname: "My Firstname",
			// }
			// if err := enc.Encode(user); err != nil {
			// 	println("error")
			// }
		}
		w.WriteString("]")
	})

	return nil
}
