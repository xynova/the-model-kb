{{ config(materialized='table') }}

with staged as (
  select * from {{ ref('stg_hosts') }}
),
flags as (
  select
    *,
    coalesce(list_contains(jobs, 'vision'), false) as job_vision,
    coalesce(list_contains(jobs, 'chat'), false) as job_chat,
    coalesce(list_contains(jobs, 'reasoning'), false) as job_reasoning,
    coalesce(list_contains(jobs, 'coding'), false) as job_coding,
    coalesce(list_contains(jobs, 'classify'), false) as job_classify
  from staged
),
cap_lists as (
  select
    *,
    list_concat(
      list_value('text'),
      case when job_vision then list_value('image') else list_value() end
    ) as input_modalities,
    list_concat(
      list_concat(
        case
          when job_chat or job_reasoning then list_value('generated_text')
          else list_value()
        end,
        case when job_coding then list_value('code') else list_value() end
      ),
      case
        when job_classify
          then list_value(
            'binary_probability',
            'categorical_distribution',
            'ordinal_distribution'
          )
        else list_value()
      end
    ) as output_contracts,
    list_value('single_response') as execution_modes
  from flags
),
cap_named as (
  select
    *,
    (
      select string_agg(x, '+' order by x)
      from unnest(input_modalities) as t(x)
    ) as input_modalities_key,
    coalesce(
      (
        select string_agg(x, '+' order by x)
        from unnest(output_contracts) as t(x)
      ),
      ''
    ) as output_contracts_key,
    (
      select string_agg(x, '+' order by x)
      from unnest(execution_modes) as t(x)
    ) as execution_modes_key
  from cap_lists
),
leaves as (
  select
    *,
    input_modalities_key
      || '__'
      || output_contracts_key
      || '__'
      || execution_modes_key as capability_cohort,
    (job_chat or job_reasoning or job_coding) as gen_applicable,
    job_vision as multimodal_applicable,
    job_coding as coding_applicable,
    job_classify as classify_applicable,
    (job_chat or job_reasoning) as ttft_applicable
  from cap_named
),
assembled as (
  select
    name || '|' || host_id as offer_id,
    name,
    title,
    kind,
    status,
    jobs,
    host_id,
    runtime,
    runtime_model_id,
    capacity,
    price_in,
    price_out,
    case
      when price_in is not null and price_out is not null
        then price_in + price_out
    end as cost_per_m,
    (capacity = 'industrial') as industrial_ok,
    capability_cohort,
    input_modalities,
    output_contracts,
    execution_modes,
    json_object(
      'quality.general_problem_solving.benchmark_composite',
      {{ quality_leaf('gen_applicable', 'aa_intelligence', 'aa_intelligence_index:v5.0.1') }},
      'quality.human_preference.pairwise_rating',
      {{ quality_leaf('gen_applicable', 'arena_elo', 'arena_board:unknown') }},
      'quality.coding.agent_benchmark_composite',
      {{ quality_leaf('coding_applicable', 'null::double', 'aa_coding_agent_index:unknown') }},
      'quality.coding.code_edit_first_attempt',
      {{ quality_leaf('coding_applicable', 'aider_polyglot', 'aider_polyglot_first_pass') }},
      'quality.multimodal_understanding.benchmark_accuracy',
      {{ quality_leaf('multimodal_applicable', 'null::double', 'mmmu_pro:unknown') }},
      'quality.typed_decision.classification_accuracy',
      {{ quality_leaf('classify_applicable', 'null::double', 'typed_decision:unknown') }},
      'quality.typed_decision.precision',
      {{ quality_leaf('classify_applicable', 'null::double', 'typed_decision:unknown') }},
      'quality.typed_decision.recall',
      {{ quality_leaf('classify_applicable', 'null::double', 'typed_decision:unknown') }},
      'quality.typed_decision.brier_loss',
      {{ quality_leaf('classify_applicable', 'null::double', 'typed_decision:unknown') }},
      'quality.typed_decision.ordinal_mae',
      {{ quality_leaf('classify_applicable', 'null::double', 'typed_decision:unknown') }}
    ) as quality,
    json_object(
      'performance.end_to_end_latency_ms',
      json_object('status', 'unknown'),
      'performance.time_to_first_token_ms',
      case
        when ttft_applicable then json_object('status', 'unknown')
        else json_object('status', 'not_applicable')
      end,
      'performance.request_success_rate',
      json_object('status', 'unknown')
    ) as performance
  from leaves
)
select * from assembled
