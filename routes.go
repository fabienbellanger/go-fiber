package main

func (s *server) routes() {
	s.router.Get("/", s.handlerHome)
	s.router.Get("/json", s.handlerBigJson)
}
