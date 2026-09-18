package genshin

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dotcchix/gonshin/internal/hoyolab"
)

func embedData() hoyolab.GenshinData {
	return hoyolab.GenshinData{
		Info: &hoyolab.AccountInfo{
			UID: 824677421, Nickname: "Zex.", Level: 60,
			GameHeadIcon: "https://cdn.test/head.png", AchievementNumber: 1362,
			ActiveDayNumber: 1893, TotalCharacters: 103,
		},
		Characters: []hoyolab.AccountCharacter{
			{Name: "Venti", Icon: "https://cdn.test/venti.png", Element: hoyolab.Anemo, Level: 90},
			{Name: "Flins", Icon: "flins.png", Element: hoyolab.Electro, Level: 90},
			{Name: "Mavuika", Icon: "https://cdn.test/mavuika.png", Element: hoyolab.Pyro, Level: 90},
			{Name: "Nefer", Icon: "nefer.png", Element: hoyolab.Dendro, Level: 90},
			{Name: "Varesa", Icon: "https://cdn.test/varesa.png", Element: hoyolab.Electro, Level: 90},
		},
		Abyss:   &hoyolab.AccountAbyss{MaxFloor: "12-3", TotalStar: 36},
		Stygian: &hoyolab.AccountStygian{Cycles: []hoyolab.StygianCycle{{Name: "Season One", Difficulty: 5, TotalClearTime: 180}}},
		Theater: &hoyolab.AccountTheater{Medal: []int{1, 1, 0}, Characters: []hoyolab.TheaterCharacter{{Name: "Furina", Level: 90}}},
	}
}

func findGallery(embed DiscordEmbed) (DiscordMediaGallery, bool) {
	for _, c := range embed.Component.Components {
		if g, ok := c.(DiscordMediaGallery); ok {
			return g, true
		}
	}
	return DiscordMediaGallery{}, false
}

func findActionRow(embed DiscordEmbed) (DiscordActionRow, bool) {
	for _, c := range embed.Component.Components {
		if r, ok := c.(DiscordActionRow); ok {
			return r, true
		}
	}
	return DiscordActionRow{}, false
}

func TestBuildDiscordEmbedRoot(t *testing.T) {
	embed := BuildDiscordEmbed(embedData(), "https://genshin.test/")
	if embed.Component.Type != discordTypeContainer {
		t.Fatalf("root type = %d, want %d", embed.Component.Type, discordTypeContainer)
	}
	if embed.Component.AccentColor != discordAccentColor {
		t.Errorf("accent = %d", embed.Component.AccentColor)
	}
	if len(embed.Component.Components) == 0 {
		t.Fatal("expected at least one component")
	}
}

func TestDiscordEmbedGallery(t *testing.T) {
	gallery, ok := findGallery(BuildDiscordEmbed(embedData(), "https://genshin.test"))
	if !ok {
		t.Fatal("expected media gallery")
	}
	want := []string{"Mavuika", "Varesa", "Nefer", "Flins"}
	if len(gallery.Items) != len(want) {
		t.Fatalf("gallery items = %d, want %d", len(gallery.Items), len(want))
	}
	for i, name := range want {
		item := gallery.Items[i]
		if !strings.HasPrefix(item.Description, name) {
			t.Errorf("item %d description = %q, want prefix %q", i, item.Description, name)
		}
		wantURL := "https://genshin.test/static/genshin/embed/" + name + ".webp"
		if item.Media.URL != wantURL {
			t.Errorf("item %d media = %q, want %q", i, item.Media.URL, wantURL)
		}
	}
}

func TestDiscordEmbedGalleryUsesBundledAssets(t *testing.T) {
	data := embedData()
	data.Characters = nil
	gallery, ok := findGallery(BuildDiscordEmbed(data, "https://genshin.test"))
	if !ok || len(gallery.Items) != 4 {
		t.Fatalf("expected four bundled items, got %+v (ok=%v)", gallery, ok)
	}
	for i, name := range discordGalleryNames {
		if got := gallery.Items[i].Description; got != name {
			t.Errorf("item %d description = %q, want name fallback %q", i, got, name)
		}
	}

	// Without a base URL the media cannot be made absolute, so the gallery is
	// omitted rather than emitting invalid relative URLs.
	if _, ok := findGallery(BuildDiscordEmbed(data, "")); ok {
		t.Errorf("gallery should be omitted without a base URL")
	}
}

