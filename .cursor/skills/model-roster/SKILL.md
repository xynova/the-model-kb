---
name: model-roster
description: >-
  Add or refresh curated model decision cards in roster/: kind (sovereign /
  provider / service), price, capability signals, personal experience notes,
  INDEX.md. Use when the user says model roster, roster add, roster refresh,
  choose a model card, or update model prices/capabilities in the roster. Card body
  prose MUST be active voice, plain English, no AI fluff.
---

# Model roster

Curate decision cards into `roster/`. This skill is a **consumer overlay** (not cursor-packs). Catalog root is repo-relative `roster/`.

## When to load

Load when the user asks to add, refresh, or update model roster cards, or to file a model into the decision shelf for choosing models by task.

## Core constraints

**CONSTRAINT:** Catalog root and taxonomy

- MUST: Read `roster/taxonomy.md` before assigning `kind`, `category`, `tags`, `jobs`, or `hosts`.
- MUST: Write entries under `roster/<category>/<slug>.md` using structure from `roster/_template.md`.
- MUST NOT: Invent kinds, categories, tags, or runtimes absent from `taxonomy.md` (unless the user first updates taxonomy).

Enforcement: Read taxonomy before any keep proposal; Grep tags against the allowed list before write.
Violation: STOP, re-Read taxonomy, fix assignments, re-verify.

CORRECT:
```text
kind: sovereign
category: coding
tags: [llm, open, coding, cloudflare]
runtime: cloudflare
```

PROHIBITED:
```text
kind: hosted
category: local
tags: [best]
```

**CONSTRAINT:** Kind is required and must match control model

- MUST: Propose exactly one of `sovereign` | `provider` | `service` on every cold read.
- MUST: Treat Cloudflare/DeepInfra **open hosted** weights as `sovereign`; put each host offer in `hosts` with `runtime` on the row.
- MUST: Treat lab-locked APIs (OpenAI, Anthropic, Gemini, …) as `provider`.
- MUST: Treat Cursor (and similar product-bundled models) as `service` with `service: cursor`.
- MUST: Treat a host as `provider` only when it **proxies** a lab-locked API; note the proxy in `price_note`.
- MUST NOT: Infer `kind` from the host name alone.
- MUST NOT: Use a frontmatter key named `provider` (use `lab` for origin).

Enforcement: Cold read states proposed `kind` and why; taxonomy CORRECT/PROHIBITED examples apply.
Violation: STOP, reclassify before write.

CORRECT:
```text
kind: sovereign | runtime: cloudflare | @cf/meta/llama-3.2-3b-instruct
kind: provider | lab: anthropic | runtime: cloudflare | proxy via AI Gateway
kind: service | service: cursor
```

PROHIBITED:
```text
kind: provider for @cf/meta/llama-3.2-3b-instruct
kind: sovereign for Cursor Composer
```

**CONSTRAINT:** Keep gate before any add write (one candidate at a time)

- MUST: Ask about **exactly one** candidate per turn.
- MUST: Show name, proposed `kind`, category, short cold read, then `keep? y/n`.
- MUST NOT: Batch-keep or auto-write without `y`.
- MUST NOT: Write until the user answers `y` for that candidate.
- MUST: On `n`, skip with no file; move to the next candidate if any.
- MUST: On `y`, enrich → write → refresh INDEX, then continue.

Enforcement: One ask → one reply → at most one new entry.
Violation: STOP, do not batch-write.

CORRECT:
```text
[1/3] llama-3.2-3b-cf | sovereign | general

Meta Llama 3.2 3B on Cloudflare Workers AI (@cf/...). Portable open weights; switch host anytime. Token-equivalent list price on CF pricing page (confirm on refresh).

keep? y/n
```

PROHIBITED: fluffy hype; datasheet bullet stacks as the default cold read; batch “reply y to file all.”

**CONSTRAINT:** Same-segment peers (apples to apples)

- MUST: When proposing competitors or “same segment” cards next to an anchor, Read `roster/taxonomy.md` same-segment rules and match **generation/cut**, **capability class**, **modality**, and (for MoE) **active-parameter band**.
- MUST: State in the cold read whether the candidate is a **peer** or an explicit **baseline** (older/weaker cut).
- MUST NOT: Treat “same host” or “similar total parameter count” alone as enough for a peer.
- MUST NOT: Present an older generation (e.g. early Qwen3 next to Gemma 4) as a silent peer; either skip, or label baseline and say why.
- MAY: Propose the same open weights on another host as a peer (host-price compare); that is always apples-to-apples for sovereignty.

Enforcement: Cold read names peer vs baseline; taxonomy CORRECT/PROHIBITED examples apply.
Violation: STOP, reclassify or drop the candidate before `keep? y/n`.

CORRECT:
```text
[2/4] gemma-4-26b-a4b | sovereign | multimodal | add host deepinfra

Append a `hosts` row on existing card `gemma-4-26b-a4b`. Same weights; host price compare vs cloudflare row.
```

