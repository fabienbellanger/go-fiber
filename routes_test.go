package main

import (
	"net/http"
	"testing"
)

func BenchmarkJsonStream(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := newServer()
		s.mode = "production" // Ne fonctionne pas car le logger est initialisé lors du newServer()

		req, _ := http.NewRequest(
			"GET",
			"/v1/json-stream",
			nil,
		)
		s.router.Test(req, -1)
	}
}
