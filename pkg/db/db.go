package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func PathDb() string {
	pathDb := os.Getenv("TODO_DBFILE")
	dbName := "scheduler.db"
	if pathDb == "" {
		pathDb = "pkg/db/"
	}
	return fmt.Sprintf("%s/%s", pathDb, dbName)
}

func Init(dbFile string) error {
	// Команда для создания и настройки таблицы
	schema := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			title VARCHAR NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
		);
		CREATE INDEX date_id ON scheduler(date);
		`
	_, err := os.Stat(dbFile)
	var install bool

	if err != nil {
		install = true
	}
	if install { // if install = true
		db, err = sql.Open("sqlite", dbFile)
		if err != nil {
			return err
		}
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
