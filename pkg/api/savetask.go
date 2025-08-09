package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
)

// Принемает Task, возвращает пустой интерфейс any{} == interface{}{}
func SaveTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("SaveTaskHandler запущен")
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
		Logger.Printf("ошибка парсинга даты Date [входные данные: %v]: %v", task.Date, err)
		sendError(w, http.StatusBadRequest, fmt.Errorf("ошибка парсинга даты Date: %w", err))
		return
	}

	// Проверяет наличие данных в Title (поле не может быть пустым)
	if task.Title == "" || len(task.Title) == 0 {
		Logger.Printf("ошибка проверки поля task.Title, поле title не может быть пустым: %v", task)
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
					Logger.Printf("ошибка рассчета даты NextDate: %v", err)
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
				Logger.Printf("ошибка рассчета даты NextDate: %v", err)
				sendError(w, http.StatusBadRequest, fmt.Errorf("ошибка рассчета даты NextDate: %w", err))
				return
			}
			task.Date = nextDate
		}
	}
	err = db.UpdateTask(&task)
	if err != nil {
		Logger.Printf("ошибка обновления task в базе данных: %v", err)
		sendError(w, http.StatusBadRequest, fmt.Errorf("ошибка обновления task в базе данных: %w", err))
		return
	}
	// Отправляет пустой интерфейс
	writeJson(w, map[string]any{})
}