func TestDiscordEmbedPlayerSection(t *testing.T) {
	embed := BuildDiscordEmbed(embedData(), "https://genshin.test")
	var content string
	for _, c := range embed.Component.Components {
		if s, ok := c.(DiscordSection); ok {
			content = s.Components[0].(DiscordTextDisplay).Content
		}
	}
	for _, want := range []string{"# Zex.", "**AR** 60", "**UID** 824677421", "**Achievements** 1362", "**Days Active** 1893"} {
		if !strings.Contains(content, want) {
			t.Errorf("player section missing %q: %q", want, content)
		}
	}
}

func TestDiscordEmbedSummary(t *testing.T) {
	embed := BuildDiscordEmbed(embedData(), "https://genshin.test")
	var summary string
	for _, c := range embed.Component.Components {
		if td, ok := c.(DiscordTextDisplay); ok && strings.Contains(td.Content, "Stygian Onslaught") {
			summary = td.Content
		}
	}
	if summary == "" {
		t.Fatal("expected summary text display")
	}
	for _, want := range []string{"Stygian Onslaught", "Season One", "180s", "Spiral Abyss", "12-3", "36★"} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary missing %q: %q", want, summary)
		}
	}
	if strings.Contains(summary, "Theater") || strings.Contains(summary, "Imaginarium") {
		t.Errorf("summary should only show stygian/abyss: %q", summary)
	}
}

func TestDiscordEmbedStygianDifficultyMatchesUI(t *testing.T) {
	cases := []struct {
		difficulty int
		want       string
	}{
		{6, "Dire"},
		{5, "Fearless"},
	}
	for _, tc := range cases {
		data := embedData()
		data.Stygian = &hoyolab.AccountStygian{Cycles: []hoyolab.StygianCycle{
			{Name: "Season One", Difficulty: tc.difficulty, TotalClearTime: 180},
		}}
		embed := BuildDiscordEmbed(data, "https://genshin.test")
		var summary string
		for _, c := range embed.Component.Components {
			if td, ok := c.(DiscordTextDisplay); ok && strings.Contains(td.Content, "Stygian Onslaught") {
				summary = td.Content
			}
		}
		if want := "· " + tc.want; !strings.Contains(summary, want) {
			t.Errorf("difficulty %d summary missing %q: %q", tc.difficulty, want, summary)
		}
	}
}

func TestDiscordEmbedButtons(t *testing.T) {
	row, ok := findActionRow(BuildDiscordEmbed(embedData(), "https://genshin.test"))
	if !ok {
		t.Fatal("expected action row")
	}
	if len(row.Components) != 5 {
		t.Fatalf("buttons = %d, want 5", len(row.Components))
	}
	var labels, urls []string
	for _, c := range row.Components {
		btn, ok := c.(DiscordButton)
		if !ok {
			t.Fatalf("component is %T, want DiscordButton", c)
		}
		if btn.Style != discordButtonStyleLink {
			t.Errorf("button %q style = %d, want link", btn.Label, btn.Style)
		}
		if !strings.HasPrefix(btn.URL, "https://") {
			t.Errorf("button %q url = %q", btn.Label, btn.URL)
		}
		labels = append(labels, btn.Label)
		urls = append(urls, btn.URL)
	}
	if strings.Contains(strings.Join(urls, " "), "?tab=theater") {
		t.Errorf("theater button should be replaced: %v", urls)
	}
	if !strings.Contains(strings.Join(labels, " "), "Akasha.cv") {
		t.Errorf("expected Akasha.cv button: %v", labels)
	}
	if !strings.Contains(strings.Join(urls, " "), "https://akasha.cv/profile/@dotcchix") {
		t.Errorf("expected akasha profile url: %v", urls)
	}
	if got := row.Components[len(row.Components)-1].(DiscordButton).URL; got != "https://enka.network/u/dotcchix/" {
		t.Errorf("enka button url = %q", got)
	}
}

