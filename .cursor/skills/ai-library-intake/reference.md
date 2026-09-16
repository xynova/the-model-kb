# AI library intake reference

LOAD-WHEN: filling entry frontmatter, choosing keep vs drop, checking INDEX shape, or running the Go `tools/library-intake` CLI during `ai-library-intake`.

## Go CLI (`tools/library-intake`)

Prep tooling is Go only. It shells out to `fabric`, `yt-dlp`, and `ffmpeg`, and calls Polypus over HTTP. It never invents project links or writes `library/` entries.

Defaults:

| Env | Default |
|-----|---------|
| `POLYPUS_BASE_URL` | `http://127.0.0.1:1320` |
| `LIBRARY_INTAKE_POLYPUS_MODEL` | `cf_local/@cf/zai-org/glm-4.7-flash` |
| `LIBRARY_INTAKE_VISION_MODEL` | `cf_local/@cf/google/gemma-4-26b-a4b-it` |
| `LIBRARY_INTAKE_WORK_ROOT` | `tmp/library-intake` |

### Full YouTube prepare

CORRECT:

```bash
cd tools/library-intake
go run ./cmd/library-intake prepare \
  --url "https://www.youtube.com/watch?v=nZYJdwM-_nI" \
  --workdir ../../tmp/library-intake
```

Writes (when successful):

- `tmp/library-intake/<id>.transcript.txt`
- `tmp/library-intake/<id>.candidates.parsed.json`
- `tmp/library-intake/<id>/frames/frame_*.jpg` (soft-fail OK)
- `tmp/library-intake/<id>.frame-map.json` (soft-fail OK)

Stdout is a JSON summary with `MediaSkipped` / `MediaSkipReason` when frames or Gemma curation fail.

### Individual commands

CORRECT:

```bash
go run ./cmd/library-intake transcript --url "URL" --out ../../tmp/library-intake/<id>.transcript.txt
go run ./cmd/library-intake candidates --transcript ../../tmp/library-intake/<id>.transcript.txt
go run ./cmd/library-intake frames --url "URL" --outdir ../../tmp/library-intake/<id>/frames --interval 8
go run ./cmd/library-intake curate \
  --frames-dir ../../tmp/library-intake/<id>/frames \
  --candidates ../../tmp/library-intake/<id>.candidates.parsed.json
```

PROHIBITED:

```bash
fabric -y "URL" -p summarize          # chat path; not transcript-only
fabric -y "URL" --visual              # OCR text; not gallery stills
.cursor/skills/.../scripts/*.sh       # removed helpers
```

**Rules:**

- MUST keep working files under `tmp/library-intake/`.
- NEVER commit transcript dumps or downloaded videos as `library/` entries.
- MUST have `fabric`, `yt-dlp`, and `ffmpeg` on PATH for full prepare.
- MUST soft-fail media; MUST fail closed on transcript or candidates.

## Candidate JSON schema

```json
{
  "candidates": [
    {
      "title": "VoiceMem",
      "name_hint": "voicemem",
      "category": "agents-memory",
      "why": "Local voice-agent memory for companion pipelines",
      "cold_read": "What it is; performance or relevance; price if known; how you run it; one caveat."
    }
  ]
}
```

**Rules:**

- MUST use only `category` values from `library/taxonomy.md`.
- MUST skip intro/outro spectacle with no tool to file.
- NEVER invent homepage, GitHub, license, or VRAM in this JSON.

## Frame-map JSON schema

```json
{
  "assignments": [
    {
      "file": "frame_0012.jpg",
      "candidate_title": "Merryold V2",
      "role": "cover",
      "reason": "depth map demo still"
    }
  ]
}
```

`role` is `cover`, `gallery`, or `skip`. `candidate_title` may be null.

### On keep: enrichment first, then media

CORRECT:

```text
1. WebSearch + WebFetch → links.homepage / links.github / license / run
2. Prefer project OG/HF card image when found → media.cover
3. Else or also: copy Gemma cover (+ up to 3 gallery) from tmp frames → library/<category>/<slug>/media/
4. Write entry with links always attempted; media optional
```

PROHIBITED:

```text
File entry with media only and skip WebSearch because the frame "looks official"
```

## Frontmatter schema

### Complete keeper entry

```yaml
---
name: glm-5.3
title: GLM 5.3
category: models
tags:
  - llm
  - open
  - local
  - moe
  - quantized
status: watch
run: local
license: open
vram_hint: multi-gpu
links:
  homepage: https://example.com/glm
  github: https://github.com/example/glm
  weights: https://huggingface.co/example/glm
media:
  cover: models/glm-5.3/media/cover.jpg
  gallery:
    - models/glm-5.3/media/01.jpg
source: "weekly roundup 2026-08-31"
added: 2026-08-31
---
```

**Rules:**

- MUST use kebab-case `name` unique across `library/`.
- MUST set exactly one `category` from `library/taxonomy.md`.
- MUST use only allowed tags from taxonomy.
- MUST NOT invent `links` values; omit keys when unknown.
- MUST set `vram_hint: unknown` when evidence is missing (or omit the field).
- MAY omit `media` entirely when no reliable still exists.

