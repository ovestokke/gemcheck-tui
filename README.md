# GemCheck

A terminal UI for analyzing Path of Exile transfigured gem prices and expected value.

Fetches live pricing from poe.ninja and gem data from the PoE wiki, then calculates EV for each gem's transfiguration pool using best-of-3 draw statistics. Helps you figure out which gems are worth farming.

## Features

- Live gem prices from poe.ninja
- Transfigured gem data scraped from poewiki.net
- Expected value calculation per gem and per color pool (best-of-3 draws)
- "Bingo" probability for hitting specific high-value gems
- Color-tabbed browsing (Red / Green / Blue)
- Fuzzy search
- Detail view with full variant breakdown
- Multi-league support
- Local caching with disk persistence

## Install

Download a binary from the [latest release](https://github.com/ovestokke/gemcheck-tui/releases/latest).

**macOS users:** clear the quarantine flag before running:

```
xattr -d com.apple.quarantine gemcheck-darwin-*
```

### Build from source

```
go build -o gemcheck ./cmd/gemcheck
```

Or using the Makefile:

```
make build
make run
```

## Usage

Launch with `./gemcheck`. Select a league, then browse gems by color tab.

### Keybindings

| Key | Action |
|-----|--------|
| `1` `2` `3` | Switch to Red / Green / Blue tab |
| `Tab` | Cycle tabs |
| `/` | Search |
| `Enter` | Open gem detail |
| `r` | Refresh prices |
| `j` / `k` | Navigate |
| `Esc` | Close overlay |
| `q` | Quit |

## How EV is calculated

Each transfigured gem belongs to a color pool. When you use a Lens on a gem, you get one of the transfigurations at random. GemCheck models a "best-of-3" scenario:

- **Gem EV** = average price across a gem's transfigurations
- **Pool EV** = weighted sum where each gem's weight is its probability of being the best result across 3 independent draws
- **Bingo chance** = probability of seeing a specific gem in at least 1 of 3 draws: `1 - ((n-1)/n)^3`

## Project structure

```
cmd/gemcheck/       Entry point
internal/
  app/              Bubble Tea top-level model
  api/              poe.ninja client + poewiki scraper
  domain/           Gem models and EV math
  cache/            In-memory TTL cache with disk persistence
  tui/              Theme, keybindings, and UI components
    components/     Table, tabs, status bar, detail, search
```

## Cache

Data is cached in `~/.cache/gemcheck/`:

| Data | TTL |
|------|-----|
| Leagues | 1 hour |
| Prices | 5 minutes |
| Wiki gems | 24 hours (disk-persisted) |

Press `r` to force-refresh prices.

---

Agentic engineered with [Claude Code](https://claude.ai/claude-code).
