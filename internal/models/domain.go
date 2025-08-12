package models

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	if domains, err := LoadDomainList(context.Background(), "domains.json"); err == nil {
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
	case strings.HasSuffix(strings.ToLower(path), ".txt"):
		return loadFromTxt(ctx, f)
	case strings.HasSuffix(strings.ToLower(path), ".csv"):
		return loadFromCSV(ctx, f)
	case strings.HasSuffix(strings.ToLower(path), ".json"):
		return loadFromJSON(f)
	default:
		return nil, fmt.Errorf("unsupported file type for %q: must be .txt, .csv, or .json", path)
	}
}

func loadFromTxt(ctx context.Context, f *os.File) ([]DomainEntry, error) {
	scanner := bufio.NewScanner(f)
	var entries []DomainEntry
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entries = append(entries, DomainEntry{Name: line})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func loadFromCSV(ctx context.Context, f *os.File) ([]DomainEntry, error) {
	r := csv.NewReader(f)
	var entries []DomainEntry
	for {
		rec, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if len(rec) == 0 || strings.TrimSpace(rec[0]) == "" {
			continue
		}
		entries = append(entries, DomainEntry{Name: strings.TrimSpace(rec[0])})
	}
	return entries, nil
}

func loadFromJSON(f *os.File) ([]DomainEntry, error) {
	var entries []DomainEntry
	dec := json.NewDecoder(f)

	var groupedData map[string]map[string][]string
	if err := dec.Decode(&groupedData); err != nil {
		if _, seekErr := f.Seek(0, 0); seekErr != nil {
			return nil, fmt.Errorf("failed to reset file pointer: %w", seekErr)
		}

		dec = json.NewDecoder(f)
		var flatEntries []DomainEntry
		if flatErr := dec.Decode(&flatEntries); flatErr != nil {
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
