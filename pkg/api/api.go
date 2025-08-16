package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/logger"
	"github.com/go-chi/chi"
)

var dateFormat string = "20060102"
var Logger *log.Logger

func staticPath() (http.Handler, error) {
	// Настройка раздачи статических файлов
	staticPath, err := filepath.Abs("./web")
	if err != nil {
		return nil, fmt.Errorf("ошибка получения статического пути: %w", err)
	}
	fs := http.FileServer(http.Dir(staticPath))
	return fs, nil
}

// InitMux инициализирует роутер chi
func InitMux(r *chi.Mux) error {
	Logger, _ = logger.NewLogger()

	fs, err := staticPath()
	if err != nil {
		return err
	}

	r.Handle("/*", fs)
	r.Get("/test", TestHandler)
	r.Get("/api/nextdate", NextDayHandler)

	r.Post("/api/task", auth(AddTaskHandler))
	r.Get("/api/task", auth(GetTaskHandler))
	r.Put("/api/task", auth(SaveTaskHandler))
	r.Delete("/api/task", auth(DeleteTaskHandler))
	r.Post("/api/task/done", auth(TaskDoneHandler))

	r.Get("/api/tasks", auth(TasksHendler))

	r.Post("/api/signin", SigninHandler)
	return nil
}

// TestHandler возвращает Metod, Host, Path
func TestHandler(w http.ResponseWriter, r *http.Request) {

	s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s",
		r.Method, r.Host, r.URL.Path)
	w.Write([]byte(s))
}

// Обработчик "/api/nextdate"
func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now, err := time.Parse(dateFormat, r.FormValue("now"))
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Errorf("ошибка парсинга дат: %w", err))
		return
	}
	naxtDate, err := NextDate(now, date, repeat)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}
	w.Write([]byte(naxtDate))
}

// sendError отправляет ошибку клиенту в формате json
func sendError(w http.ResponseWriter, statusCode int, err error) {

	// Формирует ключ:значение
	errorResponse := map[string]string{"error": fmt.Sprint(err)}

	// устанавливает заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// устанавливает HTTP статус
	w.WriteHeader(statusCode)

	// записываем мапу в json
	data, err := json.Marshal(errorResponse)
	if err != nil {
		Logger.Printf("ошибка при форматировании ответа в json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ошибка при форматировании ответа в json"))
		return
	}

	// отправляет ответ json
	w.Write(data)
}

// writeJson отправляет json ответ
func writeJson(w http.ResponseWriter, data any) {

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Маршалинг данных
	jsonData, err := json.Marshal(data)
	if err != nil {
		Logger.Printf("Ошибка при форматировании ответа в JSON: %v", err)
		http.Error(w, "Ошибка при форматировании ответа в JSON", http.StatusInternalServerError)
		return
	}

	// отправляет ответ json
	if _, err := w.Write(jsonData); err != nil {
		Logger.Printf("Ошибка при отправке JSON ответа: %v", err)
	}
}

func tasksToMaps(tasks []*db.Task) []map[string]string {
	result := make([]map[string]string, len(tasks))

	for i, task := range tasks {
		result[i] = map[string]string{
			"id":      task.ID,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}
	}
	return result
}
