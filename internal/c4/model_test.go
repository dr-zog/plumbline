package c4

// [utest->component~c4-validator~1]

import (
	"testing"

	"github.com/dr-zog/plumbline/internal/register"
)

func TestValidateCleanLadder(t *testing.T) {
	items := []register.Item{
		{ID: "feat~a~1", Type: "feat", Needs: []string{"req"}},
		{ID: "req~a~1", Type: "req", Covers: []string{"feat~a~1"}, Needs: []string{"component"}},
		{ID: "container~api~1", Type: "container", Desc: "Go service", Needs: []string{"component"}},
		{ID: "component~a~1", Type: "component", Covers: []string{"req~a~1"}, Needs: []string{"impl"}},
	}
	if vs := Errors(Validate(items)); len(vs) != 0 {
		t.Fatalf("expected no errors, got %+v", vs)
	}
}

func TestValidateViolations(t *testing.T) {
	items := []register.Item{
		{ID: "req~bad-need~1", Type: "req", Needs: []string{"impl"}},                // req may not Need impl
		{ID: "widget~x~1", Type: "widget"},                                          // unknown type
		{ID: "component~x~1", Type: "component", Covers: []string{"ghost~y~1"}},     // dangling covers
		{ID: "feat~x~1", Type: "feat", Covers: []string{"component~x~1"}},           // feat may not cover component
		{ID: "container~notech~1", Type: "container", Needs: []string{"component"}}, // warning: no tech
	}
	vs := Validate(items)

	kinds := map[string]Severity{}
	for _, v := range vs {
		kinds[v.Kind] = v.Severity
	}
	for _, want := range []string{"invalid-needs", "unknown-type", "dangling-covers", "invalid-covers"} {
		if kinds[want] != Error {
			t.Errorf("expected error kind %q, got %v", want, kinds[want])
		}
	}
	if kinds["container-missing-tech"] != Warning {
		t.Errorf("expected container-missing-tech warning, got %v", kinds["container-missing-tech"])
	}
	if got := len(Errors(vs)); got != 4 {
		t.Errorf("error count = %d, want 4 (warning excluded)", got)
	}
}

func TestDuplicateID(t *testing.T) {
	items := []register.Item{
		{ID: "req~a~1", Type: "req", Needs: []string{"component"}, File: "a.md", Line: 1},
		{ID: "req~a~1", Type: "req", Needs: []string{"component"}, File: "b.md", Line: 9},
	}
	vs := Validate(items)
	dups := 0
	for _, v := range vs {
		if v.Kind == "duplicate-id" {
			dups++
			if v.Severity != Error {
				t.Errorf("duplicate-id should be an error")
			}
		}
	}
	if dups != 1 {
		t.Fatalf("duplicate-id violations = %d, want 1 (the redefinition)", dups)
	}
}
