package ui

import (
	"strings"
	"testing"

	"github.com/mooship/blokilo/internal/models"
)

func TestNewResultsTable(t *testing.T) {
	rows := []TableRow{
		{Domain: "example.com", Status: string(models.StatusBlocked), ResponseTime: "12.34ms"},
		{Domain: "google.com", Status: string(models.StatusResolved), ResponseTime: "5.67ms"},
		{Domain: "error.com", Status: string(models.StatusError), ResponseTime: "0.00ms"},
	}

	table := NewResultsTable(rows)

	if len(table.Rows()) != 3 {
		t.Errorf("expected 3 rows, got %d", len(table.Rows()))
	}

	columns := table.Columns()
	if len(columns) != 3 {
		t.Errorf("expected 3 columns, got %d", len(columns))
	}

	expectedColumns := []string{"Domain", "Status", "Response Time"}
	for i, col := range columns {
		if col.Title != expectedColumns[i] {
			t.Errorf("column %d: expected %s, got %s", i, expectedColumns[i], col.Title)
		}
	}

	if !table.Focused() {
		t.Error("expected table to be focused")
	}
}

func TestTableView(t *testing.T) {
	rows := []TableRow{
		{Domain: "test.com", Status: string(models.StatusBlocked), ResponseTime: "10.00ms"},
	}

	table := NewResultsTable(rows)
	view := TableView(table)

	if !strings.Contains(view, "test.com") {
		t.Error("table view should contain domain name")
	}

	if !strings.Contains(view, "Domain") || !strings.Contains(view, "Status") || !strings.Contains(view, "Response Time") {
		t.Error("table view should contain column headers")
	}
}

func TestTableRowStatusColors(t *testing.T) {
	testCases := []struct {
		status   string
		expected string
	}{
		{string(models.StatusBlocked), string(models.StatusBlocked)},
		{string(models.StatusResolved), string(models.StatusResolved)},
		{string(models.StatusError), string(models.StatusError)},
	}

	for _, tc := range testCases {
		rows := []TableRow{{Domain: "test.com", Status: tc.status, ResponseTime: "10ms"}}
		table := NewResultsTable(rows)

		tableRows := table.Rows()
		if len(tableRows) != 1 {
			t.Errorf("expected 1 row, got %d", len(tableRows))
			continue
		}

		statusColumn := tableRows[0][1]
		if !strings.Contains(statusColumn, tc.expected) {
			t.Errorf("status column should contain %s, got %s", tc.expected, statusColumn)
		}
	}
}
