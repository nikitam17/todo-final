package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Функция добавления задачи в БД
func AddTask(task *Task) (int64, error) {
	var id int64

	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := TodoDB.Exec(query, sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// Функция получения задач из БД
func Tasks(limit int, search string) ([]*Task, error) {
	var s bool
	var query string
	var err error
	var rows *sql.Rows

	tasks := make([]*Task, 0, limit)
	now := time.Now()
	date := now.Format("20060102")
	if search != "" {
		t, err := time.Parse("02.01.2006", search)
		if err == nil {
			date = t.Format("20060102")
			query = `SELECT * FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit`
		} else {
			query = `SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`
			s = true
		}
	} else {
		query = `SELECT * FROM scheduler WHERE date >= :date ORDER BY date LIMIT :limit`
	}
	if s {
		rows, err = TodoDB.Query(query, sql.Named("search", "%"+search+"%"), sql.Named("limit", limit))

	} else {
		rows, err = TodoDB.Query(query, sql.Named("date", date), sql.Named("limit", limit))
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, err
}

// Функция получения задачи из БД
func GetTask(id string) (*Task, error) {
	var task Task

	query := `SELECT * FROM scheduler WHERE id = :id`
	err := TodoDB.QueryRow(query, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Функция удаления задачи из БД
func DeleteTask(id string) error {
	query := `DELETE from scheduler  WHERE id = :id`
	res, err := TodoDB.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}

// Функция обновления задачи в БД
func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := TodoDB.Exec(query, sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
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
