package ui

import (
	"context"
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mooship/blokilo/internal/dns"
	"github.com/mooship/blokilo/internal/models"
)

var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("57")).
		Padding(0, 1).
		MarginBottom(2)
)

func formatHeader(title string) string {
	switch title {
	case "Testing":
		return headerStyle.Render("🧪 Blokilo - " + title)
	case "Test Results":
		return headerStyle.Render("📊 Blokilo - " + title)
	case "Summary":
		return headerStyle.Render("📋 Blokilo - " + title)
	case "Settings":
		return headerStyle.Render("⚙️ Blokilo - " + title)
	default:
		return headerStyle.Render("Blokilo - " + title)
	}
}

type AppView int

const (
	ViewMenu AppView = iota
	ViewTest
	ViewResults
	ViewSummary
	ViewSettings
)

type AppModel struct {
	menu         MenuModel
	progress     ProgressModel
	testResults  []models.TestResult
	resultsCh    chan models.TestResult
	testRunning  bool
	testCancel   context.CancelFunc
	testCtx      context.Context
	resultsTable ResultsTableModel
	summary      SummaryModel
	settings     SettingsModel
	view         AppView
}

func (m AppModel) View() string {
	switch m.view {
	case ViewMenu:
		return m.menu.View()
	case ViewTest:
		return formatHeader("Testing") + "\n" + m.progress.View() + "\n\n[Esc/Q to Cancel]"
	case ViewResults:
		recommendationText := fmt.Sprintf("%.0f%% blocked. %s", m.summary.stats.PercentBlocked, m.summary.recommendation)

		currentDNS := m.settings.GetSelectedDNS()
		if currentDNS == "" {
			currentDNS = dns.GetSystemDNS()
		}

		return formatHeader("Test Results") + "\n" + m.resultsTable.View() + "\n\n" + recommendationText + fmt.Sprintf("\n\n🔎 Tested with: %s", currentDNS) + "\n\n[⏎ Enter: Summary, Esc/Q: Menu]"
	case ViewSummary:
		return formatHeader("Summary") + "\n" + m.summary.View() + "\n\n[Esc/Q: Menu]"
	case ViewSettings:
		return formatHeader("Settings") + "\n" + m.settings.View()
	default:
		return ""
	}
}

type SummaryModel struct {
	stats          models.Stats
	recommendation string
}

func NewSummaryModel(results []models.TestResult) SummaryModel {
	classified := make([]models.ClassifiedResult, len(results))
	for i, r := range results {
		classified[i] = models.ClassifiedResult(r)
	}
	stats := models.ComputeStats(classified)
	rec := Recommend(stats)
	return SummaryModel{stats: stats, recommendation: rec}
}

func (m SummaryModel) View() string {
	return SummaryView(m.stats, m.recommendation)
}

func Recommend(stats models.Stats) string {
	if stats.PercentBlocked == 100 {
		return "All ad/tracker domains are blocked. Excellent job!"
	} else if stats.PercentBlocked > 80 {
		return "Most ad/tracker domains are blocked. Good job!"
	} else if stats.PercentBlocked > 0 {
		return "Partial blocking detected. Consider tightening your filters."
	}
	return "No blocking detected. Check your DNS or hosts file setup."
}

type ResultsTableModel struct {
	table table.Model
}

func NewResultsTableModel(results []models.TestResult) ResultsTableModel {
	rows := make([]TableRow, len(results))
	for i, r := range results {
		rows[i] = TableRow{
			Domain:       r.Domain,
			Status:       string(r.Status),
			ResponseTime: fmt.Sprintf("%.2fms", float64(r.ResponseTime.Microseconds())/1000),
		}
	}
	t := NewResultsTable(rows)
	return ResultsTableModel{table: t}
}

func (m ResultsTableModel) View() string {
	return m.table.View()
}

func (m ResultsTableModel) Update(msg tea.Msg) (ResultsTableModel, tea.Cmd) {
	tbl, cmd := m.table.Update(msg)
	m.table = tbl
	return m, cmd
}

