package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ovestokke/gemcheck-tui/internal/domain"
	"github.com/ovestokke/gemcheck-tui/internal/tui"
)

// DetailModel displays variant details for a selected gem entry.
type DetailModel struct {
	entry  *domain.GemEntry
	active bool
	scroll int
	width  int
	height int
}

// NewDetail creates a detail popup.
func NewDetail() DetailModel {
	return DetailModel{}
}

func (m *DetailModel) SetSize(w, h int) { m.width = w; m.height = h }
func (m DetailModel) Active() bool      { return m.active }

// Show displays the detail popup for an entry.
func (m *DetailModel) Show(entry *domain.GemEntry) {
	m.entry = entry
	m.active = true
	m.scroll = 0
}

// Hide closes the detail popup.
func (m *DetailModel) Hide() {
	m.active = false
	m.entry = nil
}

func (m DetailModel) Init() tea.Cmd {
	return nil
}

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	if !m.active {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.Hide()
		case "up", "k":
			if m.scroll > 0 {
				m.scroll--
			}
		case "down", "j":
			m.scroll++
		}
	}
	return m, nil
}

func (m DetailModel) View() string {
	if !m.active || m.entry == nil {
		return ""
	}

	e := m.entry
	popupWidth := min(60, m.width-4)
	innerWidth := popupWidth - 6 // account for border + padding

	var b strings.Builder

	// Header: gem name in gem color, bold
	gemColor := tui.ColorForGem(string(e.Color))
	title := lipgloss.NewStyle().Bold(true).Foreground(gemColor).Render(e.BaseName)
	b.WriteString(title + "\n")

	// Divider under header
	b.WriteString(tui.StyleHeaderDivider.Render(strings.Repeat("\u2500", innerWidth)) + "\n")

	// EV line with price tier
	evStyle := tui.PriceStyle(e.EV)
	b.WriteString(fmt.Sprintf("%s gem  %s  EV: %s  %s  %d variants\n\n",
		e.Color.Label(),
		tui.Separator,
		evStyle.Render(domain.FormatChaos(e.EV)),
		tui.Separator,
		e.VariantCount))

	// Variants
	for _, v := range e.Variants {
		nameStyle := lipgloss.NewStyle().Foreground(tui.ColorText)
		priceStyle := tui.PriceStyle(v.SellPrice)

		if !v.Listed {
			nameStyle = nameStyle.Foreground(tui.ColorOverlay0)
			priceStyle = lipgloss.NewStyle().Foreground(tui.ColorOverlay0)
		}

		b.WriteString(fmt.Sprintf("  %s\n", nameStyle.Render(v.Name)))

		price := priceStyle.Render(domain.FormatChaos(v.SellPrice))
		prob := tui.StyleProb.Render(domain.FormatPct(v.Prob))

		unlisted := ""
		if !v.Listed {
			unlisted = tui.StyleSubtle.Render(" unlisted")
		}

		b.WriteString(fmt.Sprintf("    %s  %s%s\n", price, prob, unlisted))
	}

	// Footer hint
	b.WriteString("\n")
	b.WriteString(tui.StyleHelp.Render("esc close  \u2191\u2193 scroll"))

	content := b.String()

	// Apply scroll by trimming lines
	lines := strings.Split(content, "\n")
	maxVisible := m.height - 6
	if maxVisible < 5 {
		maxVisible = 5
	}
	if m.scroll > len(lines)-maxVisible {
		m.scroll = max(0, len(lines)-maxVisible)
	}
	if m.scroll > 0 && m.scroll < len(lines) {
		lines = lines[m.scroll:]
	}
	if len(lines) > maxVisible {
		lines = lines[:maxVisible]
	}
	content = strings.Join(lines, "\n")

	popup := tui.StyleDetailPopup.Width(popupWidth).Render(content)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, popup)
}
