package models

import (
	"context"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	wp := NewWorkerPool(2)
	if wp.workers != 2 {
		t.Errorf("Expected 2 workers, got %d", wp.workers)
	}
}

func TestWorkerPoolRun(t *testing.T) {
	wp := NewWorkerPool(2)
	domains := []string{"example.com", "google.com", "github.com"}

	testFn := func(ctx context.Context, domain string) TestResult {
		return TestResult{
			Domain:       domain,
			Status:       StatusResolved,
			ResponseTime: 10 * time.Millisecond,
			Err:          nil,
			Category:     "Test",
			Subcategory:  "Test Category",
		}
	}

	ctx := context.Background()
	results := wp.Run(ctx, domains, testFn)

	if len(results) != len(domains) {
		t.Errorf("Expected %d results, got %d", len(domains), len(results))
	}

	for i, result := range results {
		if result.Domain != domains[i] {
			t.Errorf("Expected domain %s, got %s", domains[i], result.Domain)
		}
		if result.Status != StatusResolved {
			t.Errorf("Expected status %v, got %v", StatusResolved, result.Status)
		}
	}
}

func TestWorkerPoolEmptyDomains(t *testing.T) {
	wp := NewWorkerPool(2)

	testFn := func(ctx context.Context, domain string) TestResult {
		return TestResult{}
	}

	ctx := context.Background()
	results := wp.Run(ctx, []string{}, testFn)

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty domains, got %d", len(results))
	}
}
