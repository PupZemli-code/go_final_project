/*
package server описывает структуру сервера и возвращает экземпляр Server

	type Server struct {
		Logger     *log.Logger
		HTTPServer *http.Server
	}
*/
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	handler "github.com/PupZemli-code/go-final-project/go_final_project/internal/handlers"
	"github.com/go-chi/chi"
)

// Структура сервера
type Server struct {
	Logger     *log.Logger
	HTTPServer *http.Server
}

func staticPath() (http.Handler, error) {
	// Настройка раздачи статических файлов
	staticPath, err := filepath.Abs("./web")
	if err != nil {
		return nil, err
	}
	fs := http.FileServer(http.Dir(staticPath))
	return fs, nil
}

// GetAddr принемает значение переменных окружения "TODO_HOST" и "TODO_PORT",
// возвращает Addr в формате "HOST:PORT", по умолчанию вернет "localhost:7540"
func GetAddr() string {
	host := os.Getenv("TODO_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return fmt.Sprintf("%s:%s", host, port)
}

// Создает сервер
func NewServer(logger *log.Logger) *Server {
	// Инициализация роутера
	r := chi.NewRouter()

	fs, err := staticPath()
	if err != nil {
		logger.Fatal(err)
	}

	r.Handle("/*", fs)
	r.Get("/test", handler.TestHandler)

	// Описание сервера
	server := &http.Server{
		Addr:         GetAddr(),
		Handler:      r,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	// Возвращает указатель server *http.Server и *log.Logger
	return &Server{
		Logger:     logger,
		HTTPServer: server,
	}
}
