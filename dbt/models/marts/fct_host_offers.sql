{{ config(materialized='table') }}

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
  (capacity = 'industrial') as industrial_ok
from {{ ref('stg_hosts') }}
