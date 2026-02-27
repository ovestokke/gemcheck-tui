package domain

import (
	"fmt"
	"math"
	"sort"
)

const FontDraws = 3

// ProcessGems calculates EV statistics from wiki gem data and ninja prices.
func ProcessGems(wiki WikiData, prices []GemPrice, topN int) ProcessedResult {
	// Build price lookup: name -> cheapest non-corrupted entry
	priceMap := make(map[string]GemPrice)
	totalLines := 0
	for _, p := range prices {
		if p.Corrupted {
			continue
		}
		totalLines++
		if existing, ok := priceMap[p.Name]; !ok || p.ChaosValue < existing.ChaosValue {
			priceMap[p.Name] = p
		}
	}

	// Build per-base-gem entries from the authoritative wiki list
	var gemEntries []GemEntry
	for _, c := range AllColors {
		names := wiki.TransfigGems[c]
		byBase := make(map[string][]GemVariantResult)

		for _, name := range names {
			baseName := extractBaseName(name)
			p, listed := priceMap[name]
			var sellPrice float64
			var count int
			var icon string
			if listed {
				sellPrice = p.ChaosValue
				count = p.Count
				icon = p.Icon
			}
			byBase[baseName] = append(byBase[baseName], GemVariantResult{
				Name:      name,
				SellPrice: sellPrice,
				Prob:      0, // filled below
				Count:     count,
				Icon:      icon,
				Listed:    listed,
			})
		}

		for baseName, variants := range byBase {
			n := len(variants)
			for i := range variants {
				variants[i].Prob = 1.0 / float64(n)
			}
			sort.Slice(variants, func(i, j int) bool {
				return variants[i].SellPrice > variants[j].SellPrice
			})

			var ev float64
			for _, v := range variants {
				ev += v.SellPrice
			}
			ev /= float64(n)

			gemEntries = append(gemEntries, GemEntry{
				BaseName:     baseName,
				Color:        c,
				Variants:     variants,
				EV:           ev,
				VariantCount: n,
			})
		}
	}

	// Sort gem entries by EV descending
	sort.Slice(gemEntries, func(i, j int) bool {
		return gemEntries[i].EV > gemEntries[j].EV
	})

	// Calculate color pool statistics (best-of-3 order statistic)
	colorStats := make(map[GemColor]ColorStats)
	totalTransfig := 0

	for _, c := range AllColors {
		names := wiki.TransfigGems[c]
		n := len(names)
		totalTransfig += n

		type poolGem struct {
			name      string
			sellPrice float64
			count     int
			icon      string
		}

		pool := make([]poolGem, n)
		for i, name := range names {
			p, ok := priceMap[name]
			if ok {
				pool[i] = poolGem{name, p.ChaosValue, p.Count, p.Icon}
			} else {
				pool[i] = poolGem{name: name}
			}
		}

		// Sort descending by price
		sort.Slice(pool, func(i, j int) bool {
			return pool[i].sellPrice > pool[j].sellPrice
		})

		// EV of best-of-3: P(gem at sorted-index i is max) = ((n-i)/n)^k - ((n-i-1)/n)^k
		var poolEV float64
		for i, g := range pool {
			pBest := math.Pow(float64(n-i)/float64(n), FontDraws) -
				math.Pow(float64(n-i-1)/float64(n), FontDraws)
			poolEV += g.sellPrice * pBest
		}

		// Bingo: top gems with prices
		var bingo []BingoGem
		hitProb := 1 - math.Pow(float64(n-1)/float64(n), FontDraws)
		for _, g := range pool {
			if g.sellPrice <= 0 {
				break
			}
			bingo = append(bingo, BingoGem{
				Name:      g.name,
				SellPrice: g.sellPrice,
				Prob:      hitProb,
				Count:     g.count,
				Icon:      g.icon,
			})
			if len(bingo) >= topN {
				break
			}
		}

		colorStats[c] = ColorStats{
			Color:    c,
			PoolSize: n,
			PoolEV:   poolEV,
			Bingo:    bingo,
		}
	}

	return ProcessedResult{
		ColorStats:    colorStats,
		GemPicks:      gemEntries,
		TotalLines:    totalLines,
		TotalTransfig: totalTransfig,
	}
}

// extractBaseName extracts the base gem name from a transfigured gem name.
// e.g. "Boneshatter of Carnage" -> "Boneshatter"
func extractBaseName(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if i+4 <= len(name) && name[i:i+4] == " of " {
			return name[:i]
		}
	}
	return name
}

// FormatChaos formats a chaos value for display.
func FormatChaos(v float64) string {
	if v == 0 {
		return "—"
	}
	if v >= 1000 {
		return fmt.Sprintf("%.1fk c", v/1000)
	}
	if v >= 100 {
		return fmt.Sprintf("%.0fc", v)
	}
	return fmt.Sprintf("%.1fc", v)
}

// FormatPct formats a probability as a percentage.
func FormatPct(p float64) string {
	return fmt.Sprintf("%.1f%%", p*100)
}
