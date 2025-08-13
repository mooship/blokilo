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
	domainsPath := filepath.Join("..", "..", "domains.jsonc")

	if _, err := os.Stat(domainsPath); os.IsNotExist(err) {
		t.Skipf("domains.jsonc not found at %s, skipping validation", domainsPath)
	}

	t.Run("LoadDomainsJSON", func(t *testing.T) {
		domains, err := LoadDomainList(context.Background(), domainsPath)
		if err != nil {
			t.Fatalf("Failed to load domains.jsonc: %v", err)
		}

		if len(domains) == 0 {
			t.Error("No domains loaded from domains.jsonc")
		}

		t.Logf("Successfully loaded %d domains from grouped JSON format", len(domains))
	})

	t.Run("NoDuplicateDomains", func(t *testing.T) {
		content, err := os.ReadFile(domainsPath)
		if err != nil {
			t.Fatalf("Failed to read domains.jsonc: %v", err)
		}

		jsonContent := StripJSONComments(string(content))

		var groupedData map[string]map[string][]string
		if err := json.Unmarshal([]byte(jsonContent), &groupedData); err != nil {
			t.Fatalf("Failed to decode domains.jsonc: %v", err)
		}

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
			t.Fatalf("Failed to load domains.jsonc: %v", err)
		}

		var invalidDomains []string
		for _, domain := range domains {
			name := strings.TrimSpace(domain.Name)
			if name == "" {
				invalidDomains = append(invalidDomains, "empty domain")
				continue
			}

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
		content, err := os.ReadFile(domainsPath)
		if err != nil {
			t.Fatalf("Failed to read domains.jsonc: %v", err)
		}

		jsonContent := StripJSONComments(string(content))

		var groupedData map[string]map[string][]string
		if err := json.Unmarshal([]byte(jsonContent), &groupedData); err != nil {
			t.Fatalf("Failed to decode domains.jsonc: %v", err)
		}

		if len(groupedData) == 0 {
			t.Error("No categories found in domains.jsonc")
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
	flatJSON := `[
		{"name": "example.com"},
		{"name": "test.com"}
	]`

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
		if domain.Category != "" {
			t.Errorf("Expected empty category for flat JSON, got %q", domain.Category)
		}
		if domain.Subcategory != "" {
			t.Errorf("Expected empty subcategory for flat JSON, got %q", domain.Subcategory)
		}
	}
}

func TestLoadDomainListGroupedFormat(t *testing.T) {
	groupedJSON := `{
		"Test Category": {
			"Test Subcategory": [
				"example.com",
				"test.com"
			],
			"Another Subcategory": [
				"another.com"
			]
		},
		"Another Category": {
			"Sub": [
				"final.com"
			]
		}
	}`

	tmpFile, err := os.CreateTemp("", "test_domains_grouped_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(groupedJSON); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	domains, err := LoadDomainList(context.Background(), tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load grouped JSON format: %v", err)
	}

	if len(domains) != 4 {
		t.Errorf("Expected 4 domains, got %d", len(domains))
	}

	categoryMap := make(map[string]map[string][]string)
	for _, domain := range domains {
		if categoryMap[domain.Category] == nil {
			categoryMap[domain.Category] = make(map[string][]string)
		}
		categoryMap[domain.Category][domain.Subcategory] = append(
			categoryMap[domain.Category][domain.Subcategory],
			domain.Name,
		)
	}

	if len(categoryMap) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categoryMap))
	}

	if len(categoryMap["Test Category"]) != 2 {
		t.Errorf("Expected 2 subcategories in 'Test Category', got %d", len(categoryMap["Test Category"]))
	}

	foundExampleCom := false
	for _, domain := range domains {
		if domain.Name == "example.com" {
			foundExampleCom = true
			if domain.Category != "Test Category" {
				t.Errorf("Expected category 'Test Category' for example.com, got %q", domain.Category)
			}
			if domain.Subcategory != "Test Subcategory" {
				t.Errorf("Expected subcategory 'Test Subcategory' for example.com, got %q", domain.Subcategory)
			}
			break
		}
	}

	if !foundExampleCom {
		t.Error("Expected to find 'example.com' domain")
	}
}
