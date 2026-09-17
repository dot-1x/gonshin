package genshin

import (
	"sort"

	"github.com/dotcchix/gonshin/internal/hoyolab"
)

// elementCounts returns per-element counts plus the elements ordered by count
// descending (stable by first appearance), matching the original useMemo.
func elementCounts(chars []hoyolab.AccountCharacter) (map[hoyolab.Element]int, []hoyolab.Element) {
	counts := make(map[hoyolab.Element]int)
	order := make([]hoyolab.Element, 0, len(ElementColors))
	for _, c := range chars {
		if c.Element == "" {
			continue
		}
		if _, seen := counts[c.Element]; !seen {
			order = append(order, c.Element)
		}
		counts[c.Element]++
	}
	sort.SliceStable(order, func(i, j int) bool {
		return counts[order[i]] > counts[order[j]]
	})
	return counts, order
}

// filterByElement filters characters by the active element ("All" or "" keeps
// everything).
func filterByElement(chars []hoyolab.AccountCharacter, active string) []hoyolab.AccountCharacter {
	if active == "" || active == "All" {
		return chars
	}
	filtered := make([]hoyolab.AccountCharacter, 0, len(chars))
	for _, c := range chars {
		if string(c.Element) == active {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// showcaseTarget returns the htmx filter URL, toggling off when already active.
func showcaseTarget(element, active string) string {
	if element != "All" && element == active {
		return "/showcase?element=All"
	}
	return "/showcase?element=" + element
}
