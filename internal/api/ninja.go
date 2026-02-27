package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ovestokke/gemcheck-tui/internal/domain"
)

const (
	leaguesURL = "https://api.pathofexile.com/leagues?type=main&compact=1&game=poe1"
	ninjaAPI   = "https://poe.ninja/api/data/itemoverview"
	userAgent  = "gemcheck-tui/1.0"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

func doGet(rawURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, rawURL)
	}
	return io.ReadAll(resp.Body)
}

// FetchLeagues returns the list of active PoE1 leagues.
func FetchLeagues() ([]domain.League, error) {
	body, err := doGet(leaguesURL)
	if err != nil {
		return nil, fmt.Errorf("fetching leagues: %w", err)
	}

	var raw []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing leagues: %w", err)
	}

	var leagues []domain.League
	for _, l := range raw {
		leagues = append(leagues, domain.League{ID: l.ID, Text: l.ID})
	}
	return leagues, nil
}

// FetchGemPrices fetches gem prices from poe.ninja for the given league.
func FetchGemPrices(league string) ([]domain.GemPrice, error) {
	u := fmt.Sprintf("%s?league=%s&type=SkillGem&game=poe1",
		ninjaAPI, url.QueryEscape(league))

	body, err := doGet(u)
	if err != nil {
		return nil, fmt.Errorf("fetching gem prices: %w", err)
	}

	var resp struct {
		Lines []struct {
			Name       string  `json:"name"`
			ChaosValue float64 `json:"chaosValue"`
			Count      int     `json:"count"`
			Icon       string  `json:"icon"`
			Corrupted  bool    `json:"corrupted"`
		} `json:"lines"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing gem prices: %w", err)
	}

	prices := make([]domain.GemPrice, len(resp.Lines))
	for i, l := range resp.Lines {
		prices[i] = domain.GemPrice{
			Name:       l.Name,
			ChaosValue: l.ChaosValue,
			Count:      l.Count,
			Icon:       l.Icon,
			Corrupted:  l.Corrupted,
		}
	}
	return prices, nil
}
