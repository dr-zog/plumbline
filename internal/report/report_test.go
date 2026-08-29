package report

// [utest->component~report-builder~1]
//
// oft:off — the tag literal below is a test fixture, not a real anchor.

import (
	"testing"

	"github.com/dr-zog/plumbline/internal/anchor"
	"github.com/dr-zog/plumbline/internal/register"
)

// TestFixtureEndToEnd wires the packages together over the committed fixture and
// asserts the full C4 coverage model: a clean deep chain plus one of each gap.
func TestFixtureEndToEnd(t *testing.T) {
	items, err := register.ParseFile("../../testdata/fixture/register.md")
	if err != nil {
		t.Fatal(err)
	}
	anchors, files, err := anchor.Scan([]string{"../../testdata/fixture/src"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	r := Build(items, anchors, files, GateOpts{})

	if r.Summary.OK {
		t.Fatal("expected gaps, got OK")
	}
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"registerItems", r.Summary.RegisterItems, 5},
		{"uncovered", r.Summary.UncoveredCount, 1},
		{"transitive", r.Summary.TransitiveCount, 1},
		{"broken", r.Summary.BrokenCount, 1},
		{"orphan", r.Summary.OrphanCount, 1},
		{"structuralErrors", r.Summary.StructuralErrors, 0},
		{"shallowCovered", r.Summary.ShallowCovered, 4},
		{"deepCovered", r.Summary.DeepCovered, 3},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d, want %d", c.name, c.got, c.want)
		}
	}
	if r.Summary.CoveragePct != 80.0 {
		t.Errorf("coveragePct = %.1f, want 80.0", r.Summary.CoveragePct)
	}
	if r.Summary.CompletenessPct != 60.0 {
		t.Errorf("completenessPct = %.1f, want 60.0", r.Summary.CompletenessPct)
	}
	if len(r.Uncovered) == 1 && r.Uncovered[0].ID != "req~rotate-signing-keys~1" {
		t.Errorf("uncovered ID = %q", r.Uncovered[0].ID)
	}
	if len(r.Transitive) == 1 {
		tg := r.Transitive[0]
		if tg.ID != "feat~key-management~1" {
			t.Errorf("transitive ID = %q, want feat~key-management~1", tg.ID)
		}
		if len(tg.Weak) != 1 || tg.Weak[0] != "req~rotate-signing-keys~1" {
			t.Errorf("weak coverers = %v, want [req~rotate-signing-keys~1]", tg.Weak)
		}
	}
}

// TestDeepChainClean checks a fully-covered chain passes strictly.
func TestDeepChainClean(t *testing.T) {
	items := []register.Item{
		{ID: "feat~a~1", Type: "feat", Status: "approved", Needs: []string{"req"}},
		{ID: "req~a~1", Type: "req", Status: "approved", Covers: []string{"feat~a~1"}, Needs: []string{"component"}},
		{ID: "component~a~1", Type: "component", Status: "approved", Covers: []string{"req~a~1"}, Needs: []string{"impl"}},
	}
	anchors := []anchor.Anchor{{File: "a.go", Line: 1, Covering: "impl", TargetID: "component~a~1"}}
	r := Build(items, anchors, []string{"a.go"}, GateOpts{})
	if !r.Summary.OK {
		t.Fatalf("expected OK, got %+v", r.Summary)
	}
	if r.Summary.CompletenessPct != 100.0 {
		t.Errorf("completenessPct = %.1f, want 100.0", r.Summary.CompletenessPct)
	}
}

// TestThresholdGate checks that threshold mode tolerates gaps above the floor
// but broken anchors still fail.
func TestThresholdGate(t *testing.T) {
	items := []register.Item{
		{ID: "req~a~1", Type: "req", Status: "approved", Needs: []string{"component"}},
		{ID: "component~a~1", Type: "component", Status: "approved", Covers: []string{"req~a~1"}, Needs: []string{"impl"}},
		{ID: "req~b~1", Type: "req", Status: "approved", Needs: []string{"component"}}, // uncovered
	}
	anchors := []anchor.Anchor{{File: "a.go", Line: 1, Covering: "impl", TargetID: "component~a~1"}}

	// 1 of 3 items shallow-covered directly... req~a is shallow (component covers
	// it), component~a deep, req~b uncovered => 2/3 = 66.7%.
	strict := Build(items, anchors, []string{"a.go"}, GateOpts{})
	if strict.Summary.OK {
		t.Error("strict: expected FAIL on the uncovered item")
	}
	pass := Build(items, anchors, []string{"a.go"}, GateOpts{MinCoverage: 60})
	if !pass.Summary.OK {
		t.Errorf("threshold 60%%: expected OK at %.1f%%", pass.Summary.CoveragePct)
	}
	fail := Build(items, anchors, []string{"a.go"}, GateOpts{MinCoverage: 90})
	if fail.Summary.OK {
		t.Error("threshold 90%: expected FAIL below the floor")
	}

	// A broken anchor must fail even under a lenient threshold.
	withBroken := append(anchors, anchor.Anchor{File: "b.go", Line: 1, Covering: "impl", TargetID: "component~ghost~1", Raw: "[impl->component~ghost~1]"})
	broken := Build(items, withBroken, []string{"a.go", "b.go"}, GateOpts{MinCoverage: 1})
	if broken.Summary.OK {
		t.Error("broken anchor must fail even in threshold mode")
	}
}