func TestDiscordEmbedNoBaseURL(t *testing.T) {
	embed := BuildDiscordEmbed(embedData(), "")
	if _, ok := findActionRow(embed); ok {
		t.Errorf("action row requires a base URL")
	}
	for _, c := range embed.Component.Components {
		if btn, ok := c.(DiscordButton); ok && btn.URL == "/" {
			t.Errorf("unexpected relative fallback button")
		}
	}
}

func TestDiscordEmbedJSONRules(t *testing.T) {
	for _, base := range []string{"https://genshin.test", ""} {
		raw, err := json.Marshal(BuildDiscordEmbed(embedData(), base))
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		components := collectComponents(t, raw)
		if len(components) > 40 {
			t.Errorf("component count = %d, exceeds 40", len(components))
		}
		for _, c := range components {
			typ, ok := c["type"].(float64)
			if !ok {
				t.Fatalf("component missing numeric type: %v", c)
			}
			if !allowedDiscordType(int(typ)) {
				t.Errorf("component type %d is not allowed", int(typ))
			}
			if int(typ) == discordTypeButton {
				assertButtonKeys(t, c)
			}
		}
	}
}

func TestDiscordEmbedNilData(t *testing.T) {
	embed := BuildDiscordEmbed(hoyolab.GenshinData{}, "https://genshin.test")
	if embed.Component.Type != discordTypeContainer {
		t.Fatalf("root type = %d", embed.Component.Type)
	}
	if _, err := json.Marshal(embed); err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// With no base and no avatar a section would have a null accessory, which is
	// invalid, so a bare text display is emitted instead.
	raw, _ := json.Marshal(BuildDiscordEmbed(hoyolab.GenshinData{}, ""))
	if strings.Contains(string(raw), `"accessory":null`) {
		t.Errorf("unexpected null accessory: %s", raw)
	}
	if !strings.Contains(string(raw), "Genshin Impact Profile") {
		t.Errorf("expected fallback title: %s", raw)
	}
}

func TestAbsMediaURL(t *testing.T) {
	long := "https://cdn.test/" + strings.Repeat("a", discordMaxMediaURLLen)
	cases := []struct {
		name, base, in, want string
	}{
		{"absolute", "https://x.test", "https://cdn.test/a.png", "https://cdn.test/a.png"},
		{"relative", "https://x.test", "a.png", "https://x.test/a.png"},
		{"rooted", "https://x.test", "/a.png", "https://x.test/a.png"},
		{"empty", "https://x.test", "", ""},
		{"no-base", "", "a.png", ""},
		{"too-long", "https://x.test", long, ""},
	}
	for _, tc := range cases {
		if got := absMediaURL(tc.base, tc.in); got != tc.want {
			t.Errorf("%s: absMediaURL = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func allowedDiscordType(t int) bool {
	switch t {
	case discordTypeActionRow, discordTypeButton, discordTypeSection, discordTypeTextDisplay,
		discordTypeThumbnail, discordTypeMediaGallery, discordTypeSeparator, discordTypeContainer:
		return true
	default:
		return false
	}
}

func assertButtonKeys(t *testing.T, c map[string]any) {
	t.Helper()
	for key := range c {
		switch key {
		case "type", "style", "url", "label", "emoji", "disabled":
		default:
			t.Errorf("button has disallowed key %q", key)
		}
	}
}

func collectComponents(t *testing.T, raw []byte) []map[string]any {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return walkComponents(doc)
}

func walkComponents(node any) []map[string]any {
	var out []map[string]any
	switch n := node.(type) {
	case map[string]any:
		if _, ok := n["type"]; ok {
			out = append(out, n)
		}
		for _, key := range []string{"components", "items"} {
			if list, ok := n[key].([]any); ok {
				for _, child := range list {
					out = append(out, walkComponents(child)...)
				}
			}
		}
		if acc, ok := n["accessory"]; ok {
			out = append(out, walkComponents(acc)...)
		}
	case []any:
		for _, child := range n {
			out = append(out, walkComponents(child)...)
		}
	}
	return out
}
