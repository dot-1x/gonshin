package genshin

import (
	"strconv"

	"github.com/a-h/templ"

	"github.com/dotcchix/gonshin/internal/hoyolab"
	"github.com/dotcchix/gonshin/internal/view/ui"
)

type statItem struct {
	label string
	value string
	icon  templ.Component
}

func statCardValue(v int) string {
	if v == 0 {
		return "—"
	}
	return strconv.Itoa(v)
}

func newStatItems(info *hoyolab.AccountInfo) []statItem {
	value := func(pick func(*hoyolab.AccountInfo) int) string {
		if info == nil {
			return "—"
		}
		return statCardValue(pick(info))
	}
	return []statItem{
		{"Achievements", value(func(i *hoyolab.AccountInfo) int { return i.AchievementNumber }), ui.IconTrophy("h-5 w-5 shrink-0 text-primary")},
		{"Days Active", value(func(i *hoyolab.AccountInfo) int { return i.ActiveDayNumber }), ui.IconCalendarDays("h-5 w-5 shrink-0 text-primary")},
		{"Characters", value(func(i *hoyolab.AccountInfo) int { return i.TotalCharacters }), ui.IconUsers("h-5 w-5 shrink-0 text-primary")},
		{"Friendships", value(func(i *hoyolab.AccountInfo) int { return i.TotalFriendship }), ui.IconHeart("h-5 w-5 shrink-0 text-primary")},
	}
}
