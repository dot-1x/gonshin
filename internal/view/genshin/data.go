// Package genshin renders the Genshin Impact profile page (templ components).
package genshin

import (
	"strconv"

	"github.com/dotcchix/gonshin/internal/hoyolab"
)

// ElementColors maps each element to its accent colour (from genshin-data.ts).
var ElementColors = map[hoyolab.Element]string{
	hoyolab.Pyro:    "#ef7938",
	hoyolab.Hydro:   "#4cc2f1",
	hoyolab.Electro: "#b08fc2",
	hoyolab.Cryo:    "#9fd6e3",
	hoyolab.Anemo:   "#74c2a8",
	hoyolab.Geo:     "#f2b723",
	hoyolab.Dendro:  "#a0c83b",
}

// colorOf returns the element colour, or "" when the element is unknown.
func colorOf(el *hoyolab.Element) string {
	if el == nil {
		return ""
	}
	return ElementColors[*el]
}

// FormatClearTime formats a clear time in seconds as a compact "Xs" string.
func FormatClearTime(seconds *int) string {
	if seconds == nil || *seconds <= 0 {
		return "—"
	}
	return strconv.Itoa(*seconds) + "s"
}

// FormatClearTimeN is the value-typed variant used by stygian cycles.
func FormatClearTimeN(seconds int) string {
	return FormatClearTime(&seconds)
}

// StygianClearImage returns the clear image for a difficulty index.
func StygianClearImage(index int) string {
	switch index {
	case 5:
		return "/static/genshin/stygian-clear/fearless.webp"
	default:
		return "/static/genshin/stygian-clear/dire.webp"
	}
}

// StygianDifficultyLabel maps a difficulty index to its display label.
func StygianDifficultyLabel(index int) string {
	switch index {
	case 5:
		return "Fearless"
	default:
		return "Dire"
	}
}

// navActiveClass returns the tab button styling for its active state.
func navActiveClass(active bool) string {
	if active {
		return "border-b-2 border-primary text-primary"
	}
	return "border-b-2 border-transparent text-muted-foreground hover:text-foreground"
}

// Tab describes a nav tab.
type Tab struct {
	ID    string
	Label string
	Icon  string
}

// Tabs is the ordered tab list from genshin-nav.tsx.
var Tabs = []Tab{
	{ID: "home", Label: "Detail", Icon: "/static/genshin/sumeru.png"},
	{ID: "characters", Label: "Characters", Icon: "/static/genshin/Icon_Character.webp"},
	{ID: "spiral", Label: "Spiral", Icon: "/static/genshin/spiral.webp"},
	{ID: "stygian", Label: "Stygian", Icon: "/static/genshin/stygian.png"},
	{ID: "theater", Label: "Theater", Icon: "/static/genshin/Imaginarium_Theater.webp"},
}

// IsTab reports whether id is a known tab.
func IsTab(id string) bool {
	for _, t := range Tabs {
		if t.ID == id {
			return true
		}
	}
	return false
}
