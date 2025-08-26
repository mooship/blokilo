package ui

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mooship/blokilo/internal/dns"
	"github.com/mooship/blokilo/internal/httpclient"
	"github.com/mooship/blokilo/internal/models"
	"github.com/samber/lo"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("57")).
			Padding(0, 1).
			MarginBottom(2)

	testCancelHint = "[Esc/Q to Cancel]"
	resultsHint    = "[⏎ Enter: Summary, Esc/Q: Menu]"
	summaryHint    = "[Esc/Q: Menu]"

	headerTesting     = "Testing"
	headerTestResults = "Test Results"
	headerSummary     = "Summary"
	headerSettings    = "Settings"
	headerDefault     = "Blokilo"

	systemDNSErrorFmt = "System DNS (Error: %v)"
	testedWithFmt     = "\n\n🔎 Tested with: %s"

	msgAllBlocked   = "All ad/tracker domains are blocked. Excellent job!"
	msgMostBlocked  = "Most ad/tracker domains are blocked. Good job!"
	msgPartialBlock = "Partial blocking detected. Consider tightening your filters."
	msgNoBlocking   = "No blocking detected. Check your DNS or hosts file setup."
)

func formatHeader(title string) string {
	switch title {
	case headerTesting:
		return headerStyle.Render("🧪 Blokilo - " + title)
	case headerTestResults:
		return headerStyle.Render("📊 Blokilo - " + title)
	case headerSummary:
		return headerStyle.Render("📋 Blokilo - " + title)
	case headerSettings:
		return headerStyle.Render("⚙️ Blokilo - " + title)
	default:
		return headerStyle.Render(headerDefault + " - " + title)
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
	menu           MenuModel
	progress       ProgressModel
	testResults    []models.TestResult
	resultsCh      chan models.TestResult
	testRunning    bool
	testCancel     context.CancelFunc
	testCtx        context.Context
	resultsTable   ResultsTableModel
	summary        SummaryModel
	settings       SettingsModel
	view           AppView
	categoryConfig *models.CategoryConfig
	systemDNS      string
	dnsErr         error
}

func (m AppModel) View() string {
	switch m.view {
	case ViewMenu:
		return m.menu.View()
	case ViewTest:
		return formatHeader("Testing") + "\n" + m.progress.View() + "\n\n" + testCancelHint
	case ViewResults:
		recommendationText := fmt.Sprintf("%.0f%% blocked. %s", m.summary.stats.PercentBlocked, m.summary.recommendation)

		currentDNS := m.settings.GetSelectedDNS()
		if currentDNS == "" {
			if m.dnsErr != nil {
				currentDNS = fmt.Sprintf(systemDNSErrorFmt, m.dnsErr)
			} else {
				currentDNS = m.systemDNS
			}
		}
		return formatHeader(headerTestResults) + "\n" + m.resultsTable.View() + "\n\n" + recommendationText + fmt.Sprintf(testedWithFmt, currentDNS) + "\n\n" + resultsHint
	case ViewSummary:
		return formatHeader(headerSummary) + "\n" + m.summary.View() + "\n\n" + summaryHint
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
	classified := lo.Map(results, func(r models.TestResult, _ int) models.ClassifiedResult {
		return models.ClassifiedResult(r)
	})
	stats := models.ComputeStats(classified)
	rec := Recommend(stats)
	return SummaryModel{stats: stats, recommendation: rec}
}

func (m SummaryModel) View() string {
	return SummaryView(m.stats, m.recommendation)
}

func Recommend(stats models.Stats) string {
	if stats.PercentBlocked == 100 {
		return msgAllBlocked
	} else if stats.PercentBlocked > 80 {
		return msgMostBlocked
	} else if stats.PercentBlocked > 0 {
		return msgPartialBlock
	}
	return msgNoBlocking
}

type ResultsTableModel struct {
	table table.Model
}

func NewResultsTableModel(results []models.TestResult, config *models.CategoryConfig) ResultsTableModel {
	classified := lo.Map(results, func(r models.TestResult, _ int) models.ClassifiedResult {
		return models.ClassifiedResult(r)
	})

	groups := models.GroupResultsByCategory(classified, config)

	t := NewGroupedResultsTable(groups)
	return ResultsTableModel{table: t}
}

func (m ResultsTableModel) View() string {
	return m.table.View()
}

func (m ResultsTableModel) Update(msg tea.Msg) (ResultsTableModel, tea.Cmd) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		if len(m.table.Rows()) > 0 {
			tbl, cmd := m.table.Update(sizeMsg)
			m.table = tbl
			return m, cmd
		}
		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		currentCursor := m.table.Cursor()
		rowCount := len(m.table.Rows())

		if rowCount == 0 {
			return m, nil
		}

		switch keyMsg.String() {
		case "down", "j":
			if currentCursor >= rowCount-1 {
				return m, nil
			}
		case "up", "k":
			if currentCursor <= 0 {
				return m, nil
			}
		case "pgdown", "ctrl+d":
			if currentCursor >= rowCount-1 {
				if currentCursor != rowCount-1 {
					m.table.SetCursor(rowCount - 1)
				}
				return m, nil
			}
			newCursor := currentCursor + 10
			if newCursor >= rowCount {
				m.table.SetCursor(rowCount - 1)
				return m, nil
			}
		case "pgup", "ctrl+u":
			if currentCursor <= 0 {
				return m, nil
			}
			newCursor := currentCursor - 10
			if newCursor < 0 {
				m.table.SetCursor(0)
				return m, nil
			}
		case "end", "G":
			m.table.SetCursor(rowCount - 1)
			return m, nil
		case "home", "g":
			m.table.SetCursor(0)
			return m, nil
		case "enter", " ":
			return m, nil
		}
	}

	if len(m.table.Rows()) == 0 {
		return m, nil
	}

	tbl, cmd := m.table.Update(msg)
	m.table = tbl

	if len(m.table.Rows()) > 0 {
		if m.table.Cursor() >= len(m.table.Rows()) {
			m.table.SetCursor(len(m.table.Rows()) - 1)
		}
		if m.table.Cursor() < 0 {
			m.table.SetCursor(0)
		}
	}

	return m, cmd
}

