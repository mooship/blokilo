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
		{Name: "googleads.g.doubleclick.net"},
		{Name: "pagead2.googlesyndication.com"},
		{Name: "doubleclick.net"},
		{Name: "ads.yahoo.com"},
		{Name: "ads.twitter.com"},
		{Name: "ads.facebook.com"},
		{Name: "scorecardresearch.com"},
		{Name: "securepubads.g.doubleclick.net"},
		{Name: "googletagservices.com"},
		{Name: "googletagmanager.com"},
		{Name: "amazon-adsystem.com"},
		{Name: "googlesyndication.com"},
		{Name: "googleadservices.com"},
		{Name: "connect.facebook.net"},
		{Name: "analytics.google.com"},
		{Name: "www.google-analytics.com"},
		{Name: "ads.pinterest.com"},
		{Name: "ads.linkedin.com"},
		{Name: "outbrain.com"},
		{Name: "taboola.com"},
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
	if err := dec.Decode(&entries); err != nil {
		return nil, err
	}
	return entries, nil
}
