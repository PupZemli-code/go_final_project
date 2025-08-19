package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
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

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) TaskStore {
	return TaskStore{db: db}
}

// AddTask записывает в базу данных структуру Task,
// возвращает id записанных данных, и ошибку
func (db TaskStore) AddTask(task *Task) (int64, error) {

	// Подготовленный запрос
	query := `
    INSERT INTO scheduler(date, title, comment, repeat) 
    VALUES(?, ?, ?, ?)
    `

	// Запрос с получением результатат
	res, err := db.db.Exec(query,
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
func (db TaskStore) Tasks(limit int) ([]*Task, error) {
	taskSlice := []*Task{}

	// Подготовленный запрос
	query := `
	SELECT * FROM scheduler 
	ORDER BY date LIMIT ?
	`
	// Получает *sql.Rows (строки из db в количестве limit)
	rows, err := db.db.Query(query, limit)
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
func (db TaskStore) SearchTitleComment(limit int, search string) ([]*Task, error) {
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
	rows, err := db.db.Query(query, search, search, limit)
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
func (db TaskStore) SearchDate(limit int, t time.Time) ([]*Task, error) {
	log.Printf("SearchDat: дата t == %v", t.Format(dateFormat))
	taskSlice := []*Task{}

	// Подготовленный запрос
	query := `
	SELECT * FROM scheduler 
	WHERE date = ?
	LIMIT ?
	`

	date := t.Format(dateFormat)
	log.Printf("SearchDat: дата date == %v", date)
	// Получает *sql.Rows (строки из db в количестве limit)
	rows, err := db.db.Query(query, date, limit)
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
	log.Printf("SearchDat: taskSlice == %v", taskSlice)
	return taskSlice, nil
}

// GetTask выгружает Task из db по ее id
func (db TaskStore) GetTask(id string) (*Task, error) {
	var task Task
	if id == "" {
		return &task, errors.New("ошибка GetTask: получено пустое значение id")
	}
	// Подготовленный запрос
	query := `
	SELECT * FROM scheduler
	WHERE id = ?
	`

	err := db.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)

	if err != nil {
		return &Task{}, err
	}
	return &task, nil
}

// UpdateTask обновляет запись в db
func (db TaskStore) UpdateTask(task *Task) error {
	// Проверка входных данных
	if task == nil {
		return errors.New("task не может быть nil")
	}
	if task.ID == "" {
		return errors.New("ID задачи не может быть пустым")
	}

	// Подготовленный запрос
	query := `
    UPDATE scheduler 
    SET date = ?, 
        title = ?, 
        comment = ?, 
        repeat = ?
    WHERE id = ?
    `
	res, err := db.db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)

	if err != nil {
		return err
	}

	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// Удаляет запись из базы данных
func (db TaskStore) DeleteTask(id string) error {
	// Проверка id
	if id == "" {
		return errors.New("DeleteTask: ID задачи не может быть пустым")
	}

	// Подготовленный запрос
	quary := `
	DELETE FROM scheduler 
	WHERE id = ?
	`

	_, err := db.db.Exec(quary, id)
	if err != nil {
		return err
	}
	return nil
}

// Обновляет запись в db
func (db TaskStore) UpdateDate(next string, id string) error {
	if id == "" {
		return errors.New("DeleteTask: ID задачи не может быть пустым")
	}
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("ошибка в UpdateDate: strconv.Atoi(id): %v", err)
	}

	quary := `
	UPDATE scheduler 
    SET date = ? 
    WHERE id = ?
	`

	_, err := db.db.Exec(quary, next, id)
	if err != nil {
		return fmt.Errorf("ошибка в UpdateDate: Db.Exec(quary, next, id): %v", err)
	}
	return nil
}
