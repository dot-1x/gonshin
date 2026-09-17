package genshin

import (
	"strconv"

	"github.com/dotcchix/gonshin/internal/hoyolab"
)

type statCol struct {
	label string
	key   string
}

var (
	rowCol1 = []statCol{{"MAX HP", "Max HP"}, {"ATK", "ATK"}, {"EM", "Elemental Mastery"}}
	rowCol2 = []statCol{{"CR", "CRIT Rate"}, {"CD", "CRIT DMG"}, {"ER", "Energy Recharge"}}
)

// rowStat looks up a stat value in a character's final stats.
func rowStat(stats []hoyolab.FinalStat, key string) string {
	return lookupStat(stats, key)
}

// zIndexValue renders the stacked-preview z-index (count down, matching the
// original inline `zIndex: length - i`).
func zIndexValue(total, i int) string {
	return strconv.Itoa(total - i)
}
