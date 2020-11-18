package main

import (
	"bufio"
	"encoding/json"
	"time"

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

			w.Flush()
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

	models.CachedReleases.Mux.Lock()
	defer models.CachedReleases.Mux.Unlock()
	now := time.Now()
	if len(models.CachedReleases.Releases) == 0 || models.CachedReleases.ExpireAt.Before(now) {
		releases, err := models.ReleasesProcess(projects)
		if err != nil {
			return err
		}

		models.CachedReleases.Releases = releases
		models.CachedReleases.ExpireAt = now.Local().Add(time.Hour)
	}

	return c.JSON(&models.CachedReleases.Releases)
}
