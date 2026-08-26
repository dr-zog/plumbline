// Package anchor scans a source tree for Plumbline anchors — OFT-compatible
// coverage tags embedded in code comments, of the form:
//
// oft:off — the example below documents the grammar; it is not a real anchor.
//
//	[impl->req~validate-auth-request~1]
//
// oft:on
//
// The grammar follows OpenFastTrace's Tag Importer: a covering artifact type,
// an optional name/revision on the covering side, an arrow, the full ID of the
// specification item the code-area covers, and an optional ">>" needed-coverage
// list. Like OFT, Plumbline does not parse the surrounding language; it looks
// for the pattern, which by convention lives in a comment.
//
// [impl->component~anchor-scanner~1]
package anchor

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Anchor is one coverage tag found in a source file.
type Anchor struct {
	File       string   `json:"file"`
	Line       int      `json:"line"`
	Covering   string   `json:"covering"` // covering artifact type, e.g. "impl"
	CoverName  string   `json:"coverName,omitempty"`
	CoverRev   int      `json:"coverRev,omitempty"`
	TargetID   string   `json:"targetId"` // full "type~name~rev"
	TargetType string   `json:"targetType"`
	TargetName string   `json:"targetName"`
	TargetRev  int      `json:"targetRev"`
	Needs      []string `json:"needs,omitempty"` // optional ">>" list
	Raw        string   `json:"raw"`             // the matched tag text
}

// tagRe matches the OFT coverage-tag grammar. Capture groups:
//
//	1 covering artifact type
//	2 optional covering name (may be empty for the "~~rev" form)
//	3 optional covering revision
//	4 target specification-item ID
//	5 optional ">>" needed-coverage list
var tagRe = regexp.MustCompile(
	`\[\s*` +
		`([A-Za-z][A-Za-z0-9_]*)` + // 1 covering type
		`(?:~([A-Za-z0-9_.\-]*)~(\d+))?` + // 2,3 optional cover name + rev
		`\s*->\s*` +
		`([A-Za-z][A-Za-z0-9_]*~[A-Za-z0-9_.\-]+~\d+)` + // 4 target id
		`(?:\s*>>\s*([A-Za-z0-9_,\s]+?))?` + // 5 optional needs list
		`\s*\]`,
)

// idRe splits a full specification-item ID into type, name and revision.
var idRe = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9_]*)~([A-Za-z0-9_.\-]+)~(\d+)$`)

// offRe / onRe fence a region the scanner must ignore, so a codebase can *document*
// the tag syntax (examples in doc comments, tests, skills) without those examples
// being read as real anchors. Mirrors the register parser's `oft:off`/`oft:on`.
//
// A line only toggles the fence when it *starts* (after comment punctuation) with
// the directive — so a mid-sentence mention or a quoted string literal containing
// the token doesn't trip it, but "// oft:off", "# oft:off" and "<!-- oft:off -->"
// (with or without a trailing note) do. The comment style doesn't matter.
var (
	offRe = regexp.MustCompile(`(?i)^[/#;'*<!\-\s]*oft:off\b`)
	onRe  = regexp.MustCompile(`(?i)^[/#;'*<!\-\s]*oft:on\b`)
)

// SourceExtensions are the file types Plumbline scans for anchors. The set is a
// practical subset of OpenFastTrace's Tag Importer languages. Markdown is
// included so prose artefacts that *are* the implementation — Claude Code
// skills, docs-as-code — can carry anchors in an HTML comment; scope prose you
// don't want treated as code-areas out via `roots`/`exclude`.
var SourceExtensions = map[string]bool{
	".go": true, ".py": true, ".js": true, ".mjs": true, ".cjs": true,
	".ts": true, ".tsx": true, ".jsx": true, ".java": true, ".kt": true,
	".rs": true, ".rb": true, ".c": true, ".h": true, ".cc": true,
	".cpp": true, ".hpp": true, ".cs": true, ".php": true, ".swift": true,
	".sh": true, ".bash": true, ".zsh": true, ".lua": true, ".pl": true,
	".sql": true, ".yaml": true, ".yml": true, ".toml": true,
	".md": true, ".markdown": true,
}

// skipDirs are never descended into during a scan.
var skipDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true,
}

// Scan walks the given roots and returns every anchor found, plus the list of
// scanned source files (used to detect unanchored / orphan code-areas). An
// exclude entry matches a path exactly or as a directory prefix, so excluding a
// directory skips everything under it (for example the register file, or a docs
// tree you don't want treated as code-areas).
func Scan(roots []string, exclude map[string]bool) (anchors []Anchor, files []string, err error) {
	seen := map[string]bool{}
	for _, root := range roots {
		walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			clean := filepath.Clean(path)
			if d.IsDir() {
				name := d.Name()
				if path != root && (skipDirs[name] || (strings.HasPrefix(name, ".") && name != ".") || excluded(clean, exclude)) {
					return filepath.SkipDir
				}
				return nil
			}
			if !SourceExtensions[strings.ToLower(filepath.Ext(path))] {
				return nil
			}
			if excluded(clean, exclude) || seen[clean] {
				return nil
			}
			seen[clean] = true
			fileAnchors, scanErr := scanFile(clean)
			if scanErr != nil {
				return scanErr
			}
			files = append(files, clean)
			anchors = append(anchors, fileAnchors...)
			return nil
		})
		if walkErr != nil {
			return nil, nil, walkErr
		}
	}
	return anchors, files, nil
}

// excluded reports whether clean matches an exclude entry exactly or sits under
// one as a directory prefix.
func excluded(clean string, exclude map[string]bool) bool {
	if exclude[clean] {
		return true
	}
	for e := range exclude {
		if strings.HasPrefix(clean, e+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func scanFile(path string) ([]Anchor, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []Anchor
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	line := 0
	off := false
	for sc.Scan() {
		line++
		text := sc.Text()
		switch {
		case offRe.MatchString(text):
			off = true
			continue
		case onRe.MatchString(text):
			off = false
			continue
		}
		if off {
			continue
		}
		for _, m := range tagRe.FindAllStringSubmatch(text, -1) {
			a := Anchor{
				File:     path,
				Line:     line,
				Covering: m[1],
				TargetID: m[4],
				Raw:      m[0],
			}
			if m[3] != "" {
				a.CoverName = m[2]
				a.CoverRev, _ = strconv.Atoi(m[3])
			}
			if im := idRe.FindStringSubmatch(m[4]); im != nil {
				a.TargetType = im[1]
				a.TargetName = im[2]
				a.TargetRev, _ = strconv.Atoi(im[3])
			}
			if m[5] != "" {
				for _, n := range strings.Split(m[5], ",") {
					if n = strings.TrimSpace(n); n != "" {
						a.Needs = append(a.Needs, n)
					}
				}
			}
			out = append(out, a)
		}
	}
	return out, sc.Err()
}
