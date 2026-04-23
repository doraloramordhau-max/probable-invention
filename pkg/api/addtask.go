package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"probable-invention/pkg/db"
)

var errTitleRequired = errors.New("title is required")

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)

	if task.Date == "" {
		task.Date = today
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	next := ""
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeErrorStatus(w, http.StatusBadRequest, err)
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeErrorStatus(w, http.StatusInternalServerError, err)
		return
	}

	writeJSONStatus(w, http.StatusCreated, map[string]string{"id": strconv.FormatInt(id, 10)})
}
