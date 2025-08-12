package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewMenuModel(t *testing.T) {
	menu := NewMenuModel()

	if menu.List.Title != "Blokilo - Main Menu" {
		t.Error("menu should have correct title")
	}

	items := menu.List.Items()
	if len(items) != 3 {
		t.Errorf("expected 3 menu items, got %d", len(items))
	}

	expectedLabels := []string{"Start Test", "Settings", "Exit"}
	for i, item := range items {
		if menuItem, ok := item.(MenuItem); ok {
			if menuItem.Label != expectedLabels[i] {
				t.Errorf("menu item %d: expected %s, got %s", i, expectedLabels[i], menuItem.Label)
			}
		} else {
			t.Errorf("menu item %d is not of type MenuItem", i)
		}
	}
}

func TestMenuItemInterface(t *testing.T) {
	item := MenuItem{
		Label: "Test Item",
		Desc:  "Test Description",
	}

	if item.Title() != "Test Item" {
		t.Errorf("expected Title() to return 'Test Item', got %s", item.Title())
	}

	if item.Description() != "Test Description" {
		t.Errorf("expected Description() to return 'Test Description', got %s", item.Description())
	}

	if item.FilterValue() != "Test Item" {
		t.Errorf("expected FilterValue() to return 'Test Item', got %s", item.FilterValue())
	}
}

func TestMenuUpdate(t *testing.T) {
	menu := NewMenuModel()

	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := menu.Update(sizeMsg)
	updatedMenu := updatedModel.(MenuModel)

	if updatedMenu.List.Title != "Blokilo - Main Menu" {
		t.Error("menu title should be preserved after window size update")
	}

	smallSizeMsg := tea.WindowSizeMsg{Width: 10, Height: 5}
	updatedModel, _ = menu.Update(smallSizeMsg)
	updatedMenu = updatedModel.(MenuModel)

	if updatedMenu.List.Title != "Blokilo - Main Menu" {
		t.Error("menu title should be preserved after small window size update")
	}
}

func TestMenuSelectItem(t *testing.T) {
	menu := NewMenuModel()

	cmd := menu.selectItem()
	if cmd == nil {
		t.Error("selectItem should return a command")
	}

	msg := cmd()
	if selectedMsg, ok := msg.(MenuSelectedMsg); ok {
		if selectedMsg.Item.Label != "Start Test" {
			t.Errorf("expected first item to be 'Start Test', got %s", selectedMsg.Item.Label)
		}
	} else {
		t.Error("selectItem command should return MenuSelectedMsg")
	}
}
