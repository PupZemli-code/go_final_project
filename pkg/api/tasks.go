package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Обрабатывает запрос api/tasks
func TasksHendler(w http.ResponseWriter, r *http.Request) {

	var tasks TasksResp
	var err error

	// Получаем параметр search из query string
	searchParam := r.URL.Query().Get("search")

	if searchParam != "" {
		// Если searchParam парсится
		date, ok := time.Parse("02.01.2006", searchParam)
		if ok == nil {

			// Если это дата, ищет по дате
			tasks.Tasks, err = db.SearchDate(15, date)
			if err != nil {
				Logger.Println(fmt.Errorf("ошибка поиска по дате: %w", err))
				sendError(w, http.StatusInternalServerError, fmt.Errorf("ошибка поиска по дате: %w", err))
			}
			// Если строки отсутствуют
			if tasks.Tasks == nil {
				writeJson(w, []*db.Task{})
			}
			writeJson(w, tasks)
		} else {

			// Если не дата, ищет по заголовку или комментарию
			tasks.Tasks, err = db.SearchTitleComment(15, searchParam)
			if err != nil {
				Logger.Println(fmt.Errorf("ошибка поиска по строке: %w", err))
				sendError(w, http.StatusInternalServerError, fmt.Errorf("ошибка поиска по строке: %w", err))
			}
			// отправляет ответ json
			if len(tasks.Tasks) == 0 {
				Logger.Println("db.SearchTitleComment вернул пустую структуру")
				writeJson(w, []*db.Task{})
			} else {
				writeJson(w, tasks)
			}
		}
	} else {

		// Если search пустой, возвращаем все задачи
		tasks.Tasks, err = db.Tasks(50) // в параметре максимальное количество записей
		if err != nil {
			Logger.Println(fmt.Errorf("ошибка поиска по строке: %w", err))
			sendError(w, http.StatusInternalServerError, err)
		}

		// отправляет ответ json
		if tasks.Tasks == nil {
			writeJson(w, []*db.Task{})
		} else {
			writeJson(w, tasks)
		}
	}
}
