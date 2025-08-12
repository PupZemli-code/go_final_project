package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Проверка id
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

	task, err := db.GetTask(id)
	if err != nil {
		Logger.Printf("ошибка в TaskDoneHandler: GetTask: %v", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	// Если правило повторения отсутствует,
	// удолит задачу
	if task.Repeat == "" {
		err := db.DeleteTask(task.ID)
		if err != nil {
			Logger.Printf("ошибка в TaskDoneHandler: DeleteTask: %v", err)
			sendError(w, http.StatusInternalServerError, err)
			return
		}
		writeJson(w, map[string]any{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		Logger.Printf("ошибка в TaskDoneHandler: NextDate: %v", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	err = db.UpdateDate(next, id)

	if err != nil {
		Logger.Printf("ошибка в TaskDoneHandler: UpdateDate: %v", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	writeJson(w, map[string]any{})
}
