package main

import (
	"net/http"
	"testing"

	"github.com/fabienbellanger/go-fiber/ws"
)

func BenchmarkHome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hub := ws.NewHub()
		go hub.Run()
		s := newServer("production", hub)
		s.mode = "production" // Ne fonctionne pas car le logger est initialisé lors du newServer()

		req, _ := http.NewRequest(
			"GET",
			"/v1",
			nil,
		)
		s.router.Test(req, -1)
	}
}

func BenchmarkJsonStream(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hub := ws.NewHub()
		go hub.Run()
		s := newServer("production", hub)
		s.mode = "production" // Ne fonctionne pas car le logger est initialisé lors du newServer()

		req, _ := http.NewRequest(
			"GET",
			"/v1/json-stream",
			nil,
		)
		s.router.Test(req, -1)
	}
}

func BenchmarkJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hub := ws.NewHub()
		go hub.Run()
		s := newServer("production", hub)
		s.mode = "production" // Ne fonctionne pas car le logger est initialisé lors du newServer()

		req, _ := http.NewRequest(
			"GET",
			"/v1/json",
			nil,
		)
		s.router.Test(req, -1)
	}
}
