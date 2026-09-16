package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	derrors "github.com/xynova/library-intake/internal/errors"
	"gopkg.in/yaml.v3"
)

// Entry is the catalog frontmatter fields used for indexes and weekly additions.
type Entry struct {
	Name     string   `yaml:"name"`
	Title    string   `yaml:"title"`
	Summary  string   `yaml:"summary"`
	Category string   `yaml:"category"`
	Does     string   `yaml:"does"`
	Tags     []string `yaml:"tags"`
	Status   string   `yaml:"status"`
	Run      string   `yaml:"run"`
	License  string   `yaml:"license"`
	Added    string   `yaml:"added"`
	Path     string   // repo-relative path from library root, e.g. models/auk.md
}

// SitePermalink is the Hugo public path for this entry (mounted under /library/).
func (e Entry) SitePermalink() string {
	p := strings.TrimSuffix(filepath.ToSlash(e.Path), ".md")
	return "/library/" + p + "/"
}

// AddedTime parses Added as YYYY-MM-DD in local time.
func (e Entry) AddedTime() (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(e.Added), time.Local)
	if err != nil {
		return time.Time{}, derrors.Wrap(err, derrors.CodeParse, "catalog.Entry.AddedTime", "parse added date").
			With("name", e.Name).
			With("added", e.Added)
	}
	return t, nil
}

// ISOWeekID returns YYYY-Www for the entry's added date.
func (e Entry) ISOWeekID() (string, error) {
	t, err := e.AddedTime()
	if err != nil {
		return "", err
	}
	return FormatISOWeek(t), nil
}

// FormatISOWeek returns YYYY-Www for t.
func FormatISOWeek(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}

// CurrentISOWeekID returns the ISO week id for now.
func CurrentISOWeekID(now time.Time) string {
	return FormatISOWeek(now)
}

// TableCell escapes a value for a markdown table cell.
func TableCell(s string) string {
	s = strings.ReplaceAll(s, "|", "/")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}

// LoadEntry reads one library markdown entry and returns its catalog fields.
func LoadEntry(libraryRoot, relPath string) (Entry, error) {
	op := "catalog.LoadEntry"
	abs := filepath.Join(libraryRoot, relPath)
	raw, err := os.ReadFile(abs)
	if err != nil {
		return Entry{}, derrors.Wrap(err, derrors.CodeNotFound, op, "read entry").
			With("path", abs)
	}
	fm, err := splitFrontmatter(string(raw))
	if err != nil {
		return Entry{}, err
	}
	var e Entry
	if err := yaml.Unmarshal([]byte(fm), &e); err != nil {
		return Entry{}, derrors.Wrap(err, derrors.CodeParse, op, "parse frontmatter yaml").
			With("path", relPath)
	}
	e.Name = strings.TrimSpace(e.Name)
	e.Title = strings.TrimSpace(e.Title)
	e.Summary = strings.TrimSpace(e.Summary)
	e.Category = strings.TrimSpace(e.Category)
	e.Does = strings.TrimSpace(e.Does)
	e.Status = strings.TrimSpace(e.Status)
	e.Run = strings.TrimSpace(e.Run)
	e.License = strings.TrimSpace(e.License)
	e.Added = strings.TrimSpace(e.Added)
	e.Path = filepath.ToSlash(relPath)
	if e.Summary == "" {
		if s := FirstWhatItIsSentence(string(raw)); s != "" {
			e.Summary = s
		}
	}
	if e.Name == "" || e.Title == "" || e.Category == "" || e.Does == "" || e.Run == "" || e.License == "" || e.Added == "" || e.Summary == "" {
		return Entry{}, derrors.New(derrors.CodeParse, op, "entry missing required frontmatter fields").
			With("path", relPath).
			With("name", e.Name).
			With("category", e.Category).
			With("does", e.Does).
			With("run", e.Run).
			With("license", e.License).
			With("added", e.Added).
			With("summary", e.Summary)
	}
	return e, nil
}

// FirstWhatItIsSentence returns the first prose sentence under ## What it is.
func FirstWhatItIsSentence(raw string) string {
	const heading = "## What it is"
	idx := strings.Index(raw, heading)
	if idx < 0 {
		return ""
	}
	rest := raw[idx+len(heading):]
	rest = strings.TrimLeft(rest, " \t\r\n")
	var lines []string
	for _, line := range strings.Split(rest, "\n") {
		trim := strings.TrimSpace(line)
		if trim == "" {
			if len(lines) > 0 {
				break
			}
			continue
		}
		if strings.HasPrefix(trim, "#") || strings.HasPrefix(trim, "|") || strings.HasPrefix(trim, "![") {
			if len(lines) > 0 {
				break
			}
			continue
		}
		lines = append(lines, trim)
		joined := strings.Join(lines, " ")
		if strings.Contains(joined, ". ") || strings.HasSuffix(joined, ".") {
			break
		}
	}
	if len(lines) == 0 {
		return ""
	}
	text := strings.Join(lines, " ")
	if i := strings.Index(text, ". "); i >= 0 {
		return strings.TrimSpace(text[:i+1])
	}
	return strings.TrimSpace(text)
}

func splitFrontmatter(raw string) (string, error) {
	op := "catalog.splitFrontmatter"
	s := strings.TrimPrefix(raw, "\ufeff")
	if !strings.HasPrefix(s, "---\n") && !strings.HasPrefix(s, "---\r\n") {
		return "", derrors.New(derrors.CodeParse, op, "missing opening frontmatter delimiter")
	}
	rest := s[4:]
	if strings.HasPrefix(s, "---\r\n") {
		rest = s[5:]
	}
	idx := strings.Index(rest, "\n---\n")
	if idx < 0 {
		idx = strings.Index(rest, "\r\n---\r\n")
	}
	if idx < 0 {
		idx = strings.Index(rest, "\n---\r\n")
	}
	if idx < 0 {
		return "", derrors.New(derrors.CodeParse, op, "missing closing frontmatter delimiter")
	}
	return rest[:idx], nil
}
