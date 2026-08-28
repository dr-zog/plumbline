package config

// [utest->component~config-loader~1]

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissing(t *testing.T) {
	_, found, err := Load(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatalf("missing file should not error: %v", err)
	}
	if found {
		t.Error("missing file should report found=false")
	}
}

func TestLoadAndMerge(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plumbline.json")
	if err := os.WriteFile(path, []byte(`{"register":"docs/reg.md","minCoverage":80}`), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, found, err := Load(path)
	if err != nil || !found {
		t.Fatalf("load failed: found=%v err=%v", found, err)
	}
	merged := Merge(Default(), loaded)
	if merged.Register != "docs/reg.md" {
		t.Errorf("register = %q, want docs/reg.md", merged.Register)
	}
	if merged.MinCoverage != 80 {
		t.Errorf("minCoverage = %v, want 80", merged.MinCoverage)
	}
	// Roots not set in file → keeps the default.
	if len(merged.Roots) != 1 || merged.Roots[0] != "." {
		t.Errorf("roots = %v, want [.]", merged.Roots)
	}
}

func TestMergeRegisters(t *testing.T) {
	over := Config{Registers: []string{"register/", "docs/extra.md"}}
	merged := Merge(Default(), over)
	if len(merged.Registers) != 2 || merged.Registers[0] != "register/" {
		t.Errorf("registers = %v", merged.Registers)
	}
	// Register keeps its default; the CLI decides precedence between the two.
	if merged.Register != "register.md" {
		t.Errorf("register = %q, want register.md", merged.Register)
	}
}

func TestMergeBudget(t *testing.T) {
	n, pct := 5, 12.5
	merged := Merge(Default(), Config{MaxProposed: &n, MaxProposedPct: &pct})
	if merged.MaxProposed == nil || *merged.MaxProposed != 5 {
		t.Errorf("maxProposed = %v, want 5", merged.MaxProposed)
	}
	if merged.MaxProposedPct == nil || *merged.MaxProposedPct != 12.5 {
		t.Errorf("maxProposedPct = %v, want 12.5", merged.MaxProposedPct)
	}
	// Absent in the overlay → stays nil (no limit).
	plain := Merge(Default(), Config{})
	if plain.MaxProposed != nil || plain.MaxProposedPct != nil {
		t.Error("budget should default to nil (no limit)")
	}
}
