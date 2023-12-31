package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Kbgjtn/notethingness-api.git/api/model"
	"github.com/Kbgjtn/notethingness-api.git/util"
)

type CategoryResource struct{}

func NewCategory() *CategoryResource {
	return &CategoryResource{}
}

func (rs CategoryResource) Routes(route chi.Router) {
	route.Get("/", rs.List)
	route.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get)
	})
}

// Get return a category
// @Summary Get By ID
// @Description Get a category by id
// @Tags category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} model.Category
// @Failure 400 {string} string "error: id is invalid or missing"
// @Router /categories/{id} [get]
func (rs CategoryResource) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: id is invalid or missing"))
		return
	}

	newCategory := model.Category{
		ID:    idInt,
		Label: "Category " + id,
	}

	json, err := json.Marshal(newCategory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error: failed to marshal category"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// List return a list of categories
// @Summary Get list
// @Description Get List of categories
// @Tags category
// @Accept json
// @Produce json
// @Param Pagination query string false "Pagination"
// @Success 200 {object} model.Category
// @Router /categories [get]
func (rs CategoryResource) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List Category"))
}

func (rs CategoryResource) Create(w http.ResponseWriter, r *http.Request) {
	var payload model.CategoryRequestPayload
	if err := util.ParseRequestBody(r, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: payload is invalid or missing"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Create Category"))
}
