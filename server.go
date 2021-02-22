package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/fabienbellanger/go-fiber/middlewares/timer"
	"github.com/fabienbellanger/go-fiber/ws"
	"github.com/fabienbellanger/goutils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/template/django"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
)

type server struct {
	store  store
	router *fiber.App
	mode   string
	logger *zap.Logger
}

func newServer(mode string, hub *ws.Hub) *server {
	s := &server{
		router: fiber.New(serverConfig()),
		mode:   mode,
	}

	s.initHTTPServer()

	// Pkger
	// -----
	s.router.Use("/assets", filesystem.New(filesystem.Config{
		Root: pkger.Dir("/public/assets"),
	}))

	s.initTools()
	s.routes()
	s.websocketRoutes(hub)
	s.initJWT()
	s.protectedRoutes()

	// Liste des routes
	// ----------------
	// s.displayRoutes()

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

	// Favicon
	// -------
	s.router.Use(favicon.New(favicon.Config{
		File: "favicon.png",
	}))

	// Logger
	// ------
	if s.mode != "production" {
		s.router.Use(logger.New(logger.Config{
			Next:         nil,
			Format:       "[${time}] ${status} - ${latency} - ${method} ${path}\n",
			TimeFormat:   "2006-01-02 15:04:05",
			TimeZone:     "Local",
			TimeInterval: 500 * time.Millisecond,
			Output:       os.Stderr,
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
		DisplaySeconds:      true,
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

func (s *server) initTools() {
	// Basic Auth
	// ----------
	cfg := basicauth.Config{
		Users: map[string]string{
			viper.GetString("debug.basicAuthUsername"): viper.GetString("debug.basicAuthPassword"),
		},
	}

	// Prometheus
	// ----------
	if viper.GetBool("debug.prometheus") {
		// metrics := s.router.Group("/metrics")
		// metrics.Use(basicauth.New(cfg))

		prometheus := fiberprometheus.New("go-fiber")
		prometheus.RegisterAt(s.router, "/metrics")
		s.router.Use(prometheus.Middleware)
	}

	// Pprof
	// -----
	if viper.GetBool("debug.pprof") {
		private := s.router.Group("/debug/pprof")
		private.Use(basicauth.New(cfg))
		s.router.Use(pprof.New())
	}

	if viper.GetBool("debug.monitor") {
		tools := s.router.Group("/tools")
		tools.Use(basicauth.New(cfg))
		tools.Get("/monitor", monitor.New())
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

func serverConfig() fiber.Config {
	// Initialize standard Go html template engine
	engine := django.NewFileSystem(pkger.Dir("/public/templates"), ".django")

	return fiber.Config{
		// Gestion des erreurs
		// -------------------
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			log.Printf("%+v\n", err)

			// Retreive the custom statuscode if it's an fiber.*Error
			e, ok := err.(*fiber.Error)
			if ok {
				code = e.Code
			}

			if e != nil {
				return c.JSON(e)
			}

			if code == fiber.StatusInternalServerError {
				// TODO: Logger l'erreur
				log.Printf("Error: %v\n", err)

				return c.Status(code).JSON(fiber.Map{
					"code":    code,
					"message": "Internal Server Error",
				})
			}

			return nil
		},
		Prefork:               viper.GetBool("server.prefork"),
		DisableStartupMessage: false,
		StrictRouting:         true,
		Views:                 engine,
	}
}

func (s *server) displayRoutes() {
	stask := s.router.Stack()
	for m := range stask {
		for r := range stask[m] {
			route := stask[m][r]
			if route.Method != "HEAD" && route.Method != "CONNECT" &&
				route.Method != "TRACE" && route.Method != "OPTIONS" {
				fmt.Printf("%v\t%v\t%v\n", route.Method, route.Path, route.Params)
			}
		}
	}
}
