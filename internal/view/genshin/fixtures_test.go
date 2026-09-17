package genshin

import "github.com/dotcchix/gonshin/internal/hoyolab"

func intPtr(v int) *int { return &v }

var (
	accountInfo = hoyolab.AccountInfo{
		UID:               824677421,
		Nickname:          "Zex.",
		Level:             60,
		GameHeadIcon:      "https://example.test/head.png",
		ActiveDayNumber:   1893,
		AchievementNumber: 1362,
		TotalCharacters:   103,
		TotalFriendship:   72,
	}

	accountCharacters = []hoyolab.AccountCharacter{
		{
			ID: 10000022, Icon: "https://example.test/venti.png", Name: "Venti",
			Element: hoyolab.Anemo, Level: 90, Rarity: intPtr(5), ActivedConstellationNum: 0,
			Weapon: hoyolab.Weapon{Name: "Windblume Ode", Level: 80, Icon: "https://example.test/w1.png"},
		},
		{
			ID: 10000089, Icon: "https://example.test/furina.png", Name: "Furina",
			Element: hoyolab.Hydro, Level: 90, Rarity: intPtr(5), ActivedConstellationNum: 2,
			Weapon: hoyolab.Weapon{Name: "Favonius Sword", Level: 90, Icon: "https://example.test/w2.png"},
		},
		{
			ID: 10000058, Icon: "https://example.test/yanfei.png", Name: "Yanfei",
			Element: hoyolab.Pyro, Level: 80, Rarity: intPtr(4), ActivedConstellationNum: 6,
			Weapon: hoyolab.Weapon{Name: "The Widsith", Level: 80, Icon: "https://example.test/w3.png"},
		},
	}

	accountAbyss = &hoyolab.AccountAbyss{
		StartTime: "1789502400", EndTime: "1792094399",
		TotalBattle: 33, TotalWin: 19, MaxFloor: "12-3", TotalStar: 36,
		Floors: []hoyolab.AbyssFloor{{
			MaxStar: 9,
			Levels: []hoyolab.AbyssLevel{{
				Chamber: 1, Star: 3,
				Battles: [][]hoyolab.AbyssBattleCharacter{
					{{Icon: "https://example.test/venti.png", Level: 90}, {Icon: "https://example.test/furina.png", Level: 90}},
					{{Icon: "https://example.test/unknown.png", Level: 70}},
				},
			}},
		}},
	}

	accountTheater = &hoyolab.AccountTheater{
		Medal: []int{1, 1, 0},
		Characters: []hoyolab.TheaterCharacter{
			{Avatar: "https://example.test/furina.png", Name: "Furina", Level: 90, Element: hoyolab.Hydro},
			{Avatar: "https://example.test/yanfei.png", Name: "Yanfei", Level: 80, Element: hoyolab.Pyro},
		},
	}

	accountStygian = &hoyolab.AccountStygian{
		UID: 824677421,
		Cycles: []hoyolab.StygianCycle{{
			ScheduleID: "s1", Name: "Season One", StartTime: "1", EndTime: "2",
			Difficulty: 5, TotalClearTime: 180,
			Challenges: []hoyolab.StygianChallenge{{
				Name: "First Half", Second: 90,
				Characters: []hoyolab.StygianCharacter{{
					Icon: "https://example.test/furina.png", Name: "Furina", Element: hoyolab.Hydro,
					Level: 90, ActivedConstellationNum: 2,
					Weapon:       &hoyolab.Weapon{Name: "Favonius Sword", Level: 90, Icon: "https://example.test/w2.png", Refine: intPtr(5)},
					FinalStats:   []hoyolab.FinalStat{{Name: "Max HP", Value: "40000"}, {Name: "ATK", Value: "1200"}, {Name: "Elemental Mastery", Value: "80"}, {Name: "CRIT Rate", Value: "60.0%"}, {Name: "CRIT DMG", Value: "130.0%"}, {Name: "Energy Recharge", Value: "180.0%"}, {Name: "Hydro DMG Bonus", Value: "75.0%"}},
					ArtifactSets: []string{"Golden Troupe x4", "Marechaussee Hunter x1"},
				}},
			}},
		}},
	}

	characterDetail = &hoyolab.CharacterDetail{
		ID: 10000022, Icon: "https://example.test/venti.png", Name: "Venti",
		Element: hoyolab.Anemo, Level: 90, Rarity: 5, ActivedConstellationNum: 0,
		Weapon:       &hoyolab.Weapon{Name: "Windblume Ode", Level: 80, Icon: "https://example.test/w1.png", Refine: intPtr(5)},
		FinalStats:   []hoyolab.FinalStat{{Name: "Max HP", Value: "11739"}, {Name: "Anemo DMG Bonus", Value: "61.6%"}},
		ArtifactSets: []string{"Viridescent Venerer x4", "Emblem of Severed Fate x1"},
		Constellations: []hoyolab.Constellation{
			{ID: 221, Name: "Splitting Gales", Icon: "https://example.test/c1.png", Effect: "Fires 2 additional arrows per <color=#FFD780FF>Aimed Shot</color>.", IsActived: true, Pos: 1},
			{ID: 222, Name: "Breeze of Reminiscence", Icon: "https://example.test/c2.png", Effect: "Decreases RES.", IsActived: false, Pos: 2},
		},
	}

	genshinData = hoyolab.GenshinData{
		Info:       &accountInfo,
		Characters: accountCharacters,
		Abyss:      accountAbyss,
		Theater:    accountTheater,
		Stygian:    accountStygian,
	}
)
