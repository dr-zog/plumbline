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

	r := Build(items, anchors, files, 0)

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
		{ID: "feat~a~1", Type: "feat", Needs: []string{"req"}},
		{ID: "req~a~1", Type: "req", Covers: []string{"feat~a~1"}, Needs: []string{"component"}},
		{ID: "component~a~1", Type: "component", Covers: []string{"req~a~1"}, Needs: []string{"impl"}},
	}
	anchors := []anchor.Anchor{{File: "a.go", Line: 1, Covering: "impl", TargetID: "component~a~1"}}
	r := Build(items, anchors, []string{"a.go"}, 0)
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
		{ID: "req~a~1", Type: "req", Needs: []string{"component"}},
		{ID: "component~a~1", Type: "component", Covers: []string{"req~a~1"}, Needs: []string{"impl"}},
		{ID: "req~b~1", Type: "req", Needs: []string{"component"}}, // uncovered
	}
	anchors := []anchor.Anchor{{File: "a.go", Line: 1, Covering: "impl", TargetID: "component~a~1"}}

	// 1 of 3 items shallow-covered directly... req~a is shallow (component covers
	// it), component~a deep, req~b uncovered => 2/3 = 66.7%.
	strict := Build(items, anchors, []string{"a.go"}, 0)
	if strict.Summary.OK {
		t.Error("strict: expected FAIL on the uncovered item")
	}
	pass := Build(items, anchors, []string{"a.go"}, 60)
	if !pass.Summary.OK {
		t.Errorf("threshold 60%%: expected OK at %.1f%%", pass.Summary.CoveragePct)
	}
	fail := Build(items, anchors, []string{"a.go"}, 90)
	if fail.Summary.OK {
		t.Error("threshold 90%: expected FAIL below the floor")
	}

	// A broken anchor must fail even under a lenient threshold.
	withBroken := append(anchors, anchor.Anchor{File: "b.go", Line: 1, Covering: "impl", TargetID: "component~ghost~1", Raw: "[impl->component~ghost~1]"})
	broken := Build(items, withBroken, []string{"a.go", "b.go"}, 1)
	if broken.Summary.OK {
		t.Error("broken anchor must fail even in threshold mode")
	}
}
