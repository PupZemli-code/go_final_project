package db

import (
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var dateFormat string = "20060102"

// AddTask записывает в базу данных структуру Task,
// возвращает id записанных данных, и ошибку
func AddTask(task *Task) (int64, error) {

	// Подготовленный запрос
	query := `
    INSERT INTO scheduler(date, title, comment, repeat) 
    VALUES(?, ?, ?, ?)
    `
	// Запрос с получением результатат
	res, err := Db.Exec(query,
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

// Tasks возвращает limit кол-во записей из db
func Tasks(limit int) ([]*Task, error) {
	taskSlice := []*Task{}

	// Подготовленный запрос
	query := `
	SELECT * FROM scheduler 
	ORDER BY date LIMIT ?
	`
	// Получает *sql.Rows (строки из db в количестве limit)
	rows, err := Db.Query(query, limit)
	if err != nil {
		Logger.Printf("ошибка работы db.Query: %v", err)
		return []*Task{}, fmt.Errorf("ошибка чтения базы данных: %w", err)
	}
	defer rows.Close()

	// Читаем *sql.Rows и записываем данные в taskSlice := []*Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			Logger.Printf("ошибка работы rows.Scan: %v", err)
			return []*Task{}, fmt.Errorf("ошибка чтения базы данных: %w", err)
		}
		taskSlice = append(taskSlice, &task)
	}
	// Проверка наличия ошибок при получении данных
	if err := rows.Err(); err != nil {
		Logger.Printf("ошибка rows.Err(): %v", err)
		return []*Task{}, fmt.Errorf("ошибка чтения базы данных: %w", err)
	}
	return taskSlice, nil
}

// SearchTitleComment возвращает limit значений по
// параметрам поиска подстроки title, comment.
func SearchTitleComment(limit int, search string) ([]*Task, error) {
	taskSlice := []*Task{}

	// Подготовленный запрос
	query := `
	SELECT * FROM scheduler 
	WHERE title LIKE ?
	OR comment LIKE ?
	ORDER BY date 
	LIMIT ? 
	`
	search = "%" + search + "%"
	rows, err := Db.Query(query, search, search, limit)
	if err != nil {
		Logger.Printf("ошибка работы db.Query: %v", err)
		return []*Task{}, fmt.Errorf("ошибка чтения базы данных: %w", err)
	}
	defer rows.Close()

	// Читаем *sql.Rows и записываем данные в taskSlice := []*Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			Logger.Printf("ошибка работы rows.Scan: %v", err)
			return []*Task{}, fmt.Errorf("ошибка чтения базы данных: %w", err)
		}
		taskSlice = append(taskSlice, &task) // записывает
	}
	// Проверка наличия ошибок при получении данных
	if err := rows.Err(); err != nil {
		Logger.Printf("ошибка rows.Err(): %v", err)
		return []*Task{}, fmt.Errorf("ошибка чтения базы данных: %w", err)
	}
	return taskSlice, nil
}

// SearchDate возвращает limit значений по
// параметрам поиска даты.
func SearchDate(limit int, t time.Time) ([]*Task, error) {
	taskSlice := []*Task{}

	// Подготовленный запрос
	query := `
	SELECT * FROM scheduler 
	WHERE date = ?
	LIMIT ?
	`

	date := t.Format(dateFormat)

	// Получает *sql.Rows (строки из db в количестве limit)
	rows, err := Db.Query(query, date, limit)
	if err != nil {
		Logger.Printf("ошибка работы db.Query: %v", err)
		return []*Task{}, fmt.Errorf("ошибка чтения базы данных: %w", err)
	}
	defer rows.Close()
	// Читаем *sql.Rows и записываем данные в taskSlice := []*Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			Logger.Printf("ошибка работы rows.Scan: %v", err)
			return []*Task{}, fmt.Errorf("ошибка чтения базы данных: %w", err)
		}
		taskSlice = append(taskSlice, &task)
	}
	// Проверка наличия ошибок при получении данных
	if err := rows.Err(); err != nil {
		Logger.Printf("ошибка rows.Err(): %v", err)
		return []*Task{}, fmt.Errorf("ошибка чтения базы данных: %w", err)
	}
	return taskSlice, nil
}
