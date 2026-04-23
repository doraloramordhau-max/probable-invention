package api

import (
	"net/http"

	"probable-invention/pkg/db"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, _ *http.Request) {
	tasks, err := db.Tasks(tasksLimit)
	if err != nil {
		writeErrorStatus(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
