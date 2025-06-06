package dbControl

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TomaszSkrzyp/go-lang-api/toDo/internal/models"
	"github.com/gorilla/mux"
)

// HandleGet handles the HTTP GET request to fetch a single task by its ID.
// If the task exists, it responds with status 200 and the task in JSON format.
// If the task is not found, it responds with 404 and an error message.
func (ts *TodoStorage) HandleGet(w http.ResponseWriter, r *http.Request) {
	itemId := mux.Vars(r)["id"]
	item, err := ts.getOne(itemId)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "The task with this ID has not been found",
		})
	}
}

// HandleAdd handles the HTTP POST request to create a new task.
// It expects a JSON body with "task", "status" (optional), and "due" fields.
// Validates input, assigns default status if missing, and returns status 201 on success.
// On validation or internal error, returns the appropriate error response.
func (ts *TodoStorage) HandleAdd(w http.ResponseWriter, r *http.Request) {
	var newItem map[string]string
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	task, taskOk := newItem["task"]
	status, statusOk := newItem["status"]
	due, dueOk := newItem["due"]

	if !taskOk {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing a required field: task"})
		return
	}
	if !statusOk {
		status = "Pending"
	}
	if !models.IsValidStatus(status) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid status value"})
		return
	}
	if !dueOk || due == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing or empty required field: due"})
		return
	}

	itemId, err := ts.add(task, status, due)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add task: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := map[string]string{
		"message": "Task added successfully",
		"id":      itemId,
	}
	if !statusOk {
		resp["note"] = "Status not provided, defaulted to 'Pending'"
	}
	json.NewEncoder(w).Encode(resp)
}

// HandleGetAll handles the HTTP GET request to retrieve all tasks.
// Supports optional query parameters: "status" for filtering, and "page" & "limit" for pagination.
// Responds with paginated and optionally filtered list of tasks.
func (ts *TodoStorage) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	allTasks, err := ts.getAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch tasks" + err.Error()})
		return
	}

	var filtered []models.TodoItem

	isValidStatus := false
	for _, status := range models.PossibleStatus {
		if status == statusFilter {
			isValidStatus = true
			break
		}
	}

	if isValidStatus {
		for _, item := range allTasks {
			if statusFilter == item.Status {
				filtered = append(filtered, item)
			}
		}
	} else {
		filtered = allTasks
	}

	start := (page - 1) * limit
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	paged := filtered[start:end]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": len(filtered),
		"tasks": paged,
	})
}

// HandleRemove handles the HTTP DELETE request to remove a task by its ID.
// If the task is successfully removed, it returns status 200 with a success message.
// If the ID is missing or the task cannot be removed, it returns an appropriate error.
func (ts *TodoStorage) HandleRemove(w http.ResponseWriter, r *http.Request) {
	itemId := mux.Vars(r)["id"]
	if itemId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing id in URL path"})
		return
	}
	err := ts.remove(itemId)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Task removed",
			"id":      itemId,
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove task" + err.Error()})
		return
	}
}

// HandleUpdateTask handles the HTTP PUT request to update an existing task by its ID.
// Supports a special "changeUp" flag to move the task's status up instead of regular update.
// Validates fields, applies changes, and returns success or error responses accordingly.
func (ts *TodoStorage) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	itemId := mux.Vars(r)["id"]
	if itemId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing id in URL path"})
		return
	}
	var newItem map[string]string
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return

	}

	changeStatusUpStr, changeOk := newItem["changeUp"]
	task, taskOk := newItem["task"]
	status, statusOk := newItem["status"]
	due, dueOk := newItem["due"]

	if changeOk {
		changeStatusUp, err := strconv.ParseBool(changeStatusUpStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid value for 'changeUp'; expected 'true' or 'false'",
			})
			return
		}
		if changeStatusUp {
			newStatus, err := ts.moveStatusUp(itemId)
			w.Header().Set("Content-Type", "application/json")
			if err == nil {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{
					"message":   "Status moved up",
					"id":        itemId,
					"newStatus": newStatus,
				})
			} else {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Status of this task can't be moved up",
					"id":      itemId,
				})
			}
			return
		}
	}

	if !taskOk {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing 'task' field",
		})
		return
	}
	if !statusOk {
		status = "Pending"
	}
	if !dueOk {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing 'due' field",
		})
		return
	}
	if !models.IsValidStatus(status) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid status value",
		})
		return
	}

	err := ts.changeTask(itemId, task, status, due)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]string{
			"message": "Task updated successfully",
			"id":      itemId,
		}
		if !statusOk {
			resp["note"] = "Status not provided, defaulted to 'Pending'"
		}
		json.NewEncoder(w).Encode(resp)
	} else {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to update task: " + err.Error(),
		})
		return
	}
}
