package ui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mooship/blokilo/internal/models"
)

func TestNewAppModel(t *testing.T) {
	app := NewAppModel()

	if app.view != ViewMenu {
		t.Errorf("expected initial view to be ViewMenu, got %v", app.view)
	}

	if app.testRunning {
		t.Error("testRunning should be false initially")
	}

	if len(app.testResults) != 0 {
		t.Errorf("expected empty testResults, got %d results", len(app.testResults))
	}
}

func TestAppViewRendering(t *testing.T) {
	app := NewAppModel()

	app.view = ViewMenu
	view := app.View()
	if !strings.Contains(view, "Blokilo - Main Menu") {
		t.Error("menu view should contain menu title")
	}

	app.view = ViewSettings
	view = app.View()
	if !strings.Contains(view, "Blokilo - Settings") {
		t.Error("settings view should contain settings header")
	}

	app.view = ViewTest
	app.progress = NewProgressModel(10)
	app.progress.Current = 5
	app.progress.Domain = "test.com"
	view = app.View()
	if !strings.Contains(view, "Blokilo - Testing") {
		t.Error("test view should contain testing header")
	}
	if !strings.Contains(view, "[Esc/Q to Cancel]") {
		t.Error("test view should contain cancel instructions")
	}
}

func TestAppMenuSelection(t *testing.T) {
	app := NewAppModel()
	app.view = ViewMenu

	settingsMsg := MenuSelectedMsg{
		Item: MenuItem{ID: "settings", Label: "⚙️  Settings", Desc: "Configure DNS, timeout, and more."},
	}

	updatedModel, _ := app.Update(settingsMsg)
	updatedApp := updatedModel.(AppModel)

	if updatedApp.view != ViewSettings {
		t.Error("selecting Settings should change view to ViewSettings")
	}

	if !updatedApp.settings.Focus {
		t.Error("settings should be focused after selection")
	}
}

func TestAppSettingsFinished(t *testing.T) {
	app := NewAppModel()
	app.view = ViewSettings

	finishedMsg := settingsFinishedMsg{}

	updatedModel, _ := app.Update(finishedMsg)
	updatedApp := updatedModel.(AppModel)

	if updatedApp.view != ViewMenu {
		t.Error("finishing settings should return to menu")
	}
}

func TestAppKeyHandling(t *testing.T) {
	app := NewAppModel()

	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := app.Update(ctrlCMsg)

	if cmd == nil {
		t.Error("Ctrl+C should return a quit command")
	}

	app.view = ViewResults
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := app.Update(escMsg)
	updatedApp := updatedModel.(AppModel)

	if updatedApp.view != ViewMenu {
		t.Error("Esc in results view should return to menu")
	}

	app.view = ViewResults
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = app.Update(enterMsg)
	updatedApp = updatedModel.(AppModel)

	if updatedApp.view != ViewSummary {
		t.Error("Enter in results view should go to summary")
	}
}

func TestNewSummaryModel(t *testing.T) {
	results := []models.TestResult{
		{Domain: "blocked1.com", Status: models.StatusBlocked, ResponseTime: time.Millisecond * 10},
		{Domain: "blocked2.com", Status: models.StatusBlocked, ResponseTime: time.Millisecond * 15},
		{Domain: "resolved1.com", Status: models.StatusResolved, ResponseTime: time.Millisecond * 5},
		{Domain: "error1.com", Status: models.StatusError, ResponseTime: time.Millisecond * 0},
	}

	summary := NewSummaryModel(results)

	if summary.stats.Total != 4 {
		t.Errorf("expected Total to be 4, got %d", summary.stats.Total)
	}

	expectedBlocked := 50.0
	if summary.stats.PercentBlocked != expectedBlocked {
		t.Errorf("expected PercentBlocked to be %.1f, got %.1f", expectedBlocked, summary.stats.PercentBlocked)
	}

	if summary.recommendation == "" {
		t.Error("recommendation should not be empty")
	}
}

