package genshin

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/dotcchix/gonshin/internal/hoyolab"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

func TestPlayerCardWithInfo(t *testing.T) {
	html := render(t, PlayerCard(&accountInfo))
	for _, want := range []string{`Zex.`, `AR 60`, `UID 824677421`, `h-32 w-32 border-2 border-primary`, accountInfo.GameHeadIcon} {
		if !strings.Contains(html, want) {
			t.Errorf("PlayerCard missing %q", want)
		}
	}
}

func TestPlayerCardDefaults(t *testing.T) {
	html := render(t, PlayerCard(nil))
	for _, want := range []string{`dotcchix`, `AR 60`, `UID —`} {
		if !strings.Contains(html, want) {
			t.Errorf("PlayerCard(null) missing %q", want)
		}
	}
	if !strings.Contains(html, `<svg`) {
		t.Errorf("expected fallback sparkles svg")
	}
}

func TestStatsRow(t *testing.T) {
	html := render(t, StatsRow(&accountInfo))
	for _, want := range []string{"Achievements", "Days Active", "Characters", "Friendships", "1362", "1893", "103", "72"} {
		if !strings.Contains(html, want) {
			t.Errorf("StatsRow missing %q", want)
		}
	}
	if n := strings.Count(html, `data-slot="card"`); n != 4 {
		t.Errorf("expected 4 stat cards, got %d", n)
	}
}

func TestNavActive(t *testing.T) {
	html := render(t, Nav("home", nil))
	if n := strings.Count(html, "hx-get=\"/tabs/"); n != 5 {
		t.Errorf("expected 5 tab buttons, got %d", n)
	}
	if !strings.Contains(html, `border-b-2 border-primary text-primary`) {
		t.Errorf("expected active tab styling")
	}
	if strings.Count(html, "border-b-2 border-primary text-primary") != 1 {
		t.Errorf("expected exactly one active tab")
	}
}

func TestCharacterShowcase(t *testing.T) {
	html := render(t, CharacterShowcase(accountCharacters, "All"))
	for _, want := range []string{
		"Character Showcase",
		`id="showcase-panel"`,
		`hx-get="/showcase?element=All"`,
		`hx-get="/showcase?element=Anemo"`,
		`hx-get="/characters/10000022"`,
		`hx-target="#gen-dialog"`,
		`w-64 gap-2 rounded-lg border border-border/50 bg-popover p-3`,
		"Lv.90",
		"C2",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("CharacterShowcase missing %q", want)
		}
	}
	if n := strings.Count(html, `data-slot="card"`); n < 4 {
		t.Errorf("expected header + 3 character cards, got %d", n)
	}
}

func TestCharacterShowcaseFilter(t *testing.T) {
	html := render(t, CharacterShowcasePanel(accountCharacters, "Hydro"))
	if !strings.Contains(html, "Furina") {
		t.Errorf("expected Furina in Hydro filter")
	}
	if strings.Contains(html, "Venti") || strings.Contains(html, "Yanfei") {
		t.Errorf("Hydro filter should exclude other elements")
	}
	// Clicking the already-active element toggles back to All.
	if !strings.Contains(html, `hx-get="/showcase?element=All"`) {
		t.Errorf("expected Hydro toggle to All")
	}
}

func TestCharacterShowcaseEmpty(t *testing.T) {
	html := render(t, CharacterShowcase(nil, "All"))
	if !strings.Contains(html, "No characters available.") {
		t.Errorf("expected empty state")
	}
}

func TestCharacterDialog(t *testing.T) {
	html := render(t, CharacterDialog(characterDetail))
	for _, want := range []string{
		"Venti",
		"Stats",
		"Artifact Sets",
		"Constellations",
		"Viridescent Venerer x4",
		"Splitting Gales",
		"· R5",
		"61.6%",
		"Anemo",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("CharacterDialog missing %q", want)
		}
	}
	if strings.Contains(html, "<color") {
		t.Errorf("expected colour tags stripped from effects")
	}
	if n := strings.Count(html, "fill-amber-400"); n != 5 {
		t.Errorf("expected 5 rarity stars, got %d", n)
	}
}

