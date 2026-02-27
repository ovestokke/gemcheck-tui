package components

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ovestokke/gemcheck-tui/internal/domain"
	"github.com/ovestokke/gemcheck-tui/internal/tui"
)

type leagueItem struct {
	league domain.League
}

func (i leagueItem) Title() string       { return i.league.Text }
func (i leagueItem) Description() string { return "" }
func (i leagueItem) FilterValue() string { return i.league.Text }

type LeagueSelectModel struct {
	list     list.Model
	Selected *domain.League
	width    int
	height   int
}

func NewLeagueSelect(leagues []domain.League, width, height int) LeagueSelectModel {
	items := make([]list.Item, len(leagues))
	for i, l := range leagues {
		items[i] = leagueItem{league: l}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(tui.ColorYellow).
		BorderLeftForeground(tui.ColorYellow)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(tui.ColorSubtext0).
		BorderLeftForeground(tui.ColorYellow)

	l := list.New(items, delegate, width, height-4)
	l.Title = "Select League"
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorYellow).
		Padding(0, 1)
	l.SetShowStatusBar(false)
	l.SetShowHelp(true)
	l.DisableQuitKeybindings()

	return LeagueSelectModel{list: l, width: width, height: height}
}

func (m LeagueSelectModel) Init() tea.Cmd {
	return nil
}

func (m LeagueSelectModel) Update(msg tea.Msg) (LeagueSelectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if item, ok := m.list.SelectedItem().(leagueItem); ok {
				m.Selected = &item.league
				return m, func() tea.Msg {
					return tui.LeagueSelectedMsg{League: item.league}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m LeagueSelectModel) View() string {
	return m.list.View()
}
