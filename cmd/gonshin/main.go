// Command gonshin serves the Genshin Impact profile page (GOTH stack:
// Go + Templ + HTMX + Tailwind) against the Hoyolab data API.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dotcchix/gonshin/internal/config"
	"github.com/dotcchix/gonshin/internal/handlers"
	"github.com/dotcchix/gonshin/internal/hoyolab"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	provider := hoyolab.NewHTTPProvider(cfg.APIBase, cfg.CacheTTL)

	mux := http.NewServeMux()
	mux.Handle("GET /static/", cacheStatic(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static")))))
	handlers.New(provider).Routes(mux)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on %s (api: %s)", cfg.Addr, cfg.APIBase)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

// cacheStatic sets a long-lived cache header for fingerprinted-ish assets.
func cacheStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		next.ServeHTTP(w, r)
	})
}
