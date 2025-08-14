package models

import (
	"encoding/json"
	"os"
	"slices"
	"time"
)

type ResultStatus string

const (
	StatusBlocked  ResultStatus = "BLOCKED"
	StatusResolved ResultStatus = "RESOLVED"
	StatusError    ResultStatus = "ERROR"
)

type ClassifiedResult struct {
	Domain       string
	Status       ResultStatus
	ResponseTime time.Duration
	Err          error
	Category     string
	Subcategory  string
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
	file, err := os.ReadFile(path)
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
	var s Stats
	s.Total = len(results)
	for _, r := range results {
		switch r.Status {
		case StatusBlocked:
			s.Blocked++
		case StatusResolved:
			s.Resolved++
		case StatusError:
			s.Errored++
		}
	}
	if s.Total > 0 {
		s.PercentBlocked = float64(s.Blocked) / float64(s.Total) * 100
		s.PercentResolved = float64(s.Resolved) / float64(s.Total) * 100
	}
	return s
}

func GroupResultsByCategory(results []ClassifiedResult, config *CategoryConfig) []CategoryGroup {
	categoryOrder := config.CategoryOrder
	subcategoryOrder := config.SubcategoryOrder

	categoryMap := make(map[string]map[string][]ClassifiedResult)

	for _, result := range results {
		category := result.Category
		if category == "" {
			category = "Uncategorized"
		}
		subcategory := result.Subcategory
		if subcategory == "" {
			subcategory = "Other"
		}

		if categoryMap[category] == nil {
			categoryMap[category] = make(map[string][]ClassifiedResult)
		}
		categoryMap[category][subcategory] = append(categoryMap[category][subcategory], result)
	}

	var categories []CategoryGroup

	for _, categoryName := range categoryOrder {
		subcategories, exists := categoryMap[categoryName]
		if !exists {
			continue
		}

		var subcategoryGroups []GroupedResults
		var allCategoryResults []ClassifiedResult

		subOrder, hasOrder := subcategoryOrder[categoryName]
		if !hasOrder {
			subOrder = make([]string, 0, len(subcategories))
			for subcategoryName := range subcategories {
				subOrder = append(subOrder, subcategoryName)
			}
		}

		for _, subcategoryName := range subOrder {
			subcategoryResults, exists := subcategories[subcategoryName]
			if !exists {
				continue
			}

			subcategoryStats := ComputeStats(subcategoryResults)
			subcategoryGroups = append(subcategoryGroups, GroupedResults{
				Category:    categoryName,
				Subcategory: subcategoryName,
				Results:     subcategoryResults,
				Stats:       subcategoryStats,
			})
			allCategoryResults = append(allCategoryResults, subcategoryResults...)
		}

		for subcategoryName, subcategoryResults := range subcategories {
			found := slices.Contains(subOrder, subcategoryName)
			if !found {
				subcategoryStats := ComputeStats(subcategoryResults)
				subcategoryGroups = append(subcategoryGroups, GroupedResults{
					Category:    categoryName,
					Subcategory: subcategoryName,
					Results:     subcategoryResults,
					Stats:       subcategoryStats,
				})
				allCategoryResults = append(allCategoryResults, subcategoryResults...)
			}
		}

		if len(subcategoryGroups) > 0 {
			categoryStats := ComputeStats(allCategoryResults)
			categories = append(categories, CategoryGroup{
				Category:      categoryName,
				Subcategories: subcategoryGroups,
				Stats:         categoryStats,
			})
		}
	}

	for categoryName, subcategories := range categoryMap {
		found := slices.Contains(categoryOrder, categoryName)
		if !found {
			var subcategoryGroups []GroupedResults
			var allCategoryResults []ClassifiedResult

			for subcategoryName, subcategoryResults := range subcategories {
				subcategoryStats := ComputeStats(subcategoryResults)
				subcategoryGroups = append(subcategoryGroups, GroupedResults{
					Category:    categoryName,
					Subcategory: subcategoryName,
					Results:     subcategoryResults,
					Stats:       subcategoryStats,
				})
				allCategoryResults = append(allCategoryResults, subcategoryResults...)
			}

			categoryStats := ComputeStats(allCategoryResults)
			categories = append(categories, CategoryGroup{
				Category:      categoryName,
				Subcategories: subcategoryGroups,
				Stats:         categoryStats,
			})
		}
	}

	return categories
}
