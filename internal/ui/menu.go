package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuItem struct {
	Label string
	Desc  string
}

type MenuSelectedMsg struct {
	Item MenuItem
}

func (m MenuItem) Title() string       { return m.Label }
func (m MenuItem) Description() string { return m.Desc }
func (m MenuItem) FilterValue() string { return m.Label }

var menuItems = []list.Item{
	MenuItem{"Start Test", "Run ad-block test on domains."},
	MenuItem{"Settings", "Configure DNS, timeout, and more."},
	MenuItem{"Exit", "Quit Blokilo."},
}

type MenuModel struct {
	List list.Model
}

func NewMenuModel() MenuModel {
	items := menuItems
	l := list.New(items, list.NewDefaultDelegate(), 40, 5)
	l.Title = "Blokilo - Main Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Select(0)
	return MenuModel{List: l}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		minW, minH := 30, 8
		w, h := msg.Width, msg.Height
		if w < minW {
			w = minW
		}
		if h < minH {
			h = minH
		}
		listHeight := h - 4
		if listHeight < 3 {
			listHeight = 3
		}
		m.List.SetSize(w, listHeight)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			return m, m.selectItem()
		case "q", "esc":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	return m.List.View()
}

func (m MenuModel) selectItem() tea.Cmd {
	selected := m.List.SelectedItem()
	if item, ok := selected.(MenuItem); ok {
		return func() tea.Msg {
			return MenuSelectedMsg{Item: item}
		}
	}
	return nil
}
