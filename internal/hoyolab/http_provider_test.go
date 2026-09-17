package hoyolab

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

const (
	infoJSON = `{"uid":1,"nickname":"Zex.","level":60,"game_head_icon":"icon.png",` +
		`"active_day_number":1893,"achievement_number":1362,"total_characters":103,"total_friendship":72}`
	charactersJSON = `[{"id":10000022,"icon":"a.png","name":"Venti","element":"Anemo","level":90,` +
		`"rarity":5,"actived_constellation_num":0,` +
		`"weapon":{"name":"Windblume Ode","level":80,"icon":"w.png"}}]`
	abyssJSON = `{"start_time":"1","end_time":"2","total_battle":33,"total_win":19,` +
		`"max_floor":"12-3","total_star":36,` +
		`"floors":[{"max_star":9,"levels":[{"chamber":1,"star":3,"battles":[[{"icon":"a.png","level":90}]]}]}]}`
	theaterJSON = `[{"medal":[1,1,0],"characters":[{"avatar":"t.png","name":"Cyno","level":90,"element":"Electro"}]},` +
		`{"medal":[],"characters":[]}]`
	stygianJSON = `{"uid":1,"cycles":[{"schedule_id":"s1","name":"Cycle","start_time":"1","end_time":"2",` +
		`"difficulty":5,"total_clear_time":180,"challenges":[{"name":"First Half","second":90,` +
		`"characters":[{"icon":"s.png","name":"Furina","element":"Hydro","level":90,` +
		`"actived_constellation_num":2,"weapon":{"name":"Favonius","level":90,"icon":"w.png","refine":5},` +
		`"final_stats":[{"name":"Max HP","value":"40000"}],"artifact_sets":["Golden Troupe x4"]}]}]}]}`
	detailJSON = `{"id":10000022,"icon":"a.png","name":"Venti","element":"Anemo","level":90,"rarity":5,` +
		`"actived_constellation_num":0,"weapon":{"name":"Windblume Ode","level":80,"icon":"w.png","refine":5},` +
		`"final_stats":[{"name":"Max HP","value":"11739"}],"artifact_sets":["Viridescent Venerer x4"],` +
		`"constellations":[{"id":221,"name":"Splitting Gales","icon":"c.png","effect":"x","is_actived":true,"pos":1}]}`
)

func newTestServer(t *testing.T, hits *int32) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	write := func(body string) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			atomic.AddInt32(hits, 1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(body))
		}
	}
	mux.HandleFunc("GET /api/gi/info", write(infoJSON))
	mux.HandleFunc("GET /api/gi/characters", write(charactersJSON))
	mux.HandleFunc("GET /api/gi/abyss", write(abyssJSON))
	mux.HandleFunc("GET /api/gi/theater", write(theaterJSON))
	mux.HandleFunc("GET /api/gi/stygian/detail", write(stygianJSON))
	mux.HandleFunc("GET /api/gi/characters/10000022", write(detailJSON))
	mux.HandleFunc("GET /api/gi/characters/999", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	mux.HandleFunc("GET /api/gi/characters/503", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	return httptest.NewServer(mux)
}

func TestHTTPProvider(t *testing.T) {
	var hits int32
	srv := newTestServer(t, &hits)
	defer srv.Close()

	p := NewHTTPProvider(srv.URL, time.Minute)
	ctx := context.Background()

	info, err := p.Info(ctx)
	if err != nil || info == nil || info.Nickname != "Zex." || info.Level != 60 {
		t.Fatalf("Info: %+v, %v", info, err)
	}

	chars, err := p.Characters(ctx)
	if err != nil || len(chars) != 1 || chars[0].Name != "Venti" || chars[0].Weapon.Name != "Windblume Ode" {
		t.Fatalf("Characters: %+v, %v", chars, err)
	}

	abyss, err := p.Abyss(ctx)
	if err != nil || abyss == nil || abyss.MaxFloor != "12-3" || len(abyss.Floors) != 1 {
		t.Fatalf("Abyss: %+v, %v", abyss, err)
	}

	theater, err := p.Theater(ctx)
	if err != nil || theater == nil || len(theater.Characters) != 1 || theater.Characters[0].Name != "Cyno" {
		t.Fatalf("Theater: %+v, %v", theater, err)
	}

	stygian, err := p.Stygian(ctx)
	if err != nil || stygian == nil || len(stygian.Cycles) != 1 {
		t.Fatalf("Stygian: %+v, %v", stygian, err)
	}
	if got := stygian.Cycles[0].Challenges[0].Characters[0]; got.FinalStats[0].Name != "Max HP" {
		t.Fatalf("Stygian character: %+v", got)
	}

	detail, err := p.CharacterDetail(ctx, 10000022)
	if err != nil || detail == nil || detail.Constellations[0].Name != "Splitting Gales" {
		t.Fatalf("CharacterDetail: %+v, %v", detail, err)
	}

	if _, err := p.CharacterDetail(ctx, 999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if _, err := p.CharacterDetail(ctx, 503); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected ErrUnavailable, got %v", err)
	}
}

func TestHTTPProviderCachesListEndpoints(t *testing.T) {
	var hits int32
	srv := newTestServer(t, &hits)
	defer srv.Close()

	p := NewHTTPProvider(srv.URL, time.Minute)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		if _, err := p.Info(ctx); err != nil {
			t.Fatal(err)
		}
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("expected 1 upstream hit, got %d", got)
	}

	// Detail must not be cached.
	for i := 0; i < 3; i++ {
		if _, err := p.CharacterDetail(ctx, 10000022); err != nil {
			t.Fatal(err)
		}
	}
	if got := atomic.LoadInt32(&hits); got != 4 {
		t.Fatalf("expected 4 upstream hits, got %d", got)
	}
}

func TestAllAggregates(t *testing.T) {
	var hits int32
	srv := newTestServer(t, &hits)
	defer srv.Close()

	p := NewHTTPProvider(srv.URL, time.Minute)
	data := All(context.Background(), p)

	if data.Info == nil || data.Abyss == nil || data.Theater == nil || data.Stygian == nil {
		t.Fatalf("incomplete aggregate: %+v", data)
	}
	if len(data.Characters) != 1 {
		t.Fatalf("expected 1 character, got %d", len(data.Characters))
	}
}

type errProvider struct{}

func (errProvider) Info(context.Context) (*AccountInfo, error) { return nil, errors.New("boom") }
func (errProvider) Characters(context.Context) ([]AccountCharacter, error) {
	return nil, errors.New("boom")
}
func (errProvider) Abyss(context.Context) (*AccountAbyss, error)     { return nil, errors.New("boom") }
func (errProvider) Theater(context.Context) (*AccountTheater, error) { return nil, errors.New("boom") }
func (errProvider) Stygian(context.Context) (*AccountStygian, error) { return nil, errors.New("boom") }
func (errProvider) CharacterDetail(context.Context, int) (*CharacterDetail, error) {
	return nil, errors.New("boom")
}

func TestAllToleratesFailures(t *testing.T) {
	data := All(context.Background(), errProvider{})
	if data.Info != nil || data.Abyss != nil || data.Stygian != nil {
		t.Fatalf("expected nils on failure, got %+v", data)
	}
	if data.Characters == nil || len(data.Characters) != 0 {
		t.Fatalf("expected empty characters, got %+v", data.Characters)
	}
}
