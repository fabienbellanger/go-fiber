package main

import (
	"log"

	"github.com/fabienbellanger/go-fiber/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (s *server) routes() {
	s.router.Get("static", s.handlerStatic)

	// Login
	// -----
	s.router.Post("/login", s.handlerUserLogin)

	// API v1
	// ------
	v1 := s.router.Group("/v1")

	v1.Get("/", s.handlerHome)
	v1.Get("/home/:name", s.handlerHome)
	v1.Get("/json", s.handlerBigJSON)
	v1.Get("/json-stream", s.handlerBigJSONStream)

	// Users routes
	// ------------
	users := v1.Group("/users")
	users.Get("/", s.handlerGetUser)
}

func (s *server) protectedRoutes() {
	protected := s.router.Group("/protected")

	protected.Get("/test", s.handlerProtectedTest)
}

func (s *server) websocketRoutes() {
	hub := ws.NewHub()
	go hub.Run()

	w := s.router.Group("/ws")

	s.router.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Access the websocket server: ws://localhost:8888/ws/123?v=1.0
	// https://www.websocket.org/echo.html
	w.Get("/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Printf("allowed: %v, params: %v, query: %v\n", c.Locals("allowed"), c.Params("id"), c.Query("v"))

		ws.ServeWs(hub, c.Conn)

		// cli := ws.New(c)
		// cli.Connect()

		// var (
		// 	mt  int
		// 	msg []byte
		// 	err error
		// )
		// for {
		// 	if mt, msg, err = c.ReadMessage(); err != nil {
		// 		log.Printf("[error] read: %v, type=%v, msg=%v", err, mt, msg)
		// 		break
		// 	}

		// 	log.Printf("recv: type=%v, msg=%s", mt, msg)

		// 	if err = c.WriteMessage(mt, msg); err != nil {
		// 		log.Println("[error] write:", err)
		// 		break
		// 	}
		// }
	}))
}
