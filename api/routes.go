package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Kbgjtn/notethingness-api.git/api/handler"
	"github.com/Kbgjtn/notethingness-api.git/api/repository"
	"github.com/Kbgjtn/notethingness-api.git/types"
)

func (s *Server) Routes() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/api", s.InitRoute)
	s.router = router
}

func (s *Server) InitRoute(router chi.Router) {
	var r types.Repository = repository.NewQuoteRepo(s.db)
	router.Mount("/quotes", handler.New(r).Routes(router))
}
