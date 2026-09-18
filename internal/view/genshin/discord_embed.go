package genshin

import (
	"strconv"
	"strings"

	"github.com/dotcchix/gonshin/internal/hoyolab"
)

// Discord component embeds are a read-only subset of message components. These
// type ids mirror the component reference:
// https://docs.discord.com/developers/components/reference
const (
	discordTypeActionRow    = 1
	discordTypeButton       = 2
	discordTypeSection      = 9
	discordTypeTextDisplay  = 10
	discordTypeThumbnail    = 11
	discordTypeMediaGallery = 12
	discordTypeSeparator    = 14
	discordTypeContainer    = 17
)

// discordButtonStyleLink is the only button style allowed in a component embed.
const discordButtonStyleLink = 5

// discordAccentColor is the theme primary (#306998) as a decimal colour.
const discordAccentColor = 3172760

// discordMaxMediaURLLen is Discord's unfurled media URL limit.
const discordMaxMediaURLLen = 2048

// discordGalleryNames is the fixed set of characters shown in the embed
// gallery, in display order. Each maps to a bundled webp asset under
// /static/genshin/embed.
var discordGalleryNames = []string{"Mavuika", "Varesa", "Nefer", "Flins"}

// DiscordEmbed is the document written into the discord:component-embed script.
type DiscordEmbed struct {
	Component DiscordContainer `json:"component"`
}

// DiscordContainer is component type 17, the parent component of the embed.
type DiscordContainer struct {
	Type        int   `json:"type"`
	AccentColor int   `json:"accent_color,omitempty"`
	Spoiler     bool  `json:"spoiler,omitempty"`
	Components  []any `json:"components"`
}

// DiscordTextDisplay is component type 10, a block of Discord markdown.
type DiscordTextDisplay struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// DiscordSeparator is component type 14.
type DiscordSeparator struct {
	Type    int  `json:"type"`
	Spacing *int `json:"spacing,omitempty"`
}

// DiscordUnfurledMedia is the media item accepted by thumbnails and galleries;
// only url may be set in a component embed.
type DiscordUnfurledMedia struct {
	URL string `json:"url"`
}

// DiscordThumbnail is component type 11.
type DiscordThumbnail struct {
	Type    int                  `json:"type"`
	Media   DiscordUnfurledMedia `json:"media"`
	Spoiler bool                 `json:"spoiler,omitempty"`
}

// DiscordGalleryItem is one entry in a media gallery.
type DiscordGalleryItem struct {
	Media       DiscordUnfurledMedia `json:"media"`
	Description string               `json:"description,omitempty"`
	Spoiler     bool                 `json:"spoiler,omitempty"`
}

// DiscordMediaGallery is component type 12.
type DiscordMediaGallery struct {
	Type  int                  `json:"type"`
	Items []DiscordGalleryItem `json:"items"`
}

// DiscordSection is component type 9: text components plus an accessory.
type DiscordSection struct {
	Type       int   `json:"type"`
	Components []any `json:"components"`
	Accessory  any   `json:"accessory"`
}

// DiscordButton is component type 2. Component embeds only allow link buttons.
type DiscordButton struct {
	Type  int    `json:"type"`
	Style int    `json:"style"`
	URL   string `json:"url"`
	Label string `json:"label,omitempty"`
}

// DiscordActionRow is component type 1.
type DiscordActionRow struct {
	Type       int   `json:"type"`
	Components []any `json:"components"`
}

// BuildDiscordEmbed builds the component-embed payload for the profile page.
// baseURL is the public origin used for links and relative media; when empty
// the link buttons are omitted.
func BuildDiscordEmbed(data hoyolab.GenshinData, baseURL string) DiscordEmbed {
	base := strings.TrimRight(baseURL, "/")

	container := DiscordContainer{
		Type:        discordTypeContainer,
		AccentColor: discordAccentColor,
	}

	container.Components = append(container.Components, playerSection(data.Info, base))
	if gallery, ok := characterGallery(data.Characters, base); ok {
		container.Components = append(container.Components, gallery)
	}
	if summary := progressText(data.Abyss, data.Stygian); summary != "" {
		container.Components = append(container.Components,
			DiscordSeparator{Type: discordTypeSeparator},
			DiscordTextDisplay{Type: discordTypeTextDisplay, Content: summary},
		)
	}
	if row, ok := linkButtons(base); ok {
		container.Components = append(container.Components, row)
	}

	return DiscordEmbed{Component: container}
}

