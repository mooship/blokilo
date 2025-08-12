package models

import (
	"testing"
)

func TestEnsureDNSPort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty input", "", "8.8.8.8:53"},
		{"Valid IP with port", "1.1.1.1:53", "1.1.1.1:53"},
		{"Valid IP without port", "1.1.1.1", "1.1.1.1:53"},
		{"Valid hostname with port", "dns.google:53", "dns.google:53"},
		{"Valid hostname without port", "dns.google", "dns.google:53"},
		{"Invalid port", "1.1.1.1:abc", "8.8.8.8:53"},
		{"Port out of range", "1.1.1.1:99999", "8.8.8.8:53"},
		{"Invalid format", "not-a-valid-address", "8.8.8.8:53"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnsureDNSPort(tt.input)
			if result != tt.expected {
				t.Errorf("EnsureDNSPort(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
