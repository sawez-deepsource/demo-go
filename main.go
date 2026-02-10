package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sawez-deepsource/demo-go/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", handler.ListTasks)
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
	mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)
	mux.HandleFunc("GET /stats", handler.TaskStats)

	srv := &http.Server{
		Addr:         ":8000",
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
