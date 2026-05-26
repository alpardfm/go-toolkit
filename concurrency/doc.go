// Package concurrency provides a goroutine worker pool that executes functions
// in batches with a configurable maximum number of concurrent workers.
//
// The pool collects errors from all goroutines and returns them as a single
// joined error (via errors.Join) when Do() completes. Context cancellation is
// respected between batches, allowing early termination of remaining work.
//
// Basic usage:
//
//	c := concurrency.NewConcurrency().WithMaxWorker(5)
//	c.AddFunc(func(ctx context.Context, c concurrency.Interface) {
//	    // do work
//	    if err != nil {
//	        c.AddError(err)
//	    }
//	})
//	err := c.Do(ctx)
package concurrency
