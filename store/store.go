package store

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
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

// -------------------------------------------------------
// Planted issues in store/store.go
// -------------------------------------------------------

// VET-V0008: Mutex copy
type SafeStore struct {
	mu   sync.Mutex
	data map[string]string
}

func CopySafeStore() {
	original := SafeStore{data: map[string]string{"a": "1"}}
	original.mu.Lock()
	copied := original // BAD: copies locked mutex
	copied.mu.Unlock()
	original.mu.Unlock()
}

// GO-W / SCC: WaitGroup Add inside goroutine
func ConcurrentLoad() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		go func(n int) {
			wg.Add(1) // BAD: Add must be before goroutine
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}
	wg.Wait()
}

// CRT-deferInLoop: Defer in loop
func DeferInLoop(ids []string) {
	for _, id := range ids {
		defer fmt.Println(id) // BAD: defers pile up
	}
}

// SCC-SA1000: Invalid regex
func BadSearch() {
	_ = regexp.MustCompile(`[invalid`) // BAD: unclosed bracket
}

// GO-W5006: Comparing address against nil
func NeverNilCheck() bool {
	s := SafeStore{}
	return &s != nil // BAD: always true
}

// GSC-G306: Insecure file permissions
func DumpTasks(data []byte) error {
	return os.WriteFile("/tmp/task_dump.txt", data, 0777) // BAD: world-writable
}

// SCC-SA6005: strings.ToLower comparison
func MatchID(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b) // BAD: use EqualFold
}

// GO-W5010: url.URL.Query() returns copy
func BadQueryMod(rawURL string) string {
	u, _ := url.Parse(rawURL)
	u.Query().Add("page", "1") // BAD: modifying copy
	return u.String()
}
