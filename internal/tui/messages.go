package tui

import "github.com/ovestokke/gemcheck-tui/internal/domain"

// Messages for async operations

type LeaguesFetchedMsg struct {
	Leagues []domain.League
	Err     error
}

type LeagueSelectedMsg struct {
	League domain.League
}

type WikiFetchedMsg struct {
	Wiki *domain.WikiData
	Err  error
}

type PricesFetchedMsg struct {
	Prices []domain.GemPrice
	Err    error
}

type DataReadyMsg struct {
	Result domain.ProcessedResult
}

type ErrMsg struct {
	Err error
}

type WindowSizeMsg struct {
	Width  int
	Height int
}
