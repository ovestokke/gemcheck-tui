package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ovestokke/gemcheck-tui/internal/domain"
	"github.com/ovestokke/gemcheck-tui/internal/tui"
)

type GemTabsModel struct {
	ActiveTab  int
	Tabs       []string
	Colors     []domain.GemColor
	PoolStats  *domain.ColorStats
	TotalGems  int
}

func NewGemTabs() GemTabsModel {
	return GemTabsModel{
		ActiveTab: 0,
		Tabs:      []string{"Red", "Green", "Blue"},
		Colors:    []domain.GemColor{domain.Red, domain.Green, domain.Blue},
	}
}

func (m *GemTabsModel) SetTab(idx int) {
	if idx >= 0 && idx < len(m.Tabs) {
		m.ActiveTab = idx
	}
}

func (m *GemTabsModel) NextTab() {
	m.ActiveTab = (m.ActiveTab + 1) % len(m.Tabs)
}

func (m GemTabsModel) ActiveColor() domain.GemColor {
	return m.Colors[m.ActiveTab]
}

func (m *GemTabsModel) SetPoolStats(stats *domain.ColorStats, totalGems int) {
	m.PoolStats = stats
	m.TotalGems = totalGems
}

func (m GemTabsModel) View(width int) string {
	// Logo pill
	logo := tui.StyleLogo.Render(" \u25c6 GemCheck ")

	// Tabs
	var tabs []string
	for i, t := range m.Tabs {
		c := tui.ColorForGem(string(m.Colors[i]))
		if i == m.ActiveTab {
			style := tui.StyleTabActive.
				Foreground(c).
				BorderBottomForeground(c)
			tabs = append(tabs, style.Render(t))
		} else {
			tabs = append(tabs, tui.StyleTabInactive.Render(t))
		}
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)

	// Pool stats (right of tabs)
	var statsStr string
	if m.PoolStats != nil {
		statsStr = tui.StyleSubtle.Render(
			fmt.Sprintf("%d gems%sPool EV: %s",
				m.PoolStats.PoolSize, tui.Separator, domain.FormatChaos(m.PoolStats.PoolEV)))
	}

	leftSection := lipgloss.JoinHorizontal(lipgloss.Bottom, logo, "  ", tabRow)
	if statsStr != "" {
		leftWidth := lipgloss.Width(leftSection)
		statsWidth := lipgloss.Width(statsStr)
		gap := width - leftWidth - statsWidth - 2
		if gap < 1 {
			gap = 1
		}
		leftSection = leftSection + strings.Repeat(" ", gap) + statsStr
	}

	// Full-width divider
	divider := tui.StyleHeaderDivider.Render(strings.Repeat("\u2500", width))

	return leftSection + "\n" + divider
}
