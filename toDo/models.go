package main

type todo_item struct {
	ID     string `json:"id"`
	Task   string `json:"task"`
	Status string `json:"status"`
}

var possibleStatus = []string{"Completed", "In Progress", "Pending", "Canceled"}

func isValidStatus(status string) bool {
	for _, s := range possibleStatus {
		if s == status {
			return true
		}
	}
	return false
}
