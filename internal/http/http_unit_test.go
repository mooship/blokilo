package httptest

import (
	"context"
	"testing"
	"time"
)

func TestHTTP_MockBlocked(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "nonexistentdomain.blokilotest", 100*time.Millisecond, 0)
	if res.Status != "BLOCKED" && res.Status != "ERROR" {
		t.Errorf("expected BLOCKED or ERROR, got %v", res.Status)
	}
}

func TestHTTP_MockResolved(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "example.com", 2*time.Second, 0)
	if res.Status != "RESOLVED" {
		t.Errorf("expected RESOLVED, got %v (err: %v)", res.Status, res.Err)
	}
}
