// Package register parses a Plumbline register: an OpenFastTrace-native Markdown
// document declaring the specification items (requirements, features, C4
// components) that anchors resolve up to. Items follow OFT's format — a
// backtick-wrapped ID line `type~name~revision`, optionally preceded by a
// Markdown heading that supplies a title, and followed by Needs/Covers
// attribute lines. Regions between `<!-- oft:off -->` and `<!-- oft:on -->` are
// ignored.
//
// [impl->component~register-parser~1]
package register

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Item is one specification item declared in the register.
type Item struct {
	ID     string   `json:"id"` // full "type~name~rev"
	Type   string   `json:"type"`
	Name   string   `json:"name"`
	Rev    int      `json:"rev"`
	Title  string   `json:"title,omitempty"`
	Desc   string   `json:"desc,omitempty"` // first description line after the ID
	Needs  []string `json:"needs,omitempty"`
	Covers []string `json:"covers,omitempty"`
	Status string   `json:"status,omitempty"` // OFT lifecycle: approved (default) / proposed / draft / rejected
	File   string   `json:"file"`
	Line   int      `json:"line"`
}

// statusValue returns the item's lifecycle status, defaulting to "approved" when
// none is declared — OFT's default, and also the sensible reading of an item
// constructed without a status.
func (it Item) statusValue() string {
	if it.Status == "" {
		return "approved"
	}
	return strings.ToLower(it.Status)
}

// Approved reports whether the item is gated — the default when no status is set.
func (it Item) Approved() bool { return it.statusValue() == "approved" }

// Planned reports whether the item is tracked but not gated (proposed or draft).
func (it Item) Planned() bool {
	s := it.statusValue()
	return s == "proposed" || s == "draft"
}

// Rejected reports whether the item is excluded from tracing entirely.
func (it Item) Rejected() bool { return it.statusValue() == "rejected" }

// StatusOrDefault is the normalised status, "approved" when none is declared.
func (it Item) StatusOrDefault() string { return it.statusValue() }

// KnownStatus reports whether s is a status OFT defines; an empty string is the
// approved default and therefore known.
func KnownStatus(s string) bool {
	switch strings.ToLower(s) {
	case "", "approved", "proposed", "draft", "rejected":
		return true
	}
	return false
}

// idLineRe matches a standalone spec-item ID, optionally wrapped in backticks:
//
//	`req~this-is-the-id~1`
var idLineRe = regexp.MustCompile("^\\s*`?([A-Za-z][A-Za-z0-9_]*)~([A-Za-z0-9_.\\-]+)~(\\d+)`?\\s*$")

var (
	headingRe = regexp.MustCompile(`^#{1,6}\s+(.*\S)\s*$`)
	attrRe    = regexp.MustCompile(`(?i)^\s*(Needs|Covers|Depends|Tags|Status):\s*(.*)$`)
	offRe     = regexp.MustCompile(`^\s*<!--\s*oft:off\s*-->\s*$`)
	onRe      = regexp.MustCompile(`^\s*<!--\s*oft:on\s*-->\s*$`)
)

// ParseFile reads and parses a register file.
func ParseFile(path string) ([]Item, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse(f, path)
}

func parse(r io.Reader, path string) ([]Item, error) {
	var items []Item
	var cur *Item // item currently accumulating attribute lines
	lastHeading := ""
	off := false
	line := 0

	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
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

		if m := idLineRe.FindStringSubmatch(text); m != nil {
			rev, _ := strconv.Atoi(m[3])
			items = append(items, Item{
				ID:     m[1] + "~" + m[2] + "~" + m[3],
				Type:   m[1],
				Name:   m[2],
				Rev:    rev,
				Title:  lastHeading,
				Status: "approved",
				File:   path,
				Line:   line,
			})
			cur = &items[len(items)-1]
			lastHeading = ""
			continue
		}

		if m := headingRe.FindStringSubmatch(text); m != nil {
			lastHeading = m[1]
			cur = nil // a heading opens a new section; stop accumulating
			continue
		}

		if m := attrRe.FindStringSubmatch(text); m != nil && cur != nil {
			key := strings.ToLower(m[1])
			if key == "status" {
				if v := strings.ToLower(strings.TrimSpace(m[2])); v != "" {
					cur.Status = v
				}
				continue
			}
			vals := splitList(m[2])
			switch key {
			case "needs":
				cur.Needs = append(cur.Needs, vals...)
			case "covers":
				cur.Covers = append(cur.Covers, vals...)
			}
			continue
		}

		// First non-empty prose line after an item ID becomes its description.
		if cur != nil && cur.Desc == "" {
			if s := strings.TrimSpace(text); s != "" {
				cur.Desc = s
			}
		}
	}
	return items, sc.Err()
}

func splitList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
