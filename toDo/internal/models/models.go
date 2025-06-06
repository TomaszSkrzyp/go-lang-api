package models

// TodoItem defines the structure of a task item.
type TodoItem struct {
	ID     string `json:"id"`
	Task   string `json:"task"`
	Status string `json:"status"`
	Due    string `json:"due"`
}

// PossibleStatus defines the allowed values for a task's status.
var PossibleStatus = []string{"Completed", "In Progress", "Pending", "Canceled"}

// IsValidStatus checks if a given status is valid according to PossibleStatus.
func IsValidStatus(status string) bool {
	for _, s := range PossibleStatus {
		if s == status {
			return true
		}
	}
	return false
}
