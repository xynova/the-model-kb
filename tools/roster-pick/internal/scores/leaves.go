package scores

// Keep in sync with roster/taxonomy.md Quality leaves and dbt/scripts/export_scores.py.
var QualityLeafIDs = []string{
	"quality.general_problem_solving.benchmark_composite",
	"quality.human_preference.pairwise_rating",
	"quality.coding.agent_benchmark_composite",
	"quality.coding.code_edit_first_attempt",
	"quality.multimodal_understanding.benchmark_accuracy",
	"quality.typed_decision.classification_accuracy",
	"quality.typed_decision.precision",
	"quality.typed_decision.recall",
	"quality.typed_decision.brier_loss",
	"quality.typed_decision.ordinal_mae",
}

// PerformanceLeafIDs lists performance dimension leaves on every offer.
var PerformanceLeafIDs = []string{
	"performance.end_to_end_latency_ms",
	"performance.time_to_first_token_ms",
	"performance.request_success_rate",
}

// Leaf is one quality or performance measurement slot on an offer.
type Leaf struct {
	Status        string   `json:"status"`
	Value         *float64 `json:"value,omitempty"`
	Attribution   string   `json:"attribution,omitempty"`
	ComparisonKey string   `json:"comparison_key,omitempty"`
}
