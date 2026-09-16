package catalog

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	derrors "github.com/xynova/library-intake/internal/errors"
)

var skipEntryNames = map[string]bool{
	"README.md":    true,
	"INDEX.md":     true,
	"taxonomy.md":  true,
	"_template.md": true,
	"_index.md":    true,
}

// ScanEntries walks libraryRoot for category entry markdown files.
func ScanEntries(libraryRoot string) ([]Entry, error) {
	op := "catalog.ScanEntries"
	entries := make([]Entry, 0)
	err := filepath.WalkDir(libraryRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			base := d.Name()
			if base == "media" || base == "additions" || strings.HasPrefix(base, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		if skipEntryNames[d.Name()] {
			return nil
		}
		rel, err := filepath.Rel(libraryRoot, path)
		if err != nil {
			return derrors.Wrap(err, derrors.CodeFailed, op, "rel path").With("path", path)
		}
		rel = filepath.ToSlash(rel)
		// Only one-level category entries: <category>/<slug>.md
		parts := strings.Split(rel, "/")
		if len(parts) != 2 {
			return nil
		}
		e, err := LoadEntry(libraryRoot, rel)
		if err != nil {
			return err
		}
		entries = append(entries, e)
		return nil
	})
	if err != nil {
		return nil, derrors.Wrap(err, derrors.CodeFailed, op, "walk library").
			With("library", libraryRoot)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Category != entries[j].Category {
			return entries[i].Category < entries[j].Category
		}
		return entries[i].Name < entries[j].Name
	})
	return entries, nil
}

// CategoriesFromTaxonomy lists category directory names that exist under libraryRoot
// (including empty ones) by reading taxonomy.md tables is overkill; use known dirs.
func ListCategoryDirs(libraryRoot string) ([]string, error) {
	op := "catalog.ListCategoryDirs"
	ents, err := os.ReadDir(libraryRoot)
	if err != nil {
		return nil, derrors.Wrap(err, derrors.CodeNotFound, op, "read library root").
			With("library", libraryRoot)
	}
	out := make([]string, 0)
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "additions" || name == "media" || strings.HasPrefix(name, ".") {
			continue
		}
		out = append(out, name)
	}
	sort.Strings(out)
	return out, nil
}
