package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/lo"
)

type ResultStatus string

const (
	StatusBlocked  ResultStatus = "BLOCKED"
	StatusResolved ResultStatus = "RESOLVED"
	StatusError    ResultStatus = "ERROR"
)

type ClassifiedResult struct {
	Domain         string
	Status         ResultStatus
	ResponseTime   time.Duration
	Err            error
	Category       string
	Subcategory    string
	HTTPStatusCode int
}

func ClassifyResult(dnsStatus, httpStatus ResultStatus, dnsErr, httpErr error) ResultStatus {
	if dnsStatus == StatusBlocked || httpStatus == StatusBlocked {
		return StatusBlocked
	}
	if dnsStatus == StatusError || httpStatus == StatusError {
		return StatusError
	}
	if dnsStatus == StatusResolved || httpStatus == StatusResolved {
		return StatusResolved
	}
	return StatusError
}

type Stats struct {
	Total           int
	Blocked         int
	Resolved        int
	Errored         int
	PercentBlocked  float64
	PercentResolved float64
}

type GroupedResults struct {
	Category    string
	Subcategory string
	Results     []ClassifiedResult
	Stats       Stats
}

type CategoryGroup struct {
	Category      string
	Subcategories []GroupedResults
	Stats         Stats
}

type CategoryConfig struct {
	CategoryOrder    []string            `json:"categoryOrder"`
	SubcategoryOrder map[string][]string `json:"subcategoryOrder"`
}

func LoadCategoryConfig(path string) (*CategoryConfig, error) {
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

	file, err := os.ReadFile(cleaned)
	if err != nil {
		return nil, err
	}

	var config CategoryConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func ComputeStats(results []ClassifiedResult) Stats {
	s := Stats{
		Total:    len(results),
		Blocked:  lo.CountBy(results, func(r ClassifiedResult) bool { return r.Status == StatusBlocked }),
		Resolved: lo.CountBy(results, func(r ClassifiedResult) bool { return r.Status == StatusResolved }),
		Errored:  lo.CountBy(results, func(r ClassifiedResult) bool { return r.Status == StatusError }),
	}

	if s.Total > 0 {
		s.PercentBlocked = float64(s.Blocked) / float64(s.Total) * 100
		s.PercentResolved = float64(s.Resolved) / float64(s.Total) * 100
	}
	return s
}

func GroupResultsByCategory(results []ClassifiedResult, config *CategoryConfig) []CategoryGroup {
	categoryMap := lo.GroupBy(results, func(r ClassifiedResult) string {
		if r.Category == "" {
			return "Uncategorized"
		}
		return r.Category
	})

	orderedCategories := lo.Filter(config.CategoryOrder, func(cat string, _ int) bool {
		_, exists := categoryMap[cat]
		return exists
	})

	uncategorized := lo.FilterMap(lo.Keys(categoryMap), func(cat string, _ int) (string, bool) {
		return cat, !lo.Contains(config.CategoryOrder, cat)
	})

	allCategories := append(orderedCategories, uncategorized...)

	return lo.Map(allCategories, func(categoryName string, _ int) CategoryGroup {
		subcategories := lo.GroupBy(categoryMap[categoryName], func(r ClassifiedResult) string {
			if r.Subcategory == "" {
				return "Other"
			}
			return r.Subcategory
		})

		subOrder, hasOrder := config.SubcategoryOrder[categoryName]
		if !hasOrder {
			subOrder = lo.Keys(subcategories)
		}

		orderedSubcategories := lo.Filter(subOrder, func(subcat string, _ int) bool {
			_, exists := subcategories[subcat]
			return exists
		})

		unOrderedSubcategories := lo.FilterMap(lo.Keys(subcategories), func(subcat string, _ int) (string, bool) {
			return subcat, !lo.Contains(subOrder, subcat)
		})

		allSubcategories := append(orderedSubcategories, unOrderedSubcategories...)

		subcategoryGroups := lo.Map(allSubcategories, func(subcategoryName string, _ int) GroupedResults {
			subcategoryResults := subcategories[subcategoryName]
			return GroupedResults{
				Category:    categoryName,
				Subcategory: subcategoryName,
				Results:     subcategoryResults,
				Stats:       ComputeStats(subcategoryResults),
			}
		})

		allCategoryResults := lo.Flatten(lo.Map(subcategoryGroups, func(g GroupedResults, _ int) []ClassifiedResult {
			return g.Results
		}))

		return CategoryGroup{
			Category:      categoryName,
			Subcategories: subcategoryGroups,
			Stats:         ComputeStats(allCategoryResults),
		}
	})
}
