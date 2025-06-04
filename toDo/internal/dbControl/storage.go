package main

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type todo_storage struct {
	db *sql.DB
}

func (r *todo_storage) add(task, status, due string) (string, error) {
	var id string
	err := r.db.QueryRow(
		"INSERT INTO todos (task, status, due) VALUES ($1, $2, $3) RETURNING id",
		task, status, due,
	).Scan(&id)
	return id, err
}

func (r *todo_storage) changeTask(id, task, status, due string) error {
	_, err := r.db.Exec(
		"UPDATE todos SET task=$1, status=$2, due=$3 WHERE id=$4",
		task, status, due, id,
	)
	return err
}

func (r *todo_storage) remove(id string) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id=$1", id)
	return err
}

func (r *todo_storage) getAll() ([]todo_item, error) {
	rows, err := r.db.Query("SELECT id, task, status, due FROM todos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []todo_item
	for rows.Next() {
		var t todo_item
		if err := rows.Scan(&t.ID, &t.Task, &t.Status, &t.Due); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (r *todo_storage) getOne(id string) (todo_item, error) {
	var todo todo_item
	err := r.db.QueryRow(
		"SELECT id, task, status, due FROM todos WHERE id=$1",
		id,
	).Scan(&todo.ID, &todo.Task, &todo.Status, &todo.Due)
	return todo, err
}

func (r *todo_storage) moveStatusUp(id string) (string, error) {
	var current string
	err := r.db.QueryRow("SELECT status FROM todos WHERE id = $1", id).Scan(&current)
	if err != nil {
		return "", err
	}

	next := map[string]string{
		"Canceled":    "Pending",
		"Pending":     "In Progress",
		"In Progress": "Completed",
	}

	newStatus, ok := next[current]
	if !ok {
		return "", errors.New("cannot move status forward")
	}

	_, err = r.db.Exec("UPDATE todos SET status = $1 WHERE id = $2", newStatus, id)
	return newStatus, err
}

func (r *todo_storage) seedSampleData() {
	r.add("Buy groceries", "Pending", "2025-06-10")
	r.add("Clean the house", "Pending", "2025-06-11")
	r.add("Finish project report", "Completed", "2025-06-09")
}
