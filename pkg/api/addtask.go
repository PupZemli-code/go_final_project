package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/logger"
)

var Db *db.DB

// AddTaskHandler обрабатывает запрос "/api/task" POST
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	logger, _ := logger.NewLogger()
	// Проверяет метод запроса
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, fmt.Errorf("разрешен только метод POST"))
		return
	}

	// Структура для хранения данных
	var task db.Task

	// Читает тело запроса
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	// Проверяет корректноять парсинга
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		logger.Printf("ошибка парсинга даты Date [входные данные: %v]: %v", task.Date, err)
		sendError(w, http.StatusBadRequest, fmt.Errorf("ошибка парсинга даты Date: %w", err))
		return
	}

	// Проверяет наличие данных в Title (поле не может быть пустым)
	if task.Title == "" || len(task.Title) == 0 {
		logger.Printf("ошибка проверки поля task.Title, поле title не может быть пустым: %v", task)
		sendError(w, http.StatusBadRequest, fmt.Errorf("поле title не может быть пустым"))
		return
	}

	// Если поле Date содержит данные
	if task.Date != "" {
		// Если date < текущего дня
		if !afterNow(t, time.Now()) {
			// Если параметры повтарения отсутствуют
			if task.Repeat == "" || len(task.Repeat) == 0 {
				task.Date = time.Now().Format(dateFormat)
			}
			// Если параметры повтарения присутствуют
			if task.Repeat != "" {
				nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					logger.Printf("ошибка рассчета даты NextDate: %v", err)
					sendError(w, http.StatusBadRequest, fmt.Errorf("ошибка рассчета даты NextDate: %w", err))
					return
				}
				task.Date = nextDate
			}
		}
	}
	// ели Date пусто, записывает текущую дату
	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
		if task.Repeat == "" || len(task.Repeat) == 0 {
			task.Date = time.Now().Format(dateFormat)
		}
		// Если параметры повтарения присутствуют
		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				logger.Printf("ошибка рассчета даты NextDate: %v", err)
				sendError(w, http.StatusBadRequest, fmt.Errorf("ошибка рассчета даты NextDate: %w", err))
				return
			}
			task.Date = nextDate
		}
	}
	Db, err := db.InitDb(db.PathDb())
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDb()

	// id, err := dbInstance.AddTask(&task)
	id, err := Db.AddTask(&task)
	if err != nil {
		logger.Printf("ошибка добавления task в базу данных: %v", err)
		sendError(w, http.StatusBadRequest, fmt.Errorf("ошибка добавления task в базу данных: %w", err))
		return
	}
	task.ID = strconv.Itoa(int(id))

	//
	data, err := json.Marshal(map[string]string{"id": task.ID})
	if err != nil {
		logger.Printf("ошибка при форматировании ответа в json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ошибка при форматировании ответа в json"))
		return
	}
	w.Write(data)
}
