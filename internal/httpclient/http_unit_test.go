package httpclient

import (
	"context"
	"testing"
	"time"
)

func TestHTTP_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "google.com", 8*time.Second, 1)
	if res.Status != "ERROR" && res.Status != "BLOCKED" {
		t.Errorf("expected ERROR or BLOCKED for timeout, got %v", res.Status)
	}
}

func TestHTTP_InvalidDomain(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "!!!invalid!!!", 100*time.Millisecond, 0)
	if res.Status != "ERROR" {
		t.Errorf("expected ERROR for invalid domain, got %v", res.Status)
	}
}

func TestHTTP_MockBlocked(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "nonexistentdomain.blokilotest", 100*time.Millisecond, 0)
	if res.Status != "BLOCKED" && res.Status != "ERROR" {
		t.Errorf("expected BLOCKED or ERROR, got %v", res.Status)
	}
}

func TestHTTP_MockResolved(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := CheckHTTPConnectivity(ctx, "google.com", 8*time.Second, 1)
	if res.Status != "RESOLVED" {
		t.Errorf("expected RESOLVED, got %v (err: %v)", res.Status, res.Err)
	}
}
