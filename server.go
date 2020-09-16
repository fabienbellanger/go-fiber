package main

import (
	"os"
	"strings"
	"time"

	"github.com/fabienbellanger/go-fiber/middlewares/timer"
	"github.com/fabienbellanger/goutils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/spf13/viper"
)

type server struct {
	store  store
	router *fiber.App
	mode   string
}

func newServer() *server {
	s := &server{
		router: fiber.New(errorHandling()),
		mode:   "production",
	}

	s.initHTTPServer()
	s.routes()
	s.initPprof()
	s.initJWT()
	s.protectedRoutes()

	// Custom 404 (after all routes)
	// -----------------------------
	s.router.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"message": "Resource Not Found",
		})
	})

	return s
}

func (s *server) initHTTPServer() {
	// Mode
	// ----
	s.mode = viper.GetString("environment")

	// CORS
	// ----
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(viper.GetStringSlice("server.cors.allowOrigins"), ", "),
		AllowMethods:     strings.Join(viper.GetStringSlice("server.cors.allowMethods"), ", "),
		AllowHeaders:     strings.Join(viper.GetStringSlice("server.cors.allowHeaders"), ", "),
		ExposeHeaders:    strings.Join(viper.GetStringSlice("server.cors.exposeHeaders"), ", "),
		AllowCredentials: viper.GetBool("server.cors.allowCredentials"),
		MaxAge:           int(12 * time.Hour),
	}))

	// Logger
	// ------
	if s.mode != "production" {
		s.router.Use(logger.New(logger.Config{
			Next:       nil,
			Format:     "[${time}] ${status} - ${latency} - ${method} ${path}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "Local",
			Output:     os.Stderr,
		}))
	}

	// Recover
	// -------
	s.router.Use(recover.New())

	// Request ID
	// ----------
	s.router.Use(requestid.New())

	// Timer
	// -----
	s.router.Use(timer.New(timer.Config{
		DisplayMilliseconds: true,
		DisplayHuman:        true,
	}))

	// Limiter
	// -------
	if viper.GetBool("server.limiter.enable") {
		s.router.Use(limiter.New(limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				excludedIP := viper.GetStringSlice("server.limiter.excludedIP")
				if len(excludedIP) == 0 {
					return false
				}
				return goutils.StringInSlice(c.IP(), excludedIP)
			},
			Max:      viper.GetInt("server.limiter.max"),
			Duration: viper.GetDuration("server.limiter.duration") * time.Second,
			Key: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"code":    fiber.StatusTooManyRequests,
					"message": "Too Many Requests",
				})
			},
		}))
	}
}

func (s *server) initPprof() {
	if viper.GetBool("debug.pprof") {
		private := s.router.Group("/debug/pprof")

		// Basic Auth
		// ----------
		cfg := basicauth.Config{
			Users: map[string]string{
				viper.GetString("debug.basicAuthUsername"): viper.GetString("debug.basicAuthPassword"),
			},
		}
		private.Use(basicauth.New(cfg))

		// pprof (The handled paths all begin with /debug/pprof/)
		// ------------------------------------------------------
		s.router.Use(pprof.New())
	}
}

func (s *server) initJWT() {
	s.router.Use(jwtware.New(jwtware.Config{
		SigningMethod: "HS512",
		SigningKey:    []byte(viper.GetString("jwt.secret")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "Invalid or expired JWT",
			})
		},
	}))
}

func errorHandling() fiber.Config {
	return fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// Retreive the custom statuscode if it's an fiber.*Error
			e, ok := err.(*fiber.Error)

			if ok {
				code = e.Code
			}

			if e != nil {
				return c.JSON(e)
			}

			if code == fiber.StatusInternalServerError {
				return c.JSON(fiber.Map{
					"code":    code,
					"message": "Internal Server Error",
				})
			}

			return nil
		},
	}
}
