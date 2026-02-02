package store

import (
	"fmt"
	"sync"

	"demo-go/model"
)

var (
	mu    sync.RWMutex
	tasks = map[string]model.Task{}
	nextID = 1
)

func All() []model.Task {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]model.Task, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, t)
	}
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
	tasks[t.ID] = t
	return t
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
