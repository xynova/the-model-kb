{{ config(materialized='table') }}

select
  offer_id,
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
  cost_per_m,
  industrial_ok,
  case
    when cost_per_m is null then null
    else dense_rank() over (order by cost_per_m)
  end as rank_by_cost
from {{ ref('fct_host_offers') }}