CORRECT:
```text
[3/4] qwen3-30b-a3b-cf | sovereign | general | baseline (older cut)

Older Qwen3 MoE on CF. Not a Gemma 4 peer; cheap baseline only.
```

PROHIBITED:
```text
[3/4] qwen3-30b-a3b-cf | same segment as Gemma 4 26B-A4B
# older cut, weaker class, presented as peer
```

PROHIBITED:
```text
gemma-4-26b-a4b-deepinfra.md
# second file for same open weights
```

**CONSTRAINT:** One card per model; hosts are rows

- MUST: Keep one markdown file per model (`name` slug). Same open weights on another host = append a `hosts` row on that `name`.
- MUST NOT: Create a second slug/file for the same weights (e.g. `*-deepinfra.md` beside the anchor card).
- MUST: Assign required `jobs` and per-host `capacity` from taxonomy.
- MUST: If the user names a host `id` that is not in taxonomy **Host row ids**, STOP and load skill `providers-mint` before writing a `hosts` row.
- MUST NOT: Invent a new `hosts[].id` or `runtime` inside a card-only edit.
- MUST NOT: File blended ranks or computed scores in YAML (see `roster/OBJECTIVE.md` and taxonomy **Quality leaves**).
- MUST: Map every filed benchmark number to a **leaf id** in taxonomy **Quality leaves** before write (`aa_intelligence`, `arena_elo`, `aider_polyglot` map to named leaves).
- MUST NOT: File [BenchLM](https://benchlm.ai/) overall, Estimated ranks, or copied prices as card or scores metadata.
- MUST NOT: File Perplexity output as prices, URLs, or scores.

**CONSTRAINT:** Enrichment (v1)

- MUST: Use WebSearch / WebFetch only against sources listed in `roster/taxonomy.md`, except `aa_intelligence` MAY come from `make -C site retrieve-aa` / `apply-aa` when `AA_API_KEY` is set.
- MUST: Prefer the AA Data API (`GET /api/v2/language/models/free`) over a model-page WebFetch for `aa_intelligence`.
- MUST: Prefer host-row pricing when `hosts[].runtime` is `cloudflare`, `deepinfra`, `openrouter`, or `google` (public pages, not login dashboards).
- MUST: Prefer lab pricing for `kind: provider`; product docs for `kind: service`.
- MUST: Mark uncertain fields `unknown` or omit the key.
- MUST NOT: Invent URLs, prices, or scores.
- MUST NOT: Call Perplexity MCP for enrichment in v1.
- MUST NOT: Fabricate **My experience** prose; leave a short placeholder or empty section until the user supplies notes.

Enforcement: Every `links.*` value and numeric signal must come from search/fetch or be omitted/`unknown`.
Violation: STOP, clear inventions, re-enrich.

**CONSTRAINT:** Nightly mine and capacity confirmation

- MUST: Treat `roster/retrieval/limits-review.json` as spider evidence from the **`roster-mine`** PR workflow, not as filed capacity by itself.
- MUST: Change a host row's `capacity` only after the operator adds `{name, host_id, capacity}` to `roster/retrieval/capacity-decisions.yaml` and the mine apply step runs (or the user edits the card by hand with the same intent).
- MUST NOT: Set `capacity: industrial` (or any other capacity) from page text, `suggested_capacity`, or agent hunch without a `capacity-decisions` row or explicit user instruction on the card.

Enforcement: Grep the card after apply; capacity matches decisions or prior filed value.
Violation: STOP, revert capacity, file a decisions row only after the user confirms.

**CONSTRAINT:** Dedup, compiler, INDEX

- MAY: Run `make -C site polypus-diff` to list Polypus gateway ids missing from `catalog.json`. Treat that stdout as keep/drop proposals only.
- MUST NOT: Mint `polypus` as a host id, or file Polypus route strings as `price_in` / `price_out` / `capacity`.
- MUST: Before write, Grep `roster/**/*.md` for matching `name:`.
- MUST: Skip existing `name` unless the user asks to update/refresh or add a host row.
- MUST: After any write or refresh, run `make -C site generate-roster` (writes `roster/INDEX.md`, `roster/catalog.json`, `site/static/roster.json`, and generated Hugo pages). MUST NOT hand-edit INDEX.
- MUST NOT: Create a second file when the user asked for another host on an existing model.

Enforcement: Grep before Write; INDEX row count matches entry files after generate-roster.
Violation: STOP, skip or update per user, re-run generate-roster.

**CONSTRAINT:** Refresh preserves experience

- MUST: On refresh, re-fetch every `hosts` row and card-level signals; bump `refreshed`.
- MUST: Preserve **My experience** and `status` unless the user asks to change them.
- MUST NOT: Overwrite human experience notes with generated prose.

**CONSTRAINT:** Card body prose (active, plain)

- MUST: Write sections under `#` headings in **active voice** and **plain US English** (see `roster/_template.md`).
- MUST: Use one clear idea per sentence; short paragraphs; name the **host id** and **job** when comparing offers.
- MUST: Put prices, capacity, URLs, and scores in **frontmatter** only; body explains **when to reach**, **signals in words**, and **caveats**, not a second price table.
- MUST: Leave **My experience** empty, a placeholder, or the operator's existing text only.
- MUST NOT: Use AI fluff or hype words (for example: powerful, cutting-edge, seamlessly, robust, leverage, delve, comprehensive, game-changer, state-of-the-art).
- MUST NOT: Open with throat-clearing (for example: "In today's landscape", "When it comes to", "It is worth noting").
- MAY: Use [roster-body-draft.md](roster-body-draft.md) with an external model (Gemma, etc.) when the operator asks; paste **confirmed frontmatter facts** only; MUST edit the draft before commit.

Enforcement: Scan body sections for banned fluff list and passive-only openings; facts duplicated from YAML count as violation.
Violation: STOP, rewrite body in active voice; move numbers back to frontmatter.

CORRECT:
```text
Reach for Cloudflare when you want cheap hobby inference and can live with neuron daily caps.
Pick DeepInfra on this card when you need industrial throughput and can pay the list rate.
```

PROHIBITED:
```text
This powerful, cutting-edge model seamlessly leverages best-in-class capabilities
to deliver a robust solution for users seeking to delve into multimodal AI.
```

## Flows

### Add

Copy this checklist and track progress:

```
- [ ] Read roster/taxonomy.md and roster/_template.md
- [ ] Parse candidates (or a single named model)
- [ ] For each: cold read with kind; keep? y/n; wait
- [ ] On y: dedup, enrich, write, run make -C site generate-roster
- [ ] On n: skip; next
- [ ] Summarize kept vs skipped
```

1. **Load taxonomy** — Read `roster/taxonomy.md` and `roster/_template.md`. Fail if missing.
2. **Parse** — Build an ordered candidate list. Do not write yet.
3. **Ask one** — Name, proposed `kind`, category, cold read (what it is; kind path; price/capability if known; one caveat). End with `keep? y/n`.
4. **On `n`** — Skip; continue with the next candidate.
5. **On `y`** — Dedup. Enrich. Write `roster/<category>/<slug>.md`. Run `make -C site generate-roster`. Continue.
6. **Done** — Summarize kept vs skipped.

### Refresh

```
- [ ] Read taxonomy
- [ ] Resolve slug(s) or stale-after-N-days set
- [ ] For each entry: re-fetch each host row; preserve My experience and status
- [ ] Bump refreshed; run make -C site generate-roster
```

1. User names a slug, a list, or “stale after N days.”
2. Branch enrichment by `kind` and `runtime` (below).
3. Preserve **My experience** and `status` unless asked otherwise.
4. Bump `refreshed`; run `make -C site generate-roster`.

### Enrichment branching

| host row `runtime` | Prefer |
|--------------------|--------|
| `cloudflare` | https://developers.cloudflare.com/workers-ai/platform/pricing/ then AA / Arena / Aider |
| `deepinfra` | https://deepinfra.com/pricing (and model page) then AA / Arena / Aider |
| `openrouter` | https://openrouter.ai/models then models.dev |
| `google` | Google AI / Vertex Gemma pricing docs |
| `local` | Infra / hardware `price_note`; AA / Arena when listed |
| `kind: provider` | Lab (or proxy) pricing on the host row |
| `kind: service` | Product docs in host `price_note` and `links.service_docs` |

Also pull Artificial Analysis, LMArena, OpenRouter, models.dev, Aider when the job and model coverage warrant it (see taxonomy source table).

## INDEX format

Regenerated by `make -C site generate-roster`, not by hand.

| name | title | kind | service | category | status | lab | hosts | jobs | tags | path |

`hosts` column lists comma-separated host `id` values from each card.

## Pre-completion checklist

- [ ] **Taxonomy loaded:** `roster/taxonomy.md` was Read this run
- [ ] **Kind legal:** Every write has a valid `kind`; CF open ≠ provider; Cursor ≠ sovereign/provider
- [ ] **Segment fair:** Competitor proposals labeled peer vs baseline; no silent older-cut peers
- [ ] **Keep gate held** (add flow): one ask → one reply → at most one new entry
- [ ] **Tags legal:** Every tag is in the taxonomy allowed list
- [ ] **No invented links/scores:** Traceable to WebSearch/WebFetch or omitted/`unknown`
- [ ] **Experience preserved:** Refresh did not fabricate or clobber **My experience**
- [ ] **Dedup:** No silent overwrite of existing `name`
- [ ] **INDEX fresh:** generate-roster ran; row count matches entry files
- [ ] **No Perplexity for filed metadata**
- [ ] **Prose voice:** Body is active, plain, no fluff list; My experience not invented

## Additional resources

- Frontmatter examples and kind CORRECT/PROHIBITED: [reference.md](reference.md)
- External model body draft prompt: [roster-body-draft.md](roster-body-draft.md)
- Taxonomy: [roster/taxonomy.md](../../roster/taxonomy.md)
- Template: [roster/_template.md](../../roster/_template.md)
- New host companies: skill `providers-mint`
