{# Keep in sync with roster/taxonomy.md Quality leaves and tools/roster-pick/internal/scores/leaves.go #}
{% set quality_leaf_ids = [
  'quality.general_problem_solving.benchmark_composite',
  'quality.human_preference.pairwise_rating',
  'quality.coding.agent_benchmark_composite',
  'quality.coding.code_edit_first_attempt',
  'quality.multimodal_understanding.benchmark_accuracy',
  'quality.typed_decision.classification_accuracy',
  'quality.typed_decision.precision',
  'quality.typed_decision.recall',
  'quality.typed_decision.brier_loss',
  'quality.typed_decision.ordinal_mae',
] %}
{% set performance_leaf_ids = [
  'performance.end_to_end_latency_ms',
  'performance.time_to_first_token_ms',
  'performance.request_success_rate',
] %}
