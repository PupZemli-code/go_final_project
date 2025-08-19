package main

import (
	"log"
	"os"

	"github.com/PupZemli-code/go-final-project/go_final_project/internal/server"
	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/api"
	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/db"
	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/logger"
)

var Logger *log.Logger

func main() {
	// Инициализация логгера
	var logfile os.File
	Logger, _ = logger.NewLogger()
	defer logfile.Close()

	// Инициализация базы данных
	database, err := db.InitDb(db.PathDb())
	if err != nil {
		Logger.Fatalf("ошибка инициализации db: %v", err)
	}
	defer db.CloseDb(database)

	tstore := db.NewTaskStore(database)
	tservice := api.NewTaskService(tstore)

	// Создание сервера
	srv := server.NewServer(Logger, tservice)

	// Запуск сервера
	Logger.Printf("------------------------------------------")
	Logger.Printf("Запуск сервера на %s", srv.HTTPServer.Addr)
	err = srv.HTTPServer.ListenAndServe()
	if err != nil {
		Logger.Fatalf("ошибка запуска сервера: %v", err)
	}
}
