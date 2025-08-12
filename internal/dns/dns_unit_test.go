package dns

import (
	"context"
	"testing"
	"time"

	"github.com/mooship/blokilo/internal/models"
)

func TestDNS_MockBlocked(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := TestDomainDNS(ctx, "example.com", "0.0.0.0:53")
	if res.Status != models.StatusBlocked && res.Status != models.StatusError {
		t.Errorf("expected Blocked or Error, got %v", res.Status)
	}
}

func TestDNS_MockResolved(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	res := TestDomainDNS(ctx, "example.com", "8.8.8.8:53")
	if res.Status != models.StatusResolved {
		t.Errorf("expected Resolved, got %v (err: %v)", res.Status, res.Err)
	}
}
