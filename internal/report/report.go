// Package report cross-checks anchors against the register and produces the
// bidirectional traceability report. It resolves coverage across the full
// requirement→architecture→code chain (not just direct anchors): a Need is
// satisfied by a code anchor of that type, or by a register item of that type
// that Covers the item. Deep coverage additionally requires every covering
// register item to be deeply covered itself — surfacing transitive defects.
//
// Item status governs tracing: rejected items are excluded entirely; approved
// items (the default) are gated and scored; proposed/draft items are tracked and
// surfaced as planned, but never fail the gate.
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

// Report is the full bidirectional traceability result.
type Report struct {
	Summary    Summary         `json:"summary"`
	Uncovered  []Uncovered     `json:"uncovered"`
	Transitive []TransitiveGap `json:"transitiveDefects"`
	Broken     []Broken        `json:"broken"`
	Orphans    []Orphan        `json:"orphans"`
	Structural []c4.Violation  `json:"structural"`
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
	MinCoverage      float64 `json:"minCoverage"` // gate threshold (0 = strict)
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

// Orphan is a code-area — a directory of scanned source — carrying no anchor at
// all. The anchor unit is the code-area (package / directory), one altitude up
// from a symbol, so a single anchor in any file covers the whole area.
type Orphan struct {
	Area string `json:"area"`
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
// minCoverage of 0 selects strict gating (any gap fails); a positive value
// selects threshold gating (fail below that shallow-coverage percentage, but
// broken anchors and structural errors always fail).
func Build(items []register.Item, anchors []anchor.Anchor, scannedFiles []string, minCoverage float64) Report {
	// Status governs tracing: rejected items are excluded entirely (OFT); the rest
	// are traced, but only approved items are gated and scored.
	var traced []register.Item
	rejected := 0
	for _, it := range items {
		if it.Rejected() {
			rejected++
			continue
		}
		traced = append(traced, it)
	}

	byID := make(map[string]register.Item, len(traced))
	for _, it := range traced {
		byID[it.ID] = it
	}

	// Anchors that resolve contribute a covering type; the rest are broken.
	anchorTypes := make(map[string]map[string]bool)
	var broken []Broken
	for _, a := range anchors {
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
	var planned []Planned
	shallowCount, deepCount, approvedCount := 0, 0, 0

	for _, it := range traced {
		var missing []string
		for _, need := range it.Needs {
			if !needMetShallow(it.ID, need) {
				missing = append(missing, need)
			}
		}
		isShallow := len(missing) == 0
		isDeep := isShallow && deep(it)

		// Proposed/draft: tracked and surfaced for the burndown, but never a gap
		// and never scored — you designed it, you haven't promised it yet.
		if it.Planned() {
			planned = append(planned, Planned{
				ID: it.ID, Title: it.Title, Status: it.StatusOrDefault(),
				Realised: isDeep, File: it.File, Line: it.Line,
			})
			continue
		}

		// Approved: gated as normal, and the basis for the scores.
		approvedCount++
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
	sort.Slice(planned, func(i, j int) bool { return planned[i].ID < planned[j].ID })
	sort.Slice(broken, func(i, j int) bool {
		if broken[i].File != broken[j].File {
			return broken[i].File < broken[j].File
		}
		return broken[i].Line < broken[j].Line
	})
	sort.Slice(orphans, func(i, j int) bool { return orphans[i].Area < orphans[j].Area })

	r := Report{
		Uncovered:  uncovered,
		Transitive: transitive,
		Broken:     broken,
		Orphans:    orphans,
		Structural: structural,
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
			MinCoverage:      minCoverage,
		},
	}
	r.Summary.OK = gate(r.Summary, minCoverage)
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

// gate decides the pass/fail. Broken anchors and structural errors always fail.
// Otherwise: threshold mode fails below minCoverage; strict mode fails on any gap.
// Only approved items reach the gap lists, so proposed/draft items never fail it.
func gate(s Summary, minCoverage float64) bool {
	if s.BrokenCount > 0 || s.StructuralErrors > 0 {
		return false
	}
	if minCoverage > 0 {
		return s.CoveragePct >= minCoverage
	}
	return s.UncoveredCount == 0 && s.TransitiveCount == 0 && s.OrphanCount == 0
}

func pct(n, total int) float64 {
	if total == 0 {
		return 100.0
	}
	return math.Round(float64(n)/float64(total)*1000) / 10
}
