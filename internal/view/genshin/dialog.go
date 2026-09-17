package genshin

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/dotcchix/gonshin/internal/hoyolab"
)

var colorTagRe = regexp.MustCompile(`<color[^>]*>|</color>`)

// stripColorTags removes HoYoLAB `<color=...>` markup from effect text.
func stripColorTags(s string) string {
	return colorTagRe.ReplaceAllString(s, "")
}

// normName normalises non-breaking spaces in stat names.
func normName(s string) string {
	return strings.ReplaceAll(s, "\u00a0", " ")
}

// refineSuffix renders the " · R5" weapon refinement suffix.
func refineSuffix(refine *int) string {
	if refine == nil {
		return ""
	}
	return " · R" + strconv.Itoa(*refine)
}

// splitStats splits stats into two roughly equal columns.
func splitStats(stats []hoyolab.FinalStat) (left, right []hoyolab.FinalStat) {
	mid := (len(stats) + 1) / 2
	return stats[:mid], stats[mid:]
}

// statValue looks up a stat value by normalised name.
func lookupStat(stats []hoyolab.FinalStat, key string) string {
	for _, s := range stats {
		if normName(s.Name) == key {
			return s.Value
		}
	}
	return "—"
}

// elementBonus finds the elemental DMG Bonus stat, if any.
func elementBonus(stats []hoyolab.FinalStat) (label, value string, ok bool) {
	for _, s := range stats {
		name := normName(s.Name)
		if strings.HasSuffix(name, "DMG Bonus") && s.Value != "0.0%" {
			return strings.TrimSuffix(name, " DMG Bonus"), s.Value, true
		}
	}
	return "", "", false
}
