# Model Picks long-term objective

**Choose a model for a job by fit, cost, and host limits, from facts that tools can score and route on.**

## Four layers

Only layer 1 is implemented in the nested-hosts slice. Layers 2 through 4 are follow-up work.

1. **Facts catalog.** One markdown card per model. Nested `hosts` rows hold price, capacity, and runtime ids. Operator experience stays on the card. `catalog.json` is the machine extract. Hugo Model Picks is a viewer of those facts, not a scorer.
2. **Retrieval.** Scripts pull public host catalogs into host rows. Agent WebFetch remains the fallback. No login dashboards.
3. **dbt scoring.** Extract facts, then dbt Fusion staging, then `fct_*` marts (cost-normalized scores, host rank per job, industrial gates, threshold macros). Output is a DuckDB contract. Hugo or Evidence only reads that contract.
4. **Routing service.** Coding tools and routers query scored offers (kind, jobs, capacity, price), not markdown.

## Hard rules

- Markdown and `catalog.json` hold **filed facts** only. No blended ranks in YAML.
- `kind` is portability (`sovereign`, `provider`, `service`), not who owns today's GPUs.
- Same open weights on many hosts stay **one card**.
- Uncertain numbers are `unknown` or omitted. Never invented.
