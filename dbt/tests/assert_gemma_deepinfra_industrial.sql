select 1 as fail
where not exists (
  select 1
  from {{ ref('fct_ranked_offers') }}
  where name = 'gemma-4-26b-a4b'
    and host_id = 'deepinfra'
    and industrial_ok
)
