package concurrency

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestDo_RaceCondition(t *testing.T) {
	c := NewConcurrency().WithMaxWorker(10)

	// Add 10 functions that each produce an error
	for i := 0; i < 10; i++ {
		idx := i
		c.AddFunc(func(ctx context.Context, c Interface) {
			c.AddError(fmt.Errorf("error-%d", idx))
		})
	}

	err := c.Do(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verify all 10 errors are present
	for i := 0; i < 10; i++ {
		target := fmt.Errorf("error-%d", i)
		if !errors.Is(err, target) {
			// errors.Join uses Unwrap() []error, check manually
			found := false
			for _, e := range unwrapAll(err) {
				if e.Error() == target.Error() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error-%d to be present in result", i)
			}
		}
	}
}

func TestDo_BatchErrorIsolation(t *testing.T) {
	c := NewConcurrency().WithMaxWorker(2)

	// First batch: add functions that produce errors
	c.AddFunc(func(ctx context.Context, c Interface) {
		c.AddError(fmt.Errorf("batch1-error"))
	})
	c.AddFunc(func(ctx context.Context, c Interface) {
		// no error
	})

	err1 := c.Do(context.Background())
	if err1 == nil {
		t.Fatal("expected error from first Do(), got nil")
	}

	// Second batch: add functions with no errors
	c.AddFunc(func(ctx context.Context, c Interface) {
		// no error
	})

	err2 := c.Do(context.Background())
	if err2 != nil {
		t.Fatalf("expected nil from second Do(), got: %v", err2)
	}
}

func TestDo_ReturnsAllErrors(t *testing.T) {
	c := NewConcurrency().WithMaxWorker(5)

	errA := fmt.Errorf("error-a")
	errB := fmt.Errorf("error-b")
	errC := fmt.Errorf("error-c")

	c.AddFunc(func(ctx context.Context, c Interface) {
		c.AddError(errA)
	})
	c.AddFunc(func(ctx context.Context, c Interface) {
		c.AddError(errB)
	})
	c.AddFunc(func(ctx context.Context, c Interface) {
		c.AddError(errC)
	})

	err := c.Do(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verify all errors are present using Unwrap
	errs := unwrapAll(err)
	if len(errs) != 3 {
		t.Fatalf("expected 3 errors, got %d: %v", len(errs), errs)
	}
}

func TestDo_NoErrors(t *testing.T) {
	c := NewConcurrency().WithMaxWorker(3)

	c.AddFunc(func(ctx context.Context, c Interface) {
		// no error
	})
	c.AddFunc(func(ctx context.Context, c Interface) {
		// no error
	})

	err := c.Do(context.Background())
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

// unwrapAll extracts all errors from a joined error (errors.Join result)
func unwrapAll(err error) []error {
	type unwrapper interface {
		Unwrap() []error
	}
	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}
	return []error{err}
}
