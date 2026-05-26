package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/ports/service"
)

type handlers struct {
	searchService service.SearchService
	requestTime   time.Duration
}

func (h *handlers) close() {
	h.searchService.Close()
	slog.Info("Handlers:close", slog.String("status", "successfully closed handlers"))
}

func NewHandlers(searchService service.SearchService, requestTime time.Duration) *handlers {
	return &handlers{
		searchService: searchService,
		requestTime:   requestTime,
	}
}

func (h *handlers) saveBadWord(w http.ResponseWriter, req *http.Request) {
	if !isContentTypeJSON(req) {
		writeError(w, http.StatusUnsupportedMediaType, &models.Error{
			Message: "invalid medida type, supported: JSON",
		})
		return
	}

	badWord := models.BadWord{}
	if err := json.NewDecoder(req.Body).Decode(&badWord); err != nil {
		writeError(w, http.StatusUnprocessableEntity, &models.Error{
			Message: "invalid request body",
		})
		return
	}
	defer req.Body.Close()

	ctx, cancel := context.WithTimeout(req.Context(), h.requestTime)
	defer cancel()

	err := h.searchService.SaveBadWord(ctx, badWord.Word)
	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeInternal:
			writeError(w, http.StatusInternalServerError, &e)

		case models.ErrCodeInvalidArgument:
			writeError(w, http.StatusBadRequest, &e)

		default:
			writeError(w, http.StatusBadRequest, &models.Error{
				Message: "unexpected error",
			})
		}

		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}

func (h *handlers) deleteBadWord(w http.ResponseWriter, req *http.Request) {
	if !isContentTypeJSON(req) {
		writeError(w, http.StatusUnsupportedMediaType, &models.Error{
			Message: "invalid media type, supported: JSON",
		})
		return
	}

	badWord := models.BadWord{}
	if err := json.NewDecoder(req.Body).Decode(&badWord); err != nil {
		writeError(w, http.StatusUnprocessableEntity, &models.Error{
			Message: "invalid request body",
		})
		return
	}
	defer req.Body.Close()

	ctx, cancel := context.WithTimeout(req.Context(), h.requestTime)
	defer cancel()

	err := h.searchService.DeleteBadWord(ctx, badWord.Word)
	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeInternal:
			writeError(w, http.StatusInternalServerError, &e)
		case models.ErrCodeInvalidArgument:
			writeError(w, http.StatusBadRequest, &e)

		default:
			writeError(w, http.StatusBadRequest, &models.Error{
				Message: "unexpected error",
			})
		}

		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}

func (h *handlers) findMostPopularSearches(w http.ResponseWriter, req *http.Request) {
	limitStr := req.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, &models.Error{
			Message: "requires query param 'limit' type of integer",
		})
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), h.requestTime)
	defer cancel()

	searches, err := h.searchService.FindMostPopularSearches(ctx, limit)
	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeInternal:
			writeError(w, http.StatusInternalServerError, &e)

		case models.ErrCodeInvalidArgument:
			writeError(w, http.StatusBadRequest, &e)

		default:
			writeError(w, http.StatusBadRequest, &models.Error{
				Message: "unexpected error",
			})
		}

		return
	}

	writeResponse(w, http.StatusOK, searches)
}
