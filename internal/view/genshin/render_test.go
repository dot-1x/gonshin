package genshin

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
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
