// Command plumbline is the Plumbline engine: a zero-dependency binary that
// proves anchors resolve against the register in both directions and fails the
// build when they don't. It resolves coverage across the full requirement →
// architecture → code chain, validates the register's C4 structure, scores
// coverage and completeness, and gates strictly or against a threshold. It
// exits non-zero on any gap, so it doubles as a pre-commit / CI gate with no
// LLM in the loop.
//
// [impl->component~cli~1]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dr-zog/plumbline/internal/anchor"
	"github.com/dr-zog/plumbline/internal/config"
	"github.com/dr-zog/plumbline/internal/register"
	"github.com/dr-zog/plumbline/internal/report"
)

// version is the engine's build version, stamped in at link time
// (-ldflags "-X main.version=…"); "dev" for an un-stamped local build.
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("plumbline", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var (
		cfgPath     = fs.String("config", "plumbline.json", "path to the optional JSON config file")
		regPath     = fs.String("register", "", "register file or directory (overrides config; default register.md)")
		minCov      = fs.Float64("min-coverage", -1, "fail below this coverage % (overrides config; 0 = strict)")
		asJSON      = fs.Bool("json", false, "emit the machine-readable JSON report instead of the text scorecard")
		showVersion = fs.Bool("version", false, "print the engine version and exit")
	)
	fs.Usage = func() {
		fmt.Fprint(stderr, "plumbline — light-touch requirement↔code traceability\n\n")
		fmt.Fprint(stderr, "usage: plumbline [-config PATH] [-register PATH] [-min-coverage N] [-json] [path ...]\n\n")
		fmt.Fprint(stderr, "Scans the given paths for anchors, resolves them against the register across\n")
		fmt.Fprint(stderr, "the full requirement→architecture→code chain, validates C4 structure, and\n")
		fmt.Fprint(stderr, "reports uncovered requirements, transitive defects, broken anchors, orphan\n")
		fmt.Fprint(stderr, "code-areas and structural errors. Exits 0 when clean, 1 on any gap, 2 on error.\n\n")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showVersion {
		fmt.Fprintln(stdout, version)
		return 0
	}

	// Precedence: flag > config file > baked default.
	cfg := config.Default()
	if loaded, found, err := config.Load(*cfgPath); err != nil {
		fmt.Fprintf(stderr, "plumbline: bad config %q: %v\n", *cfgPath, err)
		return 2
	} else if found {
		cfg = config.Merge(cfg, loaded)
	}
	if *minCov >= 0 {
		cfg.MinCoverage = *minCov
	}
	roots := cfg.Roots
	if fs.NArg() > 0 {
		roots = fs.Args()
	}

	// Resolve the register source(s): an explicit -register flag wins, else the
	// config's `registers` list, else its single `register` (file or directory).
	regPaths := cfg.Registers
	if len(regPaths) == 0 {
		regPaths = []string{cfg.Register}
	}
	if *regPath != "" {
		regPaths = []string{*regPath}
	}

	items, regFiles, err := register.Load(regPaths)
	if err != nil {
		fmt.Fprintf(stderr, "plumbline: cannot read register: %v\n", err)
		return 2
	}

	// Keep register prose out of the anchor scan: exclude the source paths (a
	// directory prefix skips its subtree) and every file actually parsed.
	exclude := map[string]bool{}
	for _, p := range regPaths {
		exclude[filepath.Clean(p)] = true
	}
	for _, f := range regFiles {
		exclude[f] = true
	}
	for _, e := range cfg.Exclude {
		exclude[filepath.Clean(e)] = true
	}
	anchors, files, err := anchor.Scan(roots, exclude)
	if err != nil {
		fmt.Fprintf(stderr, "plumbline: scan failed: %v\n", err)
		return 2
	}

	rep := report.Build(items, anchors, files, cfg.MinCoverage)

	if *asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(rep); err != nil {
			fmt.Fprintf(stderr, "plumbline: %v\n", err)
			return 2
		}
	} else {
		writeText(stdout, rep)
	}

	if rep.Summary.OK {
		return 0
	}
	return 1
}

