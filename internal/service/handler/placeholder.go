package handler

import (
	"application/internal/biz"
	"application/internal/service"
	"context"
	"log/slog"
	"net/http"
)

type placeholder struct {
	logger      *slog.Logger
	mux         *http.ServeMux
	placeholder biz.UsecasePlaceholder
}

var _ service.Handler = (*placeholder)(nil)

// NewPlaceholder creates a new instance of the Placeholder handler.
func NewPlaceholder(logger *slog.Logger, mux *http.ServeMux, puc biz.UsecasePlaceholder) *placeholder {
	return &placeholder{
		logger:      logger.With("layer", "Placeholder"),
		mux:         mux,
		placeholder: puc,
	}
}

// RegisterHandler registers the Placeholder handler with the given service.
func (h *placeholder) RegisterHandler(ctx context.Context) error {
	// List of placeholder endpoints
	h.mux.HandleFunc("GET /apis/mocks/placeholders", h.list)
	// Get a specific placeholder by ID
	h.mux.HandleFunc("GET /apis/mocks/placeholders/{id}", h.get)
	// Create a new placeholder
	h.mux.HandleFunc("POST /apis/mocks/placeholders", h.create)
	// Update a specific placeholder by ID
	h.mux.HandleFunc("PUT /apis/mocks/placeholders/{id}", h.update)
	// Delete a specific placeholder by ID
	h.mux.HandleFunc("DELETE /apis/mocks/placeholders/{id}", h.delete)

	return nil
}

// create implements the endpoint for creating a new placeholder.
//
//	@Summary		Create a new placeholder
//	@Description	Create a new placeholder with the provided details.
//	@Tags			Placeholders
//	@Accept			json
//	@Produce		json
//	@Param			placeholder	body		dto.CreatePlaceholderReq	true	"Placeholder details"
//	@Success		201			{object}	dto.PlaceholderResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/mocks/placeholders [post]
func (h *placeholder) create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// list implements the endpoint for listing placeholders.
//
//	@Summary		List placeholders
//	@Description	Retrieve a list of all placeholders.
//	@Tags			Placeholders
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.PlaceholderListResponse	"ok"
//	@Failure		500	{object}	dto.ErrorResponse			"Internal Server Error"
//	@Router			/apis/mocks/placeholders [get]
func (h *placeholder) list(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// update implements the endpoint for updating a specific placeholder by ID.
//
//	@Summary		Update a placeholder
//	@Description	Update the details of a specific placeholder by ID.
//	@Tags			Placeholders
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int64						true	"Placeholder ID"
//	@Param			placeholder	body		dto.CreatePlaceholderReq	true	"Updated placeholder details"
//	@Success		200			{object}	dto.PlaceholderResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404			{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/mocks/placeholders/{id} [put]
func (h *placeholder) update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// delete implements the endpoint for deleting a specific placeholder by ID.
//
//	@Summary		Delete a placeholder
//	@Description	Delete a specific placeholder by ID.
//	@Tags			Placeholders
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int64	true	"Placeholder ID"
//	@Success		204	""
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/mocks/placeholders/{id} [delete]
func (h *placeholder) delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// get implements the endpoint for retrieving a specific placeholder by ID.
//
//	@Summary		Get a placeholder
//	@Description	Retrieve the details of a specific placeholder by ID.
//	@Tags			Placeholders
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Placeholder ID"
//	@Success		200	{object}	dto.PlaceholderResp
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/mocks/placeholders/{id} [get]
func (h *placeholder) get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
