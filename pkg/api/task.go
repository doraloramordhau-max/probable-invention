package api

import (
	"errors"
	"net/http"

	"probable-invention/pkg/db"
)

var errIDRequired = errors.New("id is required")

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeErrorStatus(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, errIDRequired)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, task)
}
