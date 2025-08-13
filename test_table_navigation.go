package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mooship/blokilo/internal/models"
	"github.com/mooship/blokilo/internal/ui"
)

type testModel struct {
	table ui.ResultsTableModel
}

func (m testModel) Init() tea.Cmd {
	return nil
}

func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m testModel) View() string {
	return fmt.Sprintf("Table Navigation Test\n\n%s\n\nPress arrow keys to navigate, 'q' to quit", m.table.View())
}

func main() {
	results := []models.TestResult{
		{Domain: "example1.com", Status: models.StatusBlocked, ResponseTime: time.Millisecond * 10},
		{Domain: "example2.com", Status: models.StatusResolved, ResponseTime: time.Millisecond * 5},
		{Domain: "example3.com", Status: models.StatusError, ResponseTime: time.Millisecond * 15},
	}

	tableModel := ui.NewResultsTableModel(results)

	m := testModel{
		table: tableModel,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
