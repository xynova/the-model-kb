package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	derrors "github.com/xynova/library-intake/internal/errors"
)

// BackfillSummaries writes summary: into entry frontmatter when missing, using
// the first sentence under ## What it is. Returns how many files were updated.
func BackfillSummaries(libraryRoot string) (int, error) {
	return backfillByWalk(libraryRoot)
}

func backfillByWalk(libraryRoot string) (int, error) {
	op := "catalog.backfillByWalk"
	updated := 0
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
		if !strings.HasSuffix(d.Name(), ".md") || skipEntryNames[d.Name()] {
			return nil
		}
		rel, err := filepath.Rel(libraryRoot, path)
		if err != nil {
			return derrors.Wrap(err, derrors.CodeFailed, op, "rel path").With("path", path)
		}
		if len(strings.Split(filepath.ToSlash(rel), "/")) != 2 {
			return nil
		}
		changed, err := ensureSummaryFile(path)
		if err != nil {
			return err
		}
		if changed {
			updated++
		}
		return nil
	})
	if err != nil {
		return updated, derrors.Wrap(err, derrors.CodeFailed, op, "walk library").
			With("library", libraryRoot)
	}
	return updated, nil
}

func ensureSummaryFile(absPath string) (bool, error) {
	op := "catalog.ensureSummaryFile"
	rawBytes, err := os.ReadFile(absPath)
	if err != nil {
		return false, derrors.Wrap(err, derrors.CodeNotFound, op, "read entry").With("path", absPath)
	}
	raw := string(rawBytes)
	fm, err := splitFrontmatter(raw)
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(fm, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "summary:") {
			return false, nil
		}
	}
	summary := FirstWhatItIsSentence(raw)
	if summary == "" {
		return false, derrors.New(derrors.CodeParse, op, "could not derive summary from What it is").
			With("path", absPath)
	}
	injected, err := injectSummaryFrontmatter(raw, summary)
	if err != nil {
		return false, err
	}
	if err := os.WriteFile(absPath, []byte(injected), 0o644); err != nil {
		return false, derrors.Wrap(err, derrors.CodeFailed, op, "write entry").With("path", absPath)
	}
	return true, nil
}

func injectSummaryFrontmatter(raw, summary string) (string, error) {
	op := "catalog.injectSummaryFrontmatter"
	quoted := fmt.Sprintf("%q", summary)
	line := "summary: " + quoted + "\n"
	// Insert after title: line inside frontmatter.
	const open = "---\n"
	if !strings.HasPrefix(raw, open) && !strings.HasPrefix(raw, "---\r\n") {
		return "", derrors.New(derrors.CodeParse, op, "missing frontmatter")
	}
	rest := raw
	prefix := "---\n"
	if strings.HasPrefix(raw, "---\r\n") {
		prefix = "---\r\n"
		rest = raw[len(prefix):]
	} else {
		rest = raw[len(open):]
	}
	end := strings.Index(rest, "\n---\n")
	endLen := 5
	nl := "\n"
	if end < 0 {
		end = strings.Index(rest, "\r\n---\r\n")
		endLen = 7
		nl = "\r\n"
	}
	if end < 0 {
		end = strings.Index(rest, "\n---\r\n")
		endLen = 6
	}
	if end < 0 {
		return "", derrors.New(derrors.CodeParse, op, "missing closing frontmatter")
	}
	fm := rest[:end]
	body := rest[end+endLen:]
	lines := strings.Split(fm, nl)
	out := make([]string, 0, len(lines)+1)
	inserted := false
	for _, l := range lines {
		out = append(out, l)
		trim := strings.TrimSpace(l)
		if !inserted && strings.HasPrefix(trim, "title:") {
			out = append(out, strings.TrimSuffix(line, "\n"))
			inserted = true
		}
	}
	if !inserted {
		out = append([]string{out[0], strings.TrimSuffix(line, "\n")}, out[1:]...)
	}
	return prefix + strings.Join(out, nl) + nl + "---" + nl + body, nil
}
