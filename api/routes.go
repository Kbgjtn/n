package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Kbgjtn/notethingness-api.git/api/handler"
	"github.com/Kbgjtn/notethingness-api.git/api/repository"
	"github.com/Kbgjtn/notethingness-api.git/api/service"
)

// Routes initialize API routes
func (s *Server) Routes() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/api", s.InitRoute)
	s.router = router
}

// InitRoute initialize API routes
func (s *Server) InitRoute(router chi.Router) {
	repo := repository.New(s.db)
	service := service.New(repo)

	// Quotes resource routes
	router.Mount("/quotes", handler.New(service).Routes())
}
