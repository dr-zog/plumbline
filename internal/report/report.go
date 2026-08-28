// Package report cross-checks anchors against the register and produces the
// bidirectional traceability report. It resolves coverage across the full
// requirement→architecture→code chain (not just direct anchors): a Need is
// satisfied by a code anchor of that type, or by a register item of that type
// that Covers the item. Deep coverage additionally requires every covering
// register item to be deeply covered itself — surfacing transitive defects.
//
// Item status governs tracing (see ADR 004): rejected items are excluded, and
// code that anchors one is zombie code (a hard fail); approved items are gated
// and scored; proposed/draft items are tracked, and code against them is a
// warning (building ahead of approval), never a gate failure. Un-built spec
// (not-yet-approved feat/req) is bounded by an optional spec-debt budget.
//
// The one report yields both a gate (OK / exit code) and a scorecard (coverage
// and completeness percentages, gap counts) the skills narrate.
//
// [impl->component~report-builder~1]
package report

import (
	"math"
	"path/filepath"
	"sort"

	"github.com/dr-zog/plumbline/internal/anchor"
	"github.com/dr-zog/plumbline/internal/c4"
	"github.com/dr-zog/plumbline/internal/register"
)

// GateOpts tunes the gate: the coverage threshold and the spec-debt budget.
type GateOpts struct {
	MinCoverage    float64  // 0 = strict (any gap fails); >0 = fail below this shallow %
	MaxProposed    *int     // nil = no limit on the un-built-spec count
	MaxProposedPct *float64 // nil = no limit on the un-built-spec ratio
}

// Report is the full bidirectional traceability result.
type Report struct {
	Summary    Summary         `json:"summary"`
	Uncovered  []Uncovered     `json:"uncovered"`
	Transitive []TransitiveGap `json:"transitiveDefects"`
	Broken     []Broken        `json:"broken"`
	Orphans    []Orphan        `json:"orphans"`
	Structural []c4.Violation  `json:"structural"`
	DeadEnds   []DeadEnd       `json:"deadEnds"`
	Zombies    []Zombie        `json:"zombies"`
	Warnings   []Warning       `json:"warnings"`
	Planned    []Planned       `json:"planned"`
}

// Summary holds the headline counts, the two scores, and the pass/fail gate.
type Summary struct {
	RegisterItems    int     `json:"registerItems"` // traced items (approved + planned; rejected excluded)
	ApprovedItems    int     `json:"approvedItems"` // gated, and the basis for the scores
	PlannedItems     int     `json:"plannedItems"`  // proposed/draft — tracked, not gated
	RejectedItems    int     `json:"rejectedItems"` // excluded from tracing
	Anchors          int     `json:"anchors"`
	ShallowCovered   int     `json:"shallowCovered"`
	DeepCovered      int     `json:"deepCovered"`
	CoveragePct      float64 `json:"coveragePct"`     // shallow: direct Needs met, over approved
	CompletenessPct  float64 `json:"completenessPct"` // deep: whole chain resolves, over approved
	UncoveredCount   int     `json:"uncoveredCount"`
	TransitiveCount  int     `json:"transitiveDefectCount"`
	BrokenCount      int     `json:"brokenCount"`
	OrphanCount      int     `json:"orphanCount"`
	StructuralErrors int     `json:"structuralErrorCount"`
	DeadEndCount     int     `json:"deadEndCount"`
	ZombieCount      int     `json:"zombieCount"`    // code anchoring a rejected item — fails
	WarningCount     int     `json:"warningCount"`   // build-ahead / status-lag — never fails
	SpecDebtCount    int     `json:"specDebtCount"`  // un-built feat/req (not-yet-approved, un-realised)
	SpecTotalItems   int     `json:"specTotalItems"` // all non-rejected feat/req
	SpecDebtPct      float64 `json:"specDebtPct"`    // specDebtCount / specTotalItems
	MinCoverage      float64 `json:"minCoverage"`    // gate threshold (0 = strict)
	MaxProposed      int     `json:"maxProposed"`    // spec-debt budget (count); -1 = no limit
	MaxProposedPct   float64 `json:"maxProposedPct"` // spec-debt budget (%); -1 = no limit
	OK               bool    `json:"ok"`
}

// Uncovered is a register item with needed coverage that nothing provides.
type Uncovered struct {
	ID      string   `json:"id"`
	Title   string   `json:"title,omitempty"`
	Needs   []string `json:"needs"`
	Missing []string `json:"missing"`
	File    string   `json:"file"`
	Line    int      `json:"line"`
}

