package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	err := db.DeleteTask(id)
	if err != nil {
		Logger.Printf("получена ошибка из db.DeleteTask(id): %v", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}
	writeJson(w, map[string]any{})
}
