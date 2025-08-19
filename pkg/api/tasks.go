package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
)

var limit int = 50

// Обрабатывает запрос api/tasks
func (t TaskService) TasksHendler(w http.ResponseWriter, r *http.Request) {

	var tasks []*db.Task

	//var tasks TasksResp
	var err error

	// Получаем параметр search из query string
	searchParam := r.URL.Query().Get("search")

	if searchParam != "" {
		// Если searchParam парсится
		date, ok := time.Parse("02.01.2006", searchParam)
		if ok == nil {
			log.Printf("данные парсится %v", searchParam)
		} else {
			log.Printf("данные не парсится %v", searchParam)
		}
		if ok == nil {

			// Если это дата, ищет по дате
			tasks, err = t.store.SearchDate(15, date)
			if err != nil {
				Logger.Println(fmt.Errorf("ошибка поиска по дате: %w", err))
				sendError(w, http.StatusInternalServerError, fmt.Errorf("ошибка поиска по дате: %w", err))
				return
			}
			// Если строки отсутствуют
			if tasks == nil {
				writeJson(w, []*db.Task{})
				return
			}

			writeJson(w, map[string]any{"tasks": tasks})
			return

		} else {

			// Если не дата, ищет по заголовку или комментарию
			tasks, err = t.store.SearchTitleComment(limit, searchParam)
			if err != nil {
				Logger.Println(fmt.Errorf("ошибка поиска по строке: %w", err))
				sendError(w, http.StatusInternalServerError, fmt.Errorf("ошибка поиска по строке: %w", err))
				return
			}
			// Отправляет ответ json
			if len(tasks) == 0 {
				writeJson(w, []*db.Task{})
				return
			} else {
				writeJson(w, map[string]any{"tasks": tasks})
				return
			}
		}
	} else {

		// Если search пустой, возвращаем все задачи
		tasks, err = t.store.Tasks(limit) // в параметре максимальное количество записей
		if err != nil {
			Logger.Println(fmt.Errorf("ошибка поиска по строке: %w", err))
			sendError(w, http.StatusInternalServerError, err)
			return
		}

		// отправляет ответ json
		if tasks == nil {
			writeJson(w, []*db.Task{})
			return
		} else {
			writeJson(w, map[string]any{"tasks": tasks})
			return
		}
	}
}
