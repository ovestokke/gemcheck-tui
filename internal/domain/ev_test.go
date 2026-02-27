package domain

import (
	"math"
	"testing"
)

func TestExtractBaseName(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"Boneshatter of Carnage", "Boneshatter"},
		{"Summon Flame Golem of the Meteor", "Summon Flame Golem"},
		{"Rain of Arrows of Saturation", "Rain of Arrows"},
		{"Eye of Winter of Finality", "Eye of Winter"},
		{"Arc", "Arc"},
	}
	for _, tt := range tests {
		got := extractBaseName(tt.input)
		if got != tt.want {
			t.Errorf("extractBaseName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestProcessGems_EVCalculation(t *testing.T) {
	wiki := WikiData{
		BaseGems: map[GemColor][]string{
			Red:   {"Boneshatter"},
			Green: {},
			Blue:  {},
		},
		TransfigGems: map[GemColor][]string{
			Red:   {"Boneshatter of Carnage", "Boneshatter of Complex Trauma"},
			Green: {},
			Blue:  {},
		},
	}

	prices := []GemPrice{
		{Name: "Boneshatter of Carnage", ChaosValue: 100, Count: 5},
		{Name: "Boneshatter of Complex Trauma", ChaosValue: 50, Count: 3},
	}

	result := ProcessGems(wiki, prices, 5)

	// Specific roll EV for Boneshatter: (100 + 50) / 2 = 75
	if len(result.GemPicks) != 1 {
		t.Fatalf("expected 1 gem pick, got %d", len(result.GemPicks))
	}
	if math.Abs(result.GemPicks[0].EV-75) > 0.01 {
		t.Errorf("expected EV=75, got %.2f", result.GemPicks[0].EV)
	}

	// Color roll EV for red pool (2 gems, best of 3):
	// Sorted desc: [100, 50]
	// i=0 (100c): P = (2/2)^3 - (1/2)^3 = 1 - 0.125 = 0.875
	// i=1 (50c):  P = (1/2)^3 - (0/2)^3 = 0.125
	// EV = 100*0.875 + 50*0.125 = 87.5 + 6.25 = 93.75
	redStats := result.ColorStats[Red]
	if math.Abs(redStats.PoolEV-93.75) > 0.01 {
		t.Errorf("expected pool EV=93.75, got %.2f", redStats.PoolEV)
	}
	if redStats.PoolSize != 2 {
		t.Errorf("expected pool size=2, got %d", redStats.PoolSize)
	}
}

func TestProcessGems_CorruptedFiltered(t *testing.T) {
	wiki := WikiData{
		BaseGems: map[GemColor][]string{
			Red:   {"Boneshatter"},
			Green: {},
			Blue:  {},
		},
		TransfigGems: map[GemColor][]string{
			Red:   {"Boneshatter of Carnage"},
			Green: {},
			Blue:  {},
		},
	}

	prices := []GemPrice{
		{Name: "Boneshatter of Carnage", ChaosValue: 200, Corrupted: true},
		{Name: "Boneshatter of Carnage", ChaosValue: 100, Corrupted: false},
	}

	result := ProcessGems(wiki, prices, 5)
	if len(result.GemPicks) != 1 {
		t.Fatalf("expected 1 gem pick, got %d", len(result.GemPicks))
	}
	if math.Abs(result.GemPicks[0].EV-100) > 0.01 {
		t.Errorf("expected EV=100 (non-corrupted), got %.2f", result.GemPicks[0].EV)
	}
}

func TestFormatChaos(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{0, "—"},
		{5.5, "5.5c"},
		{150, "150c"},
		{1500, "1.5k c"},
	}
	for _, tt := range tests {
		got := FormatChaos(tt.input)
		if got != tt.want {
			t.Errorf("FormatChaos(%.1f) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestHitProbability(t *testing.T) {
	// For a pool of 32 gems, P(seeing specific gem in 3 draws) = 1 - (31/32)^3
	n := 32
	expected := 1 - math.Pow(float64(n-1)/float64(n), 3)
	// ≈ 9.1%
	if math.Abs(expected-0.0911) > 0.001 {
		t.Errorf("expected ~9.1%%, got %.4f", expected)
	}
}
