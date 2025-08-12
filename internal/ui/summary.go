package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mooship/blokilo/internal/models"
)

var (
	summaryTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render
	summaryStat  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Render
)

func SummaryView(stats models.Stats, recommendation string) string {
	var builder strings.Builder

	errorPercent := 100.0 - stats.PercentBlocked - stats.PercentResolved

	builder.WriteString(summaryTitle("📋 Summary"))
	builder.WriteString(fmt.Sprintf("\n🚫 Blocked: %s%%", summaryStat(fmt.Sprintf("%.1f", stats.PercentBlocked))))
	builder.WriteString(fmt.Sprintf("\n✅ Resolved: %s%%", summaryStat(fmt.Sprintf("%.1f", stats.PercentResolved))))

	if errorPercent > 0.1 {
		builder.WriteString(fmt.Sprintf("\n⚠️ Errors: %s%%", summaryStat(fmt.Sprintf("%.1f", errorPercent))))
	}

	builder.WriteString("\n\n")
	builder.WriteString(recommendation)

	return builder.String()
}
