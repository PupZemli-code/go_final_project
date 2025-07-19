package logger

import (
	"io"
	"log"
	"os"
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
	logger.SetOutput(io.MultiWriter(logfile, os.Stdout))
	return logger, *logfile
}
