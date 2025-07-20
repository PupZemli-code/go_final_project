package api

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
)

var dateFormat string = "20060102"

func staticPath() (http.Handler, error) {
	// Настройка раздачи статических файлов
	staticPath, err := filepath.Abs("./web")
	if err != nil {
		return nil, fmt.Errorf("ошибка получения статического пути: %w", err)
	}
	fs := http.FileServer(http.Dir(staticPath))
	return fs, nil
}

func Init(r *chi.Mux) error {

	fs, err := staticPath()
	if err != nil {
		return err
	}

	r.Handle("/*", fs)
	r.Get("/test", TestHandler)
	r.Get("/api/nextdate", NextDayHandler)
	return nil
}

// TestHandler возвращает Metod, Host, Path
func TestHandler(w http.ResponseWriter, r *http.Request) {

	s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s",
		r.Method, r.Host, r.URL.Path)
	w.Write([]byte(s))
}

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now, err := time.Parse(dateFormat, r.FormValue("now"))
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Errorf("ошибка парсинга дат: %w", err))
	}
	naxtDate, err := NextDate(now, date, repeat)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
	}
	w.Write([]byte(naxtDate))
}

func sendError(w http.ResponseWriter, statusCode int, err error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s", err.Error())
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
