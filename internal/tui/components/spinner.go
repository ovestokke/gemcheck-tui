package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ovestokke/gemcheck-tui/internal/tui"
)

type SpinnerModel struct {
	spinner spinner.Model
	Message string
}

func NewSpinner(msg string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(tui.ColorLavender)
	return SpinnerModel{spinner: s, Message: msg}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (SpinnerModel, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m SpinnerModel) View() string {
	return m.spinner.View() + " " + tui.StyleSubtle.Render(m.Message)
}
