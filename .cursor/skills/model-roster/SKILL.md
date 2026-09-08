---
name: model-roster
description: >-
  Add or refresh curated model decision cards in roster/: kind (sovereign /
  provider / service), price, capability signals, personal experience notes,
  INDEX.md. Use when the user says model roster, roster add, roster refresh,
  choose a model card, or update model prices/capabilities in the roster.
---

# Model roster

Curate decision cards into `roster/`. This skill is a **consumer overlay** (not cursor-packs). Catalog root is repo-relative `roster/`.

## When to load

Load when the user asks to add, refresh, or update model roster cards, or to file a model into the decision shelf for choosing models by task.

## Core constraints

**CONSTRAINT:** Catalog root and taxonomy

- MUST: Read `roster/taxonomy.md` before assigning `kind`, `category`, or `tags`.
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
- MUST: Treat Cloudflare/DeepInfra **open hosted** weights as `sovereign` with `runtime` set.
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
[2/4] gemma-4-26b-a4b-deepinfra | sovereign | multimodal | peer (same weights, host price)

Same Gemma 4 26B-A4B on DeepInfra. Fair sovereign host compare vs CF.
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

**CONSTRAINT:** Enrichment (v1)

- MUST: Use WebSearch / WebFetch only against sources listed in `roster/taxonomy.md`.
- MUST: Prefer runtime pricing when `runtime` is `cloudflare` or `deepinfra` (public pages, not login dashboards).
- MUST: Prefer lab pricing for `kind: provider`; product docs for `kind: service`.
- MUST: Mark uncertain fields `unknown` or omit the key.
- MUST NOT: Invent URLs, prices, or scores.
- MUST NOT: Call Perplexity MCP for enrichment in v1.
- MUST NOT: Fabricate **My experience** prose; leave a short placeholder or empty section until the user supplies notes.

Enforcement: Every `links.*` value and numeric signal must come from search/fetch or be omitted/`unknown`.
Violation: STOP, clear inventions, re-enrich.

**CONSTRAINT:** Dedup and INDEX

- MUST: Before write, Grep `roster/**/*.md` for matching `name:`.
- MUST: Skip existing `name` unless the user asks to update/refresh.
- MUST: Regenerate `roster/INDEX.md` after successful writes (all entry rows; include `kind`).
- MUST: Same weights on different hosts MAY use different slugs (e.g. `llama-3.3-70b-cf` vs `llama-3.3-70b-deepinfra`); both stay `kind: sovereign` when open.

Enforcement: Grep before Write; INDEX row count matches entry files (excluding `_template.md` and `.gitkeep`).
Violation: STOP, skip or update per user, rebuild INDEX.

**CONSTRAINT:** Refresh preserves experience

- MUST: On refresh, re-fetch price/capability fields per `kind` and `runtime`; bump `refreshed`.
- MUST: Preserve **My experience** and `status` unless the user asks to change them.
- MUST NOT: Overwrite human experience notes with generated prose.

## Flows

### Add

Copy this checklist and track progress:

```
- [ ] Read roster/taxonomy.md and roster/_template.md
- [ ] Parse candidates (or a single named model)
- [ ] For each: cold read with kind; keep? y/n; wait
- [ ] On y: dedup, enrich by kind/runtime, write, refresh INDEX
- [ ] On n: skip; next
- [ ] Summarize kept vs skipped
```

1. **Load taxonomy** — Read `roster/taxonomy.md` and `roster/_template.md`. Fail if missing.
2. **Parse** — Build an ordered candidate list. Do not write yet.
3. **Ask one** — Name, proposed `kind`, category, cold read (what it is; kind path; price/capability if known; one caveat). End with `keep? y/n`.
4. **On `n`** — Skip; continue with the next candidate.
5. **On `y`** — Dedup. Enrich (see branching below). Write `roster/<category>/<slug>.md`. Refresh INDEX. Continue.
6. **Done** — Summarize kept vs skipped.

### Refresh

```
- [ ] Read taxonomy
- [ ] Resolve slug(s) or stale-after-N-days set
- [ ] For each entry: re-fetch by kind/runtime; preserve My experience and status
- [ ] Bump refreshed; rewrite frontmatter signals; refresh INDEX
```

1. User names a slug, a list, or “stale after N days.”
2. Branch enrichment by `kind` and `runtime` (below).
3. Preserve **My experience** and `status` unless asked otherwise.
4. Bump `refreshed`; regenerate INDEX.

### Enrichment branching

| kind / runtime | Prefer |
|----------------|--------|
| `sovereign` + `cloudflare` | https://developers.cloudflare.com/workers-ai/platform/pricing/ then AA / Arena / Aider as relevant |
| `sovereign` + `deepinfra` | https://deepinfra.com/pricing (and model page) then AA / Arena / Aider |
| `sovereign` + `local` | Infra / hardware `price_note`; AA / Arena when the model is listed |
| `provider` | Lab (or proxy) pricing; optional OpenRouter / models.dev; AA / Arena |
| `service` | Product docs / plan limits in `price_note` and `links.service_docs`; do not invent `$/MTok` |

Also pull Artificial Analysis, LMArena, OpenRouter, models.dev, Aider when the job and model coverage warrant it (see taxonomy source table).

## INDEX format

```markdown
# Roster index

Regenerated by `model-roster` after writes.

| name | title | kind | category | status | lab | runtime | tags | path |
|------|-------|------|----------|--------|-----|---------|------|------|
| llama-3.2-3b-cf | Llama 3.2 3B (CF) | sovereign | general | watch | meta | cloudflare | llm, open, cloudflare | general/llama-3.2-3b-cf.md |
```

Include `service` in the table only if you extend the header; otherwise keep `kind: service` visible in the `kind` column and put the product in tags/`lab` as appropriate. Prefer adding a `service` column when any service entries exist:

| name | title | kind | service | category | status | lab | runtime | tags | path |

## Pre-completion checklist

- [ ] **Taxonomy loaded:** `roster/taxonomy.md` was Read this run
- [ ] **Kind legal:** Every write has a valid `kind`; CF open ≠ provider; Cursor ≠ sovereign/provider
- [ ] **Segment fair:** Competitor proposals labeled peer vs baseline; no silent older-cut peers
- [ ] **Keep gate held** (add flow): one ask → one reply → at most one new entry
- [ ] **Tags legal:** Every tag is in the taxonomy allowed list
- [ ] **No invented links/scores:** Traceable to WebSearch/WebFetch or omitted/`unknown`
- [ ] **Experience preserved:** Refresh did not fabricate or clobber **My experience**
- [ ] **Dedup:** No silent overwrite of existing `name`
- [ ] **INDEX fresh:** Row count matches entry files; every row shows `kind`
- [ ] **No Perplexity for filed metadata**

## Additional resources

- Frontmatter examples and kind CORRECT/PROHIBITED: [reference.md](reference.md)
- Taxonomy: [roster/taxonomy.md](../../roster/taxonomy.md)
- Template: [roster/_template.md](../../roster/_template.md)
