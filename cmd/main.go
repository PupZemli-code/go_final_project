package main

import (
	"io"
	"log"
	"os"

	"github.com/PupZemli-code/go-final-project/go_final_project/internal/server"
)

// Создает лог файл, определяет форматирование
func NewLogger() (*log.Logger, os.File) {
	// Создает файд для логов
	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ошибка создания или открытия лог файла app.log: %v", err)
	}
	// Настройка логов
	logger := log.New(logfile, "LOG ", log.Ldate|log.Ltime)
	io.MultiWriter(logfile, os.Stdout)
	return logger, *logfile
}

func main() {
	// Создание сервера
	logger, logfile := NewLogger()
	defer logfile.Close()

	srv := server.NewServer(logger)

	// Запуск сервера
	logger.Printf("запуск сервера на %s", srv.HTTPServer.Addr)

	err := srv.HTTPServer.ListenAndServe()
	if err != nil {
		logger.Fatalf("ошибка запуска сервера: %v", err)
	}
}