// TransitiveGap is an item whose direct Needs are met, but a covering item
// further down the chain is not itself deeply covered.
type TransitiveGap struct {
	ID    string   `json:"id"`
	Title string   `json:"title,omitempty"`
	Weak  []string `json:"weakCoverers"` // covering item IDs that aren't deep
	File  string   `json:"file"`
	Line  int      `json:"line"`
}

// Broken is an anchor whose target ID is absent from the register.
type Broken struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Tag      string `json:"tag"`
	TargetID string `json:"targetId"`
}

// Zombie is an anchor pointing at a rejected item — code for a requirement the
// register says was abandoned. The target exists but is rejected, so it's not a
// broken anchor; it's a hard fail of its own (remove the code, or un-reject).
type Zombie struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Tag      string `json:"tag"`
	TargetID string `json:"targetId"`
}

// Orphan is a code-area — a directory of scanned source — carrying no anchor at
// all. The anchor unit is the code-area (package / directory), one altitude up
// from a symbol, so a single anchor in any file covers the whole area.
type Orphan struct {
	Area string `json:"area"`
}

// DeadEnd is an approved item that declares no Needs. In the locked ladder every
// register-item type must need something below it, so an empty Needs terminates
// the chain prematurely — it can never be genuinely covered, and would otherwise
// read as vacuously complete and inflate the score. Only approved items qualify;
// a proposed/draft item without Needs is legitimately not-yet-armed.
type DeadEnd struct {
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
	File  string `json:"file"`
	Line  int    `json:"line"`
}

// Warning is a not-yet-approved item with code against it — surfaced, never a
// gate failure. Kind is "build-ahead" (code for an unapproved spec) or
// "status-lag" (fully covered but still unapproved → promote to approved).
type Warning struct {
	ID     string `json:"id"`
	Title  string `json:"title,omitempty"`
	Status string `json:"status"`
	Kind   string `json:"kind"`
	File   string `json:"file"`
	Line   int    `json:"line"`
}

