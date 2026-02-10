package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/sawez-deepsource/demo-go/model"
	"github.com/sawez-deepsource/demo-go/store"
)

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type statsResponse struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Pending   int `json:"pending"`
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
	doneFilter := r.URL.Query().Get("done")
	priorityFilter := r.URL.Query().Get("priority")

	var tasks []model.Task

	if doneFilter != "" {
		done := doneFilter == "true"
		tasks = store.FilterByDone(done)
	} else if priorityFilter != "" {
		p, err := strconv.Atoi(priorityFilter)
		if err != nil || !model.ValidatePriority(model.Priority(p)) {
			writeError(w, http.StatusBadRequest, "invalid priority filter")
			return
		}
		tasks = store.FilterByPriority(model.Priority(p))
	} else {
		tasks = store.All()
	}

	writeJSON(w, http.StatusOK, tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, ok := store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	if t.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if !model.ValidatePriority(t.Priority) {
		writeError(w, http.StatusBadRequest, "priority must be 0 (low), 1 (medium), or 2 (high)")
		return
	}
	created := store.Add(t)
	log.Printf("task created: id=%s title=%q", created.ID, created.Title)
	writeJSON(w, http.StatusCreated, created)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	if t.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if !model.ValidatePriority(t.Priority) {
		writeError(w, http.StatusBadRequest, "priority must be 0 (low), 1 (medium), or 2 (high)")
		return
	}
	updated, ok := store.Update(id, t)
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	log.Printf("task updated: id=%s title=%q", updated.ID, updated.Title)
	writeJSON(w, http.StatusOK, updated)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !store.Delete(id) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	log.Printf("task deleted: id=%s", id)
	w.WriteHeader(http.StatusNoContent)
}

func TaskStats(w http.ResponseWriter, r *http.Request) {
	all := store.All()
	completed := store.FilterByDone(true)
	stats := statsResponse{
		Total:     len(all),
		Completed: len(completed),
		Pending:   len(all) - len(completed),
	}
	writeJSON(w, http.StatusOK, stats)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
