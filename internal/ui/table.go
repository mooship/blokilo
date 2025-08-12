package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/mooship/blokilo/internal/models"
)

var (
	colorBlocked  = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	colorResolved = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	colorError    = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
)

type TableRow struct {
	Domain       string
	Status       string
	ResponseTime string
}

func NewResultsTable(rows []TableRow) table.Model {
	columns := []table.Column{
		{Title: "Domain", Width: 44},
		{Title: "Status", Width: 18},
		{Title: "Response Time", Width: 16},
	}
	tRows := make([]table.Row, len(rows))
	for i, r := range rows {
		status := r.Status
		switch status {
		case string(models.StatusResolved):
			status = colorResolved.Render(status)
		case string(models.StatusBlocked):
			status = colorBlocked.Render(status)
		case string(models.StatusError):
			status = colorError.Render(status)
		}
		tRows[i] = table.Row{r.Domain, status, r.ResponseTime}
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(tRows),
		table.WithFocused(true),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("15"))
	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	styles.Cell = styles.Cell.PaddingLeft(1).PaddingRight(1)
	t.SetStyles(styles)

	return t
}

func TableView(t table.Model) string {
	return t.View()
}
