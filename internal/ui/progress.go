package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mooship/blokilo/internal/dns"
)

type ProgressModel struct {
	Progress progress.Model
	Current  int
	Total    int
	Domain   string
	DNSAddr  string
}

func NewProgressModel(total int) ProgressModel {
	return ProgressModel{
		Progress: progress.New(progress.WithDefaultGradient()),
		Total:    total,
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return nil
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	updated, cmd := m.Progress.Update(msg)
	if p, ok := updated.(progress.Model); ok {
		m.Progress = p
	}
	return m, cmd
}

func (m ProgressModel) View() string {
	bar := m.Progress.ViewAs(float64(m.Current) / float64(m.Total))

	currentDNS := m.DNSAddr
	if currentDNS == "" {
		systemDNS, err := dns.GetSystemDNS()
		if err != nil {
			currentDNS = "System"
		} else {
			currentDNS = systemDNS
		}
	}

	return fmt.Sprintf("Testing: %s\n\n%s\n\n%d/%d\n\nDNS: %s",
		m.Domain, bar, m.Current, m.Total, currentDNS)
}
