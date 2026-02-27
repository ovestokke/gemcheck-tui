package app

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ovestokke/gemcheck-tui/internal/api"
	"github.com/ovestokke/gemcheck-tui/internal/cache"
	"github.com/ovestokke/gemcheck-tui/internal/domain"
	"github.com/ovestokke/gemcheck-tui/internal/tui"
	"github.com/ovestokke/gemcheck-tui/internal/tui/components"
)

// Cache TTLs
const (
	leagueTTL = 1 * time.Hour
	wikiTTL   = 24 * time.Hour
	priceTTL  = 5 * time.Minute
)

// Cache keys
const (
	keyLeagues = "leagues"
	keyWiki    = "wiki"
	keyPrices  = "prices"
)

type screenState int

const (
	screenLoading screenState = iota
	screenLeagueSelect
	screenMain
)

// Model is the top-level Bubble Tea model.
type Model struct {
	cache  *cache.Cache
	screen screenState
	width  int
	height int
	err    error

	// Sub-models
	spinner      components.SpinnerModel
	leagueSelect components.LeagueSelectModel
	tabs         components.GemTabsModel
	table        components.GemTableModel
	statusbar    components.StatusBarModel
	search       components.SearchModel
	detail       components.DetailModel

	// Data
	league     domain.League
	wiki       *domain.WikiData
	prices     []domain.GemPrice
	result     *domain.ProcessedResult
	wikiReady  bool
	priceReady bool
}

