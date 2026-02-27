package domain

// GemColor represents the attribute color of a gem.
type GemColor string

const (
	Red   GemColor = "r"
	Green GemColor = "g"
	Blue  GemColor = "b"
)

var AllColors = []GemColor{Red, Green, Blue}

func (c GemColor) Label() string {
	switch c {
	case Red:
		return "Red"
	case Green:
		return "Green"
	case Blue:
		return "Blue"
	default:
		return "Unknown"
	}
}

// BaseGem is a skill gem with its attribute color.
type BaseGem struct {
	Name  string
	Color GemColor
}

// TransfiguredGem is a variant of a base gem.
type TransfiguredGem struct {
	Name      string
	BaseName  string
	Color     GemColor
	SellPrice float64
	Listed    bool
	Count     int
	Icon      string
}

// GemPrice holds pricing info from poe.ninja.
type GemPrice struct {
	Name       string
	ChaosValue float64
	Count      int
	Icon       string
	Corrupted  bool
}

// GemVariantResult holds a transfigured gem with its probability in a specific roll.
type GemVariantResult struct {
	Name      string
	SellPrice float64
	Prob      float64
	Count     int
	Icon      string
	Listed    bool
}

// GemEntry represents a base gem and its transfigured variants with EV.
type GemEntry struct {
	BaseName     string
	Color        GemColor
	Variants     []GemVariantResult
	EV           float64
	VariantCount int
}

// BingoGem is a top gem in the color pool with its hit probability.
type BingoGem struct {
	Name      string
	SellPrice float64
	Prob      float64
	Count     int
	Icon      string
}

// ColorStats holds pool-level statistics for a gem color.
type ColorStats struct {
	Color    GemColor
	PoolSize int
	PoolEV   float64
	Bingo    []BingoGem
}

// ProcessedResult holds all computed data ready for display.
type ProcessedResult struct {
	ColorStats    map[GemColor]ColorStats
	GemPicks      []GemEntry
	TotalLines    int
	TotalTransfig int
}

// League represents a PoE league.
type League struct {
	ID   string
	Text string
}

// WikiData holds scraped gem data from poewiki.
type WikiData struct {
	BaseGems       map[GemColor][]string // color -> sorted base gem names
	TransfigGems   map[GemColor][]string // color -> sorted transfigured gem names
}
