package register

// [utest->component~register-parser~1]

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	src := "" +
		"### Validate\n" +
		"`req~validate-auth-request~1`\n" +
		"\n" +
		"Some description text on its own line.\n" +
		"\n" +
		"Needs: impl\n" +
		"\n" +
		"### Rotate\n" +
		"`req~rotate-signing-keys~1`\n" +
		"Needs: impl, dsn\n" +
		"Covers: feat~key-management~1\n"

	items, err := parse(strings.NewReader(src), "reg.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}

	first := items[0]
	if first.ID != "req~validate-auth-request~1" || first.Title != "Validate" {
		t.Errorf("first = %+v", first)
	}
	if len(first.Needs) != 1 || first.Needs[0] != "impl" {
		t.Errorf("first.Needs = %v, want [impl]", first.Needs)
	}

	second := items[1]
	if second.Rev != 1 || second.Name != "rotate-signing-keys" {
		t.Errorf("second = %+v", second)
	}
	if strings.Join(second.Needs, ",") != "impl,dsn" {
		t.Errorf("second.Needs = %v, want [impl dsn]", second.Needs)
	}
	if len(second.Covers) != 1 || second.Covers[0] != "feat~key-management~1" {
		t.Errorf("second.Covers = %v", second.Covers)
	}
}

func TestOffOn(t *testing.T) {
	src := "" +
		"`req~kept~1`\n" +
		"<!-- oft:off -->\n" +
		"`req~ignored~1`\n" +
		"<!-- oft:on -->\n" +
		"`req~also-kept~1`\n"
	items, err := parse(strings.NewReader(src), "reg.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2 (ignored region skipped)", len(items))
	}
}

func TestStatus(t *testing.T) {
	src := "" +
		"`req~default~1`\n" +
		"Some description.\n" +
		"Needs: component\n" +
		"\n" +
		"`req~planned~1`\n" +
		"Status: proposed\n" +
		"A speculative requirement.\n" +
		"Needs: component\n"
	items, err := parse(strings.NewReader(src), "reg.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}

	// An absent status defaults to approved.
	if got := items[0].StatusOrDefault(); got != "approved" {
		t.Errorf("default status = %q, want approved", got)
	}
	if !items[0].Approved() {
		t.Error("item without a status line should be Approved()")
	}

	// A Status line right under the ID must be recognised as an attribute, not
	// swallowed into the description.
	if items[1].Status != "proposed" || !items[1].Planned() {
		t.Errorf("second = %+v, want status proposed / Planned()", items[1])
	}
	if items[1].Desc != "A speculative requirement." {
		t.Errorf("desc = %q — a Status line must not become the description", items[1].Desc)
	}
}
