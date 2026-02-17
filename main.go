package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sawez-deepsource/demo-go/handler"
)

// GSC-G101: Hardcoded credentials
var adminToken = "Bearer super-secret-admin-token-12345"

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

// -------------------------------------------------------
// Planted issues in main.go
// -------------------------------------------------------

// SCC-SA1004: Suspiciously small time.Sleep
func warmupDelay() {
	time.Sleep(10) // BAD: 10 nanoseconds
}

// CRT-A0001: Shadowing builtin
func shadowLen() int {
	len := 42 // BAD: shadows builtin
	return len
}

// GO-W: Goroutine leak — channel never read
func leakyStartup() {
	ch := make(chan int)
	go func() {
		ch <- 1 // BAD: blocks forever
	}()
}

// SCC-SA2001: Empty critical section
func emptyLock() {
	var mu sync.Mutex
	mu.Lock()
	mu.Unlock() // BAD: empty critical section
}

// VET-V0013: Printf format mismatch
func badLog() {
	port := "8000"
	fmt.Printf("listening on port %d\n", port) // BAD: %d for string
}

// SCC-U1000: unused
var _ = adminToken
