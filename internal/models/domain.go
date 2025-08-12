package models

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type DomainEntry struct {
	Name         string        `json:"name"`
	Status       ResultStatus  `json:"status"`
	ResponseTime time.Duration `json:"response_time"`
}

var BuiltInDomains = loadBuiltInDomains()

func loadBuiltInDomains() []DomainEntry {
	if domains, err := LoadDomainList(context.Background(), "domains.jsonc"); err == nil {
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

	f, err := os.Open(path)
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

func loadFromJSON(f *os.File) ([]DomainEntry, error) {
	content, err := os.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	jsonContent := StripJSONComments(string(content))

	var entries []DomainEntry

	var groupedData map[string]map[string][]string
	if err := json.Unmarshal([]byte(jsonContent), &groupedData); err != nil {
		var flatEntries []DomainEntry
		if flatErr := json.Unmarshal([]byte(jsonContent), &flatEntries); flatErr != nil {
			return nil, fmt.Errorf("failed to decode JSON as either grouped or flat format: %w", err)
		}
		return flatEntries, nil
	}

	for _, subcategories := range groupedData {
		for _, domains := range subcategories {
			for _, domain := range domains {
				if strings.TrimSpace(domain) != "" {
					entries = append(entries, DomainEntry{
						Name: strings.TrimSpace(domain),
					})
				}
			}
		}
	}

	return entries, nil
}

func StripJSONComments(content string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()

		inString := false
		escaped := false
		commentStart := -1

		for i, char := range line {
			if escaped {
				escaped = false
				continue
			}

			if char == '\\' && inString {
				escaped = true
				continue
			}

			if char == '"' {
				inString = !inString
				continue
			}

			if !inString && char == '/' && i+1 < len(line) && line[i+1] == '/' {
				commentStart = i
				break
			}
		}

		if commentStart >= 0 {
			line = strings.TrimSpace(line[:commentStart])
		}

		if line != "" {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.String()
}
