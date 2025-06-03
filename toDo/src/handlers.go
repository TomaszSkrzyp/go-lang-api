package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (ts *todo_storage) handleGet(w http.ResponseWriter, r *http.Request) {
	itemId := mux.Vars(r)["id"]
	item, found := ts.getOne(itemId)
	if found {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "The task with this ID has not been found",
		})
	}
}

func (ts *todo_storage) handleAdd(w http.ResponseWriter, r *http.Request) {
	var newItem map[string]string
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, taskOk := newItem["task"]
	status, statusOk := newItem["status"]

	if !taskOk {
		http.Error(w, "Missing a required field: task", http.StatusBadRequest)
		return
	}
	if !statusOk {
		status = "Pending"
	}
	if !isValidStatus(status) {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	itemId := ts.add(task, status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := map[string]string{
		"message": "Task updated successfully",
		"id":      itemId,
	}
	if !statusOk {
		resp["note"] = "Status not provided, defaulted to 'Pending'"
	}

	json.NewEncoder(w).Encode(resp)
}

func (ts *todo_storage) handleGetAll(w http.ResponseWriter, r *http.Request) {
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
	allTasks := ts.getAll()

	var filtered []todo_item

	isValidStatus := false
	for _, status := range possibleStatus {
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

func (ts *todo_storage) handleRemove(w http.ResponseWriter, r *http.Request) {
	itemId := mux.Vars(r)["id"]
	if itemId == "" {
		http.Error(w, "Missing id in URL path", http.StatusBadRequest)
		return
	}
	ok := ts.remove(itemId)
	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Task removed",
			"id":      itemId,
		})
	} else {
		http.Error(w, "Failed to remove task", http.StatusInternalServerError)
	}
}

func (ts *todo_storage) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	itemId := mux.Vars(r)["id"]
	if itemId == "" {
		http.Error(w, "Missing id in URL path", http.StatusBadRequest)
		return
	}
	var newItem map[string]string
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	changeStatusUpStr, changeOk := newItem["changeUp"]
	task, taskOk := newItem["task"]
	status, statusOk := newItem["status"]

	if changeOk {
		changeStatusUp, err := strconv.ParseBool(changeStatusUpStr)
		if err != nil {
			http.Error(w, "Invalid value for 'changeUp'; expected 'true' or 'false'", http.StatusBadRequest)
			return
		}
		if changeStatusUp {
			newStatus, success := ts.moveStatusUp(itemId)
			w.Header().Set("Content-Type", "application/json")
			if success {
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
		http.Error(w, "Missing 'task' field", http.StatusBadRequest)
		return
	}
	if !statusOk {
		status = "Pending"
	}
	if !isValidStatus(status) {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	if ts.changeTask(itemId, task, status) {
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
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
	}
}
