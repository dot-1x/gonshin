package hoyolab

// CharacterIndex maps an icon URL to the owned roster character.
type CharacterIndex map[string]AccountCharacter

// BuildCharacterIndex indexes owned characters by icon URL so lean roster
// references (abyss battles only carry `{icon, level}`) can be enriched with
// name / element / constellation / weapon.
func BuildCharacterIndex(characters []AccountCharacter) CharacterIndex {
	index := make(CharacterIndex, len(characters))
	for _, c := range characters {
		index[c.Icon] = c
	}
	return index
}

// ResolveAbyssChar resolves a lean `{icon, level}` reference into a full
// display character, falling back to the reference itself when not owned.
func ResolveAbyssChar(ref AbyssBattleCharacter, index CharacterIndex) DisplayCharacter {
	if owned, ok := index[ref.Icon]; ok {
		element := owned.Element
		constellation := owned.ActivedConstellationNum
		weapon := owned.Weapon
		return DisplayCharacter{
			Icon:          owned.Icon,
			Name:          owned.Name,
			Element:       &element,
			Level:         owned.Level,
			Constellation: &constellation,
			Weapon:        &weapon,
		}
	}
	return DisplayCharacter{
		Icon:  ref.Icon,
		Name:  "Unknown",
		Level: ref.Level,
	}
}

// ToDisplayChar converts a theater character into a display character.
func ToDisplayChar(c TheaterCharacter) DisplayCharacter {
	element := c.Element
	return DisplayCharacter{
		Icon:    c.Avatar,
		Name:    c.Name,
		Element: &element,
		Level:   c.Level,
	}
}

// ToStygianDisplayChar converts a rich stygian character into a display
// character. Lean fallbacks (no weapon) omit the constellation so the row
// renders like the original.
func ToStygianDisplayChar(c StygianCharacter) DisplayCharacter {
	element := c.Element
	dc := DisplayCharacter{
		Icon:         c.Icon,
		Name:         c.Name,
		Element:      &element,
		Level:        c.Level,
		Weapon:       c.Weapon,
		FinalStats:   c.FinalStats,
		ArtifactSets: c.ArtifactSets,
	}
	if c.Weapon != nil {
		constellation := c.ActivedConstellationNum
		dc.Constellation = &constellation
	}
	return dc
}
