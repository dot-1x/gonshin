package handlers

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var updateGolden = flag.Bool("update-golden", false, "rewrite golden snapshot files")

func golden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", "golden", name)
	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("missing golden %s (run: go test ./... -update-golden): %v", name, err)
	}
	if got != string(want) {
		t.Errorf("golden mismatch for %s\nrun with -update-golden after reviewing the diff", name)
	}
}

func TestGoldenRoutes(t *testing.T) {
	mux := newMux(newStub())
	cases := map[string]string{
		"page.html":           "/",
		"tab_characters.html": "/tabs/characters",
		"tab_spiral.html":     "/tabs/spiral",
		"tab_theater.html":    "/tabs/theater",
		"showcase_hydro.html": "/showcase?element=Hydro",
		"dialog.html":         "/characters/1",
	}
	for file, path := range cases {
		rec := get(t, mux, path)
		if rec.Code != 200 {
			t.Fatalf("%s: status %d", path, rec.Code)
		}
		golden(t, file, rec.Body.String())
	}
}
