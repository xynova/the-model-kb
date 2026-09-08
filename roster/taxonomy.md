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

CORRECT: Gemma 4 26B-A4B (CF) peer → Qwen3.5 / 3.8 MoE in a similar active-param band with reasoning, or Gemma 4 on DeepInfra  
PROHIBITED: Gemma 4 26B-A4B peer → older Qwen3-30B-A3B treated as the same segment without calling it a baseline

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
runtime: cloudflare
runtime_model_id: "@cf/meta/llama-3.2-3b-instruct"
lab: meta
```

CORRECT:
```text
kind: provider
lab: anthropic
runtime: cloudflare
price_note: "via Cloudflare AI Gateway proxy"
```

CORRECT:
```text
kind: service
service: cursor
lab: anthropic
```

PROHIBITED:
```text
kind: provider
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

## Categories (exactly one per entry)

Job-shaped (orthogonal to `kind`):

| Category | Use for |
|----------|---------|
| `general` | Default chat / reasoning / agents |
| `coding` | Code edit, refactor, agents-on-code |
| `multimodal` | Vision / image-in or similar |
| `budget` | Explicitly cost-first picks |

Directory path: `roster/<category>/<slug>.md`.

Do **not** use a `local` category; locality is `kind: sovereign` and optional `runtime: local`.

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

## Runtime values (optional)

Preferred host, especially for `sovereign`:

- `cloudflare`
- `deepinfra`
- `openrouter`
- `local`

Other values require a taxonomy update first.

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
| `tags` | yes | List from the allowed set |
| `status` | yes | `reach-for` / `watch` / `retired` |
| `lab` | yes | Weight or API lab origin (e.g. `meta`, `openai`, `anthropic`) |
| `refreshed` | yes | ISO date `YYYY-MM-DD` of last signal/price refresh |
| `source` | yes | Why this card exists |
| `price_in` | no | USD per 1M input tokens when a list price exists |
| `price_out` | no | USD per 1M output tokens when a list price exists |
| `price_note` | no | Neurons, Flex/Priority, free allocation, Cursor caps, proxy note, hardware cost |
| `runtime` | no | Preferred host from the runtime list |
| `runtime_model_id` | no | Id on that runtime (`@cf/...`, DeepInfra path, Cursor picker id) |
| `service` | when `kind: service` | e.g. `cursor` |
| `aa_intelligence` | no | Artificial Analysis Intelligence Index (or `unknown`) |
| `arena_elo` | no | LMArena text Elo when available |
| `aider_polyglot` | no | Coding-edit score when relevant |
| `context_window` | no | Tokens if known |
| `links` | no | Map; omit missing keys |

### Allowed `links` keys

- `cloudflare_pricing`
- `deepinfra`
- `provider_pricing`
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
| Lab pricing pages | Ground truth for `kind: provider` | per vendor |

**Drop for v1:** Hugging Face Open LLM Leaderboard. Skip secondary price-aggregator blogs. Do not use login dashboards (e.g. DeepInfra dash) for agent WebFetch.

Enrichment uses WebSearch / WebFetch only. MUST NOT file Perplexity output as metadata.
