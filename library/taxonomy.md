# Library taxonomy

Source of truth for categories, tags, and keep criteria. The `ai-library-intake` skill MUST Read this file before assigning `category` or `tags`.

## Keep criteria

Keep an item only if you can name a future job for it:

- Run it locally
- Wire it into a pipeline
- Compare it when choosing a model or tool

Drop news, spectacle, closed APIs you will not pay for, and research with no attached job.

## Categories (exactly one per entry)

| Category | Use for |
|----------|---------|
| `models` | LLMs and multimodal foundation models |
| `video` | Video generation, extension, transcription tied to video |
| `3d` | Text-to-3D, scene reconstruction, 3D cleanup |
| `agents-memory` | Agent memory, voice memory, long-context agent systems |
| `robotics` | Robot models, demos, humanoid systems |
| `hardware` | Local AI boxes, chips, offline personal AI devices |
| `training-research` | Training methods, benchmarks, research techniques |
| `creative` | Agentic creative workspaces, image gen/edit products |

Directory path: `library/<category>/<slug>.md`.

## Allowed tags

Tags are lowercase kebab-case. Use only values from this list unless the user explicitly adds a new tag to this file first.

### Modality

- `llm`
- `vision`
- `video`
- `3d`
- `image`
- `audio`
- `memory`
- `robotics`
- `hardware`

### Run / access

- `local`
- `api`
- `hardware`

### License

- `open`
- `closed`
- `mixed`

### Other

- `quantized`
- `vram-heavy`
- `moe`
- `agentic`

## Status values

Exactly one of:

| Status | Meaning |
|--------|---------|
| `watch` | Interesting; revisit when a need appears |
| `reach-for` | Prefer this when the job matches |
| `drop-until-needed` | Noted; do not invest until a concrete job exists |

## Frontmatter fields

| Field | Required | Notes |
|-------|----------|-------|
| `name` | yes | Stable slug id (kebab-case), unique across the library |
| `title` | yes | Human display name |
| `category` | yes | One of the category keys above |
| `tags` | yes | List from the allowed set |
| `status` | yes | `watch` / `reach-for` / `drop-until-needed` |
| `run` | yes | `local` / `api` / `hardware` |
| `license` | yes | `open` / `closed` / `mixed` |
| `vram_hint` | no | Short hint, e.g. `consumer`, `multi-gpu`, `unknown` |
| `links` | no | Map: `homepage`, `github`, `weights` (omit missing keys) |
| `source` | yes | Where the candidate came from |
| `added` | yes | ISO date `YYYY-MM-DD` |

Uncertain enrichment fields: use `unknown` or omit the key. MUST NOT invent URLs or licenses.
