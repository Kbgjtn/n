package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Kbgjtn/notethingness-api.git/api/model"
	repo "github.com/Kbgjtn/notethingness-api.git/api/repository"
	"github.com/Kbgjtn/notethingness-api.git/types"
	"github.com/Kbgjtn/notethingness-api.git/util"
)

type QuotesResource struct {
	repo *repo.QuoteRepository
}

func New(r *repo.QuoteRepository) *QuotesResource {
	return &QuotesResource{r}
}

func (rs QuotesResource) Routes(route chi.Router) {
	route.Get("/", rs.List)
	route.Post("/", rs.Create)
	route.Route("/{id}",
		func(r chi.Router) {
			r.Get("/", rs.Get)
			r.Delete("/", rs.Delete)
			r.Put("/", rs.Update)
		})
}

// Get returns a quote by id
// @Summary Get a quote
// @Description Get a quote by id
// @Tags quote
// @Accept  json
// @Produce  json
// @Param id path string true "Quote ID"
// @Success 200 {object} types.JSONResult{data=model.Quote}
// @Failure 400 {string} string "error: id is invalid"
// @Failure 404 {string} string "error: quote not found"
// @Router /quotes/{id} [get]
func (rs QuotesResource) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var reqDTO model.QuoteURLParams
	err := reqDTO.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := rs.repo.Get(r.Context(), reqDTO)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	jsonData, err := json.Marshal(data.CreateResponseDto())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// List returns a list of quotes
// @Summary List quotes

// @Description Get List quotes
// @Tags quote
// @Accept  json
// @Produce  json
// @Param offset query string false "string default example" default(0) example(1)
// @Param limit query string false "string default example" default(10) example(20)
// @Success 200 {object} types.JSONResult{data=model.Quotes,paginate=types.Pageable,length=int}
// @Failure 400 {string} string "error: offset or limit is invalid"
// @Router /quotes [get]
func (rs QuotesResource) List(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")
	p := types.Pageable{}.Parse(limit, offset)

	data, err := rs.repo.List(r.Context(), &p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonData, err := json.Marshal(data.CreateResponseDto(&p))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonData)
}

// Delete deletes a quote by id
// @Summary Delete a quote
// @Description Delete a quote by id
// @Tags quote
// @Accept  json
// @Produce  json
// @Param id path string true "Quote ID"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "error: id is invalid"
// @Failure 404 {string} string "error: quote not found"
// @Router /quotes/{id} [delete]
func (rs QuotesResource) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var reqDTO model.QuoteURLParams
	err := reqDTO.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = rs.repo.Delete(r.Context(), reqDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(
			[]byte(
				"error: Cannot delete the quote with id " + id + " because it is not a valid",
			),
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Create creates a quote
// @Summary Create a quote
// @Description Create a quote
// @Tags quote
// @Accept  json
// @Produce  json
// @Param request body model.QuoteRequestPayload true "default"
// @Success 201 {object} types.JSONResult{data=model.Quote}
// @Failure 400 {string} string "Bad Request: Invalid payload"
// @Router /quotes [post]
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

	data, err := rs.repo.Create(r.Context(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonData, err := json.Marshal(data.CreateResponseDto())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// Create creates a quote
// @Summary Create a quote
// @Description Create a quote
// @Tags quote
// @Accept  json
// @Produce  json
// @Param request body model.QuoteRequestPayload true "default"
// @Success 200 {object} types.JSONResult{data=model.Quote}
// @Failure 400 {string} string "Bad Request: Invalid payload"
// @Router /quotes/{id} [put]
func (rs QuotesResource) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var reqDTO model.QuoteURLParams
	err := reqDTO.Parse(id)
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

	data, err := rs.repo.Update(r.Context(), reqDTO, payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	jsonData, err := json.Marshal(data.CreateResponseDto())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
