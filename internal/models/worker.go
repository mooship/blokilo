package models

import (
	"context"
	"sync"
	"time"
)

type TestResult struct {
	Domain       string
	Status       ResultStatus
	ResponseTime time.Duration
	Err          error
	Category     string
	Subcategory  string
}

type TestFunc func(ctx context.Context, domain string) TestResult

type WorkerPool struct {
	workers int
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{workers: workers}
}

func (wp *WorkerPool) Run(ctx context.Context, domains []string, testFn TestFunc) []TestResult {
	if len(domains) == 0 {
		return []TestResult{}
	}

	results := make([]TestResult, len(domains))
	var wg sync.WaitGroup
	resultsCh := make(chan struct {
		idx int
		res TestResult
	}, min(len(domains), 100))

	sem := make(chan struct{}, wp.workers)

	for i, domain := range domains {
		wg.Add(1)
		go func(idx int, d string) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			res := testFn(ctx, d)

			select {
			case resultsCh <- struct {
				idx int
				res TestResult
			}{idx, res}:
			case <-ctx.Done():
				return
			}
		}(i, domain)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for r := range resultsCh {
		results[r.idx] = r.res
	}

	return results
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
