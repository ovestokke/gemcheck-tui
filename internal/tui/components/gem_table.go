package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ovestokke/gemcheck-tui/internal/domain"
	"github.com/ovestokke/gemcheck-tui/internal/tui"
)

// gemEntryItem adapts domain.GemEntry to list.Item.
type gemEntryItem struct {
	entry domain.GemEntry
}

func (i gemEntryItem) FilterValue() string { return i.entry.BaseName }

// itemDelegate renders gem entries with left-border selection and price tiers.
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(gemEntryItem)
	if !ok {
		return
	}

	e := item.entry
	selected := index == m.Index()
	width := m.Width()

	gemColor := tui.ColorForGem(string(e.Color))
	evStr := domain.FormatChaos(e.EV) + " EV"
	priceStyle := tui.PriceStyle(e.EV)

	var border, name, line2prefix string
	if selected {
		border = lipgloss.NewStyle().Foreground(gemColor).Render("\u2503 ")
		name = lipgloss.NewStyle().Bold(true).Foreground(tui.ColorText).Render(e.BaseName)
		line2prefix = lipgloss.NewStyle().Foreground(gemColor).Render("\u2503 ")
	} else {
		border = "  "
		name = lipgloss.NewStyle().Foreground(tui.ColorSubtext0).Render(e.BaseName)
		line2prefix = "  "
	}

	// Line 1: border + name ... EV (right-aligned)
	evRendered := priceStyle.Render(evStr)
	nameWidth := lipgloss.Width(border) + lipgloss.Width(name)
	evWidth := lipgloss.Width(evRendered)
	gap := width - nameWidth - evWidth - 1
	if gap < 1 {
		gap = 1
	}
	line1 := border + name + strings.Repeat(" ", gap) + evRendered

	// Line 2: border + variant count + best price
	var bestPrice float64
	for _, v := range e.Variants {
		if v.SellPrice > bestPrice {
			bestPrice = v.SellPrice
		}
	}
	detailText := fmt.Sprintf("  %d variants%sbest: %s",
		e.VariantCount, tui.Separator, domain.FormatChaos(bestPrice))
	line2 := line2prefix + tui.StyleSubtle.Render(detailText)

	fmt.Fprintf(w, "%s\n%s", line1, line2)
}

// GemTableModel is a scrollable gem list for a single color tab.
type GemTableModel struct {
	list   list.Model
	color  domain.GemColor
	width  int
	height int
}

// NewGemTable creates an empty gem table.
func NewGemTable(width, height int) GemTableModel {
	delegate := itemDelegate{}

	l := list.New(nil, delegate, width, height)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	l.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up", "k"))
	l.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down", "j"))

	return GemTableModel{list: l, width: width, height: height}
}

// SetEntries populates the table with gem entries for a given color.
func (m *GemTableModel) SetEntries(entries []domain.GemEntry, color domain.GemColor) {
	m.color = color
	items := make([]list.Item, 0, len(entries))
	for _, e := range entries {
		if e.Color == color {
			items = append(items, gemEntryItem{entry: e})
		}
	}
	m.list.SetItems(items)
}

// SelectedEntry returns the currently highlighted gem entry, if any.
func (m GemTableModel) SelectedEntry() *domain.GemEntry {
	item, ok := m.list.SelectedItem().(gemEntryItem)
	if !ok {
		return nil
	}
	return &item.entry
}

// SetSize updates the table dimensions.
func (m *GemTableModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

func (m GemTableModel) Init() tea.Cmd {
	return nil
}

func (m GemTableModel) Update(msg tea.Msg) (GemTableModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m GemTableModel) View() string {
	return m.list.View()
}
