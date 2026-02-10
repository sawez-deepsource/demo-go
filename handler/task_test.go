package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sawez-deepsource/demo-go/handler"
	"github.com/sawez-deepsource/demo-go/model"
	"github.com/sawez-deepsource/demo-go/store"
)

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("GET /tasks", handler.ListTasks)
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
	mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)
	mux.HandleFunc("GET /stats", handler.TaskStats)
	return mux
}

func TestCreateAndListTasks(t *testing.T) {
	store.Clear()
	mux := setupMux()

	body := `{"title":"Test Task","description":"A test","priority":1}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)
	if created.ID == "" {
		t.Fatal("expected task to have an ID")
	}
	if created.CreatedAt == "" {
		t.Fatal("expected task to have a created_at timestamp")
	}

	req = httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Test Task") {
		t.Fatal("expected response to contain 'Test Task'")
	}
}

func TestCreateTaskValidation(t *testing.T) {
	store.Clear()
	mux := setupMux()

	body := `{"description":"no title"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing title, got %d", w.Code)
	}
}

func TestCreateTaskInvalidPriority(t *testing.T) {
	store.Clear()
	mux := setupMux()

	body := `{"title":"Bad Priority","priority":5}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid priority, got %d", w.Code)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	store.Clear()
	mux := setupMux()

	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	store.Clear()
	mux := setupMux()

	body := `{"title":"Original","description":"original desc","priority":0}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)

	updateBody := `{"title":"Updated","description":"updated desc","done":true,"priority":2}`
	req = httptest.NewRequest(http.MethodPut, "/tasks/"+created.ID, strings.NewReader(updateBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var updated model.Task
	json.NewDecoder(w.Body).Decode(&updated)
	if updated.Title != "Updated" {
		t.Fatalf("expected title 'Updated', got %q", updated.Title)
	}
	if !updated.Done {
		t.Fatal("expected task to be marked done")
	}
	if updated.CreatedAt != created.CreatedAt {
		t.Fatal("expected created_at to be preserved")
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	store.Clear()
	mux := setupMux()

	body := `{"title":"Ghost","priority":0}`
	req := httptest.NewRequest(http.MethodPut, "/tasks/999", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	store.Clear()
	mux := setupMux()

	body := `{"title":"To Delete","priority":0}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)

	req = httptest.NewRequest(http.MethodDelete, "/tasks/"+created.ID, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/tasks/"+created.ID, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", w.Code)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	store.Clear()
	mux := setupMux()

	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestTaskStats(t *testing.T) {
	store.Clear()
	mux := setupMux()

	store.Add(model.NewTask("Task 1", "desc", model.PriorityLow))
	store.Add(model.NewTask("Task 2", "desc", model.PriorityHigh))

	task3 := store.Add(model.NewTask("Task 3", "desc", model.PriorityMedium))
	task3.MarkDone()
	store.Update(task3.ID, task3)

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var stats struct {
		Total     int `json:"total"`
		Completed int `json:"completed"`
		Pending   int `json:"pending"`
	}
	json.NewDecoder(w.Body).Decode(&stats)

	if stats.Total != 3 {
		t.Fatalf("expected total 3, got %d", stats.Total)
	}
	if stats.Completed != 1 {
		t.Fatalf("expected completed 1, got %d", stats.Completed)
	}
	if stats.Pending != 2 {
		t.Fatalf("expected pending 2, got %d", stats.Pending)
	}
}

func TestFilterByDone(t *testing.T) {
	store.Clear()
	mux := setupMux()

	store.Add(model.NewTask("Pending Task", "desc", model.PriorityLow))
	doneTask := store.Add(model.NewTask("Done Task", "desc", model.PriorityLow))
	doneTask.MarkDone()
	store.Update(doneTask.ID, doneTask)

	req := httptest.NewRequest(http.MethodGet, "/tasks?done=true", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var tasks []model.Task
	json.NewDecoder(w.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 done task, got %d", len(tasks))
	}
	if tasks[0].Title != "Done Task" {
		t.Fatalf("expected 'Done Task', got %q", tasks[0].Title)
	}
}

func TestFilterByPriority(t *testing.T) {
	store.Clear()
	mux := setupMux()

	store.Add(model.NewTask("Low", "desc", model.PriorityLow))
	store.Add(model.NewTask("High 1", "desc", model.PriorityHigh))
	store.Add(model.NewTask("High 2", "desc", model.PriorityHigh))

	req := httptest.NewRequest(http.MethodGet, "/tasks?priority=2", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var tasks []model.Task
	json.NewDecoder(w.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Fatalf("expected 2 high priority tasks, got %d", len(tasks))
	}
}

func TestInvalidJSON(t *testing.T) {
	store.Clear()
	mux := setupMux()

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("{bad json"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
