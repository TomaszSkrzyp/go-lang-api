package validate

import (
	"errors"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"

	"github.com/TomaszSkrzyp/go-lang-api/toDo/internal/models"
)

var sanitizer = bluemonday.UGCPolicy()

func sanitizeString(input string) string {
	return sanitizer.Sanitize(strings.TrimSpace(input))
}

func SanitizeAndValidateTaskInput(data map[string]string) (string, string, string, error) {
	task, ok := data["task"]
	if !ok || strings.TrimSpace(task) == "" {
		return "", "", "", errors.New("missing or empty 'task' field")
	}
	task = sanitizeString(task)

	status := data["status"]
	if status == "" {
		status = "Pending"
	} else {
		status = sanitizeString(status)
		if !models.IsValidStatus(status) {
			return "", "", "", errors.New("invalid status value")
		}
	}

	due, ok := data["due"]
	if !ok || strings.TrimSpace(due) == "" {
		return "", "", "", errors.New("missing or empty 'due' field")
	}
	due = sanitizeString(due)

	// Validate date format (optional)
	_, err := time.Parse("2006-01-02", due)
	if err != nil {
		return "", "", "", errors.New("invalid due date format, expected YYYY-MM-DD")
	}

	return task, status, due, nil
}
