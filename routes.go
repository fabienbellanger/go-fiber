package main

func (s *server) routes() {
	// Login
	// -----
	s.router.Post("/login", s.handlerUserLogin)

	// Github projects
	// ---------------
	s.router.Get("releases", s.handlerReleases)

	// API v1
	// ------
	v1 := s.router.Group("/v1")

	v1.Get("/", s.handlerHome)
	v1.Get("/home/:name", s.handlerHome)
	v1.Get("/json", s.handlerBigJSON)
	v1.Get("/json-stream", s.handlerBigJSONStream)
	v1.Get("/github", s.handlerGithub)

	// Users routes
	// ------------
	users := v1.Group("/users")
	users.Get("/", s.handlerGetUser)
}

func (s *server) protectedRoutes() {
	protected := s.router.Group("/protected")

	protected.Get("/test", s.handlerProtectedTest)
}
