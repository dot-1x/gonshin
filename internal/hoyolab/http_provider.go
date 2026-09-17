package hoyolab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HTTPProvider talks to the existing Hoyolab HTTP API. It is the current
// implementation of Provider; swapping in a native Go implementation later
// requires no handler or view changes.
type HTTPProvider struct {
	baseURL string
	client  *http.Client
	ttl     time.Duration
	cache   *ttlCache
}

// NewHTTPProvider builds a provider against baseURL (e.g.
// "https://hoyo.dotcchix.dev"). ttl controls caching of list endpoints.
func NewHTTPProvider(baseURL string, ttl time.Duration) *HTTPProvider {
	return &HTTPProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 20 * time.Second},
		ttl:     ttl,
		cache:   newTTLCache(),
	}
}

func (p *HTTPProvider) getJSON(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL+"/api/gi"+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch {
	case res.StatusCode == http.StatusNotFound:
		return ErrNotFound
	case res.StatusCode == http.StatusServiceUnavailable:
		return ErrUnavailable
	case res.StatusCode >= http.StatusBadRequest:
		return fmt.Errorf("hoyolab: %s: unexpected status %d", path, res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(out)
}

func fetchCached[T any](ctx context.Context, p *HTTPProvider, path string) (*T, error) {
	v, err := p.cache.do(path, p.ttl, func() (any, error) {
		var out T
		if err := p.getJSON(ctx, path, &out); err != nil {
			return nil, err
		}
		return &out, nil
	})
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return v.(*T), nil
}

// Info implements Provider.
func (p *HTTPProvider) Info(ctx context.Context) (*AccountInfo, error) {
	return fetchCached[AccountInfo](ctx, p, "/info")
}

// Characters implements Provider.
func (p *HTTPProvider) Characters(ctx context.Context) ([]AccountCharacter, error) {
	list, err := fetchCached[[]AccountCharacter](ctx, p, "/characters")
	if err != nil || list == nil {
		return nil, err
	}
	return *list, nil
}

// Abyss implements Provider.
func (p *HTTPProvider) Abyss(ctx context.Context) (*AccountAbyss, error) {
	return fetchCached[AccountAbyss](ctx, p, "/abyss")
}

// Theater implements Provider, returning the current (first) season.
func (p *HTTPProvider) Theater(ctx context.Context) (*AccountTheater, error) {
	list, err := fetchCached[[]AccountTheater](ctx, p, "/theater")
	if err != nil || list == nil || len(*list) == 0 {
		return nil, err
	}
	return &(*list)[0], nil
}

// Stygian implements Provider.
func (p *HTTPProvider) Stygian(ctx context.Context) (*AccountStygian, error) {
	return fetchCached[AccountStygian](ctx, p, "/stygian/detail")
}

// CharacterDetail implements Provider. It is intentionally uncached, matching
// the original `cache: "no-store"` dialog fetch.
func (p *HTTPProvider) CharacterDetail(ctx context.Context, id int) (*CharacterDetail, error) {
	var out CharacterDetail
	if err := p.getJSON(ctx, fmt.Sprintf("/characters/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

var _ Provider = (*HTTPProvider)(nil)