// TestStatusLifecycle checks the OFT status semantics: approved gates and scores,
// proposed is tracked-not-gated, rejected is excluded, and scores are over the
// approved set.
func TestStatusLifecycle(t *testing.T) {
	items := []register.Item{
		// approved, fully covered
		{ID: "req~done~1", Type: "req", Status: "approved", Needs: []string{"component"}},
		{ID: "component~done~1", Type: "component", Status: "approved", Covers: []string{"req~done~1"}, Needs: []string{"impl"}},
		// proposed, uncovered — tracked, must NOT fail the gate
		{ID: "req~planned~1", Type: "req", Status: "proposed", Needs: []string{"component"}},
		// rejected, uncovered — excluded from tracing entirely
		{ID: "req~dropped~1", Type: "req", Status: "rejected", Needs: []string{"component"}},
	}
	anchors := []anchor.Anchor{{File: "a.go", Line: 1, Covering: "impl", TargetID: "component~done~1"}}
	r := Build(items, anchors, []string{"a.go"}, GateOpts{})

	if !r.Summary.OK {
		t.Fatalf("expected OK — proposed uncovered must not fail, rejected excluded; got %+v", r.Summary)
	}
	if r.Summary.RegisterItems != 3 {
		t.Errorf("registerItems = %d, want 3 (rejected excluded)", r.Summary.RegisterItems)
	}
	if r.Summary.RejectedItems != 1 {
		t.Errorf("rejectedItems = %d, want 1", r.Summary.RejectedItems)
	}
	if r.Summary.PlannedItems != 1 {
		t.Errorf("plannedItems = %d, want 1", r.Summary.PlannedItems)
	}
	if r.Summary.ApprovedItems != 2 {
		t.Errorf("approvedItems = %d, want 2", r.Summary.ApprovedItems)
	}
	// Scores are over the approved set (both covered) — planned/rejected don't dilute.
	if r.Summary.CoveragePct != 100.0 || r.Summary.CompletenessPct != 100.0 {
		t.Errorf("scores = %.1f/%.1f, want 100.0/100.0 over approved", r.Summary.CoveragePct, r.Summary.CompletenessPct)
	}
	if r.Summary.UncoveredCount != 0 {
		t.Errorf("uncovered = %d, want 0 (a proposed item is not a gap)", r.Summary.UncoveredCount)
	}
	if len(r.Planned) != 1 || r.Planned[0].ID != "req~planned~1" || r.Planned[0].Realised {
		t.Errorf("planned = %+v, want [req~planned~1, realised=false]", r.Planned)
	}
}

// TestDeadEnd checks that an approved item declaring no Needs is a dead-end (a
// hard fail), not a vacuous 100% — while a proposed item with no Needs is not.
func TestDeadEnd(t *testing.T) {
	items := []register.Item{
		{ID: "req~naked~1", Type: "req", Status: "approved"},                            // approved, no Needs → dead-end
		{ID: "req~ok~1", Type: "req", Status: "approved", Needs: []string{"component"}}, // approved, needs a component
		{ID: "component~ok~1", Type: "component", Status: "approved", Covers: []string{"req~ok~1"}, Needs: []string{"impl"}},
		{ID: "req~future~1", Type: "req", Status: "proposed"}, // proposed, no Needs → NOT a dead-end
	}
	anchors := []anchor.Anchor{{File: "a.go", Line: 1, Covering: "impl", TargetID: "component~ok~1"}}
	r := Build(items, anchors, []string{"a.go"}, GateOpts{})

	if r.Summary.OK {
		t.Fatal("expected FAIL — an approved item with no Needs is a dead-end")
	}
	if r.Summary.DeadEndCount != 1 {
		t.Errorf("deadEndCount = %d, want 1", r.Summary.DeadEndCount)
	}
	if len(r.DeadEnds) != 1 || r.DeadEnds[0].ID != "req~naked~1" {
		t.Errorf("deadEnds = %+v, want [req~naked~1]", r.DeadEnds)
	}
	// The dead-end must NOT count as covered — no vacuous inflation. Of 3 approved
	// items (naked, ok, component~ok), only the two real ones are shallow-covered.
	if r.Summary.ApprovedItems != 3 || r.Summary.ShallowCovered != 2 {
		t.Errorf("approved=%d shallow=%d, want 3 and 2 (dead-end excluded from covered)", r.Summary.ApprovedItems, r.Summary.ShallowCovered)
	}
	// A proposed item with no Needs is planned, never a dead-end.
	if r.Summary.PlannedItems != 1 {
		t.Errorf("plannedItems = %d, want 1 (the proposed item)", r.Summary.PlannedItems)
	}
}

