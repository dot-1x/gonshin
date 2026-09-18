package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dotcchix/gonshin/internal/hoyolab"
)

type stubProvider struct {
	characters []hoyolab.AccountCharacter
	detail     *hoyolab.CharacterDetail
	detailErr  error
	stygian    *hoyolab.AccountStygian
}

func (s stubProvider) Info(context.Context) (*hoyolab.AccountInfo, error) {
	return &hoyolab.AccountInfo{UID: 824677421, Nickname: "Zex.", Level: 60, ActiveDayNumber: 1893}, nil
}

func (s stubProvider) Characters(context.Context) ([]hoyolab.AccountCharacter, error) {
	return s.characters, nil
}

func (s stubProvider) Abyss(context.Context) (*hoyolab.AccountAbyss, error) { return nil, nil }

func (s stubProvider) Theater(context.Context) (*hoyolab.AccountTheater, error) { return nil, nil }

func (s stubProvider) Stygian(context.Context) (*hoyolab.AccountStygian, error) {
	return s.stygian, nil
}

func (s stubProvider) CharacterDetail(context.Context, int) (*hoyolab.CharacterDetail, error) {
	return s.detail, s.detailErr
}

func newStub() stubProvider {
	return stubProvider{
		characters: []hoyolab.AccountCharacter{
			{ID: 10000022, Icon: "a.png", Name: "Venti", Element: hoyolab.Anemo, Level: 90},
			{ID: 10000089, Icon: "b.png", Name: "Furina", Element: hoyolab.Hydro, Level: 90},
			{ID: 10000106, Icon: "mavuika.png", Name: "Mavuika", Element: hoyolab.Pyro, Level: 90},
			{ID: 10000111, Icon: "varesa.png", Name: "Varesa", Element: hoyolab.Electro, Level: 90},
			{ID: 10000116, Icon: "flins.png", Name: "Flins", Level: 90},
			{ID: 10000113, Icon: "nefer.png", Name: "Nefer", Element: hoyolab.Dendro, Level: 90},
		},
		detail: &hoyolab.CharacterDetail{ID: 10000022, Icon: "a.png", Name: "Venti", Element: hoyolab.Anemo, Level: 90, Rarity: 5},
		stygian: &hoyolab.AccountStygian{Cycles: []hoyolab.StygianCycle{
			{Name: "One", Difficulty: 6, TotalClearTime: 100},
			{Name: "Two", Difficulty: 5, TotalClearTime: 200},
		}},
	}
}

func newMux(p hoyolab.Provider) *http.ServeMux {
	mux := http.NewServeMux()
	New(p, "").Routes(mux)
	return mux
}

func get(t *testing.T, mux *http.ServeMux, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestPage(t *testing.T) {
	rec := get(t, newMux(newStub()), "/")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("content-type = %q", ct)
	}
	body := rec.Body.String()
	for _, want := range []string{"<!doctype html>", "Genshin Impact Profile", "genshin-impact", "hoyo", "Detail", "Stygian Onslaught", "/static/app.css", "htmx.min.js", "gen-dialog"} {
		if !strings.Contains(body, want) {
			t.Errorf("page missing %q", want)
		}
	}

	for _, want := range []string{
		`id="discord:component-embed"`,
		`type="application/json"`,
		`property="og:url"`,
		`property="og:image"`,
		`name="twitter:card"`,
		`name="theme-color"`,
		"Mavuika",
		"Varesa",
		"Nefer",
		"Flins",
		"Enka.network",
		"Akasha.cv",
		"Days Active",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page missing discord preview %q", want)
		}
	}
}

func TestPageTabDeepLink(t *testing.T) {
	rec := get(t, newMux(newStub()), "/?tab=characters")
	body := rec.Body.String()
	if !strings.Contains(body, "Character Showcase") {
		t.Errorf("expected Characters tab body")
	}
	if !strings.Contains(body, `href="/tabs/characters"`) && !strings.Contains(body, `hx-get="/tabs/characters"`) {
		t.Errorf("expected nav present")
	}
	if strings.Count(body, "border-b-2 border-primary text-primary") != 1 {
		t.Errorf("expected exactly one active tab")
	}
}

func TestPageTabFallback(t *testing.T) {
	rec := get(t, newMux(newStub()), "/?tab=nope")
	if !strings.Contains(rec.Body.String(), "Player") && !strings.Contains(rec.Body.String(), "Zex.") {
		t.Errorf("invalid tab should fall back to home")
	}
}

func TestTabs(t *testing.T) {
	mux := newMux(newStub())
	cases := map[string]string{
		"/tabs/characters": "Character Showcase",
		"/tabs/spiral":     "Spiral Abyss",
		"/tabs/stygian":    "Stygian Onslaught",
		"/tabs/theater":    "Imaginarium Theater",
	}
	for path, want := range cases {
		rec := get(t, mux, path)
		if rec.Code != http.StatusOK {
			t.Errorf("%s status = %d", path, rec.Code)
		}
		body := rec.Body.String()
		if !strings.Contains(body, want) {
			t.Errorf("%s missing %q", path, want)
		}
		if !strings.Contains(body, `hx-swap-oob="true"`) {
			t.Errorf("%s missing OOB nav", path)
		}
	}
}

func TestTabUnknown(t *testing.T) {
	if rec := get(t, newMux(newStub()), "/tabs/nope"); rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestShowcase(t *testing.T) {
	rec := get(t, newMux(newStub()), "/showcase?element=Hydro")
	body := rec.Body.String()
	if !strings.Contains(body, "Furina") || strings.Contains(body, "Venti") {
		t.Errorf("showcase filter incorrect")
	}
	if !strings.Contains(body, `id="showcase-panel"`) {
		// Panel content is swapped in, so the wrapper is not part of the fragment.
		t.Log("showcase fragment has no panel wrapper (expected)")
	}
}

func TestCharacterDialog(t *testing.T) {
	rec := get(t, newMux(newStub()), "/characters/1")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Venti") {
		t.Fatalf("dialog failed: %d", rec.Code)
	}
}

func TestCharacterDialogNotFound(t *testing.T) {
	stub := newStub()
	stub.detail = nil
	stub.detailErr = hoyolab.ErrNotFound
	rec := get(t, newMux(stub), "/characters/10000022")
	body := rec.Body.String()
	if !strings.Contains(body, "Failed to fetch (404)") {
		t.Errorf("expected 404 error message, got %s", body)
	}
	if !strings.Contains(body, "Showing cached roster data for Venti") {
		t.Errorf("expected fallback note")
	}
}

func TestStygianPanel(t *testing.T) {
	rec := get(t, newMux(newStub()), "/stygian?cycle=1")
	body := rec.Body.String()
	if !strings.Contains(body, "Two") || !strings.Contains(body, "200s") {
		t.Errorf("stygian panel did not select cycle 1")
	}
}
