// Package config loads Plumbline's optional JSON config file and merges it over
// the locked defaults. The artifact-type ladder is not configurable (that is the
// point of adopting the prior art); config only tunes paths, exclusions and the
// coverage threshold. Precedence is flag > config file > baked default.
//
// [impl->component~config-loader~1]
package config

import (
	"encoding/json"
	"os"
)

// Config tunes a Plumbline run. All fields are optional in the file.
type Config struct {
	Register    string   `json:"register"`    // register file or directory (default "register.md")
	Registers   []string `json:"registers"`   // explicit list of register files/dirs/globs; takes precedence over Register
	Roots       []string `json:"roots"`       // source roots to scan (default ["."])
	Exclude     []string `json:"exclude"`     // cleaned paths to skip
	MinCoverage float64  `json:"minCoverage"` // 0 = strict gating; >0 = fail below this %

	// Spec-debt budget over the requirements/features axis — nil = no limit.
	MaxProposed    *int     `json:"maxProposed,omitempty"`    // max un-built feat/req allowed (count)
	MaxProposedPct *float64 `json:"maxProposedPct,omitempty"` // max un-built feat/req allowed (%)
}

// Default returns the locked defaults.
func Default() Config {
	return Config{
		Register: "register.md",
		Roots:    []string{"."},
	}
}

// Load reads a config file if present. found is false (with no error) when the
// path does not exist, so callers can fall back to defaults silently.
func Load(path string) (cfg Config, found bool, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, false, nil
		}
		return Config{}, false, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, true, err
	}
	return cfg, true, nil
}

// Merge overlays a loaded file's non-zero fields onto the base (the defaults).
func Merge(base, over Config) Config {
	if over.Register != "" {
		base.Register = over.Register
	}
	if len(over.Registers) > 0 {
		base.Registers = over.Registers
	}
	if len(over.Roots) > 0 {
		base.Roots = over.Roots
	}
	if len(over.Exclude) > 0 {
		base.Exclude = over.Exclude
	}
	if over.MinCoverage != 0 {
		base.MinCoverage = over.MinCoverage
	}
	if over.MaxProposed != nil {
		base.MaxProposed = over.MaxProposed
	}
	if over.MaxProposedPct != nil {
		base.MaxProposedPct = over.MaxProposedPct
	}
	return base
}
