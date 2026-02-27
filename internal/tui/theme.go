package tui

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha palette
var (
	ColorBase     = lipgloss.Color("#1e1e2e")
	ColorMantle   = lipgloss.Color("#181825")
	ColorCrust    = lipgloss.Color("#11111b")
	ColorSurface0 = lipgloss.Color("#313244")
	ColorSurface1 = lipgloss.Color("#45475a")
	ColorSurface2 = lipgloss.Color("#585b70")
	ColorOverlay0 = lipgloss.Color("#6c7086")
	ColorOverlay1 = lipgloss.Color("#7f849c")
	ColorText     = lipgloss.Color("#cdd6f4")
	ColorSubtext0 = lipgloss.Color("#a6adc8")
	ColorSubtext1 = lipgloss.Color("#bac2de")
	ColorRed      = lipgloss.Color("#f38ba8")
	ColorGreen    = lipgloss.Color("#a6e3a1")
	ColorBlue     = lipgloss.Color("#89b4fa")
	ColorYellow   = lipgloss.Color("#f9e2af")
	ColorPeach    = lipgloss.Color("#fab387")
	ColorMauve    = lipgloss.Color("#cba6f7")
	ColorLavender = lipgloss.Color("#b4befe")
	ColorTeal     = lipgloss.Color("#94e2d5")
	ColorGold     = lipgloss.Color("#f5c211")
)

// Separator used between status items
var Separator = lipgloss.NewStyle().Foreground(ColorOverlay0).Render(" \u2022 ")

var (
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorYellow)

	StyleSubtle = lipgloss.NewStyle().
			Foreground(ColorOverlay1)

	StylePrice = lipgloss.NewStyle().
			Foreground(ColorYellow).
			Bold(true)

	StyleProb = lipgloss.NewStyle().
			Foreground(ColorGreen)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorOverlay0)

	// --- Branded logo ---
	StyleLogo = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCrust).
			Background(ColorGold).
			Padding(0, 1)

	// --- Tab styles ---
	StyleTabActive = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2).
			Border(lipgloss.NormalBorder(), false, false, true, false)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(ColorOverlay1).
				Padding(0, 2)

	// --- Header divider ---
	StyleHeaderDivider = lipgloss.NewStyle().
				Foreground(ColorSurface1)

	// --- Card styles ---
	StyleCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSurface1).
			Padding(0, 1)

	StyleCardTitle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	// --- Search overlay ---
	StyleSearchInput = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorLavender).
				Padding(0, 1)

	// --- Detail popup ---
	StyleDetailPopup = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorMauve).
				Padding(1, 2)

	// --- Status bar segments ---
	StyleStatusLeague = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorCrust).
				Background(ColorMauve).
				Padding(0, 1)

	StyleStatusInfo = lipgloss.NewStyle().
			Foreground(ColorSubtext0).
			Background(ColorMantle).
			Padding(0, 1)

	StyleStatusHelp = lipgloss.NewStyle().
			Foreground(ColorOverlay0).
			Background(ColorCrust).
			Padding(0, 1)

	StyleStatusBar = lipgloss.NewStyle().
			Background(ColorMantle)

	// --- Left-border selection indicators ---
	StyleSelectedBorder = lipgloss.NewStyle().
				Foreground(ColorYellow).
				SetString("\u2503 ")

	StyleNormalBorder = lipgloss.NewStyle().
				SetString("  ")

	// --- Price tier styles ---
	StylePriceHigh = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	StylePriceMid = lipgloss.NewStyle().
			Foreground(ColorYellow)

	StylePriceLow = lipgloss.NewStyle().
			Foreground(ColorOverlay1)
)

// PriceStyle returns a tier-colored style based on chaos value.
func PriceStyle(chaos float64) lipgloss.Style {
	switch {
	case chaos >= 50:
		return StylePriceHigh
	case chaos >= 10:
		return StylePriceMid
	default:
		return StylePriceLow
	}
}

func ColorForGem(color string) lipgloss.Color {
	switch color {
	case "r":
		return ColorRed
	case "g":
		return ColorGreen
	case "b":
		return ColorBlue
	default:
		return ColorOverlay1
	}
}
