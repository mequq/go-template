package handler

import (
	"application/internal/biz"
	"application/internal/service"
	"application/internal/service/dto"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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
	logger := h.logger.With("method", "Create")
	ctx := r.Context()
	logger.DebugContext(ctx, "Create placeholder")

	placeholder := new(dto.CreatePlaceholderReq)
	if err := json.NewDecoder(r.Body).Decode(placeholder); err != nil {
		logger.WarnContext(ctx, "failed to decode request body", "error", err)
		dto.HandleError(err, w)

		return
	}

	id, err := h.placeholder.Create(ctx, placeholder.Name)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create placeholder", "error", err)
		dto.HandleError(err, w)

		return
	}

	resp := dto.PlaceholderResp{
		ID:   id.String(),
		Name: placeholder.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
		dto.HandleError(err, w)

		return
	}
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
	logger := h.logger.With("method", "List")
	ctx := r.Context()
	logger.DebugContext(ctx, "List placeholders")

	placeholders, err := h.placeholder.List(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to list placeholders", "error", err)
		dto.HandleError(err, w)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := dto.ToPlaceholderResps(placeholders)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
		panic(err)
	}

	logger.InfoContext(ctx, "Listed placeholders", "count", resp.Count)
}

// update implements the endpoint for updating a specific placeholder by ID.
//
//	@Summary		Update a placeholder
//	@Description	Update the details of a specific placeholder by ID.
//	@Tags			Placeholders
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string						true	"Placeholder UUID"
//	@Param			placeholder	body		dto.CreatePlaceholderReq	true	"Updated placeholder details"
//	@Success		200			{object}	dto.PlaceholderResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404			{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/mocks/placeholders/{id} [put]
func (h *placeholder) update(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Update")
	ctx := r.Context()
	logger.DebugContext(ctx, "Update placeholder")

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		logger.WarnContext(ctx, "invalid UUID format", "error", err)
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	placeholder := new(dto.CreatePlaceholderReq)
	if err := json.NewDecoder(r.Body).Decode(placeholder); err != nil {
		logger.WarnContext(ctx, "failed to decode request body", "error", err)
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	if err := h.placeholder.Update(ctx, id, placeholder.Name); err != nil {
		logger.ErrorContext(ctx, "failed to update placeholder", "error", err)
		dto.HandleError(err, w)

		return
	}

	resp := dto.PlaceholderResp{
		ID:   id.String(),
		Name: placeholder.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
		panic(err)
	}
}

// delete implements the endpoint for deleting a specific placeholder by ID.
//
//	@Summary		Delete a placeholder
//	@Description	Delete a specific placeholder by ID.
//	@Tags			Placeholders
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Placeholder UUID"
//	@Success		204	""
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/mocks/placeholders/{id} [delete]
func (h *placeholder) delete(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Delete")
	ctx := r.Context()
	logger.DebugContext(ctx, "Delete placeholder")

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		logger.WarnContext(ctx, "invalid UUID format", "error", err)
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	if err := h.placeholder.Delete(ctx, id); err != nil {
		logger.ErrorContext(ctx, "failed to delete placeholder", "error", err)
		dto.HandleError(err, w)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// get implements the endpoint for retrieving a specific placeholder by ID.
//
//	@Summary		Get a placeholder
//	@Description	Retrieve the details of a specific placeholder by ID.
//	@Tags			Placeholders
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Placeholder UUID"
//	@Success		200	{object}	dto.PlaceholderResp
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/mocks/placeholders/{id} [get]
func (h *placeholder) get(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Get")
	ctx := r.Context()
	logger.DebugContext(ctx, "Get placeholder")

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		logger.WarnContext(ctx, "invalid UUID format", "error", err)
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	placeholder, err := h.placeholder.Get(ctx, id)
	if err != nil {
		logger.ErrorContext(ctx, "failed to get placeholder", "error", err)
		dto.HandleError(err, w)

		return
	}

	resp := dto.PlaceholderResp{
		ID:   placeholder.ID.String(),
		Name: placeholder.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
		panic(err)
	}
}
