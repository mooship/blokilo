package ui

import (
	"fmt"
	"net"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type settingsFinishedMsg struct{}

type SettingsModel struct {
	DNSInput textinput.Model
	Focus    bool
	Value    string
	Err      error
}

func NewSettingsModel(currentDNS string) SettingsModel {
	ti := textinput.New()
	ti.Placeholder = "1.1.1.1:53 or empty for system DNS"
	ti.Width = 50
	ti.SetValue(currentDNS)
	ti.Focus()
	return SettingsModel{
		DNSInput: ti,
		Focus:    true,
		Value:    currentDNS,
	}
}

func (m SettingsModel) GetSelectedDNS() string {
	return m.Value
}

func (m SettingsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			val := m.DNSInput.Value()
			if val == "" {
				m.Err = nil
				m.Value = ""
				m.Focus = false
				return m, nil
			}

			host, port, err := net.SplitHostPort(val)
			if err != nil {
				if net.ParseIP(val) != nil {
					m.Err = nil
					m.Value = net.JoinHostPort(val, "53")
				} else {
					m.Err = fmt.Errorf("invalid format, must be host:port or IP")
				}
			} else {
				m.Err = nil
				m.Value = net.JoinHostPort(host, port)
			}

			if m.Err == nil {
				m.Focus = false
				return m, func() tea.Msg { return settingsFinishedMsg{} }
			}
			return m, nil
		case "esc":
			m.Err = nil
			return m, func() tea.Msg { return settingsFinishedMsg{} }
		}
	}

	if m.Focus {
		m.DNSInput, cmd = m.DNSInput.Update(msg)
	}

	return m, cmd
}

func (m SettingsModel) View() string {
	errMsg := ""
	if m.Err != nil {
		errMsg = "\n\n❌ Error: " + m.Err.Error()
	}

	helpText := "\n\n💡 Examples: 1.1.1.1, 8.8.8.8:53, 9.9.9.9:853"
	controls := "\n\n⌨️  [Enter] Save  [Esc] Cancel"

	return fmt.Sprintf("🔧 Custom DNS Server (host:port):\n%s%s%s%s",
		m.DNSInput.View(), helpText, controls, errMsg)
}
