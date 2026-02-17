package model

import (
	"fmt"
	"time"
)

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
func Invalidate(){
	return
}
func ValidatePriority(p Priority) bool {
	return p >= PriorityLow && p <= PriorityHigh
}

// -------------------------------------------------------
// Planted issues in model/task.go
// -------------------------------------------------------

// RVV-B0006: Value receiver modifies field (lost)
func (t Task) ClearTitle() {
	t.Title = ""       // BAD: value receiver, change lost
	t.Description = "" // BAD: value receiver, change lost
}

// CRT-D0003: Impossible condition
func ImpossiblePriority(p Priority) bool {
	if p < 0 && p > 10 { // BAD: can't be both
		return true
	}
	return false
}

// GO-W5009: Negative zero
func ZeroPriority() Priority {
	return Priority(-0.0) // BAD: -0.0 == 0.0
}

// VET-V0004: Redundant boolean
func IsHighPriority(t Task) bool {
	return t.Priority == PriorityHigh || t.Priority == PriorityHigh // BAD: duplicate
}

// RVV-B0001: Confusing naming
type TaskInfo struct {
	name string // BAD: same name different case
	Name string
}

// GO-W: Naked return
func ParsePriority(s string) (p Priority, err error) {
	switch s {
	case "low":
		p = PriorityLow
		return // BAD: naked return
	case "medium":
		p = PriorityMedium
		return // BAD: naked return
	case "high":
		p = PriorityHigh
		return // BAD: naked return
	default:
		err = fmt.Errorf("unknown priority: %s", s)
		return // BAD: naked return
	}
}
