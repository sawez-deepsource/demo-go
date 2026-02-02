package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sawez-deepsource/demo-go/model"
	"github.com/sawez-deepsource/demo-go/store"
)

func ListTasks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, store.All())
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, ok := store.Get(id)
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	created := store.Add(t)
	writeJSON(w, http.StatusCreated, created)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !store.Delete(id) {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
