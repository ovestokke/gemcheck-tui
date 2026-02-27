package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ovestokke/gemcheck-tui/internal/app"
	"github.com/ovestokke/gemcheck-tui/internal/cache"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	cacheDir := filepath.Join(home, ".cache", "gemcheck")
	c := cache.New(cacheDir)

	m := app.NewModel(c)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
