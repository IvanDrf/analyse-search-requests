package http

import (
	"encoding/json"
	"net/http"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

func isContentTypeJSON(req *http.Request) bool {
	return req.Header.Get("Content-Type") == "application/json"
}

func writeError(w http.ResponseWriter, status int, err *models.Error) {
	http.Error(w, err.Error(), status)
}

func writeResponse(w http.ResponseWriter, status int, response any) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	if response != nil {
		json.NewEncoder(w).Encode(response)
	}
}
