package ui

import (
	"strings"
	"testing"

	"github.com/mooship/blokilo/internal/models"
)

func TestSummaryView(t *testing.T) {
	stats := models.Stats{
		Total:           10,
		Blocked:         7,
		Resolved:        2,
		Errored:         1,
		PercentBlocked:  70.0,
		PercentResolved: 20.0,
	}

	recommendation := "Most ad/tracker domains are blocked. Good job!"
	view := SummaryView(stats, recommendation)

	if !strings.Contains(view, "Blocked: 70.0%") {
		t.Error("summary view should contain blocked percentage")
	}

	if !strings.Contains(view, "Resolved: 20.0%") {
		t.Error("summary view should contain resolved percentage")
	}

	if !strings.Contains(view, recommendation) {
		t.Error("summary view should contain recommendation")
	}
}

func TestSummaryViewWithErrors(t *testing.T) {
	stats := models.Stats{
		Total:           10,
		Blocked:         6,
		Resolved:        2,
		Errored:         2,
		PercentBlocked:  60.0,
		PercentResolved: 20.0,
	}

	recommendation := "Some errors occurred"
	view := SummaryView(stats, recommendation)

	if !strings.Contains(view, "Errors: 20.0%") {
		t.Error("summary view should contain error percentage when errors > 0.1%")
	}
}

func TestSummaryViewNoErrors(t *testing.T) {
	stats := models.Stats{
		Total:           10,
		Blocked:         8,
		Resolved:        2,
		Errored:         0,
		PercentBlocked:  80.0,
		PercentResolved: 20.0,
	}

	recommendation := "Great blocking!"
	view := SummaryView(stats, recommendation)

	if strings.Contains(view, "Errors:") {
		t.Error("summary view should not contain error percentage when errors = 0%")
	}
}

func TestRecommend(t *testing.T) {
	testCases := []struct {
		name           string
		stats          models.Stats
		expectedPhrase string
	}{
		{
			name:           "Perfect blocking",
			stats:          models.Stats{PercentBlocked: 100.0, PercentResolved: 0.0},
			expectedPhrase: "All ad/tracker domains are blocked",
		},
		{
			name:           "Good blocking",
			stats:          models.Stats{PercentBlocked: 85.0, PercentResolved: 15.0},
			expectedPhrase: "Most ad/tracker domains are blocked",
		},
		{
			name:           "Partial blocking",
			stats:          models.Stats{PercentBlocked: 50.0, PercentResolved: 50.0},
			expectedPhrase: "Partial blocking detected",
		},
		{
			name:           "No blocking",
			stats:          models.Stats{PercentBlocked: 0.0, PercentResolved: 100.0},
			expectedPhrase: "No blocking detected",
		},
	}

	for _, tc := range testCases {
		recommendation := Recommend(tc.stats)
		if !strings.Contains(recommendation, tc.expectedPhrase) {
			t.Errorf("expected recommendation to contain %q, got %q", tc.expectedPhrase, recommendation)
		}
	}
}

func TestSummaryView_ZeroTotal(t *testing.T) {
	stats := models.Stats{}
	recommendation := "No data"
	view := SummaryView(stats, recommendation)
	if !strings.Contains(view, "0/0") && !strings.Contains(view, "0%") {
		t.Error("summary view should handle zero total domains")
	}
}

func TestSummaryView_AllErrors(t *testing.T) {
	stats := models.Stats{
		Total:           5,
		Blocked:         0,
		Resolved:        0,
		Errored:         5,
		PercentBlocked:  0.0,
		PercentResolved: 0.0,
	}
	recommendation := "All errors"
	view := SummaryView(stats, recommendation)
	if !strings.Contains(view, "Errors: 100.0%") {
		t.Error("summary view should show 100% errors when all errored")
	}
}
