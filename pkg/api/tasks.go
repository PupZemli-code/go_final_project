package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
)

// type TasksResp struct {
// 	Tasks []*db.Task `json:"tasks"`
// }

var Tasks []*db.Task

// Обрабатывает запрос api/tasks
func TasksHendler(w http.ResponseWriter, r *http.Request) {

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
			Tasks, err = db.SearchDate(15, date)
			if err != nil {
				Logger.Println(fmt.Errorf("ошибка поиска по дате: %w", err))
				sendError(w, http.StatusInternalServerError, fmt.Errorf("ошибка поиска по дате: %w", err))
				return
			}
			// Если строки отсутствуют
			if Tasks == nil {
				writeJson(w, []*db.Task{})
				return
			}
			//writeJson(w, tasks)
			//Tasks = tasks.Tasks
			log.Printf("длина Tasks == %v", len(Tasks))
			log.Printf("Tasks == %v", Tasks)
			//resp := tasksToMaps(Tasks)
			writeJson(w, map[string]any{"tasks": Tasks})
			return

		} else {

			// Если не дата, ищет по заголовку или комментарию
			Tasks, err = db.SearchTitleComment(15, searchParam)
			if err != nil {
				Logger.Println(fmt.Errorf("ошибка поиска по строке: %w", err))
				sendError(w, http.StatusInternalServerError, fmt.Errorf("ошибка поиска по строке: %w", err))
				return
			}
			// Отправляет ответ json
			if len(Tasks) == 0 {
				writeJson(w, []*db.Task{})
				return
			} else {
				writeJson(w, map[string]any{"tasks": Tasks})
				return
			}
		}
	} else {

		// Если search пустой, возвращаем все задачи
		Tasks, err = db.Tasks(50) // в параметре максимальное количество записей
		if err != nil {
			Logger.Println(fmt.Errorf("ошибка поиска по строке: %w", err))
			sendError(w, http.StatusInternalServerError, err)
			return
		}

		// отправляет ответ json
		if Tasks == nil {
			writeJson(w, []*db.Task{})
			return
		} else {
			writeJson(w, map[string]any{"tasks": Tasks})
			return
		}
	}
}
