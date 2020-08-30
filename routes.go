package main

func (s *server) routes() {
	v1 := s.router.Group("/v1")

	v1.Get("/", s.handlerHome)
	v1.Get("/home", s.handlerHome)
	v1.Get("/json", s.handlerBigJSON)

	// Users routes
	// ------------
	users := v1.Group("/users")
	users.Get("/", s.handlerGetUser)
}
