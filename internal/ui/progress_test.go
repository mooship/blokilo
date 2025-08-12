package ui

import (
	"strings"
	"testing"
)

func TestNewProgressModel(t *testing.T) {
	total := 100
	progress := NewProgressModel(total)

	if progress.Total != total {
		t.Errorf("expected Total to be %d, got %d", total, progress.Total)
	}

	if progress.Current != 0 {
		t.Errorf("expected Current to be 0, got %d", progress.Current)
	}

	if progress.Domain != "" {
		t.Errorf("expected Domain to be empty, got %s", progress.Domain)
	}

	if progress.DNSAddr != "" {
		t.Errorf("expected DNSAddr to be empty, got %s", progress.DNSAddr)
	}
}

func TestProgressView(t *testing.T) {
	progress := NewProgressModel(10)
	progress.Current = 3
	progress.Domain = "example.com"
	progress.DNSAddr = "1.1.1.1:53"

	view := progress.View()

	expectedElements := []string{
		"Testing: example.com",
		"3/10",
		"DNS: 1.1.1.1:53",
	}

	for _, element := range expectedElements {
		if !strings.Contains(view, element) {
			t.Errorf("progress view should contain '%s'", element)
		}
	}
}

func TestProgressViewWithSystemDNS(t *testing.T) {
	progress := NewProgressModel(5)
	progress.Current = 2
	progress.Domain = "test.com"
	progress.DNSAddr = ""

	view := progress.View()

	if !strings.Contains(view, "DNS:") {
		t.Error("progress view should contain DNS information")
	}

	if !strings.Contains(view, "Testing: test.com") {
		t.Error("progress view should contain current domain")
	}

	if !strings.Contains(view, "2/5") {
		t.Error("progress view should contain progress counter")
	}
}

func TestProgressViewProgressBar(t *testing.T) {
	progress := NewProgressModel(4)
	progress.Current = 2
	progress.Domain = "halfway.com"

	view := progress.View()

	lines := strings.Split(view, "\n")

	if len(lines) < 5 {
		t.Errorf("progress view should have at least 5 lines, got %d", len(lines))
	}

	if !strings.Contains(lines[0], "Testing: halfway.com") {
		t.Error("first line should contain testing information")
	}

	foundCounter := false
	for _, line := range lines {
		if strings.Contains(line, "2/4") {
			foundCounter = true
			break
		}
	}
	if !foundCounter {
		t.Error("should contain progress counter 2/4")
	}
}

func TestProgressViewEdgeCases(t *testing.T) {
	progress := NewProgressModel(10)
	progress.Current = 0
	progress.Domain = "first.com"
	progress.DNSAddr = "8.8.8.8:53"

	view := progress.View()
	if !strings.Contains(view, "0/10") {
		t.Error("should handle zero progress correctly")
	}

	progress.Current = 10
	progress.Domain = "last.com"

	view = progress.View()
	if !strings.Contains(view, "10/10") {
		t.Error("should handle complete progress correctly")
	}

	if !strings.Contains(view, "Testing: last.com") {
		t.Error("should show current domain even when complete")
	}
}

func TestProgressViewLongDomain(t *testing.T) {
	progress := NewProgressModel(1)
	progress.Current = 1
	progress.Domain = "very-long-domain-name-that-might-cause-display-issues.example.com"
	progress.DNSAddr = "1.1.1.1:53"

	view := progress.View()

	if !strings.Contains(view, progress.Domain) {
		t.Error("should display full domain name even if long")
	}
}
