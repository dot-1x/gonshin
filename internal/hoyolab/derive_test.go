package hoyolab

import "testing"

func TestBuildCharacterIndexAndResolve(t *testing.T) {
	chars := []AccountCharacter{
		{
			ID: 1, Icon: "a.png", Name: "Venti", Element: Anemo, Level: 90,
			ActivedConstellationNum: 3,
			Weapon:                  Weapon{Name: "Elegy", Level: 90, Icon: "w.png"},
		},
	}
	index := BuildCharacterIndex(chars)

	owned := ResolveAbyssChar(AbyssBattleCharacter{Icon: "a.png", Level: 90}, index)
	if owned.Name != "Venti" || owned.Element == nil || *owned.Element != Anemo {
		t.Fatalf("owned resolve: %+v", owned)
	}
	if owned.Constellation == nil || *owned.Constellation != 3 {
		t.Fatalf("owned constellation: %+v", owned.Constellation)
	}
	if owned.Weapon == nil || owned.Weapon.Name != "Elegy" {
		t.Fatalf("owned weapon: %+v", owned.Weapon)
	}

	unknown := ResolveAbyssChar(AbyssBattleCharacter{Icon: "x.png", Level: 80}, index)
	if unknown.Name != "Unknown" || unknown.Element != nil || unknown.Constellation != nil || unknown.Weapon != nil {
		t.Fatalf("unknown resolve: %+v", unknown)
	}
	if unknown.Icon != "x.png" || unknown.Level != 80 {
		t.Fatalf("unknown fallback fields: %+v", unknown)
	}
}

func TestToStygianDisplayChar(t *testing.T) {
	full := StygianCharacter{
		Icon: "s.png", Name: "Furina", Element: Hydro, Level: 90,
		ActivedConstellationNum: 2,
		Weapon:                  &Weapon{Name: "Favonius", Level: 90, Icon: "w.png"},
		FinalStats:              []FinalStat{{Name: "Max HP", Value: "40000"}},
		ArtifactSets:            []string{"Golden Troupe x4"},
	}
	got := ToStygianDisplayChar(full)
	if got.Weapon == nil || got.Constellation == nil || *got.Constellation != 2 {
		t.Fatalf("full build: %+v", got)
	}
	if len(got.FinalStats) != 1 || len(got.ArtifactSets) != 1 {
		t.Fatalf("full build stats: %+v", got)
	}

	lean := StygianCharacter{Icon: "l.png", Name: "Unknown", Element: Pyro, Level: 70}
	gotLean := ToStygianDisplayChar(lean)
	if gotLean.Weapon != nil || gotLean.Constellation != nil {
		t.Fatalf("lean fallback should omit build: %+v", gotLean)
	}
}

func TestToDisplayChar(t *testing.T) {
	dc := ToDisplayChar(TheaterCharacter{Avatar: "t.png", Name: "Cyno", Level: 90, Element: Electro})
	if dc.Icon != "t.png" || dc.Name != "Cyno" || dc.Element == nil || *dc.Element != Electro {
		t.Fatalf("ToDisplayChar: %+v", dc)
	}
	if dc.Weapon != nil || dc.Constellation != nil {
		t.Fatalf("theater char should have no build: %+v", dc)
	}
}