// TestDeadEndCoverageAware checks ADR 007: an approved item with no Needs is a
// dead-end only when nothing covers it. A top-of-axis node covered from below is
// scored through its coverers, not flagged.
func TestDeadEndCoverageAware(t *testing.T) {
	// A context with no Needs, covered from below by a fully-deep chain.
	deepChain := []register.Item{
		{ID: "context~sys~1", Type: "context", Status: "approved"}, // no Needs, covered below
		{ID: "container~svc~1", Type: "container", Status: "approved", Covers: []string{"context~sys~1"}, Needs: []string{"component"}},
		{ID: "component~w~1", Type: "component", Status: "approved", Covers: []string{"container~svc~1"}, Needs: []string{"impl"}},
	}
	anchors := []anchor.Anchor{{File: "a.go", Line: 1, Covering: "impl", TargetID: "component~w~1"}}
	r := Build(deepChain, anchors, []string{"a.go"}, GateOpts{})
	if !r.Summary.OK {
		t.Fatalf("a context covered from below must not fail; summary %+v deadEnds %+v", r.Summary, r.DeadEnds)
	}
	if r.Summary.DeadEndCount != 0 {
		t.Errorf("deadEndCount = %d, want 0 (a covered top-node is not a dead-end)", r.Summary.DeadEndCount)
	}
	if r.Summary.ShallowCovered != 3 || r.Summary.DeepCovered != 3 {
		t.Errorf("shallow=%d deep=%d, want 3 and 3 (context scored through its coverers)", r.Summary.ShallowCovered, r.Summary.DeepCovered)
	}

	// Same context, but its container is not itself covered → context is a transitive
	// defect (shallow but not deep), still NOT a dead-end.
	weakChain := []register.Item{
		{ID: "context~sys~1", Type: "context", Status: "approved"},
		{ID: "container~svc~1", Type: "container", Status: "approved", Covers: []string{"context~sys~1"}, Needs: []string{"component"}}, // uncovered
	}
	r = Build(weakChain, nil, nil, GateOpts{})
	if r.Summary.DeadEndCount != 0 {
		t.Errorf("weak-chain: deadEndCount = %d, want 0 (context has a coverer)", r.Summary.DeadEndCount)
	}
	if r.Summary.OK {
		t.Fatal("weak-chain must fail (container uncovered, context transitive)")
	}
	if r.Summary.ShallowCovered != 1 || r.Summary.DeepCovered != 0 {
		t.Errorf("weak-chain: shallow=%d deep=%d, want 1 and 0 (context shallow via its coverer, not deep)", r.Summary.ShallowCovered, r.Summary.DeepCovered)
	}
	if r.Summary.UncoveredCount != 1 || len(r.Uncovered) != 1 || r.Uncovered[0].ID != "container~svc~1" {
		t.Errorf("weak-chain: uncovered = %+v, want [container~svc~1]", r.Uncovered)
	}

	// A top-of-axis node with no Needs AND no coverers is still a dead-end.
	terminal := []register.Item{{ID: "context~alone~1", Type: "context", Status: "approved"}}
	r = Build(terminal, nil, nil, GateOpts{})
	if r.Summary.DeadEndCount != 1 || len(r.DeadEnds) != 1 || r.DeadEnds[0].ID != "context~alone~1" {
		t.Errorf("terminal: deadEnds = %+v, want [context~alone~1]", r.DeadEnds)
	}
}

