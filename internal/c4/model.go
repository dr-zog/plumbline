// Package c4 defines Plumbline's locked artifact-type ladder — the C4 model plus
// the requirements axis — and validates a register's structure against it. The
// ladder is deliberately NOT configurable: Plumbline adopts the prior art
// (C4 + OpenFastTrace) without bespoke variation.
//
// Two axes meet at the component:
//
//	requirements:   feat ── needs ──▶ req ──┐
//	                                         ├─ needs ─▶ component ─ needs ─▶ impl | utest | itest
//	architecture:   context ─▶ container ────┘                                     (code anchors)
//
// "Needs" flows from the abstract down to the concrete. Code anchors attach at
// the component (the C4 Code level, one altitude up from a symbol). Because the
// ladder edges only ever point downward, the coverage graph is acyclic by
// construction — no cycle detection is required.
//
// [impl->component~c4-validator~1]
package c4

import (
	"fmt"

	"github.com/dr-zog/plumbline/internal/register"
)

// Type is a known artifact type and the types it may require coverage from.
type Type struct {
	Name  string
	Axis  string   // "requirements", "architecture" or "code"
	Needs []string // artifact types this type is allowed to Need
	Leaf  bool     // a code leaf: covered by anchors, never by register items
}

// ladder is the locked type model. A type may only Need the types listed here.
var ladder = map[string]Type{
	"feat":      {Name: "feat", Axis: "requirements", Needs: []string{"req"}},
	"req":       {Name: "req", Axis: "requirements", Needs: []string{"component"}},
	"context":   {Name: "context", Axis: "architecture", Needs: []string{"container"}},
	"container": {Name: "container", Axis: "architecture", Needs: []string{"component"}},
	"component": {Name: "component", Axis: "architecture", Needs: []string{"impl", "utest", "itest"}},
	"impl":      {Name: "impl", Axis: "code", Leaf: true},
	"utest":     {Name: "utest", Axis: "code", Leaf: true},
	"itest":     {Name: "itest", Axis: "code", Leaf: true},
}

// Known reports whether t is a type in the locked ladder.
func Known(t string) bool { _, ok := ladder[t]; return ok }

// IsLeaf reports whether t is a code leaf (covered by anchors, not by register
// items). Anchors may only use leaf covering types.
func IsLeaf(t string) bool { return ladder[t].Leaf }

// mayNeed reports whether a covering of type cov satisfies a Need of type need.
func mayNeed(parent, child string) bool {
	for _, n := range ladder[parent].Needs {
		if n == child {
			return true
		}
	}
	return false
}

// Severity distinguishes hard errors (fail the build) from advisory warnings.
type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
)

// Violation is one structural problem found in the register.
type Violation struct {
	ItemID   string   `json:"itemId"`
	Kind     string   `json:"kind"`
	Detail   string   `json:"detail"`
	Severity Severity `json:"severity"`
	File     string   `json:"file"`
	Line     int      `json:"line"`
}

// Validate checks every register item against the locked ladder: unknown types,
// disallowed Needs edges, and Covers links whose target is missing or of a type
// this item may not cover. Containers that name no technology yield a warning.
func Validate(items []register.Item) []Violation {
	byID := make(map[string]register.Item, len(items))
	for _, it := range items {
		byID[it.ID] = it
	}

	var vs []Violation

	// Duplicate IDs — the hazard when a register spans multiple files. Report
	// every redefinition, pointing back at where the ID was first declared.
	firstSeen := make(map[string]register.Item, len(items))
	for _, it := range items {
		if prev, ok := firstSeen[it.ID]; ok {
			vs = append(vs, Violation{it.ID, "duplicate-id",
				fmt.Sprintf("already defined at %s:%d", prev.File, prev.Line), Error, it.File, it.Line})
			continue
		}
		firstSeen[it.ID] = it
	}

	for _, it := range items {
		if !register.KnownStatus(it.Status) {
			vs = append(vs, Violation{it.ID, "unknown-status",
				fmt.Sprintf("%q is not a known status (approved/proposed/draft/rejected)", it.Status), Error, it.File, it.Line})
		}
		if !Known(it.Type) {
			vs = append(vs, Violation{it.ID, "unknown-type",
				fmt.Sprintf("%q is not a Plumbline artifact type", it.Type), Error, it.File, it.Line})
			continue
		}
		for _, n := range it.Needs {
			if !Known(n) {
				vs = append(vs, Violation{it.ID, "invalid-needs",
					fmt.Sprintf("%s may not Need unknown type %q", it.Type, n), Error, it.File, it.Line})
			} else if !mayNeed(it.Type, n) {
				vs = append(vs, Violation{it.ID, "invalid-needs",
					fmt.Sprintf("%s may not Need %s", it.Type, n), Error, it.File, it.Line})
			}
		}
		for _, cov := range it.Covers {
			target, ok := byID[cov]
			if !ok {
				vs = append(vs, Violation{it.ID, "dangling-covers",
					fmt.Sprintf("Covers %s, which is not in the register", cov), Error, it.File, it.Line})
				continue
			}
			if !mayNeed(target.Type, it.Type) {
				vs = append(vs, Violation{it.ID, "invalid-covers",
					fmt.Sprintf("%s may not cover %s", it.Type, target.Type), Error, it.File, it.Line})
			}
		}
		// C4 Container level names chosen technology (advisory).
		if it.Type == "container" && it.Desc == "" {
			vs = append(vs, Violation{it.ID, "container-missing-tech",
				"container should name its chosen technology in a description", Warning, it.File, it.Line})
		}
	}
	return vs
}

// Errors returns only the fail-the-build violations.
func Errors(vs []Violation) []Violation {
	var out []Violation
	for _, v := range vs {
		if v.Severity == Error {
			out = append(out, v)
		}
	}
	return out
}
