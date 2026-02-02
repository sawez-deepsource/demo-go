package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sawez-deepsource/demo-go/handler"
)

func TestCreateAndListTasks(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("GET /tasks", handler.ListTasks)

	// Create a task
	body := `{"title":"Test Task","description":"A test"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// List tasks
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
