package hoyolab

import (
	"context"
	"errors"
	"sync"
)

// ErrNotFound is returned when a resource (e.g. a character detail) does not
// exist for the account.
var ErrNotFound = errors.New("hoyolab: not found")

// ErrUnavailable is returned when the upstream data is not cached yet and the
// API responds with 503.
var ErrUnavailable = errors.New("hoyolab: data unavailable")

// Provider is the migration seam between the web app and the Hoyolab data
// source. Today HTTPProvider proxies the Python FastAPI service; a future
// native Go port only needs to implement this interface.
type Provider interface {
	// Info returns the flattened player profile.
	Info(ctx context.Context) (*AccountInfo, error)
	// Characters returns the owned character roster.
	Characters(ctx context.Context) ([]AccountCharacter, error)
	// Abyss returns the spiral abyss summary.
	Abyss(ctx context.Context) (*AccountAbyss, error)
	// Theater returns the current imaginarium theater season.
	Theater(ctx context.Context) (*AccountTheater, error)
	// Stygian returns stygian onslaught cycles enriched with character builds.
	Stygian(ctx context.Context) (*AccountStygian, error)
	// CharacterDetail returns the detailed build for a single character.
	CharacterDetail(ctx context.Context, id int) (*CharacterDetail, error)
}

// All loads every cached dataset in parallel, mirroring the original page's
// Promise.all aggregation. Per-endpoint failures yield nil/empty values rather
// than failing the whole page.
func All(ctx context.Context, p Provider) GenshinData {
	var (
		wg         sync.WaitGroup
		info       *AccountInfo
		characters []AccountCharacter
		abyss      *AccountAbyss
		theater    *AccountTheater
		stygian    *AccountStygian
	)

	wg.Add(5)
	go func() { defer wg.Done(); info, _ = p.Info(ctx) }()
	go func() { defer wg.Done(); characters, _ = p.Characters(ctx) }()
	go func() { defer wg.Done(); abyss, _ = p.Abyss(ctx) }()
	go func() { defer wg.Done(); theater, _ = p.Theater(ctx) }()
	go func() { defer wg.Done(); stygian, _ = p.Stygian(ctx) }()
	wg.Wait()

	if characters == nil {
		characters = []AccountCharacter{}
	}
	return GenshinData{
		Info:       info,
		Characters: characters,
		Abyss:      abyss,
		Theater:    theater,
		Stygian:    stygian,
	}
}
