// Package hoyolab models the Hoyolab data API and provides a Provider
// abstraction so the Python FastAPI backend can later be swapped for a native
// Go implementation without touching the HTTP handlers or views.
package hoyolab

// Element is a Genshin elemental type.
type Element string

const (
	Pyro    Element = "Pyro"
	Hydro   Element = "Hydro"
	Electro Element = "Electro"
	Cryo    Element = "Cryo"
	Anemo   Element = "Anemo"
	Geo     Element = "Geo"
	Dendro  Element = "Dendro"
)

// AccountInfo is the flattened player profile (/api/gi/info).
type AccountInfo struct {
	UID               int    `json:"uid"`
	Nickname          string `json:"nickname"`
	Level             int    `json:"level"`
	GameHeadIcon      string `json:"game_head_icon"`
	ActiveDayNumber   int    `json:"active_day_number"`
	AchievementNumber int    `json:"achievement_number"`
	TotalCharacters   int    `json:"total_characters"`
	TotalFriendship   int    `json:"total_friendship"`
}

// Weapon is a character's equipped weapon.
type Weapon struct {
	Name   string `json:"name"`
	Level  int    `json:"level"`
	Icon   string `json:"icon"`
	Refine *int   `json:"refine,omitempty"`
}

// AccountCharacter is a roster entry (/api/gi/characters).
type AccountCharacter struct {
	ID                      int     `json:"id"`
	Icon                    string  `json:"icon"`
	Name                    string  `json:"name"`
	Element                 Element `json:"element"`
	Level                   int     `json:"level"`
	Rarity                  *int    `json:"rarity,omitempty"`
	ActivedConstellationNum int     `json:"actived_constellation_num"`
	Weapon                  Weapon  `json:"weapon"`
}

// FinalStat is a single resolved character stat.
type FinalStat struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Constellation is a character constellation.
type Constellation struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Effect    string `json:"effect"`
	IsActived bool   `json:"is_actived"`
	Pos       int    `json:"pos"`
}

// CharacterDetail is the detailed build (/api/gi/characters/{id}).
type CharacterDetail struct {
	ID                      int             `json:"id"`
	Icon                    string          `json:"icon"`
	Name                    string          `json:"name"`
	Element                 Element         `json:"element"`
	Level                   int             `json:"level"`
	Rarity                  int             `json:"rarity"`
	ActivedConstellationNum int             `json:"actived_constellation_num"`
	Weapon                  *Weapon         `json:"weapon"`
	FinalStats              []FinalStat     `json:"final_stats"`
	ArtifactSets            []string        `json:"artifact_sets"`
	Constellations          []Constellation `json:"constellations"`
}

// AbyssBattleCharacter is a lean `{icon, level}` abyss roster reference.
type AbyssBattleCharacter struct {
	Icon  string `json:"icon"`
	Level int    `json:"level"`
}

// AbyssLevel is a single chamber.
type AbyssLevel struct {
	Chamber int                      `json:"chamber"`
	Star    int                      `json:"star"`
	Battles [][]AbyssBattleCharacter `json:"battles"`
}

// AbyssFloor is a spiral abyss floor (the API returns floor 12 only).
type AbyssFloor struct {
	MaxStar int          `json:"max_star"`
	Levels  []AbyssLevel `json:"levels"`
}

// AccountAbyss is the spiral abyss summary (/api/gi/abyss).
type AccountAbyss struct {
	StartTime   string       `json:"start_time"`
	EndTime     string       `json:"end_time"`
	TotalBattle int          `json:"total_battle"`
	TotalWin    int          `json:"total_win"`
	MaxFloor    string       `json:"max_floor"`
	TotalStar   int          `json:"total_star"`
	Floors      []AbyssFloor `json:"floors"`
}

// TheaterCharacter is an imaginarium theater roster entry.
type TheaterCharacter struct {
	Avatar  string  `json:"avatar"`
	Name    string  `json:"name"`
	Level   int     `json:"level"`
	Element Element `json:"element"`
}

// AccountTheater is a theater season (/api/gi/theater).
type AccountTheater struct {
	Medal      []int              `json:"medal"`
	Characters []TheaterCharacter `json:"characters"`
}

// StygianCharacter is a stygian team member, either a full build or a lean
// `{icon, name, element, level}` fallback.
type StygianCharacter struct {
	Icon                    string      `json:"icon"`
	Name                    string      `json:"name"`
	Element                 Element     `json:"element"`
	Level                   int         `json:"level"`
	ActivedConstellationNum int         `json:"actived_constellation_num"`
	Weapon                  *Weapon     `json:"weapon"`
	FinalStats              []FinalStat `json:"final_stats"`
	ArtifactSets            []string    `json:"artifact_sets"`
}

// StygianChallenge is a single stygian challenge with its clearing team.
type StygianChallenge struct {
	Name       string             `json:"name"`
	Second     int                `json:"second"`
	Characters []StygianCharacter `json:"characters"`
}

// StygianCycle is one stygian onslaught cycle.
type StygianCycle struct {
	ScheduleID     string             `json:"schedule_id"`
	Name           string             `json:"name"`
	StartTime      string             `json:"start_time"`
	EndTime        string             `json:"end_time"`
	Difficulty     int                `json:"difficulty"`
	TotalClearTime int                `json:"total_clear_time"`
	Challenges     []StygianChallenge `json:"challenges"`
}

// AccountStygian is the stygian detail payload (/api/gi/stygian/detail).
type AccountStygian struct {
	UID    int            `json:"uid"`
	Cycles []StygianCycle `json:"cycles"`
}

// DisplayCharacter is a character resolved to a display-ready shape.
type DisplayCharacter struct {
	Icon          string
	Name          string
	Element       *Element
	Level         int
	Constellation *int
	Weapon        *Weapon
	FinalStats    []FinalStat
	ArtifactSets  []string
}

// DisplayTeam is a team rendered by the team card.
type DisplayTeam struct {
	Name       string
	Characters []DisplayCharacter
	ClearTime  *int
}

// GenshinData aggregates everything the page needs.
type GenshinData struct {
	Info       *AccountInfo
	Characters []AccountCharacter
	Abyss      *AccountAbyss
	Theater    *AccountTheater
	Stygian    *AccountStygian
}
