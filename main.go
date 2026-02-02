package main

import (
	"log"
	"net/http"

	"github.com/sawez-deepsource/demo-go/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", handler.ListTasks)
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)

	log.Println("Server starting on :8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}
