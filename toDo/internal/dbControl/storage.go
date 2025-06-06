package dbControl

import (
	"database/sql"
	"errors"

	"github.com/TomaszSkrzyp/go-lang-api/toDo/internal/models"
	_ "github.com/lib/pq"
)

// TodoStorage represents a wrapper around the database connection
// used to operate on the `todos` table.
type TodoStorage struct {
	DB *sql.DB
}

// add inserts a new task into the `todos` table and returns the generated ID.
func (r *TodoStorage) add(task, status, due string) (string, error) {
	var id string
	err := r.DB.QueryRow(
		"INSERT INTO todos (task, status, due) VALUES ($1, $2, $3) RETURNING id",
		task, status, due,
	).Scan(&id)
	return id, err
}

// changeTask updates the task, status, and due date for a given task ID.
func (r *TodoStorage) changeTask(id, task, status, due string) error {
	_, err := r.DB.Exec(
		"UPDATE todos SET task=$1, status=$2, due=$3 WHERE id=$4",
		task, status, due, id,
	)
	return err
}

// remove deletes a task from the `todos` table by its ID.
func (r *TodoStorage) remove(id string) error {
	_, err := r.DB.Exec("DELETE FROM todos WHERE id=$1", id)
	return err
}

// getAll retrieves all tasks from the database ordered by ID.
func (r *TodoStorage) getAll() ([]models.TodoItem, error) {
	rows, err := r.DB.Query("SELECT id, task, status, due FROM todos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.TodoItem
	for rows.Next() {
		var t models.TodoItem
		if err := rows.Scan(&t.ID, &t.Task, &t.Status, &t.Due); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

// getOne fetches a single task by its ID.
func (r *TodoStorage) getOne(id string) (models.TodoItem, error) {
	var todo models.TodoItem
	err := r.DB.QueryRow(
		"SELECT id, task, status, due FROM todos WHERE id=$1",
		id,
	).Scan(&todo.ID, &todo.Task, &todo.Status, &todo.Due)
	return todo, err
}

// moveStatusUp transitions a task to the next status in the predefined flow.
// Returns an error if the task can't be transitioned or does not exist.
func (r *TodoStorage) moveStatusUp(id string) (string, error) {
	var current string
	err := r.DB.QueryRow("SELECT status FROM todos WHERE id = $1", id).Scan(&current)
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

	res, err := r.DB.Exec("UPDATE todos SET status = $1 WHERE id = $2", newStatus, id)
	if err != nil {
		return "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", errors.New("no task found with given id")
	}

	return newStatus, nil
}

// seedSampleData populates the database with sample tasks for testing/demo.
func (r *TodoStorage) seedSampleData() {
	r.add("Buy groceries", "Pending", "2025-06-10")
	r.add("Clean the house", "Pending", "2025-06-11")
	r.add("Finish project report", "Completed", "2025-06-09")
}