func TestNewResultsTableModel(t *testing.T) {
	results := []models.TestResult{
		{Domain: "test1.com", Status: models.StatusBlocked, ResponseTime: time.Millisecond * 10},
		{Domain: "test2.com", Status: models.StatusResolved, ResponseTime: time.Millisecond * 5},
	}

	tableModel := NewResultsTableModel(results)

	rows := tableModel.table.Rows()
	if len(rows) != 4 {
		t.Errorf("expected 4 rows (including headers), got %d", len(rows))
	}

	if !strings.Contains(rows[0][0], "Uncategorized") {
		t.Errorf("expected first row to be category header with 'Uncategorized', got %s", rows[0][0])
	}

	if !strings.Contains(rows[1][0], "Other") {
		t.Errorf("expected second row to be subcategory header with 'Other', got %s", rows[1][0])
	}

	if !strings.Contains(rows[2][0], "test1.com") {
		t.Errorf("expected third row to contain 'test1.com', got %s", rows[2][0])
	}

	if !strings.Contains(rows[2][2], "10ms") {
		t.Errorf("expected third row response time to contain '10ms', got %s", rows[2][2])
	}
}

func TestFormatHeader(t *testing.T) {
	header := formatHeader("Test Page")

	if !strings.Contains(header, "Blokilo - Test Page") {
		t.Error("formatHeader should include 'Blokilo -' prefix")
	}
}

func TestAppTestResultHandling(t *testing.T) {
	app := NewAppModel()
	app.view = ViewTest
	app.progress = NewProgressModel(2)
	app.testResults = make([]models.TestResult, 2)

	result := models.TestResult{
		Domain:       "example.com",
		Status:       models.StatusBlocked,
		ResponseTime: time.Millisecond * 10,
	}

	testResultMsg := testResultMsg{Result: result}
	updatedModel, cmd := app.Update(testResultMsg)
	updatedApp := updatedModel.(AppModel)

	if updatedApp.progress.Current != 1 {
		t.Errorf("expected progress.Current to be 1, got %d", updatedApp.progress.Current)
	}

	if updatedApp.progress.Domain != "example.com" {
		t.Errorf("expected progress.Domain to be 'example.com', got %s", updatedApp.progress.Domain)
	}

	if updatedApp.testResults[0].Domain != "example.com" {
		t.Errorf("expected first test result domain to be 'example.com', got %s", updatedApp.testResults[0].Domain)
	}

	if cmd == nil {
		t.Error("should return a command to continue listening for results")
	}
}

func TestAppAllTestsComplete(t *testing.T) {
	app := NewAppModel()
	app.view = ViewTest
	app.testRunning = true

	completeMsg := allTestsCompleteMsg{}
	updatedModel, _ := app.Update(completeMsg)
	updatedApp := updatedModel.(AppModel)

	if updatedApp.testRunning {
		t.Error("testRunning should be false after all tests complete")
	}

	if updatedApp.view != ViewResults {
		t.Error("view should change to ViewResults after all tests complete")
	}

	if updatedApp.testCancel != nil {
		t.Error("testCancel should be nil after all tests complete")
	}
}

func TestResultsTableBoundaryHandling(t *testing.T) {
	results := []models.TestResult{
		{Domain: "test1.com", Status: models.StatusBlocked, ResponseTime: time.Millisecond * 10},
		{Domain: "test2.com", Status: models.StatusResolved, ResponseTime: time.Millisecond * 5},
	}

	tableModel := NewResultsTableModel(results)

	emptyTableModel := NewResultsTableModel([]models.TestResult{})
	downKeyMsg := tea.KeyMsg{Type: tea.KeyDown}
	_, cmd := emptyTableModel.Update(downKeyMsg)
	if cmd != nil {
		t.Error("empty table should return nil command for navigation keys")
	}

	endKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("end")}
	updatedModel, _ := tableModel.Update(endKeyMsg)

	updatedModel, _ = updatedModel.Update(downKeyMsg)

	for range 5 {
		updatedModel, _ = updatedModel.Update(downKeyMsg)
	}

	homeKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("home")}
	updatedModel, _ = updatedModel.Update(homeKeyMsg)

	upKeyMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ = updatedModel.Update(upKeyMsg)

	for range 5 {
		updatedModel, _ = updatedModel.Update(upKeyMsg)
	}

	pgDownMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("pgdown")}
	updatedModel, _ = updatedModel.Update(endKeyMsg)
	updatedModel, _ = updatedModel.Update(pgDownMsg)

	pgUpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("pgup")}
	updatedModel, _ = updatedModel.Update(homeKeyMsg)
	updatedModel, _ = updatedModel.Update(pgUpMsg)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd = updatedModel.Update(enterMsg)
	if cmd != nil {
		t.Error("enter key should return nil command (handled by parent)")
	}

	spaceMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")}
	_, cmd = updatedModel.Update(spaceMsg)
	if cmd != nil {
		t.Error("space key should return nil command (handled by parent)")
	}
}
