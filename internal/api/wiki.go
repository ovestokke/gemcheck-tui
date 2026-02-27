package api

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/net/html"

	"github.com/ovestokke/gemcheck-tui/internal/domain"
)

const (
	baseGemsURL     = "https://www.poewiki.net/wiki/List_of_skill_gems"
	transfigGemsURL = "https://www.poewiki.net/wiki/Transfigured_skill_gem"
)

// FetchWikiData scrapes poewiki for base gem colors and transfigured gem lists.
func FetchWikiData() (*domain.WikiData, error) {
	baseGems, err := scrapeGemTable(baseGemsURL, 3)
	if err != nil {
		return nil, fmt.Errorf("scraping base gems: %w", err)
	}

	transfigGems, err := scrapeGemTable(transfigGemsURL, 3)
	if err != nil {
		return nil, fmt.Errorf("scraping transfigured gems: %w", err)
	}

	return &domain.WikiData{
		BaseGems:     baseGems,
		TransfigGems: transfigGems,
	}, nil
}

// scrapeGemTable fetches a poewiki page and extracts gem names from the first
// nTables item-tables. Tables are in order: Strength (r), Dexterity (g), Intelligence (b).
func scrapeGemTable(pageURL string, nTables int) (map[domain.GemColor][]string, error) {
	body, err := doGet(pageURL)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	tables := findItemTables(doc)
	colors := []domain.GemColor{domain.Red, domain.Green, domain.Blue}
	result := make(map[domain.GemColor][]string)

	for i := 0; i < nTables && i < len(tables); i++ {
		gems := extractGemsFromTable(tables[i])
		sort.Strings(gems)
		result[colors[i]] = gems
	}

	// Ensure all colors have entries
	for _, c := range colors {
		if _, ok := result[c]; !ok {
			result[c] = []string{}
		}
	}

	return result, nil
}

// findItemTables finds all <table> elements with class "wikitable sortable item-table".
func findItemTables(n *html.Node) []*html.Node {
	var tables []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			if hasClass(n, "item-table") {
				tables = append(tables, n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return tables
}

// hasClass checks if an HTML node has a specific class.
func hasClass(n *html.Node, class string) bool {
	for _, a := range n.Attr {
		if a.Key == "class" {
			for _, c := range strings.Fields(a.Val) {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

// extractGemsFromTable extracts gem names from a wikitable's first column links.
func extractGemsFromTable(table *html.Node) []string {
	var gems []string
	seen := make(map[string]bool)

	rows := findElements(table, "tr")
	for _, row := range rows {
		tds := findElements(row, "td")
		if len(tds) == 0 {
			continue
		}
		// Get the first <a> in the first <td> that looks like a gem link
		links := findElements(tds[0], "a")
		for _, link := range links {
			href := getAttr(link, "href")
			if strings.HasPrefix(href, "/wiki/File:") || strings.HasPrefix(href, "/wiki/Special:") {
				continue
			}
			title := getAttr(link, "title")
			text := getTextContent(link)
			if text != "" && text == title && !seen[text] {
				seen[text] = true
				gems = append(gems, text)
				break
			}
		}
	}

	return gems
}

// findElements finds all descendant elements with the given tag name.
func findElements(n *html.Node, tag string) []*html.Node {
	var result []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			result = append(result, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return result
}

// getAttr returns the value of an attribute.
func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// getTextContent returns the text content of a node.
func getTextContent(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(sb.String())
}
