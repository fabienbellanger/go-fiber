package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/fabienbellanger/go-fiber/models"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (s *server) handlerStatic(c *fiber.Ctx) error {
	s.logger.Debug("Index route", zap.String("url", "https://www.google.com"))
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

func (s *server) handlerGithub(c *fiber.Ctx) error {
	projects, err := models.LoadProjectsFromFile("projects.json")
	if err != nil {
		return err
	}

	// Version non concurrente
	// -----------------------
	releases := make([]models.Release, 0)
	for _, project := range projects {
		release, err := project.GetInformation()
		if err == nil {
			releases = append(releases, release)
		}
	}

	// Version concurrente
	// -------------------
	numCPU := runtime.NumCPU()
	fmt.Printf("NumCPU=%v\n", numCPU)

	// Pour le pool de workers
	// https://www.prakharsrivastav.com/posts/golang-concurrent-worker-pool/
	// --> https://github.com/PrakharSrivastav/workers
	// https://brandur.org/go-worker-pool
	// https://medium.com/@j.d.livni/write-a-go-worker-pool-in-15-minutes-c9b42f640923

	return c.JSON(&releases)
}
