// Package handlers wires the Genshin page routes to the templ views and the
// Hoyolab provider.
package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"github.com/dotcchix/gonshin/internal/hoyolab"
	"github.com/dotcchix/gonshin/internal/view/genshin"
)

// Handlers holds the route handlers' dependencies.
type Handlers struct {
	provider hoyolab.Provider
	baseURL  string
}

// New creates a Handlers backed by the given provider. baseURL is the public
// site origin used for Discord link previews; when empty it is derived from the
// request.
func New(p hoyolab.Provider, baseURL string) *Handlers {
	return &Handlers{provider: p, baseURL: strings.TrimRight(baseURL, "/")}
}

// Routes registers all page routes on mux.
func (h *Handlers) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", h.page)
	mux.HandleFunc("GET /tabs/{tab}", h.tab)
	mux.HandleFunc("GET /showcase", h.showcase)
	mux.HandleFunc("GET /characters/{id}", h.characterDialog)
	mux.HandleFunc("GET /stygian", h.stygianPanel)
}

func render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.Render(r.Context(), w); err != nil {
		log.Printf("handlers: render: %v", err)
	}
}

func (h *Handlers) page(w http.ResponseWriter, r *http.Request) {
	active := r.URL.Query().Get("tab")
	if !genshin.IsTab(active) {
		active = "home"
	}
	data := hoyolab.All(r.Context(), h.provider)
	baseURL := h.publicBaseURL(r)
	render(w, r, genshin.Layout(genshin.Page(data, active), genshin.BuildDiscordEmbed(data, baseURL), baseURL))
}

// publicBaseURL resolves the site origin for links and og tags, preferring the
// configured value and otherwise reconstructing it from the request.
func (h *Handlers) publicBaseURL(r *http.Request) string {
	if h.baseURL != "" {
		return h.baseURL
	}
	host := r.Host
	if host == "" {
		return ""
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		if comma := strings.IndexByte(proto, ','); comma >= 0 {
			proto = proto[:comma]
		}
		if proto = strings.TrimSpace(proto); proto != "" {
			scheme = proto
		}
	}
	return scheme + "://" + host
}

func (h *Handlers) tab(w http.ResponseWriter, r *http.Request) {
	tab := r.PathValue("tab")
	if !genshin.IsTab(tab) {
		http.NotFound(w, r)
		return
	}
	data := hoyolab.All(r.Context(), h.provider)
	render(w, r, genshin.TabResponse(data, tab))
}

func (h *Handlers) showcase(w http.ResponseWriter, r *http.Request) {
	element := r.URL.Query().Get("element")
	if element == "" {
		element = "All"
	}
	characters, err := h.provider.Characters(r.Context())
	if err != nil && characters == nil {
		log.Printf("handlers: showcase: %v", err)
	}
	render(w, r, genshin.CharacterShowcasePanel(characters, element))
}

func (h *Handlers) characterDialog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	detail, err := h.provider.CharacterDetail(r.Context(), id)
	if err != nil || detail == nil {
		render(w, r, genshin.CharacterDialogError(dialogErrorMessage(err), h.fallbackName(r, id)))
		return
	}
	render(w, r, genshin.CharacterDialog(detail))
}

func (h *Handlers) stygianPanel(w http.ResponseWriter, r *http.Request) {
	cycle := 0
	if raw := r.URL.Query().Get("cycle"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			cycle = n
		}
	}
	stygian, err := h.provider.Stygian(r.Context())
	if err != nil || stygian == nil || len(stygian.Cycles) == 0 {
		http.NotFound(w, r)
		return
	}
	render(w, r, genshin.StygianPanel(stygian, cycle))
}

// fallbackName finds a roster character's name for the dialog error state.
func (h *Handlers) fallbackName(r *http.Request, id int) string {
	characters, err := h.provider.Characters(r.Context())
	if err != nil {
		return ""
	}
	for _, c := range characters {
		if c.ID == id {
			return c.Name
		}
	}
	return ""
}

func dialogErrorMessage(err error) string {
	switch {
	case errors.Is(err, hoyolab.ErrNotFound):
		return "Failed to fetch (404)"
	case errors.Is(err, hoyolab.ErrUnavailable):
		return "Failed to fetch (503)"
	default:
		return "Failed to fetch"
	}
}