func TestCharacterDialogStates(t *testing.T) {
	errHTML := render(t, CharacterDialogError("Failed to fetch (503)"))
	if !strings.Contains(errHTML, "Failed to fetch (503)") {
		t.Errorf("error dialog missing message")
	}
	loadHTML := render(t, CharacterDialogLoading("Venti"))
	if !strings.Contains(loadHTML, "Loading Venti...") || !strings.Contains(loadHTML, "animate-spin") {
		t.Errorf("loading dialog malformed")
	}
}

func TestSpiralAbyssSection(t *testing.T) {
	html := render(t, SpiralAbyssSection(accountAbyss, accountCharacters))
	for _, want := range []string{
		"Spiral Abyss",
		"Floor 12-3",
		"Chamber 1",
		"First Half",
		"Second Half",
		"Venti",
		"Unknown",
		`<details open`,
		"group-open:rotate-180",
		"Hydro",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("SpiralAbyssSection missing %q", want)
		}
	}
	if strings.Contains(html, `<details open="false"`) {
		t.Errorf("second team should not be open")
	}
}

func TestSpiralAbyssEmpty(t *testing.T) {
	html := render(t, SpiralAbyssSection(nil, nil))
	if !strings.Contains(html, "No Spiral Abyss records available.") {
		t.Errorf("expected empty state")
	}
}

func TestCharacterRowUnknown(t *testing.T) {
	html := render(t, CharacterRow(hoyolab.DisplayCharacter{Icon: "x.png", Name: "Unknown", Level: 70}))
	if !strings.Contains(html, "Unknown") || !strings.Contains(html, "Lv.70") {
		t.Errorf("unknown character row malformed")
	}
	if !strings.Contains(html, "—") {
		t.Errorf("expected em dash placeholders for missing weapon/stats")
	}
}

func TestStygianOnslaughtSection(t *testing.T) {
	html := render(t, StygianOnslaughtSection(accountStygian, 0))
	for _, want := range []string{
		"Stygian Onslaught",
		"Season One",
		"Fearless",
		"180s",
		"Total Clear Time",
		"First Half",
		"Furina",
		"Golden Troupe x4",
		"MAX HP",
		"40000",
		"90s",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("StygianOnslaughtSection missing %q", want)
		}
	}
	if strings.Contains(html, `hx-get="/stygian?cycle=`) {
		t.Errorf("single cycle should not render a cycle selector")
	}
}

func TestStygianCycleSelector(t *testing.T) {
	multi := &hoyolab.AccountStygian{Cycles: []hoyolab.StygianCycle{
		{Name: "One", Difficulty: 6, TotalClearTime: 100},
		{Name: "Two", Difficulty: 5, TotalClearTime: 200},
	}}
	html := render(t, StygianOnslaughtSection(multi, 1))
	if !strings.Contains(html, `hx-get="/stygian?cycle=0"`) || !strings.Contains(html, `hx-get="/stygian?cycle=1"`) {
		t.Errorf("expected cycle selector URLs")
	}
	if !strings.Contains(html, "Two") || !strings.Contains(html, "200s") {
		t.Errorf("expected selected cycle content")
	}
	if strings.Count(html, "border-primary/60 bg-primary/10 text-foreground") != 1 {
		t.Errorf("expected exactly one active cycle")
	}
}

func TestStygianEmpty(t *testing.T) {
	html := render(t, StygianOnslaughtSection(nil, 0))
	if !strings.Contains(html, "No Stygian Onslaught records available.") {
		t.Errorf("expected empty state")
	}
}

func TestImaginariumTheaterSection(t *testing.T) {
	html := render(t, ImaginariumTheaterSection(accountTheater))
	for _, want := range []string{
		"Imaginarium Theater",
		"Suli Could Never",
		"2 / 3",
		"Acts Cleared",
		"Medals across acts",
		"Furina",
		"Yanfei",
		"fill-primary",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("ImaginariumTheaterSection missing %q", want)
		}
	}
}

func TestImaginariumTheaterEmpty(t *testing.T) {
	html := render(t, ImaginariumTheaterSection(nil))
	if !strings.Contains(html, "No Imaginarium Theater records available.") {
		t.Errorf("expected empty state")
	}
}
