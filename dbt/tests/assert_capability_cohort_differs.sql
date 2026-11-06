select 1 as fail
where not exists (
  select 1
  from {{ ref('fct_host_offers') }} as jev
  cross join {{ ref('fct_host_offers') }} as ds
  where jev.name = 'jev'
    and jev.host_id = 'local'
    and ds.name = 'deepseek-v4-flash'
    and jev.capability_cohort is distinct from ds.capability_cohort
)
