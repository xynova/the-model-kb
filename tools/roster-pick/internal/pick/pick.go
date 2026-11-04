package pick

import (
	"context"
	"sort"

	"github.com/xynova/roster-pick/internal/config"
	derrors "github.com/xynova/roster-pick/internal/errors"
	"github.com/xynova/roster-pick/internal/lexicon"
	"github.com/xynova/roster-pick/internal/observability"
	"github.com/xynova/roster-pick/internal/scores"
)

// Request is input for Pick.
type Request struct {
	Task     string
	Job      string
	Capacity string
	Kind     string
	Limit    int
}

// Result is the pick response.
type Result struct {
	MatchedJobs     []string       `json:"matched_jobs"`
	MatchedCapacity string         `json:"matched_capacity"`
	MatchedKind     string         `json:"matched_kind"`
	Offers          []scores.Offer `json:"offers"`
}

// Picker ranks offers from scores.json for a task.
type Picker struct {
	doc    *scores.Document
	logger observability.Logger
	limit  int
}

// CreatePicker builds a Picker from config.
func CreatePicker(cfg config.Config, logger observability.Logger) (*Picker, error) {
	if logger == nil {
		panic("pick.CreatePicker: logger is required")
	}
	if cfg.ScoresPath == "" {
		return nil, derrors.New(derrors.CodeUnavailable, "pick.CreatePicker", "scores path is empty; run make -C dbt build")
	}
	doc, err := scores.Load(cfg.ScoresPath)
	if err != nil {
		return nil, err
	}
	limit := cfg.Limit
	if limit <= 0 {
		limit = 5
	}
	return &Picker{doc: doc, logger: logger, limit: limit}, nil
}

// Pick filters and sorts offers for the request.
func (p *Picker) Pick(ctx context.Context, req Request) (Result, error) {
	if p == nil {
		return Result{}, derrors.New(derrors.CodeFailed, "pick.Pick", "picker is nil")
	}
	if ctx == nil {
		return Result{}, derrors.New(derrors.CodeInvalidArgument, "pick.Pick", "context is required")
	}
	select {
	case <-ctx.Done():
		return Result{}, ctx.Err()
	default:
	}

	match, err := lexicon.ParseTask(req.Task)
	if err != nil {
		return Result{}, err
	}
	if req.Job != "" {
		match.Jobs = []string{req.Job}
	}
	if req.Capacity != "" {
		match.Capacity = req.Capacity
	}
	if req.Kind != "" {
		match.Kind = req.Kind
	}

	limit := req.Limit
	if limit <= 0 {
		limit = p.limit
	}

	filtered := filterOffers(p.doc.Offers, match)
	sortOffers(filtered)

	if len(filtered) > limit {
		filtered = filtered[:limit]
	}

	res := Result{
		MatchedJobs:     match.Jobs,
		MatchedCapacity: match.Capacity,
		MatchedKind:     match.Kind,
		Offers:          filtered,
	}
	if len(filtered) == 0 {
		return res, derrors.New(derrors.CodeNotFound, "pick.Pick", "no offers match the task filters").
			With("matched_jobs", match.Jobs).
			With("matched_capacity", match.Capacity).
			With("matched_kind", match.Kind)
	}
	return res, nil
}

func filterOffers(offers []scores.Offer, match lexicon.Match) []scores.Offer {
	filtered := make([]scores.Offer, 0, len(offers))
	for _, o := range offers {
		if o.Status == "retired" {
			continue
		}
		if match.Capacity == "industrial" {
			if !o.IndustrialOK {
				continue
			}
		} else if match.Capacity != "" && o.Capacity != match.Capacity {
			continue
		}
		if match.Kind != "" && o.Kind != match.Kind {
			continue
		}
		if !jobIntersect(o.Jobs, match.Jobs) {
			continue
		}
		filtered = append(filtered, o)
	}
	return filtered
}

func sortOffers(filtered []scores.Offer) {
	sort.SliceStable(filtered, func(i, j int) bool {
		a, b := filtered[i], filtered[j]
		if cmpCost(a.CostPerM, b.CostPerM) != 0 {
			return cmpCost(a.CostPerM, b.CostPerM) < 0
		}
		if statusRank(a.Status) != statusRank(b.Status) {
			return statusRank(a.Status) < statusRank(b.Status)
		}
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		return a.HostID < b.HostID
	})
}

func jobIntersect(offerJobs, matched []string) bool {
	set := map[string]struct{}{}
	for _, j := range offerJobs {
		set[j] = struct{}{}
	}
	for _, m := range matched {
		if _, ok := set[m]; ok {
			return true
		}
	}
	return false
}

func cmpCost(a, b *float64) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return 1
	}
	if b == nil {
		return -1
	}
	if *a < *b {
		return -1
	}
	if *a > *b {
		return 1
	}
	return 0
}

func statusRank(s string) int {
	switch s {
	case "reach-for":
		return 0
	case "watch":
		return 1
	default:
		return 2
	}
}
