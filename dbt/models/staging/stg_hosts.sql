{{ config(materialized='view') }}

with catalog as (
  select * from read_json_auto('{{ var("catalog_path") }}')
),
cards as (
  select unnest(cards) as card from catalog
)

select
  card.name::varchar as name,
  card.title::varchar as title,
  card.kind::varchar as kind,
  card.status::varchar as status,
  card.jobs as jobs,
  h.id::varchar as host_id,
  h.runtime::varchar as runtime,
  h.runtime_model_id::varchar as runtime_model_id,
  h.capacity::varchar as capacity,
  try_cast(h.price_in as double) as price_in,
  try_cast(h.price_out as double) as price_out
from cards
cross join unnest(card.hosts) as t(h)
