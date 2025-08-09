package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("GetTaskHandler запущен")
	var task *db.Task
	var err error
	// Получаем значение id
	id := r.URL.Query().Get("id")

	// Если поле id не пустое
	if id == "" {
		Logger.Printf("получено пустое поле id")
		sendError(w, http.StatusBadRequest, fmt.Errorf("поле id пустое"))
		return
	}

	// Проверка формата id
	if _, err = strconv.Atoi(id); err != nil {
		Logger.Printf("id не соответствует формату числа")
		sendError(w, http.StatusBadRequest, fmt.Errorf("id не соответствует формату числа"))
		return
	}

	// Поиск задачи
	task, err = db.GetTask(id)
	if err != nil {
		Logger.Printf("ошибка в GetTask: %v", err)
		sendError(w, http.StatusInternalServerError, fmt.Errorf("ошибка в GetTask: %w", err))
		return
	}

	// Проверка существования задачи
	if task == nil {
		sendError(w, http.StatusNotFound, errors.New("задача не найдена"))
		return
	}
	writeJson(w, task)
}
