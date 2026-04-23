package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"probable-invention/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {
	writeJSONStatus(w, http.StatusOK, data)
}

func writeJSONStatus(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, db.ErrTaskNotFound) {
		status = http.StatusNotFound
	}
	writeErrorStatus(w, status, err)
}

func writeErrorStatus(w http.ResponseWriter, status int, err error) {
	writeJSONStatus(w, status, map[string]string{"error": err.Error()})
}
