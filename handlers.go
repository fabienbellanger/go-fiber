package main

import (
	"bufio"
	"encoding/json"
	"runtime"
	"sync"

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

	numProjects := len(projects)
	jobs := make(chan models.Project, numProjects)
	results := make(chan models.Release, numProjects)

	// Nombre de workers
	// -----------------
	numWorkers := runtime.NumCPU()
	if numProjects < numWorkers {
		numWorkers = numProjects
	}

	// Lancement des workers
	// ---------------------
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			models.ReleaseWorker(jobs, results)
		}()
	}

	// Fermeture du channel results quand tous les workers ont terminés
	// ----------------------------------------------------------------
	go func() {
		defer close(results)
		wg.Wait()
	}()

	// Envoi des jobs
	// --------------
	go func() {
		defer close(jobs)
		for _, project := range projects {
			jobs <- project
		}
	}()

	// Traitement des résultats
	// ------------------------
	releases := make([]models.Release, 0)
	for r := range results {
		releases = append(releases, r)
	}

	return c.JSON(&releases)
}
