package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"probable-invention/pkg/db"
)

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeErrorStatus(w, http.StatusBadRequest, err)
		return
	}

	if task.ID == "" {
		writeErrorStatus(w, http.StatusBadRequest, errIDRequired)
		return
	}

	if task.Title == "" {
		writeErrorStatus(w, http.StatusBadRequest, errTitleRequired)
		return
	}

	if err := checkDate(&task); err != nil {
		writeErrorStatus(w, http.StatusBadRequest, err)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, errIDRequired)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]string{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorStatus(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

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

	if task.Repeat == "" {
		if err = db.DeleteTask(id); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, err)
		return
	}

	if err = db.UpdateDate(next, id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]string{})
}
