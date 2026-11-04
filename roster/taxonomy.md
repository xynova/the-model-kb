# Roster taxonomy

Source of truth for kind, categories, tags, status, and keep criteria. The `model-roster` skill MUST Read this file before assigning `kind`, `category`, or `tags`.

## Keep criteria

Keep a card only if you can name a job where you would open this entry to choose a model. Drop spectacle and cards with no decision value.

## Same-segment peers (apples to apples)

When adding competitors or “same segment” cards next to an anchor model, peers MUST match the decision class, not merely the same host or rough parameter count.

Match on:

- **Generation / cut date** — same model generation or newer (e.g. Gemma 4 peers, not an older Qwen3 cut sold as equivalent)
- **Capability class** — comparable reasoning / tools / agent strength when that is why the anchor exists
- **Modality** — do not pair a vision/multimodal anchor with a text-only peer unless the cold read says the modality gap is intentional
- **Active size band** for MoE — compare active-parameter class, not only total parameters

Older or weaker cuts MAY be offered only as an explicit **baseline** (cold read must say baseline, not peer). Same open weights on another host are always a fair peer (host price compare), not a generation compare.

CORRECT: Gemma 4 26B-A4B peer → Qwen3.5 / 3.8 MoE in a similar active-param band with reasoning  
PROHIBITED: Gemma 4 26B-A4B peer → older Qwen3-30B-A3B treated as the same segment without calling it a baseline

Same open weights on another host are **not** a second card. Append a row to `hosts` on the existing card (host price compare).

## Kind (exactly one per entry; required)

| Kind | Meaning | Typical examples |
|------|---------|------------------|
| `sovereign` | Open / self-hostable weights you can **switch hosts** for | Local GGUF; Docker images in this repo; Cloud Run llama.cpp; Cloudflare-hosted or DeepInfra-hosted open models |
| `provider` | Lab-locked API; weights **not** relocatable. A host counts only when it **proxies** that API | OpenAI, Anthropic, Gemini; CF/DeepInfra only if proxying those |
| `service` | Model access **bundled by a product** | Cursor-supplied models |

### Kind rules

- Portability defines sovereignty, not who owns today's GPUs.
- Cloudflare (and similar gateways) are **not** a kind. Same account can back a `sovereign` card (hosted open model) and a `provider` card (proxy to Claude/OpenAI/Gemini).
- Do **not** label Cloudflare or DeepInfra open hosted weights as `provider`.
- Do **not** label a host as `sovereign` when it only proxies a locked lab API.
- Cursor (and similar) are always `service`, even when the underlying lab model is Claude or GPT.
- Weight/lab origin field is `lab` (e.g. `openai`, `anthropic`, `meta`). MUST NOT use a frontmatter key named `provider` (collides with `kind: provider`).

### CORRECT / PROHIBITED

CORRECT:
```text
kind: sovereign
lab: meta
jobs: [chat, cheap]
preferred_host: cloudflare
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/meta/llama-3.2-3b-instruct"
    price_in: 0.051
    price_out: 0.335
    capacity: hobby
```

CORRECT:
```text
kind: provider
lab: anthropic
jobs: [coding]
hosts:
  - id: cloudflare
    runtime: cloudflare
    price_note: "via Cloudflare AI Gateway proxy"
    capacity: unknown
```

CORRECT:
```text
kind: service
service: cursor
lab: anthropic
jobs: [coding]
hosts:
  - id: cursor
    runtime_model_id: composer-2.5
    price_in: 3
    price_out: 15
    capacity: unknown
```

PROHIBITED:
```text
kind: provider
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/meta/llama-3.2-3b-instruct"
# open hosted weights are sovereign
```

PROHIBITED:
```text
kind: sovereign
service: cursor
# product-bundled access is service
```

PROHIBITED:
```text
# blended rank or score computed in YAML
pick_score: 0.87
```

Scores and host ranks belong in a later dbt layer, not in markdown or `catalog.json` schema 1.

## Categories (exactly one per entry)

Job-shaped (orthogonal to `kind`):

| Category | Use for |
|----------|---------|
| `general` | Default chat / reasoning / agents |
| `coding` | Code edit, refactor, agents-on-code |
| `multimodal` | Vision / image-in or similar |
| `budget` | Explicitly cost-first picks |

Directory path: `roster/<category>/<slug>.md`.

Do **not** use a `local` category; locality is `kind: sovereign` and a host row with `id: local` or `runtime: local`.

## Jobs (required on every card)

Routing keys. Use only values from this list:

- `chat`
- `reasoning`
- `coding`
- `vision`
- `long-context`
- `cheap`
- `high-throughput`

Assign every kept card at least one job. Category and tags help choose jobs but `jobs` is the machine key.

## Allowed tags

Tags are lowercase kebab-case. Use only values from this list unless the user explicitly adds a new tag to this file first.

### Modality

- `llm`
- `vision`
- `video`
- `image`
- `audio`

### Run / access flavor

- `local`
- `api`
- `open`
- `closed`
- `moe`
- `agentic`
- `quantized`

### Roster-specific

- `reasoning`
- `fast`
- `cheap`
- `coding`
- `cloudflare`
- `deepinfra`
- `cursor`

Do **not** encode `kind` only in tags; `kind` is the source of truth.

## Status values

Exactly one of:

| Status | Meaning |
|--------|---------|
| `reach-for` | Prefer this when the job matches |
| `watch` | On the shortlist; not primary yet |
| `retired` | Superseded; keep for history, do not recommend |

## Runtime values (on host rows)

Use on each `hosts` entry when `kind` is `sovereign` or `provider`:

