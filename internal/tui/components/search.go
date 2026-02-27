package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ovestokke/gemcheck-tui/internal/domain"
	"github.com/ovestokke/gemcheck-tui/internal/tui"
)

const maxSearchResults = 10

// SearchModel is a fuzzy search overlay for gem entries.
type SearchModel struct {
	input   textinput.Model
	allGems []domain.GemEntry
	results []domain.GemEntry
	cursor  int
	active  bool
	width   int
	height  int
}

// NewSearch creates a search overlay.
func NewSearch() SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.CharLimit = 64
	ti.Width = 40
	ti.PromptStyle = lipgloss.NewStyle().Foreground(tui.ColorLavender)
	ti.TextStyle = lipgloss.NewStyle().Foreground(tui.ColorText)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(tui.ColorOverlay0)

	return SearchModel{input: ti}
}

func (m *SearchModel) SetGems(gems []domain.GemEntry) { m.allGems = gems }
func (m *SearchModel) SetSize(w, h int)               { m.width = w; m.height = h }
func (m SearchModel) Active() bool                     { return m.active }

// Open activates the search overlay.
func (m *SearchModel) Open() tea.Cmd {
	m.active = true
	m.input.SetValue("")
	m.results = nil
	m.cursor = 0
	m.input.Focus()
	return textinput.Blink
}

// Close deactivates the search overlay.
func (m *SearchModel) Close() {
	m.active = false
	m.input.Blur()
}

// SelectedEntry returns the entry under the cursor, if any.
func (m SearchModel) SelectedEntry() *domain.GemEntry {
	if m.cursor >= 0 && m.cursor < len(m.results) {
		return &m.results[m.cursor]
	}
	return nil
}

func (m SearchModel) Init() tea.Cmd {
	return nil
}

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	if !m.active {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Close()
			return m, nil
		case "enter":
			return m, nil // caller checks SelectedEntry
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.filterResults()
	return m, cmd
}

func (m *SearchModel) filterResults() {
	query := strings.ToLower(m.input.Value())
	if query == "" {
		m.results = nil
		m.cursor = 0
		return
	}

	m.results = nil
	for _, g := range m.allGems {
		if strings.Contains(strings.ToLower(g.BaseName), query) {
			m.results = append(m.results, g)
			if len(m.results) >= maxSearchResults {
				break
			}
		}
	}
	if m.cursor >= len(m.results) {
		m.cursor = max(0, len(m.results)-1)
	}
}

func (m SearchModel) View() string {
	if !m.active {
		return ""
	}

	var b strings.Builder
	b.WriteString(m.input.View())
	b.WriteString("\n")

	accentCursor := lipgloss.NewStyle().Foreground(tui.ColorLavender)

	for i, g := range m.results {
		gemColor := tui.ColorForGem(string(g.Color))
		dot := lipgloss.NewStyle().Foreground(gemColor).Render("\u25cf ")

		var prefix string
		if i == m.cursor {
			prefix = accentCursor.Render("\u276f ")
		} else {
			prefix = "  "
		}

		evStyle := tui.PriceStyle(g.EV)
		line := prefix + dot + g.BaseName + "  " +
			evStyle.Render(domain.FormatChaos(g.EV)+" EV")
		b.WriteString(line + "\n")
	}

	if m.input.Value() == "" && len(m.results) == 0 {
		b.WriteString(tui.StyleSubtle.Render("  Type to search..."))
	} else if m.input.Value() != "" && len(m.results) == 0 {
		b.WriteString(tui.StyleSubtle.Render("  No results"))
	}

	popup := tui.StyleSearchInput.Width(46).Render(b.String())

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, popup)
}
