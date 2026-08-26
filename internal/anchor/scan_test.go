package anchor

// [utest->component~anchor-scanner~1]
//
// oft:off — the tag literals below are test fixtures, not real anchors.

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTagGrammar(t *testing.T) {
	cases := []struct {
		line      string
		cover     string
		coverName string
		coverRev  int
		target    string
		needs     []string
	}{
		{line: "// [impl->req~validate-auth-request~1]", cover: "impl", target: "req~validate-auth-request~1"},
		{line: "# [impl~~2->dsn~foo~1]", cover: "impl", coverRev: 2, target: "dsn~foo~1"},
		{line: "// [impl~validate-password~2->dsn~foo~1]", cover: "impl", coverName: "validate-password", coverRev: 2, target: "dsn~foo~1"},
		{line: "  <!-- [doc->req~user-guide~1] -->", cover: "doc", target: "req~user-guide~1"},
		{line: "' [dsn->req~1password-login~1>>impl,test]", cover: "dsn", target: "req~1password-login~1", needs: []string{"impl", "test"}},
	}
	for _, c := range cases {
		m := tagRe.FindStringSubmatch(c.line)
		if m == nil {
			t.Fatalf("no match for %q", c.line)
		}
		if m[1] != c.cover {
			t.Errorf("%q: cover = %q, want %q", c.line, m[1], c.cover)
		}
		if m[4] != c.target {
			t.Errorf("%q: target = %q, want %q", c.line, m[4], c.target)
		}
	}
}

func TestScanFixture(t *testing.T) {
	anchors, files, err := Scan([]string{"../../testdata/fixture/src"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Fatalf("scanned files = %d, want 3", len(files))
	}
	if len(anchors) != 2 {
		t.Fatalf("anchors = %d, want 2 (good + broken)", len(anchors))
	}
	targets := map[string]bool{}
	for _, a := range anchors {
		targets[a.TargetID] = true
	}
	if !targets["component~auth-validator~1"] || !targets["component~does-not-exist~1"] {
		t.Errorf("unexpected anchor targets: %v", targets)
	}
}

// TestMarkdownScanned proves Markdown files carry anchors (in an HTML comment),
// so prose-as-implementation like Claude Code skills can be traced.
func TestMarkdownScanned(t *testing.T) {
	dir := t.TempDir()
	md := "---\nname: onboard\n---\n\n<!-- [impl->component~onboard-skill~1] -->\n\n# Onboard\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}
	anchors, files, err := Scan([]string{dir}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || len(anchors) != 1 {
		t.Fatalf("files=%d anchors=%d, want 1 and 1", len(files), len(anchors))
	}
	if anchors[0].TargetID != "component~onboard-skill~1" || anchors[0].Covering != "impl" {
		t.Errorf("anchor = %+v", anchors[0])
	}
}

// TestOffOnFence proves oft:off/oft:on suppress example tags, so a codebase can
// document the anchor syntax without those examples becoming real anchors.
func TestOffOnFence(t *testing.T) {
	dir := t.TempDir()
	src := "package x\n" +
		"// [impl->component~real~1]\n" +
		"// oft:off\n" +
		"// [impl->component~example~1]\n" +
		"// oft:on\n" +
		"// [impl->component~real2~1]\n"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	anchors, _, err := Scan([]string{dir}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, a := range anchors {
		got[a.TargetID] = true
	}
	if !got["component~real~1"] || !got["component~real2~1"] {
		t.Errorf("real anchors missing: %v", got)
	}
	if got["component~example~1"] {
		t.Error("fenced example tag was recorded as an anchor")
	}
}

// TestExcludeDirPrefix proves an exclude entry skips a whole directory subtree.
func TestExcludeDirPrefix(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "keep"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "skip", "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(p string) {
		if err := os.WriteFile(filepath.Join(dir, p), []byte("// [impl->component~x~1]\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join("keep", "a.go"))
	write(filepath.Join("skip", "b.go"))
	write(filepath.Join("skip", "nested", "c.go"))

	exclude := map[string]bool{filepath.Clean(filepath.Join(dir, "skip")): true}
	_, files, err := Scan([]string{dir}, exclude)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || filepath.Base(files[0]) != "a.go" {
		t.Fatalf("scanned %v, want only keep/a.go", files)
	}
}
