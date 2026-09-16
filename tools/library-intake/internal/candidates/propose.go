package candidates

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/xynova/library-intake/internal/clients/polypus"
	derrors "github.com/xynova/library-intake/internal/errors"
)

const systemPrompt = `You extract library intake candidates from an AI news video transcript. Reply with JSON only matching {"candidates":[{"title","name_hint","category","why","cold_read"}]}. Categories must be one of: models, video, 3d, agents-memory, robotics, hardware, training-research, creative. Skip intro/outro with no product to file. Skip sponsors unless they are the product under review. One product/model/system per candidate. Do not invent URLs, licenses, or VRAM. Prefer real product names as spoken; correct obvious ASR mangling when the intended name is clear.`

// Candidate is one library intake proposal.
type Candidate struct {
	Title    string `json:"title"`
	NameHint string `json:"name_hint"`
	Category string `json:"category"`
	Why      string `json:"why"`
	ColdRead string `json:"cold_read"`
}

// List is the Polypus candidate response envelope.
type List struct {
	Candidates []Candidate `json:"candidates"`
}

// ProposeConfig holds candidate proposal construction inputs.
type ProposeConfig struct {
	Client *polypus.Client
	Model  string
}

// Proposer asks Polypus for a candidate list.
type Proposer struct {
	client *polypus.Client
	model  string
}

// CreateProposer builds a Proposer.
func (c ProposeConfig) CreateProposer() (*Proposer, error) {
	if c.Client == nil {
		return nil, derrors.New(derrors.CodeInvalidArgument, "candidates.CreateProposer", "client is nil")
	}
	model := strings.TrimSpace(c.Model)
	if model == "" {
		return nil, derrors.New(derrors.CodeInvalidArgument, "candidates.CreateProposer", "model is empty")
	}
	return &Proposer{client: c.Client, model: model}, nil
}

// ProposeFromFile reads a transcript file and writes parsed candidates JSON.
func (p *Proposer) ProposeFromFile(ctx context.Context, transcriptPath, outPath string) (*List, error) {
	if p == nil {
		return nil, derrors.New(derrors.CodeFailed, "candidates.ProposeFromFile", "proposer is nil")
	}
	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		return nil, derrors.Wrap(err, derrors.CodeNotFound, "candidates.ProposeFromFile", "read transcript").
			With("path", transcriptPath)
	}
	text := strings.TrimSpace(string(data))
	if text == "" {
		return nil, derrors.New(derrors.CodeInvalidArgument, "candidates.ProposeFromFile", "transcript is empty")
	}
	list, err := p.Propose(ctx, text)
	if err != nil {
		return nil, err
	}
	if err := writeJSON(outPath, list); err != nil {
		return nil, err
	}
	return list, nil
}

// Propose asks Polypus for candidates from transcript text.
func (p *Proposer) Propose(ctx context.Context, transcript string) (*List, error) {
	content, err := p.client.ChatCompletions(ctx, p.model, 0.2, []polypus.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: transcript},
	})
	if err != nil {
		return nil, err
	}
	raw, err := polypus.ExtractJSONObject(content)
	if err != nil {
		return nil, err
	}
	var list List
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, derrors.Wrap(err, derrors.CodeParse, "candidates.Propose", "unmarshal candidates")
	}
	if len(list.Candidates) == 0 {
		return nil, derrors.New(derrors.CodeParse, "candidates.Propose", "candidates list is empty")
	}
	for i, c := range list.Candidates {
		if strings.TrimSpace(c.Title) == "" {
			return nil, derrors.New(derrors.CodeParse, "candidates.Propose", "candidate title is empty").
				With("index", i)
		}
		if strings.TrimSpace(c.Category) == "" {
			return nil, derrors.New(derrors.CodeParse, "candidates.Propose", "candidate category is empty").
				With("index", i).
				With("title", c.Title)
		}
	}
	return &list, nil
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "candidates.writeJSON", "mkdir").
			With("path", path)
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "candidates.writeJSON", "marshal")
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "candidates.writeJSON", "write").
			With("path", path)
	}
	return nil
}
