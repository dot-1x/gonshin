package genshin

import (
	"strconv"

	"github.com/dotcchix/gonshin/internal/hoyolab"
)

// stygianTeams resolves a cycle's challenges into display teams.
func stygianTeams(cycle hoyolab.StygianCycle) []hoyolab.DisplayTeam {
	teams := make([]hoyolab.DisplayTeam, 0, len(cycle.Challenges))
	for _, ch := range cycle.Challenges {
		chars := make([]hoyolab.DisplayCharacter, 0, len(ch.Characters))
		for _, c := range ch.Characters {
			chars = append(chars, hoyolab.ToStygianDisplayChar(c))
		}
		seconds := ch.Second
		teams = append(teams, hoyolab.DisplayTeam{
			Name:       ch.Name,
			Characters: chars,
			ClearTime:  &seconds,
		})
	}
	return teams
}

// stygianTarget builds the cycle-selector htmx URL.
func stygianTarget(i int) string {
	return "/stygian?cycle=" + strconv.Itoa(i)
}

// clampCycle keeps a requested cycle index within range.
func clampCycle(i, length int) int {
	if i < 0 || i >= length {
		return 0
	}
	return i
}
