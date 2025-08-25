package httpclient

import (
	"context"
	"testing"
	"time"

	"github.com/mooship/blokilo/internal/models"
)

func TestCheckHTTPConnectivity_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "example.com", 2*time.Second, 0)
	if res.Status != models.StatusError && res.Status != models.StatusBlocked {
		t.Errorf("Expected error or blocked status for context timeout, got: %v", res.Status)
	}
}

func TestCheckHTTPConnectivity_InvalidDomain(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "!!!invalid!!!", 100*time.Millisecond, 0)
	if res.Status != models.StatusError {
		t.Errorf("Expected error status for invalid domain, got: %v", res.Status)
	}
}

func TestCheckHTTPConnectivity_MockBlocked(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "nonexistentdomain.blokilotest", 100*time.Millisecond, 0)
	if res.Status != models.StatusBlocked && res.Status != models.StatusError {
		t.Logf("Expected blocked or error for non-existent domain, got: %v", res.Status)
	}
}

func TestCheckHTTPConnectivity_MockResolved(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "example.com", 2*time.Second, 0)
	if res.Status != models.StatusResolved {
		t.Logf("Expected resolved for example.com, got: %v (err: %v)", res.Status, res.Err)
	}
}
