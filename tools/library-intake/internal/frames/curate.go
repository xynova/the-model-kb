package frames

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xynova/library-intake/internal/candidates"
	"github.com/xynova/library-intake/internal/clients/polypus"
	derrors "github.com/xynova/library-intake/internal/errors"
)

const curateSystem = `You assign YouTube roundup frames to library intake candidates. Reply with JSON only: {"assignments":[{"file":"frame_0001.jpg","candidate_title":null_or_exact_title,"role":"cover|gallery|skip","reason":"short"}]}. candidate_title MUST be exactly one of the provided titles, or null for intro/outro/sponsor/unrelated. Prefer cover for a clear product demo or title card; gallery for useful secondary stills; skip for blurry/text-only or duplicate frames. Do not invent project URLs.`

// Assignment maps one frame file to a candidate or skip.
type Assignment struct {
	File           string  `json:"file"`
	CandidateTitle *string `json:"candidate_title"`
	Role           string  `json:"role"`
	Reason         string  `json:"reason"`
}

// MapFile is the durable frame-map artifact.
type MapFile struct {
	Model           string       `json:"model"`
	FramesDir       string       `json:"frames_dir"`
	CandidateTitles []string     `json:"candidate_titles"`
	Assignments     []Assignment `json:"assignments"`
}

// CurateConfig holds frame curation construction inputs.
type CurateConfig struct {
	Client    *polypus.Client
	Model     string
	BatchSize int
}

// Curator maps frames to candidates via Gemma vision.
type Curator struct {
	client    *polypus.Client
	model     string
	batchSize int
}

// CreateCurator builds a Curator.
func (c CurateConfig) CreateCurator() (*Curator, error) {
	if c.Client == nil {
		return nil, derrors.New(derrors.CodeInvalidArgument, "frames.CreateCurator", "client is nil")
	}
	model := strings.TrimSpace(c.Model)
	if model == "" {
		return nil, derrors.New(derrors.CodeInvalidArgument, "frames.CreateCurator", "model is empty")
	}
	batch := c.BatchSize
	if batch <= 0 {
		batch = 4
	}
	return &Curator{client: c.Client, model: model, batchSize: batch}, nil
}

// Curate maps frames under framesDir to candidates and writes outPath.
func (c *Curator) Curate(ctx context.Context, framesDir string, list *candidates.List, outPath string) (*MapFile, error) {
	if c == nil {
		return nil, derrors.New(derrors.CodeFailed, "frames.Curate", "curator is nil")
	}
	if list == nil || len(list.Candidates) == 0 {
		return nil, derrors.New(derrors.CodeInvalidArgument, "frames.Curate", "candidates are empty")
	}
	frames, err := listFrames(framesDir)
	if err != nil {
		return nil, err
	}
	titles := make([]string, 0, len(list.Candidates))
	for _, cand := range list.Candidates {
		title := strings.TrimSpace(cand.Title)
		if title == "" {
			title = strings.TrimSpace(cand.NameHint)
		}
		titles = append(titles, title)
	}

	var assignments []Assignment
	for i := 0; i < len(frames); i += c.batchSize {
		end := i + c.batchSize
		if end > len(frames) {
			end = len(frames)
		}
		batch := frames[i:end]
		batchAsg, err := c.curateBatch(ctx, titles, batch)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, batchAsg...)
	}

	out := &MapFile{
		Model:           c.model,
		FramesDir:       framesDir,
		CandidateTitles: titles,
		Assignments:     assignments,
	}
	if err := writeJSON(outPath, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Curator) curateBatch(ctx context.Context, titles []string, batch []string) ([]Assignment, error) {
	var b strings.Builder
	b.WriteString("Candidates:\n")
	for _, t := range titles {
		b.WriteString("- ")
		b.WriteString(t)
		b.WriteByte('\n')
	}
	b.WriteString("\nAssign each of these ")
	b.WriteString(fmt.Sprintf("%d", len(batch)))
	b.WriteString(" frames. Filenames: ")
	names := make([]string, len(batch))
	for i, p := range batch {
		names[i] = filepath.Base(p)
	}
	b.WriteString(strings.Join(names, ", "))

	parts := []any{
		polypus.TextPart{Type: "text", Text: b.String()},
	}
	for _, path := range batch {
		dataURL, err := jpegDataURL(path)
		if err != nil {
			return nil, err
		}
		parts = append(parts,
			polypus.ImagePart{Type: "image_url", ImageURL: polypus.ImageURL{URL: dataURL}},
			polypus.TextPart{Type: "text", Text: "(filename: " + filepath.Base(path) + ")"},
		)
	}

	content, err := c.client.ChatCompletions(ctx, c.model, 0.1, []polypus.Message{
		{Role: "system", Content: curateSystem},
		{Role: "user", Content: parts},
	})
	if err != nil {
		return nil, err
	}
	raw, err := polypus.ExtractJSONObject(content)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Assignments []Assignment `json:"assignments"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, derrors.Wrap(err, derrors.CodeParse, "frames.curateBatch", "unmarshal assignments")
	}
	if len(parsed.Assignments) == 0 {
		return nil, derrors.New(derrors.CodeParse, "frames.curateBatch", "empty assignments")
	}
	return parsed.Assignments, nil
}

func listFrames(dir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "frame_*.jpg"))
	if err != nil {
		return nil, derrors.Wrap(err, derrors.CodeFailed, "frames.listFrames", "glob")
	}
	if len(matches) == 0 {
		return nil, derrors.New(derrors.CodeNotFound, "frames.listFrames", "no frame_*.jpg found").
			With("dir", dir)
	}
	sort.Strings(matches)
	return matches, nil
}

func jpegDataURL(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", derrors.Wrap(err, derrors.CodeFailed, "frames.jpegDataURL", "read frame").
			With("path", path)
	}
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(data), nil
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "frames.writeJSON", "mkdir").
			With("path", path)
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "frames.writeJSON", "marshal")
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "frames.writeJSON", "write").
			With("path", path)
	}
	return nil
}
