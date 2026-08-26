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
