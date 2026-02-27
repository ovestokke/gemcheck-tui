package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/ovestokke/gemcheck-tui/internal/tui"
)

// StatusBarModel is a view-only status bar.
type StatusBarModel struct {
	league   string
	cacheAge time.Duration
	gemCount int
	width    int
}

// NewStatusBar creates a status bar.
func NewStatusBar() StatusBarModel {
	return StatusBarModel{}
}

func (m *StatusBarModel) SetLeague(name string)    { m.league = name }
func (m *StatusBarModel) SetCacheAge(d time.Duration) { m.cacheAge = d }
func (m *StatusBarModel) SetGemCount(n int)        { m.gemCount = n }
func (m *StatusBarModel) SetWidth(w int)           { m.width = w }

func (m StatusBarModel) View() string {
	// Segment 1: League pill
	leagueSeg := tui.StyleStatusLeague.Render(m.league)

	// Segment 2: Cache + gem count
	infoText := fmt.Sprintf("Cache: %s", formatAge(m.cacheAge))
	if m.gemCount > 0 {
		infoText += fmt.Sprintf("  %d gems", m.gemCount)
	}
	infoSeg := tui.StyleStatusInfo.Render(infoText)

	// Segment 3: Help keys (right-aligned)
	helpSeg := tui.StyleStatusHelp.Render("1-3 tab  / search  r refresh  q quit")

	// Calculate gap fill
	leftWidth := lipgloss.Width(leagueSeg) + lipgloss.Width(infoSeg)
	rightWidth := lipgloss.Width(helpSeg)
	gap := m.width - leftWidth - rightWidth
	if gap < 1 {
		gap = 1
	}
	gapFill := tui.StyleStatusInfo.Render(strings.Repeat(" ", gap))

	content := leagueSeg + infoSeg + gapFill + helpSeg
	return lipgloss.NewStyle().Width(m.width).Render(content)
}

func formatAge(d time.Duration) string {
	if d <= 0 {
		return "fresh"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm ago", int(d.Minutes()))
}
