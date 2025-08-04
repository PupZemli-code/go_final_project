package db

import "fmt"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask записывает в базу данных структуру Task,
// возвращает id записанных данных, и ошибку

func (db *DB) AddTask(task *Task) (int64, error) {

	// Подготовленный запрос
	query := `
    INSERT INTO scheduler(date, title, comment, repeat) 
    VALUES(?, ?, ?, ?)
    `
	// Запрос с получением результатат
	res, err := db.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	// Извлекает id из результата
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %w", err)
	}

	return id, nil
}
