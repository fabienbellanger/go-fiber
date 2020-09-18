package main

import (
	"net/http"
	"testing"
)

func BenchmarkJsonStream(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// // http.Request
		// req := httptest.NewRequest("GET", "http://google.com", nil)
		// req.Header.Set("X-Custom-Header", "hi")

		// // http.Response
		// resp, _ := app.Test(req)

		// // Do something with results:
		// if resp.StatusCode == 200 {
		// 	body, _ := ioutil.ReadAll(resp.Body)
		// 	fmt.Println(string(body)) // => Hello, World!
		// }
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