- `cloudflare`
- `deepinfra`
- `openrouter`
- `google`
- `xinference`
- `local`

Do **not** use `cursor` as a runtime. Cursor offers use `id: cursor` on `kind: service` cards without `runtime`.

## Capacity (required on every host row)

Exactly one of:

| Value | Meaning |
|-------|---------|
| `hobby` | Default hosted limits, neurons/day, low RPM, dev use |
| `industrial` | High throughput or operator-confirmed production limits |
| `local` | Self-hosted hardware |
| `unknown` | Limits not filed yet |

Public pages rarely publish RPM/TPM. File `unknown` or an operator note in `limits_note`; do not invent numeric limits.

## Host row ids

Stable `id` on each `hosts` entry:

- `cloudflare`
- `deepinfra`
- `openrouter`
- `google`
- `xinference`
- `cursor`
- `local`

Other ids require a taxonomy update first. Use project skill `providers-mint` (one public pricing URL, `keep? y/n`) so `ALLOWED_HOST_IDS` and retrieve-limits `LINK_KEYS` stay in sync. MUST NOT add an id only in markdown.

## Minting a host

A new serving company is a **mint**, not a nightly mine fetch. Mine only GETs hostnames already in `LINK_KEYS`.

- MUST: Fetch the host's own public pricing or model page (WebSearch / WebFetch).
- MUST NOT: File Perplexity output as `id`, URL, price, or capacity.
- MUST: Add `id` here, in `site/scripts/generate-roster.py` allowlists, and in `site/scripts/retrieve-limits.py` `LINK_KEYS` on the same keep.
- MUST NOT: Reuse `provider_pricing` for a non-Google hostname.

## Service values

Required when `kind: service`:

- `cursor`

## Frontmatter fields

| Field | Required | Notes |
|-------|----------|-------|
| `name` | yes | Stable slug id (kebab-case), unique across the roster |
| `title` | yes | Human display name |
| `kind` | yes | `sovereign` / `provider` / `service` |
| `category` | yes | One of the category keys above |
| `jobs` | yes | List from the jobs list above |
| `tags` | yes | List from the allowed set |
| `status` | yes | `reach-for` / `watch` / `retired` |
| `lab` | yes | Weight or API lab origin (e.g. `meta`, `openai`, `anthropic`) |
| `refreshed` | yes | ISO date `YYYY-MM-DD` of last signal/price refresh |
| `source` | yes | Why this card exists |
| `hosts` | yes | List of host offer maps (see below); at least one row |
| `preferred_host` | no | Must match a host `id` when set |
| `aliases` | no | Hugo redirect paths when a slug was renamed |
| `service` | when `kind: service` | e.g. `cursor` |
| `aa_intelligence` | no | Artificial Analysis Intelligence Index (or `unknown`) |
| `arena_elo` | no | LMArena text Elo when available |
| `aider_polyglot` | no | Coding-edit score when relevant |
| `context_window` | no | Tokens if known |
| `links` | no | Card-level map; omit missing keys |

### Host row fields (`hosts` list)

| Field | Required | Notes |
|-------|----------|-------|
| `id` | yes | Stable host id from the list above |
| `runtime` | when card `kind` is `sovereign` or `provider` | From runtime list |
| `runtime_model_id` | no | Id on that runtime (`@cf/...`, DeepInfra path, Cursor picker id) |
| `price_in` | no | USD per 1M input tokens when a list price exists |
| `price_out` | no | USD per 1M output tokens when a list price exists |
| `price_note` | no | Neurons, plan caps, proxy note |
| `capacity` | yes | `hobby` / `industrial` / `local` / `unknown` |
| `limits_note` | no | Neurons/day, RPM caveats, marketplace notes |
| `links` | no | Host-specific keys from the links list below |

MUST NOT put `kind` on a host row. MUST NOT put card-level `runtime`, `price_in`, `price_out`, or `price_note` on the card (use `hosts`).

### Allowed `links` keys

- `cloudflare_pricing`
- `deepinfra`
- `provider_pricing`
- `xinference`
- `service_docs`
- `artificial_analysis`
- `openrouter`
- `models_dev`
- `aider`

Uncertain enrichment fields: use `unknown` or omit the key. MUST NOT invent URLs, prices, or scores.

## Allowed enrichment sources (v1)

| Source | Role | Primary URL |
|--------|------|-------------|
| Artificial Analysis | Intelligence / cost / speed | https://artificialanalysis.ai/leaderboards/models |
| LMArena (Arena AI) | Human-preference Elo | https://lmarena.ai/ |
| OpenRouter | Cross-host pricing / ids | https://openrouter.ai/models |
| models.dev | Specs, pricing, capabilities | https://models.dev/ |
| Aider Polyglot | Coding-edit signal | https://aider.chat/docs/leaderboards/ |
| Cloudflare Workers AI | Price when `runtime: cloudflare` | https://developers.cloudflare.com/workers-ai/platform/pricing/ |
| DeepInfra | Price when `runtime: deepinfra` | https://deepinfra.com/pricing |
| Xinference | Price when `runtime: xinference` | https://xinference.co/pricing.html |
| Lab pricing pages | Ground truth for `kind: provider` | per vendor |
| Minted host pricing pages | Price when `runtime` matches a `providers-mint` id | that host's public page (skill `providers-mint`) |

**Drop for v1:** Hugging Face Open LLM Leaderboard. Skip secondary price-aggregator blogs. Do not use login dashboards (e.g. DeepInfra dash) for agent WebFetch.

Enrichment uses WebSearch / WebFetch only. MUST NOT file Perplexity output as metadata.
