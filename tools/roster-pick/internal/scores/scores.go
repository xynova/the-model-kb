package scores

import (
	"encoding/json"
	"fmt"
	"os"

	derrors "github.com/xynova/roster-pick/internal/errors"
)

const schemaVersion = 1

// Document is the on-disk scores.json envelope.
type Document struct {
	ScoresSchema        int     `json:"scores_schema"`
	GeneratedAt         string  `json:"generated_at"`
	SourceCatalogSchema int     `json:"source_catalog_schema"`
	Offers              []Offer `json:"offers"`
}

// Offer is one host row with derived score fields.
type Offer struct {
	Name             string          `json:"name"`
	Title            string          `json:"title"`
	Kind             string          `json:"kind"`
	Status           string          `json:"status"`
	Jobs             []string        `json:"jobs"`
	HostID           string          `json:"host_id"`
	Runtime          string          `json:"runtime,omitempty"`
	RuntimeModelID   string          `json:"runtime_model_id,omitempty"`
	Capacity         string          `json:"capacity"`
	PriceIn          *float64        `json:"price_in,omitempty"`
	PriceOut         *float64        `json:"price_out,omitempty"`
	CostPerM         *float64        `json:"cost_per_m,omitempty"`
	RankByCost       *int            `json:"rank_by_cost,omitempty"`
	IndustrialOK     bool            `json:"industrial_ok"`
	CapabilityCohort string          `json:"capability_cohort,omitempty"`
	InputModalities  []string        `json:"input_modalities,omitempty"`
	OutputContracts  []string        `json:"output_contracts,omitempty"`
	ExecutionModes   []string        `json:"execution_modes,omitempty"`
	LeafCohorts      []string        `json:"leaf_cohorts,omitempty"`
	Quality          map[string]Leaf `json:"quality,omitempty"`
	Performance      map[string]Leaf `json:"performance,omitempty"`
}

// Load reads and validates scores.json.
func Load(path string) (*Document, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- operator-supplied scores path.
	if err != nil {
		if os.IsNotExist(err) {
			return nil, derrors.Wrap(err, derrors.CodeUnavailable, "scores.Load", "scores file missing; run make -C dbt build")
		}
		return nil, derrors.Wrap(err, derrors.CodeFailed, "scores.Load", "read scores file")
	}
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, derrors.Wrap(err, derrors.CodeFailed, "scores.Load", "parse scores json")
	}
	if doc.ScoresSchema != schemaVersion {
		return nil, derrors.New(derrors.CodeFailed, "scores.Load", fmt.Sprintf("unsupported scores_schema %d (want %d)", doc.ScoresSchema, schemaVersion))
	}
	return &doc, nil
}