// playerSection renders the nickname and headline stats with an avatar
// thumbnail, falling back to a link button when no avatar is available. A
// section requires an accessory, so a bare text display is returned when
// neither an absolute avatar nor a base URL is available.
func playerSection(info *hoyolab.AccountInfo, base string) any {
	nickname := "Genshin Impact Profile"
	var stats []string
	if info != nil {
		if info.Nickname != "" {
			nickname = info.Nickname
		}
		if info.Level != 0 {
			stats = append(stats, "**AR** "+strconv.Itoa(info.Level))
		}
		if info.UID != 0 {
			stats = append(stats, "**UID** "+strconv.Itoa(info.UID))
		}
		if info.AchievementNumber != 0 {
			stats = append(stats, "**Achievements** "+strconv.Itoa(info.AchievementNumber))
		}
		if info.ActiveDayNumber != 0 {
			stats = append(stats, "**Days Active** "+strconv.Itoa(info.ActiveDayNumber))
		}
		if info.TotalCharacters != 0 {
			stats = append(stats, "**Characters** "+strconv.Itoa(info.TotalCharacters))
		}
	}

	content := "# " + escapeDiscordMarkdown(nickname)
	if len(stats) > 0 {
		content += "\n" + strings.Join(stats, " · ")
	}

	text := DiscordTextDisplay{Type: discordTypeTextDisplay, Content: content}
	section := DiscordSection{
		Type:       discordTypeSection,
		Components: []any{text},
	}

	if url := absMediaURL(base, headIcon(info)); url != "" {
		section.Accessory = DiscordThumbnail{
			Type:  discordTypeThumbnail,
			Media: DiscordUnfurledMedia{URL: url},
		}
		return section
	}
	if base != "" {
		section.Accessory = DiscordButton{
			Type:  discordTypeButton,
			Style: discordButtonStyleLink,
			URL:   base + "/",
			Label: "Open profile",
		}
		return section
	}
	return text
}

// characterGallery renders the fixed gallery names from bundled webp assets,
// enriching the description with roster data when the character is owned.
func characterGallery(characters []hoyolab.AccountCharacter, base string) (DiscordMediaGallery, bool) {
	byName := make(map[string]hoyolab.AccountCharacter, len(characters))
	for _, c := range characters {
		byName[c.Name] = c
	}

	var items []DiscordGalleryItem
	for _, name := range discordGalleryNames {
		url := absMediaURL(base, discordGalleryImage(name))
		if url == "" {
			continue
		}
		description := name
		if c, ok := byName[name]; ok {
			description = galleryDescription(c)
		}
		items = append(items, DiscordGalleryItem{
			Media:       DiscordUnfurledMedia{URL: url},
			Description: description,
		})
	}
	if len(items) == 0 {
		return DiscordMediaGallery{}, false
	}
	return DiscordMediaGallery{Type: discordTypeMediaGallery, Items: items}, true
}

// discordGalleryImage is the bundled asset path for a gallery character.
func discordGalleryImage(name string) string {
	return "/static/genshin/embed/" + name + ".webp"
}

func galleryDescription(c hoyolab.AccountCharacter) string {
	if c.Element == "" {
		return c.Name
	}
	return c.Name + " · " + string(c.Element) + " · " + lv(c.Level)
}

// progressText summarizes Stygian Onslaught and Spiral Abyss only.
func progressText(abyss *hoyolab.AccountAbyss, stygian *hoyolab.AccountStygian) string {
	var lines []string

	if stygian != nil && len(stygian.Cycles) > 0 {
		cycle := stygian.Cycles[0]
		line := "**Stygian Onslaught**"
		if cycle.Name != "" {
			line += " · " + escapeDiscordMarkdown(cycle.Name)
		}
		if cycle.Difficulty > 0 {
			line += " · " + StygianDifficultyLabel(cycle.Difficulty)
		}
		if cycle.TotalClearTime > 0 {
			line += " · " + FormatClearTimeN(cycle.TotalClearTime)
		}
		lines = append(lines, line)
	}

	if abyss != nil && (abyss.MaxFloor != "" || abyss.TotalStar > 0) {
		line := "**Spiral Abyss**"
		if abyss.MaxFloor != "" {
			line += " · " + escapeDiscordMarkdown(abyss.MaxFloor)
		}
		if abyss.TotalStar > 0 {
			line += " · " + strconv.Itoa(abyss.TotalStar) + "★"
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// linkButtons builds the action row of deep links plus the Enka profile.
func linkButtons(base string) (DiscordActionRow, bool) {
	if base == "" {
		return DiscordActionRow{}, false
	}
	return DiscordActionRow{
		Type: discordTypeActionRow,
		Components: []any{
			DiscordButton{Type: discordTypeButton, Style: discordButtonStyleLink, URL: base + "/?tab=characters", Label: "Characters"},
			DiscordButton{Type: discordTypeButton, Style: discordButtonStyleLink, URL: base + "/?tab=spiral", Label: "Spiral Abyss"},
			DiscordButton{Type: discordTypeButton, Style: discordButtonStyleLink, URL: base + "/?tab=stygian", Label: "Stygian"},
			DiscordButton{Type: discordTypeButton, Style: discordButtonStyleLink, URL: "https://akasha.cv/profile/@dotcchix", Label: "Akasha.cv"},
			DiscordButton{Type: discordTypeButton, Style: discordButtonStyleLink, URL: "https://enka.network/u/dotcchix/", Label: "Enka.network"},
		},
	}, true
}

func headIcon(info *hoyolab.AccountInfo) string {
	if info == nil {
		return ""
	}
	return info.GameHeadIcon
}

// absMediaURL resolves media against base, enforcing the scheme and length
// rules Discord applies to unfurled media.
func absMediaURL(base, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || len(raw) > discordMaxMediaURLLen {
		return ""
	}
	if strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "http://") {
		return raw
	}
	if base == "" {
		return ""
	}
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	return base + raw
}

// escapeDiscordMarkdown neutralises markdown syntax in user-supplied text.
func escapeDiscordMarkdown(s string) string {
	return strings.NewReplacer(
		"\\", "\\\\",
		"*", "\\*",
		"_", "\\_",
		"~", "\\~",
		"`", "\\`",
		"|", "\\|",
		"[", "\\[",
		"]", "\\]",
	).Replace(s)
}