// Planned is a proposed or draft item: tracked and surfaced for the planned-vs-
// realised view, but never gated. Realised is true when its chain already
// resolves to code.
type Planned struct {
	ID       string `json:"id"`
	Title    string `json:"title,omitempty"`
	Status   string `json:"status"`
	Realised bool   `json:"realised"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Build cross-checks anchors against register items and scanned source files.
// Gating: broken anchors, structural errors, dead-ends and zombie code always
// fail; strict mode fails on any gap; a coverage threshold or a spec-debt budget
// (GateOpts) adds further conditions.
func Build(items []register.Item, anchors []anchor.Anchor, scannedFiles []string, opts GateOpts) Report {
	// Status governs tracing: rejected items are excluded from coverage (OFT), but
	// remembered so code anchoring one reads as zombie code, not a broken anchor.
	var traced []register.Item
	rejected := 0
	rejectedIDs := map[string]bool{}
	for _, it := range items {
		if it.Rejected() {
			rejected++
			rejectedIDs[it.ID] = true
			continue
		}
		traced = append(traced, it)
	}

	byID := make(map[string]register.Item, len(traced))
	for _, it := range traced {
		byID[it.ID] = it
	}

	// Anchors that resolve contribute a covering type; an anchor on a rejected item
	// is zombie code; the rest are broken.
	anchorTypes := make(map[string]map[string]bool)
	var broken []Broken
	var zombies []Zombie
	for _, a := range anchors {
		if rejectedIDs[a.TargetID] {
			zombies = append(zombies, Zombie{File: a.File, Line: a.Line, Tag: a.Raw, TargetID: a.TargetID})
			continue
		}
		if _, ok := byID[a.TargetID]; !ok {
			broken = append(broken, Broken{File: a.File, Line: a.Line, Tag: a.Raw, TargetID: a.TargetID})
			continue
		}
		if anchorTypes[a.TargetID] == nil {
			anchorTypes[a.TargetID] = map[string]bool{}
		}
		anchorTypes[a.TargetID][a.Covering] = true
	}

	// coverers[X] = register items that declare `Covers: X`.
	coverers := make(map[string][]register.Item)
	for _, it := range traced {
		for _, tgt := range it.Covers {
			coverers[tgt] = append(coverers[tgt], it)
		}
	}

	// needMetShallow: a Need type is met by an anchor or by any covering item.
	needMetShallow := func(id, need string) bool {
		if anchorTypes[id][need] {
			return true
		}
		for _, y := range coverers[id] {
			if y.Type == need {
				return true
			}
		}
		return false
	}

	// deep is memoised: an item is deep iff every Need is met by an anchor, or
	// by a covering item that is itself deep.
	memo := make(map[string]bool)
	visiting := make(map[string]bool)
	var deep func(it register.Item) bool
	deep = func(it register.Item) bool {
		if v, ok := memo[it.ID]; ok {
			return v
		}
		if visiting[it.ID] {
			return false // defensive: ladder is acyclic, but never loop
		}
		visiting[it.ID] = true
		res := true
		for _, need := range it.Needs {
			if !needMetDeep(it.ID, need, anchorTypes, coverers, deep) {
				res = false
				break
			}
		}
		visiting[it.ID] = false
		memo[it.ID] = res
		return res
	}

	var uncovered []Uncovered
	var transitive []TransitiveGap
	var deadEnds []DeadEnd
	var warnings []Warning
	var planned []Planned
	shallowCount, deepCount, approvedCount := 0, 0, 0
	specDebtCount, specTotal := 0, 0

	for _, it := range traced {
		var missing []string
		for _, need := range it.Needs {
			if !needMetShallow(it.ID, need) {
				missing = append(missing, need)
			}
		}
		isShallow := len(missing) == 0
		isDeep := isShallow && deep(it)

		// Spec-debt budget is measured over the requirements/features axis only:
		// a not-yet-approved feat/req that isn't fully realised is un-built spec.
		if it.Type == "feat" || it.Type == "req" {
			specTotal++
			if it.Planned() && !isDeep {
				specDebtCount++
			}
		}

		// Proposed/draft: tracked, never gated. Code against one is a warning —
		// building ahead of approval (status-lag once fully covered) — not a gap.
		if it.Planned() {
			planned = append(planned, Planned{
				ID: it.ID, Title: it.Title, Status: it.StatusOrDefault(),
				Realised: isDeep, File: it.File, Line: it.Line,
			})
			if metNeeds := len(it.Needs) - len(missing); metNeeds > 0 {
				kind := "build-ahead"
				if isDeep {
					kind = "status-lag" // fully built but unapproved → promote to approved
				}
				warnings = append(warnings, Warning{
					ID: it.ID, Title: it.Title, Status: it.StatusOrDefault(),
					Kind: kind, File: it.File, Line: it.Line,
				})
			}
			continue
		}

		// Approved: gated as normal, and the basis for the scores.
		approvedCount++
		if len(it.Needs) == 0 {
			// A dead-end: an approved item that declares no downward need. It can
			// never be genuinely covered, so it's a gap — not a vacuous 100%.
			deadEnds = append(deadEnds, DeadEnd{ID: it.ID, Title: it.Title, File: it.File, Line: it.Line})
			continue
		}
		if isShallow {
			shallowCount++
		}
		if isDeep {
			deepCount++
		}
		switch {
		case !isShallow:
			uncovered = append(uncovered, Uncovered{
				ID: it.ID, Title: it.Title, Needs: it.Needs,
				Missing: missing, File: it.File, Line: it.Line,
			})
		case !isDeep:
			transitive = append(transitive, TransitiveGap{
				ID: it.ID, Title: it.Title,
				Weak: weakCoverers(it, coverers, deep),
				File: it.File, Line: it.Line,
			})
		}
	}

	// Orphans: code-areas (directories) that contain scanned source but carry no
	// anchor in any of their files. A single anchor covers the whole area, so
	// adding files to an already-anchored package doesn't demand new anchors —
	// that coarseness is what keeps Plumbline light-touch.
	dirHasAnchor := make(map[string]bool)
	for _, a := range anchors {
		dirHasAnchor[filepath.Dir(a.File)] = true
	}
	seenDir := make(map[string]bool)
	var orphans []Orphan
	for _, f := range scannedFiles {
		d := filepath.Dir(f)
		if seenDir[d] {
			continue
		}
		seenDir[d] = true
		if !dirHasAnchor[d] {
			orphans = append(orphans, Orphan{Area: d})
		}
	}

	structural := c4.Validate(items)
	structuralErrors := len(c4.Errors(structural))

	sort.Slice(uncovered, func(i, j int) bool { return uncovered[i].ID < uncovered[j].ID })
	sort.Slice(transitive, func(i, j int) bool { return transitive[i].ID < transitive[j].ID })
	sort.Slice(deadEnds, func(i, j int) bool { return deadEnds[i].ID < deadEnds[j].ID })
	sort.Slice(warnings, func(i, j int) bool { return warnings[i].ID < warnings[j].ID })
	sort.Slice(planned, func(i, j int) bool { return planned[i].ID < planned[j].ID })
	sort.Slice(broken, func(i, j int) bool {
		if broken[i].File != broken[j].File {
			return broken[i].File < broken[j].File
		}
		return broken[i].Line < broken[j].Line
	})
	sort.Slice(zombies, func(i, j int) bool {
		if zombies[i].File != zombies[j].File {
			return zombies[i].File < zombies[j].File
		}
		return zombies[i].Line < zombies[j].Line
	})
	sort.Slice(orphans, func(i, j int) bool { return orphans[i].Area < orphans[j].Area })

	r := Report{
		Uncovered:  uncovered,
		Transitive: transitive,
		Broken:     broken,
		Orphans:    orphans,
		Structural: structural,
		DeadEnds:   deadEnds,
		Zombies:    zombies,
		Warnings:   warnings,
		Planned:    planned,
		Summary: Summary{
			RegisterItems:    len(traced),
			ApprovedItems:    approvedCount,
			PlannedItems:     len(planned),
			RejectedItems:    rejected,
			Anchors:          len(anchors),
			ShallowCovered:   shallowCount,
			DeepCovered:      deepCount,
			CoveragePct:      pct(shallowCount, approvedCount),
			CompletenessPct:  pct(deepCount, approvedCount),
			UncoveredCount:   len(uncovered),
			TransitiveCount:  len(transitive),
			BrokenCount:      len(broken),
			OrphanCount:      len(orphans),
			StructuralErrors: structuralErrors,
			DeadEndCount:     len(deadEnds),
			ZombieCount:      len(zombies),
			WarningCount:     len(warnings),
			SpecDebtCount:    specDebtCount,
			SpecTotalItems:   specTotal,
			SpecDebtPct:      pctOrZero(specDebtCount, specTotal),
			MinCoverage:      opts.MinCoverage,
			MaxProposed:      derefOr(opts.MaxProposed, -1),
			MaxProposedPct:   derefFloatOr(opts.MaxProposedPct, -1),
		},
	}
	r.Summary.OK = gate(r.Summary, opts)
	return r
}

// needMetDeep reports whether a Need is met by an anchor (a leaf, always deep)
// or by a covering item that is itself deeply covered.
func needMetDeep(id, need string, anchorTypes map[string]map[string]bool,
	coverers map[string][]register.Item, deep func(register.Item) bool) bool {
	if anchorTypes[id][need] {
		return true
	}
	for _, y := range coverers[id] {
		if y.Type == need && deep(y) {
			return true
		}
	}
	return false
}

// weakCoverers lists the covering items of it that are not themselves deep.
func weakCoverers(it register.Item, coverers map[string][]register.Item, deep func(register.Item) bool) []string {
	var out []string
	for _, y := range coverers[it.ID] {
		if !deep(y) {
			out = append(out, y.ID)
		}
	}
	sort.Strings(out)
	return out
}

// gate decides the pass/fail. Broken anchors, structural errors, dead-ends and
// zombie code always fail. Then the spec-debt budget (if set) fails when un-built
// spec exceeds it. Then: threshold mode fails below minCoverage; strict mode fails
// on any gap. Warnings never fail; only approved items reach the gap lists, so
// proposed/draft items never fail it.
func gate(s Summary, opts GateOpts) bool {
	if s.BrokenCount > 0 || s.StructuralErrors > 0 || s.DeadEndCount > 0 || s.ZombieCount > 0 {
		return false
	}
	if opts.MaxProposed != nil && s.SpecDebtCount > *opts.MaxProposed {
		return false
	}
	if opts.MaxProposedPct != nil && s.SpecDebtPct > *opts.MaxProposedPct {
		return false
	}
	if opts.MinCoverage > 0 {
		return s.CoveragePct >= opts.MinCoverage
	}
	return s.UncoveredCount == 0 && s.TransitiveCount == 0 && s.OrphanCount == 0
}

func pct(n, total int) float64 {
	if total == 0 {
		return 100.0
	}
	return math.Round(float64(n)/float64(total)*1000) / 10
}

// pctOrZero is like pct but returns 0 for an empty set (no spec = no debt, not
// "100% debt").
func pctOrZero(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return math.Round(float64(n)/float64(total)*1000) / 10
}

func derefOr(p *int, d int) int {
	if p != nil {
		return *p
	}
	return d
}

func derefFloatOr(p *float64, d float64) float64 {
	if p != nil {
		return *p
	}
	return d
}
