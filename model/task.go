package model

import "time"

type Priority int

const (
	PriorityLow    Priority = 0
	PriorityMedium Priority = 1
	PriorityHigh   Priority = 2
)

type Task struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Done        bool     `json:"done"`
	Priority    Priority `json:"priority"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func NewTask(title, description string, priority Priority) Task {
	now := time.Now().UTC().Format(time.RFC3339)
	return Task{
		Title:       title,
		Description: description,
		Priority:    priority,
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) MarkDone() {
	t.Done = true
	t.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}

func (t *Task) SetPriority(p Priority) {
	t.Priority = p
	t.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}

func ValidatePriority(p Priority) bool {
	return p >= PriorityLow && p <= PriorityHigh
}
