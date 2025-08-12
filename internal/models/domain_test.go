package models

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDomainsJSONValidation(t *testing.T) {
	// Get the path to domains.json relative to the project root
	domainsPath := filepath.Join("..", "..", "domains.json")

	// Check if the file exists
	if _, err := os.Stat(domainsPath); os.IsNotExist(err) {
		t.Skipf("domains.json not found at %s, skipping validation", domainsPath)
	}

	t.Run("LoadDomainsJSON", func(t *testing.T) {
		domains, err := LoadDomainList(context.Background(), domainsPath)
		if err != nil {
			t.Fatalf("Failed to load domains.json: %v", err)
		}

		if len(domains) == 0 {
			t.Error("No domains loaded from domains.json")
		}

		t.Logf("Successfully loaded %d domains from grouped JSON format", len(domains))
	})

	t.Run("NoDuplicateDomains", func(t *testing.T) {
		f, err := os.Open(domainsPath)
		if err != nil {
			t.Fatalf("Failed to open domains.json: %v", err)
		}
		defer f.Close()

		var groupedData map[string]map[string][]string
		dec := json.NewDecoder(f)
		if err := dec.Decode(&groupedData); err != nil {
			t.Fatalf("Failed to decode domains.json: %v", err)
		}

		// Track all domains and their locations for duplicate detection
		domainMap := make(map[string][]string)
		totalDomains := 0

		for category, subcategories := range groupedData {
			for subcategory, domains := range subcategories {
				for _, domain := range domains {
					domain = strings.TrimSpace(domain)
					if domain == "" {
						continue
					}

					location := category + " -> " + subcategory
					if existing, found := domainMap[domain]; found {
						existing = append(existing, location)
						domainMap[domain] = existing
					} else {
						domainMap[domain] = []string{location}
					}
					totalDomains++
				}
			}
		}

		// Check for duplicates
		var duplicates []string
		for domain, locations := range domainMap {
			if len(locations) > 1 {
				duplicates = append(duplicates, domain+" (found in: "+strings.Join(locations, ", ")+")")
			}
		}

		if len(duplicates) > 0 {
			t.Errorf("Found %d duplicate domain(s):\n%s", len(duplicates), strings.Join(duplicates, "\n"))
		}

		t.Logf("Validated %d unique domains across %d total entries", len(domainMap), totalDomains)
	})

	t.Run("ValidDomainFormat", func(t *testing.T) {
		domains, err := LoadDomainList(context.Background(), domainsPath)
		if err != nil {
			t.Fatalf("Failed to load domains.json: %v", err)
		}

		var invalidDomains []string
		for _, domain := range domains {
			name := strings.TrimSpace(domain.Name)
			if name == "" {
				invalidDomains = append(invalidDomains, "empty domain")
				continue
			}

			// Basic domain validation - check for basic format
			if strings.Contains(name, " ") {
				invalidDomains = append(invalidDomains, name+" (contains spaces)")
			}
			if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") {
				invalidDomains = append(invalidDomains, name+" (starts or ends with dot)")
			}
			if strings.Contains(name, "..") {
				invalidDomains = append(invalidDomains, name+" (contains consecutive dots)")
			}
		}

		if len(invalidDomains) > 0 {
			t.Errorf("Found %d invalid domain(s):\n%s", len(invalidDomains), strings.Join(invalidDomains, "\n"))
		}
	})

	t.Run("GroupedStructureIntegrity", func(t *testing.T) {
		f, err := os.Open(domainsPath)
		if err != nil {
			t.Fatalf("Failed to open domains.json: %v", err)
		}
		defer f.Close()

		var groupedData map[string]map[string][]string
		dec := json.NewDecoder(f)
		if err := dec.Decode(&groupedData); err != nil {
			t.Fatalf("Failed to decode domains.json: %v", err)
		}

		if len(groupedData) == 0 {
			t.Error("No categories found in domains.json")
		}

		categoryCount := 0
		subcategoryCount := 0
		for category, subcategories := range groupedData {
			categoryCount++
			if category == "" {
				t.Error("Found empty category name")
			}
			if len(subcategories) == 0 {
				t.Errorf("Category %q has no subcategories", category)
			}

			for subcategory, domains := range subcategories {
				subcategoryCount++
				if subcategory == "" {
					t.Errorf("Found empty subcategory name in category %q", category)
				}
				if len(domains) == 0 {
					t.Errorf("Subcategory %q in category %q has no domains", subcategory, category)
				}
			}
		}

		t.Logf("Validated structure: %d categories, %d subcategories", categoryCount, subcategoryCount)
	})
}

func TestLoadDomainListBackwardCompatibility(t *testing.T) {
	// Test that the function can still handle flat JSON arrays
	flatJSON := `[
		{"name": "example.com"},
		{"name": "test.com"}
	]`

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_domains_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(flatJSON); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	domains, err := LoadDomainList(context.Background(), tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load flat JSON format: %v", err)
	}

	if len(domains) != 2 {
		t.Errorf("Expected 2 domains, got %d", len(domains))
	}

	expectedDomains := []string{"example.com", "test.com"}
	for i, domain := range domains {
		if domain.Name != expectedDomains[i] {
			t.Errorf("Expected domain %q, got %q", expectedDomains[i], domain.Name)
		}
	}
}
