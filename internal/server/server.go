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
	"time"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/api"
	"github.com/go-chi/chi"
)

// Структура сервера
type Server struct {
	Logger     *log.Logger
	HTTPServer *http.Server
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
	// Инициализация роутера и хендлеров
	r := chi.NewRouter()
	err := api.Init(r)
	if err != nil {
		logger.Fatal(err)
	}

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
