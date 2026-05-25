package concurrency_test

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/alpardfm/go-toolkit/concurrency"
)

func ExampleNewConcurrency() {
	// Create a worker pool with 2 max workers
	pool := concurrency.NewConcurrency().WithMaxWorker(2)

	var mu sync.Mutex
	results := make([]int, 0, 3)

	// Add tasks to the pool
	pool.AddFunc(func(ctx context.Context, c concurrency.Interface) {
		mu.Lock()
		results = append(results, 1)
		mu.Unlock()
	})
	pool.AddFunc(func(ctx context.Context, c concurrency.Interface) {
		mu.Lock()
		results = append(results, 2)
		mu.Unlock()
	})
	pool.AddFunc(func(ctx context.Context, c concurrency.Interface) {
		mu.Lock()
		results = append(results, 3)
		mu.Unlock()
	})

	// Execute all tasks concurrently
	err := pool.Do(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Sort for deterministic output
	sort.Ints(results)
	fmt.Println(results)

	// Output:
	// [1 2 3]
}