// NewModel creates the application model.
func NewModel(c *cache.Cache) Model {
	return Model{
		cache:     c,
		screen:    screenLoading,
		spinner:   components.NewSpinner("Fetching leagues..."),
		tabs:      components.NewGemTabs(),
		table:     components.NewGemTable(80, 20),
		statusbar: components.NewStatusBar(),
		search:    components.NewSearch(),
		detail:    components.NewDetail(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Init(),
		fetchLeaguesCmd(m.cache),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetSize(msg.Width, msg.Height-4) // tabs + statusbar
		m.statusbar.SetWidth(msg.Width)
		m.search.SetSize(msg.Width, msg.Height)
		m.detail.SetSize(msg.Width, msg.Height)
		if m.screen == screenLeagueSelect {
			m.leagueSelect, _ = m.leagueSelect.Update(msg)
		}
		return m, nil

	case tui.LeaguesFetchedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.leagueSelect = components.NewLeagueSelect(msg.Leagues, m.width, m.height)
		m.screen = screenLeagueSelect
		return m, nil

	case tui.LeagueSelectedMsg:
		m.league = msg.League
		m.screen = screenLoading
		m.spinner = components.NewSpinner("Loading gem data...")
		m.wikiReady = false
		m.priceReady = false
		m.statusbar.SetLeague(m.league.Text)
		return m, tea.Batch(
			m.spinner.Init(),
			fetchWikiCmd(m.cache),
			fetchPricesCmd(m.cache, m.league.ID),
		)

	case tui.WikiFetchedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.wiki = msg.Wiki
		m.wikiReady = true
		return m, m.tryProcessGems()

	case tui.PricesFetchedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.prices = msg.Prices
		m.priceReady = true
		return m, m.tryProcessGems()

	case tui.DataReadyMsg:
		m.result = &msg.Result
		m.search.SetGems(msg.Result.GemPicks)
		m.populateTable()
		m.screen = screenMain
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Forward to active sub-model
	var cmd tea.Cmd
	switch m.screen {
	case screenLoading:
		m.spinner, cmd = m.spinner.Update(msg)
	case screenLeagueSelect:
		m.leagueSelect, cmd = m.leagueSelect.Update(msg)
	case screenMain:
		if m.search.Active() {
			m.search, cmd = m.search.Update(msg)
		} else if m.detail.Active() {
			m.detail, cmd = m.detail.Update(msg)
		} else {
			m.table, cmd = m.table.Update(msg)
		}
	}
	return m, cmd
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit
	if key.Matches(msg, tui.Keys.Quit) && !m.search.Active() {
		return m, tea.Quit
	}

	if m.screen == screenLeagueSelect {
		var cmd tea.Cmd
		m.leagueSelect, cmd = m.leagueSelect.Update(msg)
		return m, cmd
	}

	if m.screen != screenMain {
		return m, nil
	}

	// Search overlay takes priority
	if m.search.Active() {
		if msg.String() == "enter" {
			if entry := m.search.SelectedEntry(); entry != nil {
				m.search.Close()
				m.detail.Show(entry)
			}
			return m, nil
		}
		var cmd tea.Cmd
		m.search, cmd = m.search.Update(msg)
		return m, cmd
	}

	// Detail overlay
	if m.detail.Active() {
		m.detail, _ = m.detail.Update(msg)
		return m, nil
	}

	// Normal main screen keys
	switch {
	case key.Matches(msg, tui.Keys.Tab1):
		m.tabs.SetTab(0)
		m.populateTable()
	case key.Matches(msg, tui.Keys.Tab2):
		m.tabs.SetTab(1)
		m.populateTable()
	case key.Matches(msg, tui.Keys.Tab3):
		m.tabs.SetTab(2)
		m.populateTable()
	case key.Matches(msg, tui.Keys.NextTab):
		m.tabs.NextTab()
		m.populateTable()
	case key.Matches(msg, tui.Keys.Search):
		return m, m.search.Open()
	case key.Matches(msg, tui.Keys.Refresh):
		m.cache.Clear(keyPrices)
		m.screen = screenLoading
		m.spinner = components.NewSpinner("Refreshing prices...")
		m.priceReady = false
		m.wikiReady = true // wiki is still valid
		return m, tea.Batch(
			m.spinner.Init(),
			fetchPricesCmd(m.cache, m.league.ID),
		)
	case key.Matches(msg, tui.Keys.Select):
		if entry := m.table.SelectedEntry(); entry != nil {
			m.detail.Show(entry)
		}
	default:
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) tryProcessGems() tea.Cmd {
	if !m.wikiReady || !m.priceReady {
		return nil
	}
	wiki := m.wiki
	prices := m.prices
	return func() tea.Msg {
		result := domain.ProcessGems(*wiki, prices, 10)
		return tui.DataReadyMsg{Result: result}
	}
}

func (m *Model) populateTable() {
	if m.result == nil {
		return
	}
	activeColor := m.tabs.ActiveColor()
	m.table.SetEntries(m.result.GemPicks, activeColor)
	age := m.cache.Age(keyPrices, priceTTL)
	m.statusbar.SetCacheAge(age)

	// Pass stats to tabs and status bar
	if stats, ok := m.result.ColorStats[activeColor]; ok {
		m.tabs.SetPoolStats(&stats, len(m.result.GemPicks))
	}
	m.statusbar.SetGemCount(len(m.result.GemPicks))
}

func (m Model) View() string {
	if m.err != nil {
		return tui.StyleError.Render("Error: "+m.err.Error()) + "\n\n" +
			tui.StyleHelp.Render("Press q to quit")
	}

	switch m.screen {
	case screenLoading:
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			m.spinner.View())

	case screenLeagueSelect:
		return m.leagueSelect.View()

	case screenMain:
		tabBar := m.tabs.View(m.width)
		tableView := m.table.View()
		statusBar := m.statusbar.View()

		main := lipgloss.JoinVertical(lipgloss.Left, tabBar, tableView, statusBar)

		// Render main to full terminal size
		mainPlaced := lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, main)

		// Overlay search or detail on top using Place compositing
		if m.search.Active() {
			overlay := m.search.View()
			return overlayCenter(mainPlaced, overlay, m.width, m.height)
		}
		if m.detail.Active() {
			overlay := m.detail.View()
			return overlayCenter(mainPlaced, overlay, m.width, m.height)
		}
		return mainPlaced
	}
	return ""
}

// overlayCenter composites an overlay (which already contains Place positioning)
// on top of a background. The overlay's View already handles its own centering,
// so we just return it directly (it fills the full terminal).
func overlayCenter(_, overlay string, _, _ int) string {
	return overlay
}

// Async commands

func fetchLeaguesCmd(c *cache.Cache) tea.Cmd {
	return func() tea.Msg {
		if data, ok := c.Get(keyLeagues); ok {
			if leagues, ok := data.([]domain.League); ok {
				return tui.LeaguesFetchedMsg{Leagues: leagues}
			}
		}
		leagues, err := api.FetchLeagues()
		if err != nil {
			return tui.LeaguesFetchedMsg{Err: err}
		}
		c.Set(keyLeagues, leagues, leagueTTL)
		return tui.LeaguesFetchedMsg{Leagues: leagues}
	}
}

func fetchWikiCmd(c *cache.Cache) tea.Cmd {
	return func() tea.Msg {
		// Try disk cache first
		var wiki domain.WikiData
		if c.LoadFromDisk(keyWiki, &wiki) {
			return tui.WikiFetchedMsg{Wiki: &wiki}
		}
		w, err := api.FetchWikiData()
		if err != nil {
			return tui.WikiFetchedMsg{Err: err}
		}
		c.SaveToDisk(keyWiki, w, wikiTTL)
		return tui.WikiFetchedMsg{Wiki: w}
	}
}

func fetchPricesCmd(c *cache.Cache, league string) tea.Cmd {
	return func() tea.Msg {
		if data, ok := c.Get(keyPrices); ok {
			if prices, ok := data.([]domain.GemPrice); ok {
				return tui.PricesFetchedMsg{Prices: prices}
			}
		}
		prices, err := api.FetchGemPrices(league)
		if err != nil {
			return tui.PricesFetchedMsg{Err: err}
		}
		c.Set(keyPrices, prices, priceTTL)
		return tui.PricesFetchedMsg{Prices: prices}
	}
}
