# gonshin

A rewrite of the **Genshin Impact profile page** from
`dotcchix-porto26/app/hoyo/genshin-impact` (Next.js/React) using the **GOTH
stack**: **Go + Templ + HTMX + Tailwind CSS**. The UI is a 1:1 port of the
original — same markup, same Tailwind classes, same theme tokens, same assets.

Data comes from the existing Hoyolab HTTP API (`https://hoyo.dotcchix.dev`)
behind a provider interface, so the Python FastAPI backend in
`~/Experiments/hoyolab-api` can later be replaced by a native Go implementation
without touching the handlers or views.

## Stack

| Layer | Choice |
|---|---|
| Server | Go (standard `net/http`, Go 1.22+ routing) |
| Templates | [templ](https://templ.guide) |
| Interactivity | [HTMX](https://htmx.org) + a few lines of vanilla JS for the dialog |
| Styling | Tailwind CSS v4 via `@tailwindcss/cli`, same theme CSS as the React app |
| Icons | Hand-ported lucide SVGs |
| Fonts | Self-hosted Geist / Geist Mono variable fonts |

## Layout

```
cmd/gonshin/main.go            # server + graceful shutdown
internal/config/               # env config
internal/hoyolab/              # models, Provider interface, HTTP provider, cache, derive
internal/handlers/             # routes + render glue
internal/view/ui/              # shadcn-equivalent templ primitives + icons
internal/view/genshin/         # page, nav, tabs and section components
internal/style/                # tailwind-merge-equivalent class merger
web/css/input.css              # Tailwind entry (theme mirrors app/globals.css)
web/static/                    # fonts, htmx, genshin assets, compiled app.css
```

## Migration seam

All data access goes through `hoyolab.Provider`:

```go
type Provider interface {
    Info(ctx) (*AccountInfo, error)
    Characters(ctx) ([]AccountCharacter, error)
    Abyss(ctx) (*AccountAbyss, error)
    Theater(ctx) (*AccountTheater, error)
    Stygian(ctx) (*AccountStygian, error)
    CharacterDetail(ctx, id) (*CharacterDetail, error)
}
```

`hoyolab.HTTPProvider` is the current implementation (proxies the existing API,
with a TTL cache matching the original `revalidate: 1200`). To move the backend
to Go, implement the same interface and swap it in `cmd/gonshin/main.go`.

## Configuration

Copy `.env.example` to `.env` (or export the vars):

| Variable | Default | Description |
|---|---|---|
| `HOYOLAB_API_BASE` | `https://hoyo.dotcchix.dev` | Hoyolab API root |
| `ADDR` | `:8080` | HTTP listen address |
| `CACHE_TTL_SECONDS` | `1200` | Cache TTL for list endpoints (detail is never cached) |

## Docker

```sh
docker compose up -d --build   # serves on http://localhost:8510
```

The image is a multi-stage build (Tailwind CSS → Go/templ build → minimal
Alpine runtime) and runs as a non-root user. Override the API base or TTL via
environment variables:

```sh
HOYOLAB_API_BASE=http://host.docker.internal:8504 docker compose up -d
```

## Development

```sh
pnpm install      # tailwind + fonts
make build        # templ generate + css + go build
make run          # build and run
make check        # templ generate + go vet + go test + go build
```

For live CSS during development: `pnpm css:watch`. Regenerate templates with
`make generate` (requires `go install github.com/a-h/templ/cmd/templ@latest`).

## Tests

```sh
go test ./...
```

Covers the class merger (validated against `tailwind-merge`), the HTTP provider
(httptest upstream, caching, error mapping), derive helpers, component
rendering, and HTTP routes. Golden snapshots for the page and tab routes live in
`internal/handlers/testdata/golden`; regenerate with
`go test ./internal/handlers -update-golden` after reviewing the diff.

## Routes

| Route | Description |
|---|---|
| `GET /` | Full page (Detail tab active) |
| `GET /tabs/{tab}` | HTMX tab partial + out-of-band nav (`home`, `characters`, `spiral`, `stygian`, `theater`) |
| `GET /showcase?element=` | Character grid filter partial |
| `GET /stygian?cycle=` | Stygian cycle panel partial |
| `GET /characters/{id}` | Character detail dialog body |
| `GET /static/` | Fonts, CSS, HTMX, Genshin assets |

## Notes

- The original `next/image` usage is replaced with plain `<img>` plus the
  equivalent `absolute inset-0 h-full w-full object-cover` classes; remote
  hoyoverse images load directly.
- Radix tooltips are reimplemented as CSS `group-hover` tooltips; the dialog is
  a fixed overlay toggled by a small inline script. Visual output matches the
  original.
