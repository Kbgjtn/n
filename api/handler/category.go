package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Kbgjtn/notethingness-api.git/api/model"
	"github.com/Kbgjtn/notethingness-api.git/api/repository"
	"github.com/Kbgjtn/notethingness-api.git/types"
	"github.com/Kbgjtn/notethingness-api.git/util"
)

type CategoryResource struct {
	repo *repository.CategoryRepository
}

func NewCategory(repo *repository.CategoryRepository) *CategoryResource {
	return &CategoryResource{repo}
}

func (rs CategoryResource) Routes(route chi.Router) {
	route.Get("/", rs.List)
	route.Post("/", rs.Create)
	route.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Put("/", rs.Update)
		r.Delete("/", rs.Delete)
	})
}

// Get return a category
// @Summary Get By ID
// @Description Get a category by id
// @Tags category
// @Accept json
// @Produce json
// @Param id path string true "Quote ID"
// @Success 200 {object} model.Category
// @Failure 400 {string} string "Bad Request: id is invalid or missing"
// @Router /categories/{id} [get]
// !curl localhost:3000/api/categories/1 | jq
func (rs CategoryResource) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	args, err := model.ParseParams(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := rs.repo.Get(r.Context(), args)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if result.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json, err := json.Marshal(result.ToJSON(200, "Success"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error: failed to marshal category"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// List return a list of categories
// @Summary Get list
// @Description Get List of categories
// @Tags category
// @Accept json
// @Produce json
// @Param offset query string false "string default example" default(0) example(1)
// @Param limit query string false "string default example" default(10) example(20)
// @Success 200 {object} types.JSONResult{data=model.Categories}
// @Failure 400 {string} string "Bad Request: error message"
// @Router /categories [get]
// !curl localhost:3000/api/categories | jq
func (rs CategoryResource) List(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	p := types.Pageable{}.Parse(limit, offset)

	result, err := rs.repo.List(r.Context(), &p)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error: failed to get categories"))
		return
	}

	data, err := json.Marshal(result.ToJSON(p))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("error: failed to parsing to JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Create a category
// @Summary Create a new category
// @Description Create a new category
// @Tags category
// @Accept json
// @Produce json
// @Param request body model.CategoryRequestPayload true "default"
// @Success 201 {object} types.JSONResult{data=model.Category}
// @Failure 400 {string} string "Bad Request: label is invalid or missing"
// @Router /categories [post]
// !curl -v 'POST' localhost:3000/api/categories -d '{"label":"test"}' -H "Content-Type: application/json" | jq
func (rs CategoryResource) Create(w http.ResponseWriter, r *http.Request) {
	var payload model.CategoryRequestPayload
	err := util.ParseRequestBody(r, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err = payload.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := rs.repo.Create(r.Context(), payload.Label)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(result.ToJSON(201, "Created"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error: failed to marshal category"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

// Delete a category
// @Summary Delete a category
// @Description Delete a category
// @Tags category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request: id is invalid or missing"
// @Router /categories/{id} [delete]
// !curl -v -X DELETE localhost:3000/api/categories/1 | jq
func (rs CategoryResource) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	args, err := model.ParseParams(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = rs.repo.Delete(r.Context(), args); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// Update a category
// @Summary Update a category
// @Description Update a category
// @Tags category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body model.CategoryRequestPayload true "default"
// @Success 200 {object} types.JSONResult{data=model.Category}
// @Failure 400 {string} string "Bad Request: id is invalid or missing"
// @Router /categories/{id} [put]
// !curl -v -X PUT localhost:3000/api/categories/1 -d '{"label":"test"}' -H "Content-Type: application/json" | jq
func (rs CategoryResource) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	args, err := model.ParseParams(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var payload model.CategoryRequestPayload
	if err = util.ParseRequestBody(r, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err = payload.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := rs.repo.Update(r.Context(), args, payload.Label)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(result.ToJSON(200, "Success"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error: failed to marshal category"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
