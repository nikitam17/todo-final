package api

import (
	"net/http"
	"todo-final/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	tasks, err := db.Tasks(50, search) // в параметре максимальное количество записей
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
