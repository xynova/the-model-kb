# Model Picks long-term objective

**Choose a model for a job by fit, cost, and host limits, from facts that tools can score and route on.**

## Four layers

Layer 1 (facts catalog) and the first agent pick slice (layer 3 output plus layer 4 API) are in progress on branch `feat/roster-pick`. Layer 2 retrieval is wired: nightly **`roster-mine`** PR plus operator **`capacity-decisions`**.

1. **Facts catalog.** One markdown card per model. Nested `hosts` rows hold price, capacity, and runtime ids. Operator experience stays on the card. `catalog.json` is the machine extract. Hugo Model Picks is a viewer of those facts, not a scorer.
2. **Retrieval.** `make -C site mine-hosts` runs `retrieve-limits.py` (allowlisted host URLs from `catalog.json`) and `apply-mine.py`, then `retrieve-aa.py` / `apply-aa.py` when `AA_API_KEY` is set. Outputs: `roster/retrieval/limits-review.json` and `roster/retrieval/aa-intelligence-review.json`. **`capacity` changes only** for rows listed in `roster/retrieval/capacity-decisions.yaml`. **`aa_intelligence`** is copied from the Artificial Analysis free Data API (thinking cut: max, xhigh, high, or reasoning; never non-reasoning; index major 5) onto matching non-service cards; unmatched or null AA rows that already have a number clear to `unknown`. The same pull records this card's AA cost-per-task plus off-shelf **Alternative** (at least this index, lower AA cost-per-task, then best index per cost) and **Similar** (closest thinking index) into `site/data/aa_compare.json`; those hints do not mint cards. New host companies enter that allowlist through skill `providers-mint` (public page, `keep? y/n`), not through Perplexity text or the nightly PR. `make -C site polypus-diff` compares live Polypus OpenAI and Systemone models lists to `catalog.json` and prints keep/drop proposals only; it does not mint Polypus as a host or file prices. Scheduled workflow `.github/workflows/roster-mine.yml` pushes branch `roster-mine` and opens/updates that PR; it does not merge or push `main`. Agent WebFetch of an AA model page remains the fallback when the API key is unset. No login dashboards.
3. **dbt scoring.** `make -C dbt build` reads `catalog.json`, runs dbt Fusion staging and `fct_*` marts, and writes committed `scores.json` (`scores_schema: 1`). Each offer carries a full **quality** and **performance** leaf grid (see [taxonomy.md](taxonomy.md) heading **Quality leaves**): `measured`, `unknown`, or `not_applicable` per leaf, plus `capability_cohort` and `leaf_cohorts`. DuckDB under `artifacts/` is gitignored. Operators refresh scores after catalog changes; agents do not run dbt.
4. **Routing service.** `Picker.Pick` in `tools/roster-pick` filters and sorts `scores.json` offers for a task by **fit, capacity, kind, then ascending `cost_per_m`**. Quality leaves are exposed on each offer for inspect and future compare; they do not replace cost sort in v1. Agents call MCP `roster_pick` (stdio). The Cobra CLI (`roster-pick pick`, `inspect`, `mcp`, `version`) is for tests and operators, not for inventing picks from markdown.

## Hard rules

- Markdown and `catalog.json` hold **filed facts** only. No blended ranks in YAML.
- `kind` is portability (`sovereign`, `provider`, `service`), not who owns today's GPUs.
- Same open weights on many hosts stay **one card**.
- Uncertain numbers are `unknown` or omitted. Never invented.
