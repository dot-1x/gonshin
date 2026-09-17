package genshin

import "github.com/dotcchix/gonshin/internal/hoyolab"

// abyssTeams resolves a chamber's two halves into display teams.
func abyssTeams(level hoyolab.AbyssLevel, index hoyolab.CharacterIndex) []hoyolab.DisplayTeam {
	teams := make([]hoyolab.DisplayTeam, 0, len(level.Battles))
	for half, battle := range level.Battles {
		name := "Second Half"
		if half == 0 {
			name = "First Half"
		}
		chars := make([]hoyolab.DisplayCharacter, 0, len(battle))
		for _, c := range battle {
			chars = append(chars, hoyolab.ResolveAbyssChar(c, index))
		}
		teams = append(teams, hoyolab.DisplayTeam{Name: name, Characters: chars})
	}
	return teams
}
