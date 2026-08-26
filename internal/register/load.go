package register

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Load parses one or more register sources — files, directories (recursed for
// .md/.markdown), or globs — and returns the aggregated specification items plus
// the cleaned list of register files actually parsed. Callers pass the file list
// to the anchor scanner's exclude set so register prose isn't also treated as a
// code-area. This is how OpenFastTrace works: a register is a *set* of documents
// across a tree, not one file, so it scales to a per-level or per-container split.
func Load(paths []string) (items []Item, files []string, err error) {
	seen := map[string]bool{}
	for _, raw := range paths {
		resolved := []string{raw}
		if strings.ContainsAny(raw, "*?[") {
			m, gErr := filepath.Glob(raw)
			if gErr != nil {
				return nil, nil, fmt.Errorf("register glob %q: %w", raw, gErr)
			}
			if len(m) == 0 {
				return nil, nil, fmt.Errorf("register glob %q matched nothing", raw)
			}
			resolved = m
		}
		for _, p := range resolved {
			if err := loadPath(p, &items, &files, seen); err != nil {
				return nil, nil, err
			}
		}
	}
	return items, files, nil
}

func loadPath(path string, items *[]Item, files *[]string, seen map[string]bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return parseInto(path, items, files, seen)
	}
	return filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !isRegisterFile(p) {
			return nil
		}
		return parseInto(p, items, files, seen)
	})
}

func parseInto(path string, items *[]Item, files *[]string, seen map[string]bool) error {
	clean := filepath.Clean(path)
	if seen[clean] {
		return nil
	}
	seen[clean] = true
	got, err := ParseFile(path)
	if err != nil {
		return err
	}
	*items = append(*items, got...)
	*files = append(*files, clean)
	return nil
}

func isRegisterFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown":
		return true
	}
	return false
}
