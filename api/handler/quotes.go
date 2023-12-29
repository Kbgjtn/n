package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Kbgjtn/notethingness-api.git/api/model"
	"github.com/Kbgjtn/notethingness-api.git/types"
	"github.com/Kbgjtn/notethingness-api.git/util"
)

type QuotesResource struct {
	repo types.Repository
}

func New(r types.Repository) *QuotesResource {
	return &QuotesResource{r}
}

func (rs QuotesResource) Routes(route chi.Router) chi.Router {
	route.Get("/", rs.List)
	route.Post("/", rs.Create)
	route.Route("/{id}",
		func(r chi.Router) {
			r.Get("/", rs.Get)
			r.Delete("/", rs.Delete)
			r.Put("/", rs.Update)
		})

	return route
}

func (rs QuotesResource) List(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")
	p := types.Pageable{}.Parse(offset, limit)

	data, page, err := rs.repo.List(r.Context(), p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonData, err := json.Marshal(data.(model.Quotes).CreateResponseDto(page.(types.Pageable)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonData)
}

func (rs QuotesResource) Delete(w http.ResponseWriter, r *http.Request) {
	requestIDParam := r.Context().Value(chi.RouteCtxKey).(*chi.Context).URLParam("id")
	reqDTO, err := model.QuoteURLParams{}.Parse(requestIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = rs.repo.Delete(r.Context(), reqDTO)
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

func (rs QuotesResource) Get(w http.ResponseWriter, r *http.Request) {
	requestIDParam := r.Context().Value(chi.RouteCtxKey).(*chi.Context).URLParam("id")
	reqDTO, err := model.QuoteURLParams{}.Parse(requestIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := rs.repo.Get(r.Context(), reqDTO)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	jsonData, err := json.Marshal(data.(model.Quote).CreateResponseDto())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

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

	jsonData, err := json.Marshal(data.(model.Quote).CreateResponseDto())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

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

	data, err := rs.repo.Update(r.Context(), params, payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	jsonData, err := json.Marshal(data.(model.Quote).CreateResponseDto())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
