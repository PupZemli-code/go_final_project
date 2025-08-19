package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PupZemli-code/go-final-project/go_final_project/pkg/logger"
	_ "modernc.org/sqlite"
)

var Logger *log.Logger

// Возвращает путь до базы данных, берет его из переменной
// окружения "TODO_DBFILE", или использует путь по умолчанию
func PathDb() string {
	// Путь по переменной окружения
	pathDb := os.Getenv("TODO_DBFILE")
	// Имя db
	dbName := "scheduler.db"
	// Если переменная окружения пуста (не используется)
	if pathDb == "" {
		// путь по умолчанию
		pathDb = "pkg/db/"
		// Создает деррикторию если она не сузествует
		if err := os.MkdirAll(pathDb, 0755); err != nil {
			log.Fatalf("ошибка создания директория: %v", err)
		}
	}
	return fmt.Sprintf("%s/%s", pathDb, dbName)
}

// Создает подключение к базе данных, если базы
// нет, создает ее по указанной схеме
// Инициализация базы данных
func InitDb(pathDb string) (*sql.DB, error) {
	Logger, _ = logger.NewLogger()
	schema := `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date CHAR(8) NOT NULL DEFAULT '',
        title TEXT NOT NULL DEFAULT '',
        comment TEXT NOT NULL DEFAULT '',
        repeat VARCHAR(128) NOT NULL DEFAULT ''
    );
    CREATE INDEX IF NOT EXISTS date_id ON scheduler(date);
    `

	// Открываем соединение
	db, err := sql.Open("sqlite", pathDb)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Настройка пула подключений
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 10)

	// Создаем таблицы
	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	return db, nil
}

// Закрытие подключения
func CloseDb(db *sql.DB) {
	db.Close()
}