### Minimal legal entry

```yaml
---
name: fibo-1.5
title: Fibo 1.5
category: creative
tags:
  - image
  - open
  - local
status: watch
run: local
license: open
source: "weekly roundup 2026-08-31"
added: 2026-08-31
---
```

**Rules:**

- MUST include all required fields from taxonomy even when `links` and `vram_hint` are absent.
- NEVER leave `tags:` empty.

## Keep vs drop

### Interactive one-at-a-time gate

CORRECT:

```text
[10/21] GLM-5.3-Flash | models

Open multimodal model from Z.ai (text and images). MIT weights plus a cheap API (cents per million tokens vs full GLM-5.3). Pitched for strong coding/agent work; local full-size needs multi-GPU.

keep? y/n
```

PROHIBITED: fluffy hype paragraphs; default bold Performance/Price/Run stacks; paper jargon before `y`; batch keep sets.

**Rules:**

- MUST ask exactly one candidate per turn.
- MUST include performance and price (or say unknown) when available.
- MUST be direct but readable: facts in short sentences, no filler.
- MUST NOT treat a lone `y`/`yes` as approval of a multi-item keep set.
- MUST still ask on drop-lean items so the user can override.
- NEVER write entry files for `n` replies.

### Keep: future job is nameable

CORRECT (lean keep in the one-liner):

```text
[7/21] VoiceMem | agents-memory | Local voice-agent memory for companion/voice pipelines
keep? y/n
```

PROHIBITED (auto-keep with no job):

```text
[15/21] World Humanoid Robot Games | robotics | Cool robots broke records
keep? y/n
```
(Ask anyway if it is a real candidate, but the why MUST say spectacle / no tool to run.)

### Drop lean: no shelf value yet

CORRECT:

```text
[15/21] World Humanoid Robot Games | robotics | Spectacle; nothing to run or wire
keep? y/n
```

PROHIBITED:

```text
Skip asking about drops; only ask about interesting keepers
```

**Rules:**

- MUST still ask `y/n` on drop-lean candidates so the user can override.
- NEVER write entry files for `n`.

## Enrichment

### WebSearch then write

CORRECT:

```text
1. WebSearch "VoiceMem AI voice agent memory GitHub"
2. Confirm repo URL and license from results/fetch
3. Set links.github, license: open, run: local
4. Write agents-memory/voicemem.md
```

PROHIBITED:

```text
links:
  github: https://github.com/someone/voicemem  # guessed, not searched
license: open  # assumed because "open-source" in a YouTube summary alone
```

**Rules:**

- MUST search before setting homepage/github/weights.
- MUST NOT call Perplexity MCP in v1.
- MUST prefer primary project pages over secondary blog recaps when both appear.
- MUST run enrichment on every `y` even when curated stills already exist.
- MAY download a verified project Open Graph / Hugging Face card image into `media/` after links are confirmed.

## Dedup

### Existing name

CORRECT:

```text
Grep name: voicemem → found agents-memory/voicemem.md
→ Skip write; tell user; update only if they ask
```

PROHIBITED:

```text
Overwrite agents-memory/voicemem.md without asking
```

**Rules:**

- MUST Grep `^name: <slug>` under `library/` before Write.
- NEVER silently overwrite.

## INDEX regeneration

### After two writes

CORRECT:

```markdown
# Library index

Regenerated by `ai-library-intake` after writes.

| name | title | category | status | tags | path |
|------|-------|----------|--------|------|------|
| block-3d | Block 3D | 3d | watch | 3d, open, local | 3d/block-3d.md |
| voicemem | VoiceMem | agents-memory | watch | memory, open, local, audio | agents-memory/voicemem.md |
```

**Rules:**

- MUST include every `library/<category>/*.md` except empty dirs.
- MUST NOT list `_template.md`, `README.md`, `taxonomy.md`, or `INDEX.md` as entries.
- MUST sort by category then name (stable scan order is enough if documented in the rewrite).

## Smoke path (pasted dump / VoiceMem)

Mental walkthrough only; do not file unless the user says `y`:

1. Parse → ordered candidates including VoiceMem.
2. Ask → `[n] VoiceMem | agents-memory | Local voice memory for companion agents` + `keep? y/n`.
3. Wait → user `y`.
4. Enrich → WebSearch; fill links/license/run.
5. Write → `library/agents-memory/voicemem.md`; refresh INDEX.
6. Ask → next candidate.

## Smoke path (YouTube)

Mental walkthrough only; do not file unless the user says `y`:

1. `library-intake prepare` → transcript + candidates (+ optional frames/frame-map).
2. Ask → first candidate cold read (+ optional cover image) + `keep? y/n`.
3. On `y` → WebSearch enrich (required) → attach media if any → write → INDEX; then next candidate.
4. On `n` → skip; next candidate.
