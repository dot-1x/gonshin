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
