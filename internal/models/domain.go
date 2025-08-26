package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/lo"
)

type DomainEntry struct {
	Name         string        `json:"name"`
	Status       ResultStatus  `json:"status"`
	ResponseTime time.Duration `json:"response_time"`
	Category     string        `json:"category,omitempty"`
	Subcategory  string        `json:"subcategory,omitempty"`
}

var BuiltInDomains = loadBuiltInDomains()

func loadBuiltInDomains() []DomainEntry {
	if domains, err := LoadDomainList(context.Background(), "data/domains.jsonc"); err == nil {
		return domains
	}

	return []DomainEntry{
		{Name: "pagead2.googlesyndication.com"},
		{Name: "googletagmanager.com"},
		{Name: "amazon-adsystem.com"},
		{Name: "ad.doubleclick.net"},
		{Name: "outbrain.com"},
		{Name: "taboola.com"},
		{Name: "analytics.google.com"},
		{Name: "omtrdc.net"},
		{Name: "mixpanel.com"},
		{Name: "script.hotjar.com"},
		{Name: "pixel.facebook.com"},
		{Name: "ads.linkedin.com"},
		{Name: "ads.pinterest.com"},
		{Name: "ads.youtube.com"},
		{Name: "ads.tiktok.com"},
		{Name: "notify.bugsnag.com"},
		{Name: "browser.sentry-cdn.com"},
		{Name: "ads.yahoo.com"},
		{Name: "api.ad.xiaomi.com"},
		{Name: "samsungads.com"},
	}
}

func LoadDomainList(ctx context.Context, path string) ([]DomainEntry, error) {
	if path == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	cleaned := filepath.Clean(path)

	allowed := false
	tempDir := filepath.Clean(os.TempDir())
	if filepath.IsAbs(cleaned) {
		s := filepath.ToSlash(cleaned)
		if strings.HasPrefix(s, filepath.ToSlash(tempDir)+"/") || strings.Contains(s, "/data/") {
			allowed = true
		}
	} else {
		s := filepath.ToSlash(cleaned)
		if strings.HasPrefix(s, "data/") || strings.Contains(s, "/data/") {
			allowed = true
		}
	}

	if !allowed {
		return nil, fmt.Errorf("only files under the data/ directory or temp dir are allowed: %q", path)
	}

	f, err := os.Open(cleaned)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer f.Close()

	switch {
	case strings.HasSuffix(strings.ToLower(path), ".json"), strings.HasSuffix(strings.ToLower(path), ".jsonc"):
		return loadFromJSON(f)
	default:
		return nil, fmt.Errorf("unsupported file type for %q: must be .json or .jsonc", path)
	}
}

func loadFromJSON(f io.Reader) ([]DomainEntry, error) {
	contentBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	content := string(contentBytes)

	jsonContent := StripJSONComments(string(content))

	var groupedData map[string]map[string][]string
	if err := json.Unmarshal([]byte(jsonContent), &groupedData); err != nil {
		var flatEntries []DomainEntry
		if flatErr := json.Unmarshal([]byte(jsonContent), &flatEntries); flatErr != nil {
			return nil, fmt.Errorf("failed to decode JSON as either grouped or flat format: %w", err)
		}
		return flatEntries, nil
	}

	return lo.Flatten(lo.MapToSlice(groupedData, func(category string, subcategories map[string][]string) []DomainEntry {
		return lo.Flatten(lo.MapToSlice(subcategories, func(subcategory string, domains []string) []DomainEntry {
			return lo.Map(domains, func(domain string, _ int) DomainEntry {
				return DomainEntry{
					Name:        strings.TrimSpace(domain),
					Category:    category,
					Subcategory: subcategory,
				}
			})
		}))
	})), nil
}

func StripJSONComments(content string) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "//") {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}
	return result.String()
}