func NewAppModel() AppModel {
	config, err := models.LoadCategoryConfig("data/categories.jsonc")
	if err != nil {
		config = &models.CategoryConfig{}
	}

	systemDNS, dnsErr := dns.GetSystemDNS()

	return AppModel{
		menu:           NewMenuModel(),
		progress:       ProgressModel{},
		testResults:    []models.TestResult{},
		resultsTable:   NewResultsTableModel([]models.TestResult{}, config),
		summary:        NewSummaryModel([]models.TestResult{}),
		settings:       NewSettingsModel(""),
		view:           ViewMenu,
		categoryConfig: config,
		systemDNS:      systemDNS,
		dnsErr:         dnsErr,
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
				domainList, err := models.LoadDomainList(ctx, "data/domains.jsonc")
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
			m.resultsTable = NewResultsTableModel(m.testResults, m.categoryConfig)
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
			wg.Go(func() {
				for domainEntry := range jobs {
					dnsResult := dns.TestDomainDNS(ctx, domainEntry.Name, dnsServer)
					httpResult := httpclient.CheckHTTPConnectivity(ctx, domainEntry.Name, 5*time.Second, 2)

					merged := dnsResult
					merged.HTTPStatusCode = httpResult.HTTPStatusCode
					merged.Category = domainEntry.Category
					merged.Subcategory = domainEntry.Subcategory

					if dnsResult.Err != nil {
						merged.Err = dnsResult.Err
					} else if httpResult.Err != nil {
						merged.Err = httpResult.Err
					} else {
						merged.Err = nil
					}

					select {
					case <-ctx.Done():
						return
					default:
						select {
						case resultsCh <- merged:
						case <-ctx.Done():
						}
					}
				}
			})
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