func NewAppModel() AppModel {
	return AppModel{
		menu:         NewMenuModel(),
		progress:     ProgressModel{},
		testResults:  []models.TestResult{},
		resultsTable: NewResultsTableModel([]models.TestResult{}),
		summary:      NewSummaryModel([]models.TestResult{}),
		settings:     NewSettingsModel(""),
		view:         ViewMenu,
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

type testResultMsg struct {
	Result models.TestResult
}

type allTestsCompleteMsg struct{}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case settingsFinishedMsg:
		m.view = ViewMenu
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.testCancel != nil {
				m.testCancel()
			}
			return m, tea.Quit
		}
	}

	switch m.view {
	case ViewMenu:
		if sel, ok := msg.(MenuSelectedMsg); ok {
			switch sel.Item.ID {
			case "start":
				ctx := context.Background()
				domainList, err := models.LoadDomainList(ctx, "domains.jsonc")
				if err != nil {
					domainList = models.BuiltInDomains
				}
				m.progress = NewProgressModel(len(domainList))
				m.progress.DNSAddr = m.settings.GetSelectedDNS()
				m.testResults = make([]models.TestResult, len(domainList))
				m.testRunning = true
				testCtx, cancel := context.WithCancel(context.Background())
				m.testCancel = cancel
				m.testCtx = testCtx
				m.view = ViewTest
				m.resultsCh = make(chan models.TestResult, len(domainList))
				return m, tea.Batch(
					runParallelTests(testCtx, domainList, m.resultsCh, m.progress.DNSAddr),
					listenForTestResults(m.resultsCh),
				)
			case "settings":
				m.view = ViewSettings
				m.settings = NewSettingsModel(m.settings.GetSelectedDNS())
				return m, nil
			case "exit":
				return m, tea.Quit
			}
		}
		menuModel, cmd := m.menu.Update(msg)
		m.menu = menuModel.(MenuModel)
		return m, cmd

	case ViewTest:
		switch msg := msg.(type) {
		case testResultMsg:
			if m.progress.Current < len(m.testResults) {
				m.testResults[m.progress.Current] = msg.Result
				m.progress.Current++
				m.progress.Domain = msg.Result.Domain
			}
			return m, listenForTestResults(m.resultsCh)
		case allTestsCompleteMsg:
			m.testRunning = false
			m.testCancel = nil
			m.testCtx = nil
			m.resultsTable = NewResultsTableModel(m.testResults)
			m.summary = NewSummaryModel(m.testResults)
			m.view = ViewResults
			return m, nil
		case tea.KeyMsg:
			if msg.String() == "esc" || msg.String() == "q" {
				if m.testCancel != nil {
					m.testCancel()
				}
				m.testRunning = false
				m.view = ViewMenu
				return m, nil
			}
		}
		progModel, cmd := m.progress.Update(msg)
		m.progress = progModel.(ProgressModel)
		return m, cmd

	case ViewResults:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter", " ":
				m.view = ViewSummary
				return m, nil
			case "esc", "q":
				m.view = ViewMenu
				return m, nil
			}
		}
		tbl, cmd := m.resultsTable.Update(msg)
		m.resultsTable = tbl
		return m, cmd

	case ViewSummary:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "esc", "q":
				m.view = ViewMenu
				return m, nil
			}
		}
		return m, nil

	case ViewSettings:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "esc" || key.String() == "q" {
				m.view = ViewMenu
				return m, nil
			}
		}
		settingsModel, cmd := m.settings.Update(msg)
		m.settings = settingsModel.(SettingsModel)
		return m, cmd

	default:
		return m, nil
	}
}

func listenForTestResults(resultsCh chan models.TestResult) tea.Cmd {
	return func() tea.Msg {
		result, ok := <-resultsCh
		if !ok {
			return allTestsCompleteMsg{}
		}
		return testResultMsg{Result: result}
	}
}

func runParallelTests(ctx context.Context, domainList []models.DomainEntry, resultsCh chan models.TestResult, dnsAddr string) tea.Cmd {
	return func() tea.Msg {
		var wg sync.WaitGroup
		const numWorkers = 5
		jobs := make(chan models.DomainEntry, len(domainList))

		dnsServer := ""
		if dnsAddr != "" {
			dnsServer = dnsAddr
		}

		for range numWorkers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for domain := range jobs {
					result := dns.TestDomainDNS(ctx, domain.Name, dnsServer)

					select {
					case <-ctx.Done():
						return
					default:
						select {
						case resultsCh <- result:
						case <-ctx.Done():
						}
					}
				}
			}()
		}

		go func() {
			defer close(jobs)
			for _, domain := range domainList {
				select {
				case jobs <- domain:
				case <-ctx.Done():
					return
				}
			}
		}()

		go func() {
			wg.Wait()
			close(resultsCh)
		}()

		return nil
	}
}
