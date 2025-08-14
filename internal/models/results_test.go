package models

import (
	"testing"
	"time"
)

func TestClassifyResult(t *testing.T) {
	cases := []struct {
		dns, http       ResultStatus
		dnsErr, httpErr error
		expect          ResultStatus
	}{
		{StatusBlocked, StatusBlocked, nil, nil, StatusBlocked},
		{StatusResolved, StatusBlocked, nil, nil, StatusBlocked},
		{StatusBlocked, StatusResolved, nil, nil, StatusBlocked},
		{StatusResolved, StatusResolved, nil, nil, StatusResolved},
		{StatusError, StatusBlocked, nil, nil, StatusBlocked},
		{StatusError, StatusError, nil, nil, StatusError},
	}
	for _, c := range cases {
		got := ClassifyResult(c.dns, c.http, c.dnsErr, c.httpErr)
		if got != c.expect {
			t.Errorf("ClassifyResult(%v, %v) = %v, want %v", c.dns, c.http, got, c.expect)
		}
	}
}

func TestComputeStats(t *testing.T) {
	results := []ClassifiedResult{
		{Status: StatusBlocked},
		{Status: StatusResolved},
		{Status: StatusResolved},
		{Status: StatusError},
	}

	stats := ComputeStats(results)

	if stats.Total != 4 {
		t.Errorf("Expected total 4, got %d", stats.Total)
	}
	if stats.Blocked != 1 {
		t.Errorf("Expected blocked 1, got %d", stats.Blocked)
	}
	if stats.Resolved != 2 {
		t.Errorf("Expected resolved 2, got %d", stats.Resolved)
	}
	if stats.Errored != 1 {
		t.Errorf("Expected errored 1, got %d", stats.Errored)
	}

	expectedBlocked := 25.0
	if stats.PercentBlocked != expectedBlocked {
		t.Errorf("Expected blocked percent %.1f, got %.1f", expectedBlocked, stats.PercentBlocked)
	}

	expectedResolved := 50.0
	if stats.PercentResolved != expectedResolved {
		t.Errorf("Expected resolved percent %.1f, got %.1f", expectedResolved, stats.PercentResolved)
	}
}

func TestGroupResultsByCategory(t *testing.T) {
	results := []ClassifiedResult{
		{
			Domain:       "example.com",
			Status:       StatusBlocked,
			ResponseTime: 10 * time.Millisecond,
			Category:     "Ads",
			Subcategory:  "Google Ads",
		},
		{
			Domain:       "test.com",
			Status:       StatusResolved,
			ResponseTime: 5 * time.Millisecond,
			Category:     "Ads",
			Subcategory:  "Google Ads",
		},
		{
			Domain:       "analytics.com",
			Status:       StatusBlocked,
			ResponseTime: 15 * time.Millisecond,
			Category:     "Analytics",
			Subcategory:  "Google Analytics",
		},
		{
			Domain:       "uncategorized.com",
			Status:       StatusError,
			ResponseTime: 0,
			Category:     "",
			Subcategory:  "",
		},
	}

	config := &CategoryConfig{
		CategoryOrder: []string{"Ads", "Analytics", "Uncategorized"},
	}
	groups := GroupResultsByCategory(results, config)

	if len(groups) != 3 {
		t.Errorf("Expected 3 category groups, got %d", len(groups))
	}

	expectedOrder := []string{"Ads", "Analytics", "Uncategorized"}
	for i, group := range groups {
		if group.Category != expectedOrder[i] {
			t.Errorf("Expected category %d to be %q, got %q", i, expectedOrder[i], group.Category)
		}
	}

	adsGroup := &groups[0]
	if adsGroup.Category != "Ads" {
		t.Fatalf("Expected first category to be 'Ads', got %q", adsGroup.Category)
	}

	if len(adsGroup.Subcategories) != 1 {
		t.Errorf("Expected 1 subcategory in Ads group, got %d", len(adsGroup.Subcategories))
	}

	googleAdsSubgroup := adsGroup.Subcategories[0]
	if googleAdsSubgroup.Subcategory != "Google Ads" {
		t.Errorf("Expected subcategory 'Google Ads', got %q", googleAdsSubgroup.Subcategory)
	}

	if len(googleAdsSubgroup.Results) != 2 {
		t.Errorf("Expected 2 results in Google Ads subcategory, got %d", len(googleAdsSubgroup.Results))
	}

	uncategorizedGroup := &groups[2]
	if uncategorizedGroup.Category != "Uncategorized" {
		t.Fatalf("Expected third category to be 'Uncategorized', got %q", uncategorizedGroup.Category)
	}

	if len(uncategorizedGroup.Subcategories) != 1 {
		t.Errorf("Expected 1 subcategory in Uncategorized group, got %d", len(uncategorizedGroup.Subcategories))
	}

	if uncategorizedGroup.Subcategories[0].Subcategory != "Other" {
		t.Errorf("Expected subcategory 'Other' for uncategorized, got %q", uncategorizedGroup.Subcategories[0].Subcategory)
	}
}

func TestGroupResultsByCategoryStats(t *testing.T) {
	results := []ClassifiedResult{
		{
			Domain:      "blocked1.com",
			Status:      StatusBlocked,
			Category:    "Ads",
			Subcategory: "Google Ads",
		},
		{
			Domain:      "blocked2.com",
			Status:      StatusBlocked,
			Category:    "Ads",
			Subcategory: "Google Ads",
		},
		{
			Domain:      "resolved1.com",
			Status:      StatusResolved,
			Category:    "Ads",
			Subcategory: "Google Ads",
		},
	}

	config := &CategoryConfig{
		CategoryOrder: []string{"Ads"},
	}
	groups := GroupResultsByCategory(results, config)

	if len(groups) != 1 {
		t.Errorf("Expected 1 category group, got %d", len(groups))
	}

	adsGroup := groups[0]
	if adsGroup.Stats.Total != 3 {
		t.Errorf("Expected total 3 in category stats, got %d", adsGroup.Stats.Total)
	}
	if adsGroup.Stats.Blocked != 2 {
		t.Errorf("Expected blocked 2 in category stats, got %d", adsGroup.Stats.Blocked)
	}
	if adsGroup.Stats.Resolved != 1 {
		t.Errorf("Expected resolved 1 in category stats, got %d", adsGroup.Stats.Resolved)
	}

	subcategoryStats := adsGroup.Subcategories[0].Stats
	if subcategoryStats.Total != 3 {
		t.Errorf("Expected total 3 in subcategory stats, got %d", subcategoryStats.Total)
	}
	if subcategoryStats.Blocked != 2 {
		t.Errorf("Expected blocked 2 in subcategory stats, got %d", subcategoryStats.Blocked)
	}
}