// TestGatePolicy checks ADR 004's maturity policy: code against a not-yet-approved
// item is a warning (never a failure), and an anchor on a rejected item is zombie
// code (a hard fail, distinct from a broken anchor).
func TestGatePolicy(t *testing.T) {
	// Status-lag: a proposed item fully covered — a warning, never a failure.
	{
		items := []register.Item{
			{ID: "req~lag~1", Type: "req", Status: "proposed", Needs: []string{"component"}},
			{ID: "component~lag~1", Type: "component", Status: "approved", Covers: []string{"req~lag~1"}, Needs: []string{"impl"}},
		}
		anchors := []anchor.Anchor{{File: "a.go", Line: 1, Covering: "impl", TargetID: "component~lag~1"}}
		r := Build(items, anchors, []string{"a.go"}, GateOpts{})
		if !r.Summary.OK {
			t.Fatalf("status-lag must not fail the gate; got %+v", r.Summary)
		}
		if r.Summary.WarningCount != 1 || len(r.Warnings) != 1 || r.Warnings[0].Kind != "status-lag" {
			t.Errorf("warnings = %+v, want one status-lag", r.Warnings)
		}
	}

	// Build-ahead: a proposed item with shallow-but-not-deep coverage.
	{
		items := []register.Item{
			{ID: "req~ba~1", Type: "req", Status: "proposed", Needs: []string{"component"}},
			{ID: "component~ba~1", Type: "component", Status: "proposed", Covers: []string{"req~ba~1"}, Needs: []string{"impl"}},
		}
		r := Build(items, nil, nil, GateOpts{})
		if !r.Summary.OK {
			t.Fatalf("build-ahead must not fail the gate; got %+v", r.Summary)
		}
		if r.Summary.WarningCount != 1 || r.Warnings[0].Kind != "build-ahead" || r.Warnings[0].ID != "req~ba~1" {
			t.Errorf("warnings = %+v, want one build-ahead on req~ba~1", r.Warnings)
		}
	}

	// Zombie: an anchor on a rejected item — a hard fail, not a broken anchor.
	{
		items := []register.Item{
			{ID: "req~gone~1", Type: "req", Status: "rejected", Needs: []string{"component"}},
		}
		anchors := []anchor.Anchor{{File: "z.go", Line: 3, Covering: "impl", TargetID: "req~gone~1", Raw: "[impl->req~gone~1]"}}
		r := Build(items, anchors, []string{"z.go"}, GateOpts{})
		if r.Summary.OK {
			t.Fatal("zombie code must fail the gate")
		}
		if r.Summary.ZombieCount != 1 || r.Summary.BrokenCount != 0 {
			t.Errorf("zombie=%d broken=%d, want 1 and 0 (not a broken anchor)", r.Summary.ZombieCount, r.Summary.BrokenCount)
		}
	}
}

// TestSpecDebtBudget checks the spec-debt count/ratio over the feat/req axis and
// the optional budget gate (count and percentage).
func TestSpecDebtBudget(t *testing.T) {
	items := []register.Item{
		{ID: "req~a~1", Type: "req", Status: "approved", Needs: []string{"component"}}, // approved
		{ID: "component~a~1", Type: "component", Status: "approved", Covers: []string{"req~a~1"}, Needs: []string{"impl"}},
		{ID: "req~b~1", Type: "req", Status: "proposed", Needs: []string{"component"}}, // un-built spec
		{ID: "req~c~1", Type: "req", Status: "draft", Needs: []string{"component"}},    // un-built spec
	}
	anchors := []anchor.Anchor{{File: "a.go", Line: 1, Covering: "impl", TargetID: "component~a~1"}}

	// feat/req axis = req~a, req~b, req~c (3); component~a is off-axis. Un-built spec
	// = req~b + req~c (not-yet-approved, not realised) = 2 → 66.7%.
	base := Build(items, anchors, []string{"a.go"}, GateOpts{})
	if base.Summary.SpecDebtCount != 2 || base.Summary.SpecTotalItems != 3 {
		t.Fatalf("specDebt = %d/%d, want 2/3", base.Summary.SpecDebtCount, base.Summary.SpecTotalItems)
	}
	if !base.Summary.OK {
		t.Error("no budget set → must not fail on spec-debt")
	}

	one, two, fifty := 1, 2, 50.0
	if Build(items, anchors, []string{"a.go"}, GateOpts{MaxProposed: &one}).Summary.OK {
		t.Error("spec-debt 2 over budget ≤1 must fail")
	}
	if !Build(items, anchors, []string{"a.go"}, GateOpts{MaxProposed: &two}).Summary.OK {
		t.Error("spec-debt 2 within budget ≤2 must pass")
	}
	if Build(items, anchors, []string{"a.go"}, GateOpts{MaxProposedPct: &fifty}).Summary.OK {
		t.Error("spec-debt 66.7% over budget ≤50% must fail")
	}
}
