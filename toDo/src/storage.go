package main

import (
	"sort"
	"strconv"
	"sync"
)

type todo_storage struct {
	items  map[string]*todo_item
	mu     sync.Mutex
	nextID int
}

func (ts *todo_storage) add(task string, status string) string {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	id := strconv.Itoa(ts.nextID)
	ts.items[id] = &todo_item{id, task, status}
	ts.nextID++

	return id
}

func (ts *todo_storage) changeTask(id, task, status string) bool {

	ts.mu.Lock()
	defer ts.mu.Unlock()
	oldTask, ok := ts.items[id]
	if !ok {
		return false
	}
	oldTask.Task = task
	oldTask.Status = status
	return true
}

func (ts *todo_storage) remove(id string) bool {

	ts.mu.Lock()
	defer ts.mu.Unlock()
	_, ok := ts.items[id]
	if ok {
		delete(ts.items, id)
	} else {
		return false
	}

	return true
}

func (ts *todo_storage) getAll() []todo_item {
	tasks := make([]todo_item, 0)
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for _, item := range ts.items {
		tasks = append(tasks, *item)
	}

	// Sort by ID (assuming IDs are numeric strings)
	sort.Slice(tasks, func(i, j int) bool {
		id1, _ := strconv.Atoi(tasks[i].ID)
		id2, _ := strconv.Atoi(tasks[j].ID)
		return id1 < id2
	})

	return tasks
}

func (ts *todo_storage) getOne(id string) (todo_item, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	v, ok := ts.items[id]
	if ok {
		return *v, true
	}
	return todo_item{}, false
}

func (ts *todo_storage) moveStatusUp(id string) (string, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	v, ok := ts.items[id]
	if !ok {
		return "", false
	}
	currentStatus := (*v).Status
	nextStatus := map[string]string{"Canceled": "Pending", "Pending": "In Progress", "In Progress": "Completed"}
	newStatus, newStatusOk := nextStatus[currentStatus]
	if !newStatusOk {
		return "", false
	}
	v.Status = newStatus
	return newStatus, true
}
func (ts *todo_storage) seedSampleData() {
	ts.add("Buy groceries", "Pending")
	ts.add("Clean the house", "Pending")
	ts.add("Finish project report", "Completed")
}
