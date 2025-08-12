package dns

import (
	"context"
	"testing"
	"time"

	"github.com/mooship/blokilo/internal/models"
)

func TestDNS(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := TestDomainDNS(ctx, "example.com", "8.8.8.8:53")
	if result.Status == models.StatusError {
		t.Errorf("unexpected error: %v", result.Err)
	}

	if result.Domain != "example.com" {
		t.Errorf("expected domain 'example.com', got '%s'", result.Domain)
	}
}
