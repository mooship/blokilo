package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewSettingsModel(t *testing.T) {
	currentDNS := "1.1.1.1:53"
	settings := NewSettingsModel(currentDNS)

	if settings.Value != currentDNS {
		t.Errorf("expected Value to be %s, got %s", currentDNS, settings.Value)
	}

	if !settings.Focus {
		t.Error("settings should be focused by default")
	}

	if settings.DNSInput.Value() != currentDNS {
		t.Errorf("expected input value to be %s, got %s", currentDNS, settings.DNSInput.Value())
	}
}

func TestGetSelectedDNS(t *testing.T) {
	settings := NewSettingsModel("8.8.8.8:53")
	settings.Value = "9.9.9.9:53"

	if settings.GetSelectedDNS() != "9.9.9.9:53" {
		t.Errorf("expected GetSelectedDNS to return %s, got %s", "9.9.9.9:53", settings.GetSelectedDNS())
	}
}

func TestSettingsUpdateEnterKey(t *testing.T) {
	settings := NewSettingsModel("")

	settings.DNSInput.SetValue("1.1.1.1")
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, cmd := settings.Update(keyMsg)
	updatedSettings := updatedModel.(SettingsModel)

	if updatedSettings.Value != "1.1.1.1:53" {
		t.Errorf("expected Value to be '1.1.1.1:53', got %s", updatedSettings.Value)
	}

	if updatedSettings.Err != nil {
		t.Errorf("expected no error, got %v", updatedSettings.Err)
	}

	if updatedSettings.Focus {
		t.Error("focus should be false after successful input")
	}

	if cmd == nil {
		t.Error("should return a command after successful input")
	}
}

func TestSettingsUpdateEnterKeyWithPort(t *testing.T) {
	settings := NewSettingsModel("")

	settings.DNSInput.SetValue("8.8.8.8:53")
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, _ := settings.Update(keyMsg)
	updatedSettings := updatedModel.(SettingsModel)

	if updatedSettings.Value != "8.8.8.8:53" {
		t.Errorf("expected Value to be '8.8.8.8:53', got %s", updatedSettings.Value)
	}

	if updatedSettings.Err != nil {
		t.Errorf("expected no error, got %v", updatedSettings.Err)
	}
}

func TestSettingsUpdateInvalidInput(t *testing.T) {
	settings := NewSettingsModel("")

	settings.DNSInput.SetValue("invalid-dns")
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, cmd := settings.Update(keyMsg)
	updatedSettings := updatedModel.(SettingsModel)

	if updatedSettings.Err == nil {
		t.Error("expected an error for invalid input")
	}

	if !updatedSettings.Focus {
		t.Error("focus should remain true after invalid input")
	}

	if cmd != nil {
		t.Error("should not return a command after invalid input")
	}
}

func TestSettingsUpdateEmptyInput(t *testing.T) {
	settings := NewSettingsModel("1.1.1.1:53")

	settings.DNSInput.SetValue("")
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, _ := settings.Update(keyMsg)
	updatedSettings := updatedModel.(SettingsModel)

	if updatedSettings.Value != "" {
		t.Errorf("expected Value to be empty, got %s", updatedSettings.Value)
	}

	if updatedSettings.Err != nil {
		t.Errorf("expected no error for empty input, got %v", updatedSettings.Err)
	}

	if updatedSettings.Focus {
		t.Error("focus should be false after empty input")
	}
}

func TestSettingsUpdateEscapeKey(t *testing.T) {
	settings := NewSettingsModel("")
	settings.Err = &testError{"some error"}

	keyMsg := tea.KeyMsg{Type: tea.KeyEsc}

	updatedModel, cmd := settings.Update(keyMsg)
	updatedSettings := updatedModel.(SettingsModel)

	if updatedSettings.Err != nil {
		t.Error("error should be cleared after escape")
	}

	if cmd == nil {
		t.Error("should return a command after escape")
	}
}

func TestSettingsView(t *testing.T) {
	settings := NewSettingsModel("1.1.1.1:53")

	view := settings.View()

	expectedElements := []string{
		"🔧 Custom DNS Server",
		"💡 Examples:",
		"⌨️  [⏎ Enter: Save, Esc/Q: Cancel]",
	}

	for _, element := range expectedElements {
		if !containsString(view, element) {
			t.Errorf("view should contain '%s'", element)
		}
	}
}

func TestSettingsViewWithError(t *testing.T) {
	settings := NewSettingsModel("")
	settings.Err = &testError{"invalid format"}

	view := settings.View()

	if !containsString(view, "❌ Error: invalid format") {
		t.Error("view should contain error message")
	}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

func containsString(text, substr string) bool {
	return len(text) >= len(substr) && indexOf(text, substr) >= 0
}

func indexOf(text, substr string) int {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
