package store

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/sawez-deepsource/demo-go/model"
)

var (
	mu     sync.RWMutex
	tasks  = map[string]model.Task{}
	nextID = 1
)

func All() []model.Task {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]model.Task, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}

func Get(id string) (model.Task, bool) {
	mu.RLock()
	defer mu.RUnlock()
	t, ok := tasks[id]
	return t, ok
}

func Add(t model.Task) model.Task {
	mu.Lock()
	defer mu.Unlock()
	t.ID = fmt.Sprintf("%d", nextID)
	nextID++
	now := time.Now().UTC().Format(time.RFC3339)
	if t.CreatedAt == "" {
		t.CreatedAt = now
	}
	t.UpdatedAt = now
	tasks[t.ID] = t
	return t
}

func Update(id string, updated model.Task) (model.Task, bool) {
	mu.Lock()
	defer mu.Unlock()
	existing, ok := tasks[id]
	if !ok {
		return model.Task{}, false
	}
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	tasks[id] = updated
	return updated, true
}

func Delete(id string) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := tasks[id]; !ok {
		return false
	}
	delete(tasks, id)
	return true
}

func Count() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(tasks)
}

func Clear() {
	mu.Lock()
	defer mu.Unlock()
	tasks = map[string]model.Task{}
	nextID = 1
}

func FilterByDone(done bool) []model.Task {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]model.Task, 0)
	for _, t := range tasks {
		if t.Done == done {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}

func FilterByPriority(p model.Priority) []model.Task {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]model.Task, 0)
	for _, t := range tasks {
		if t.Priority == p {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}
