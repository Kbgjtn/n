package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Kbgjtn/notethingness-api.git/api/model"
	"github.com/Kbgjtn/notethingness-api.git/api/service"
	"github.com/Kbgjtn/notethingness-api.git/types"
	"github.com/Kbgjtn/notethingness-api.git/util"
)

// QuotesResource implements REST API for quotes
type QuotesResource struct {
	service *service.QuotesService
}

// New creates new QuotesResource
func New(service *service.QuotesService) *QuotesResource {
	return &QuotesResource{service}
}

// Routes returns quotes resource routes
func (rs QuotesResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", rs.List)
	r.Post("/", rs.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Delete("/", rs.Delete)
		r.Put("/", rs.Update)
	})
	return r
}

// List returns list of quotes GET /quotes
func (rs QuotesResource) List(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")
	page := types.Pageable{}.Parse(limit, offset)
	quotes, err := rs.service.List(r.Context(), page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonData, err := json.Marshal(quotes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonData)
}

// Delete deletes quote DELETE /quotes/{id}
func (rs QuotesResource) Delete(w http.ResponseWriter, r *http.Request) {
	requestIDParam := r.Context().Value(chi.RouteCtxKey).(*chi.Context).URLParam("id")
	reqDTO, err := model.QuoteURLParams{}.Parse(requestIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = rs.service.Delete(r.Context(), &reqDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(
			[]byte(
				"error: Cannot delete the quote with id " + requestIDParam + " because it is not a valid",
			),
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Get returns quote by id GET /quotes/{id}
func (rs QuotesResource) Get(w http.ResponseWriter, r *http.Request) {
	requestIDParam := r.Context().Value(chi.RouteCtxKey).(*chi.Context).URLParam("id")
	reqDTO, err := model.QuoteURLParams{}.Parse(requestIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rs.service.Get(r.Context(), &reqDTO)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Create creates new quote POST /quotes
func (rs QuotesResource) Create(w http.ResponseWriter, r *http.Request) {
	var payload model.QuoteRequestPayload
	if err := util.ParseRequestBody(r, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: payload is invalid or missing"))
		return
	}

	if _, err := payload.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := rs.service.Create(r.Context(), &payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// Update updates quote PUT /quotes/{id}
func (rs QuotesResource) Update(w http.ResponseWriter, r *http.Request) {
	requestIDParam := r.Context().Value(chi.RouteCtxKey).(*chi.Context).URLParam("id")
	params, err := model.QuoteURLParams{}.Parse(requestIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var payload model.QuoteRequestPayload
	if err = util.ParseRequestBody(r, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: payload is invalid or missing"))
		return
	}

	if _, err = payload.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := rs.service.Update(r.Context(), &params, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
