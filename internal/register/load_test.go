package register

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAggregatesTree(t *testing.T) {
	dir := t.TempDir()
	write := func(rel, body string) {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("requirements.md", "`req~a~1`\nNeeds: component\n")
	write("components/engine.md", "`component~a~1`\nCovers: req~a~1\nNeeds: impl\n")
	write("README.md", "Just prose, no spec items.\n")

	// Directory recursion aggregates both spec files (README contributes nothing).
	items, files, err := Load([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
	if len(files) != 3 {
		t.Fatalf("register files = %d, want 3 (README parsed, yields no items)", len(files))
	}

	// A glob targets one level.
	only, _, err := Load([]string{filepath.Join(dir, "components", "*.md")})
	if err != nil {
		t.Fatal(err)
	}
	if len(only) != 1 || only[0].ID != "component~a~1" {
		t.Fatalf("glob load = %+v", only)
	}
}

func TestLoadMissing(t *testing.T) {
	if _, _, err := Load([]string{filepath.Join(t.TempDir(), "nope")}); err == nil {
		t.Error("expected error for a missing register path")
	}
}
