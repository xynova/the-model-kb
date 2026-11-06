package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/xynova/roster-mill/internal/aa"
	"github.com/xynova/roster-mill/internal/card"
	derrors "github.com/xynova/roster-mill/internal/errors"
)

type loaded struct {
	fm   map[string]any
	file *card.File
	rel  string
}

// Result is generate output paths and card count.
type Result struct {
	GeneratedDir string
	CatalogPath  string
	IndexPath    string
	Cards        int
}

// Run validates cards, writes Hugo pages, catalog.json, INDEX.md, and aa_compare.json.
func Run(repoRoot string, stdout io.Writer) (*Result, error) {
	src := filepath.Join(repoRoot, "roster")
	siteDir := filepath.Join(repoRoot, "site")
	dst := filepath.Join(siteDir, "generated", "roster")
	catalogPath := filepath.Join(src, "catalog.json")
	staticJSON := filepath.Join(siteDir, "static", "roster.json")
	indexPath := filepath.Join(src, "INDEX.md")

	info, err := os.Stat(src)
	if err != nil || !info.IsDir() {
		return nil, derrors.New(derrors.CodeNotFound, "generate", "roster source missing: "+src)
	}

	var cards []loaded
	namesSeen := map[string]struct{}{}

	err = filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		parts := strings.Split(rel, string(filepath.Separator))
		if d.IsDir() {
			if rel != "." && (parts[0] == "retrieval" || strings.HasPrefix(parts[0], ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if _, skip := card.SkipNames[d.Name()]; skip {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		f := card.Split(string(raw))
		fm := f.AsMap()
		if fm == nil {
			fmt.Fprintf(os.Stderr, "skip (no frontmatter): %s\n", filepath.ToSlash(rel))
			return nil
		}
		slashRel := filepath.ToSlash(rel)
		if vErr := card.ValidateCard(fm, slashRel); vErr != nil {
			return vErr
		}
		name, _ := fm["name"].(string)
		if _, dup := namesSeen[name]; dup {
			return fmt.Errorf("duplicate card name: %s", name)
		}
		namesSeen[name] = struct{}{}
		cards = append(cards, loaded{fm: fm, file: f, rel: slashRel})
		return nil
	})
	if err != nil {
		return nil, derrors.Wrap(err, derrors.CodeFailed, "generate", "walk roster")
	}

	if err := os.RemoveAll(dst); err != nil {
		return nil, derrors.Wrap(err, derrors.CodeFailed, "generate", "clear generated/roster")
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return nil, derrors.Wrap(err, derrors.CodeFailed, "generate", "mkdir generated/roster")
	}

	jsonCards := make([]map[string]any, 0, len(cards))
	indexRows := make([]loaded, 0, len(cards))
	for _, c := range cards {
		out := filepath.Join(dst, filepath.FromSlash(c.rel))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return nil, err
		}
		hugoFM, err := card.DumpHugo(c.file.Front)
		if err != nil {
			return nil, derrors.Wrap(err, derrors.CodeFailed, "generate", "dump hugo "+c.rel)
		}
		text := "---\n" + hugoFM + "\n---\n\n" + c.file.Body
		if err := os.WriteFile(out, []byte(text), 0o644); err != nil {
			return nil, err
		}
		jsonCards = append(jsonCards, card.ToJSON(c.fm))
		indexRows = append(indexRows, c)
	}

	sort.Slice(jsonCards, func(i, j int) bool {
		return str(jsonCards[i]["name"]) < str(jsonCards[j]["name"])
	})
	generatedAt := time.Now().UTC().Format("2006-01-02")
	catalog := map[string]any{
		"catalog_schema": 1,
		"generated_at":   generatedAt,
		"cards":          jsonCards,
	}
	catalogBytes, err := marshalJSON(catalog)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(catalogPath, catalogBytes, 0o644); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(staticJSON), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(staticJSON, catalogBytes, 0o644); err != nil {
		return nil, err
	}
	if err := writeIndex(indexRows, indexPath); err != nil {
		return nil, err
	}

	reviewPath := filepath.Join(src, "retrieval", "aa-intelligence-review.json")
	comparePath := filepath.Join(siteDir, "data", "aa_compare.json")
	if raw, err := os.ReadFile(reviewPath); err == nil {
		var review map[string]any
		if json.Unmarshal(raw, &review) == nil {
			if err := os.MkdirAll(filepath.Dir(comparePath), 0o755); err != nil {
				return nil, err
			}
			cmp := aa.CompareData(review, catalog)
			b, err := marshalJSON(cmp)
			if err != nil {
				return nil, err
			}
			if err := os.WriteFile(comparePath, b, 0o644); err != nil {
				return nil, err
			}
		}
	}

	fmt.Fprintf(stdout, "generated roster pages under %s\n", dst)
	fmt.Fprintf(stdout, "wrote %s and %s\n", catalogPath, staticJSON)
	fmt.Fprintf(stdout, "wrote %s (%d cards)\n", indexPath, len(jsonCards))
	return &Result{
		GeneratedDir: dst,
		CatalogPath:  catalogPath,
		IndexPath:    indexPath,
		Cards:        len(jsonCards),
	}, nil
}

func writeIndex(cards []loaded, indexPath string) error {
	sort.Slice(cards, func(i, j int) bool {
		return str(cards[i].fm["name"]) < str(cards[j].fm["name"])
	})
	var b strings.Builder
	b.WriteString("# Roster index\n\n")
	b.WriteString("Regenerated by `tools/roster-mill`.\n\n")
	b.WriteString("| name | title | kind | service | category | status | lab | hosts | jobs | tags | path |\n")
	b.WriteString("|------|-------|------|---------|----------|--------|-----|-------|------|------|------|\n")
	for _, c := range cards {
		fm := c.fm
		hostIDs := joinIDs(fm["hosts"])
		jobs := joinList(fm["jobs"])
		tags := joinList(fm["tags"])
		service := str(fm["service"])
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			str(fm["name"]), str(fm["title"]), str(fm["kind"]), service,
			str(fm["category"]), str(fm["status"]), str(fm["lab"]),
			hostIDs, jobs, tags, c.rel)
	}
	b.WriteString("\n")
	return os.WriteFile(indexPath, []byte(b.String()), 0o644)
}

func joinIDs(hosts any) string {
	list, _ := hosts.([]any)
	ids := make([]string, 0, len(list))
	for _, h := range list {
		hm, ok := h.(map[string]any)
		if !ok {
			continue
		}
		ids = append(ids, str(hm["id"]))
	}
	return strings.Join(ids, ",")
}

func joinList(v any) string {
	list, ok := v.([]any)
	if !ok {
		return ""
	}
	parts := make([]string, 0, len(list))
	for _, item := range list {
		parts = append(parts, str(item))
	}
	return strings.Join(parts, ",")
}

func str(v any) string {
	s, _ := v.(string)
	return s
}

func marshalJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