func writeText(w io.Writer, r report.Report) {
	s := r.Summary
	fmt.Fprint(w, "Plumbline traceability report\n")
	fmt.Fprintf(w, "  register items : %d\n", s.RegisterItems)
	fmt.Fprintf(w, "  anchors found  : %d\n", s.Anchors)
	fmt.Fprintf(w, "  coverage       : %d/%d items (%.1f%%)   [direct needs met]\n", s.ShallowCovered, s.RegisterItems, s.CoveragePct)
	fmt.Fprintf(w, "  completeness   : %d/%d items (%.1f%%)   [whole chain resolves]\n", s.DeepCovered, s.RegisterItems, s.CompletenessPct)
	if s.MinCoverage > 0 {
		fmt.Fprintf(w, "  gate           : threshold — fail below %.1f%%\n", s.MinCoverage)
	} else {
		fmt.Fprint(w, "  gate           : strict — any gap fails\n")
	}
	fmt.Fprint(w, "  measures       : traceability, not design quality — a green gate is not a verdict on your architecture\n")
	fmt.Fprintln(w)

	if len(r.Structural) > 0 {
		fmt.Fprintf(w, "STRUCTURAL (%d) — register vs the C4 ladder:\n", len(r.Structural))
		fmt.Fprint(w, "  → fix: correct the type, Needs edge, or Covers link at each location below.\n")
		for _, v := range r.Structural {
			fmt.Fprintf(w, "  [%s] %s: %s  (%s:%d)\n", v.Severity, v.ItemID, v.Detail, v.File, v.Line)
		}
		fmt.Fprintln(w)
	}
	if len(r.Uncovered) > 0 {
		fmt.Fprintf(w, "UNCOVERED (%d) — needed coverage with nothing to provide it:\n", len(r.Uncovered))
		fmt.Fprint(w, "  → fix: anchor a code-area to each item, or retire the item from the register if it's no longer required.\n")
		for _, u := range r.Uncovered {
			fmt.Fprintf(w, "  %s%s  missing: %s  (%s:%d)\n", u.ID, titleOf(u.Title), strings.Join(u.Missing, ","), u.File, u.Line)
		}
		fmt.Fprintln(w)
	}
	if len(r.Transitive) > 0 {
		fmt.Fprintf(w, "TRANSITIVE DEFECTS (%d) — direct needs met, but a coverer below isn't:\n", len(r.Transitive))
		fmt.Fprint(w, "  → fix: cover the weak item(s) named below so the whole chain resolves.\n")
		for _, t := range r.Transitive {
			fmt.Fprintf(w, "  %s%s  weak: %s  (%s:%d)\n", t.ID, titleOf(t.Title), strings.Join(t.Weak, ","), t.File, t.Line)
		}
		fmt.Fprintln(w)
	}
	if len(r.Broken) > 0 {
		fmt.Fprintf(w, "BROKEN anchors (%d) — tag target not in register:\n", len(r.Broken))
		fmt.Fprint(w, "  → fix: add the target to the register, or point the tag at an existing item.\n")
		for _, b := range r.Broken {
			fmt.Fprintf(w, "  %s  (%s:%d)\n", b.Tag, b.File, b.Line)
		}
		fmt.Fprintln(w)
	}
	if len(r.Orphans) > 0 {
		fmt.Fprintf(w, "ORPHAN code-areas (%d) — directories with no anchor:\n", len(r.Orphans))
		fmt.Fprint(w, "  → fix: add an anchor in each directory (pointing at the component it implements), or exclude it in plumbline.json.\n")
		for _, o := range r.Orphans {
			fmt.Fprintf(w, "  %s\n", o.Area)
		}
		fmt.Fprintln(w)
	}

	if s.OK {
		fmt.Fprint(w, "OK — traceability holds.\n")
	} else {
		// Lead with progress, not just a verdict — during a large move you want to
		// see how far along you are, not a bare red.
		fmt.Fprintf(w, "FAIL — %d/%d items covered; still to clear: %d uncovered, %d transitive, %d broken, %d orphan, %d structural.\n",
			s.ShallowCovered, s.RegisterItems, s.UncoveredCount, s.TransitiveCount, s.BrokenCount, s.OrphanCount, s.StructuralErrors)
	}
}

func titleOf(t string) string {
	if t == "" {
		return ""
	}
	return " — " + t
}
