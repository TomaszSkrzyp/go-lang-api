package models

type TodoItem struct {
	ID     string `json:"id"`
	Task   string `json:"task"`
	Status string `json:"status"`
	Due    string `json:"due"`
}

var PossibleStatus = []string{"Completed", "In Progress", "Pending", "Canceled"}

func IsValidStatus(status string) bool {
	for _, s := range PossibleStatus {
		if s == status {
			return true
		}
	}
	return false
}
