package api

import (
	"errors"
	"net/http"
	"strconv"
)

// Реализует запрос r.Delete"/api/task"
func (t TaskService) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		err := errors.New("DeleteTaskHandler получено пустое значение для id")
		Logger.Println(err)
		sendError(w, http.StatusBadRequest, err)
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		Logger.Println(err)
		sendError(w, http.StatusBadRequest, err)
		return
	}

	err := t.store.DeleteTask(id)
	if err != nil {
		Logger.Printf("получена ошибка из db.DeleteTask(id): %v", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}
	writeJson(w, map[string]any{})
}
