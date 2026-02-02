package main

import (
	"log"
	"net/http"

	"demo-go/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", handler.ListTasks)
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
