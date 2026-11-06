select 1 as fail
where not exists (
  select 1
  from {{ ref('fct_host_offers') }}
  where name = 'deepseek-v4-flash'
    and json_extract_string(
      quality,
      '$."quality.general_problem_solving.benchmark_composite".status'
    ) = 'unknown'
)
