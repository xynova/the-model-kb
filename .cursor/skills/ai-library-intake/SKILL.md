---
name: ai-library-intake
description: >-
  Intake AI news dumps or YouTube roundups into the curated library/: prepare
  transcripts, candidates, and stills with the Go tools/library-intake CLI
  (Fabric + Polypus GLM/Gemma), ask keep/drop one at a time, enrich keepers via
  WebSearch for project links,   write tagged markdown entries with optional media
  gallery, refresh INDEX.md and weekly additions. Use when the user says library
  intake, catalog roundup, keep into library, ingest this into the library,
  YouTube intake, fabric transcript, frame gallery, or attaches a weekly AI news
  dump or video URL for filing.
---

# AI library intake

Curate keepers into `library/`. This skill is a **consumer overlay** (not cursor-packs). Catalog root is repo-relative `library/`.

YouTube prep tooling lives in `tools/library-intake` (Go). The CLI shells out to installed `fabric`, `yt-dlp`, and `ffmpeg`, and calls Polypus over HTTP. It does **not** invent project links or write library entries.

## When to load

Load when the user asks to ingest, catalog, or keep items into the library, or attaches a roundup / chapter notes / YouTube URL for filing.

## Core constraints

**CONSTRAINT:** Catalog root and taxonomy

- MUST: Read `library/taxonomy.md` before assigning `category`, `does`, or `tags`.
- MUST: Write entries under `library/<category>/<slug>.md` using structure from `library/_template.md`.
- MUST NOT: Invent categories, `does` values, or tags absent from `taxonomy.md` (unless the user first updates taxonomy).

Enforcement: Read taxonomy before the proposal table; Grep `does` and tags against the allowed lists before write.
Violation: STOP, re-Read taxonomy, fix assignments, re-verify.

CORRECT:
```text
category: agents-memory
does: memory
tags: [memory, open, local, audio]
```

PROHIBITED:
```text
category: misc
does: cool-stuff
tags: [cool, new]
```

**CONSTRAINT:** YouTube prep via Go CLI

- MUST: When the intake source is a YouTube URL, prepare artifacts with `tools/library-intake` before keep asks.
- MUST: Prefer `go run ./cmd/library-intake -- prepare --url "<url>"` from `tools/library-intake`, or the built `bin/library-intake prepare` binary.
- MUST: Write working artifacts under `tmp/library-intake/` (transcript, candidates JSON, optional frames + frame-map). MUST NOT file those as library entries.
- MUST NOT: Invent transcript text, paste a guessed summary as the transcript, or call Perplexity for the transcript.
- MUST NOT: Use shell or Python helpers under `.cursor/skills/ai-library-intake/scripts/` (removed).
- MUST NOT: Require Nix; `fabric`, `yt-dlp`, and `ffmpeg` on PATH plus a running Polypus gateway are enough.

Enforcement: Shell history shows `library-intake` (or `go run ... library-intake`); transcript + candidates files exist before the first keep ask.
Violation: STOP, run prepare (or transcript + candidates), then continue.

CORRECT:
```bash
cd tools/library-intake
go run ./cmd/library-intake prepare --url "https://www.youtube.com/watch?v=nZYJdwM-_nI" --workdir ../../tmp/library-intake
```

PROHIBITED:
```text
Skip the Go CLI; paraphrase the video title into fake chapter notes
```

**CONSTRAINT:** Candidate parse via Polypus (YouTube and long dumps)

- MUST: For YouTube intakes, obtain the ordered candidate list from `library-intake` (`prepare` or `candidates`), which calls Polypus at `POLYPUS_BASE_URL` (default `http://127.0.0.1:1320`).
- MUST: Use model id `LIBRARY_INTAKE_POLYPUS_MODEL` when set; otherwise default `cf_local/@cf/zai-org/glm-4.7-flash`.
- MUST: Fail closed if Polypus is unreachable or returns non-JSON: tell the user, do not invent candidates.
- MUST: For short pasted dumps (not YouTube), the agent MAY parse candidates directly without Polypus.
- MUST NOT: Treat Polypus output as filed metadata (no invented `links`, licenses, or VRAM). Enrichment stays WebSearch after `y`.

Enforcement: `tmp/library-intake/<id>.candidates.parsed.json` exists (or tool history shows `library-intake candidates`) before the first keep ask on YouTube intakes.
Violation: STOP, fix Polypus reachability or model id, re-run candidates.

CORRECT:
```text
library-intake prepare → GLM candidates JSON → ask keep one-by-one
```

PROHIBITED:
```text
YouTube URL → agent invents a 12-item list with no library-intake and no Polypus
```

**CONSTRAINT:** Identifying stills (YouTube frames + Gemma)

- MUST: For YouTube intakes, prefer `library-intake prepare` (or `frames` then `curate`) so sparse frames land in `tmp/library-intake/<video-id>/frames/` and Gemma maps them via Polypus (`LIBRARY_INTAKE_VISION_MODEL`, default `cf_local/@cf/google/gemma-4-26b-a4b-it`).
- MUST: Keep frame work under `tmp/library-intake/` until a keeper gets `y`. On `y`, copy at most one `cover` and up to three `gallery` stills into `library/<category>/<slug>/media/`.
- MUST: Prefer project Open Graph / README / Hugging Face card images from WebSearch enrichment when they identify the project better than a roundup frame; stills from the video are a fallback or supplement.
- MUST NOT: Treat stills as a substitute for enrichment. Links, license, and run mode still come from WebSearch/WebFetch after `y`.
- MUST NOT: Use Fabric `--visual` (OCR text) as the gallery pipeline; that path is text-only.
- MUST NOT: Fail the whole intake if frame extract or Gemma mapping fails; `prepare` soft-fails media and continues keep/drop + enrichment without `media`.

Enforcement: For YouTube, prepare result or frame-map exists when media succeeded; on `y`, enrichment still runs.
Violation: STOP if stills replaced enrichment; restore WebSearch enrich, then continue.

CORRECT:
```text
library-intake prepare → candidates + optional frame-map
→ keep? y → WebSearch links + copy cover/gallery → write entry
```

PROHIBITED:
```text
Skip WebSearch because a nice demo frame exists; invent links from the video overlay
```

**CONSTRAINT:** Keep gate before any write (one candidate at a time)

- MUST: Parse all candidates first (skip intro/outro), then ask about **exactly one** candidate per turn.
- MUST: For each candidate, show name, suggested category, one-line why (future job or drop reason), and ask `keep? y/n` (or Yes/No). Wait for the reply before the next candidate.
- MUST NOT: Present a batch keep list and treat a single `y`/`yes` as approval of every candidate.
- MUST NOT: Auto-keep the whole dump.
- MUST NOT: Write that candidate's entry (or refresh INDEX for it) until the user answers `y` for that candidate.
- MUST: On `n`, skip with no file; move to the next candidate.
- MUST: On `y`, enrich (WebSearch links required) → attach curated stills when available → write → `library-intake index refresh`, then ask about the next candidate (unless the user stops the intake).

Enforcement: Each keep decision in chat names one candidate; no multi-keeper Write before one-per-candidate `y`.
Violation: STOP, do not batch-write; revert unsolicited batch files if the user asks; resume one-at-a-time.

CORRECT:
```text
[10/21] GLM-5.3-Flash | models

Open multimodal model from Z.ai (text and images). Downloadable MIT weights; also an API. Claims strong coding and agent results at much lower API cost than full GLM-5.3 (roughly cents per million tokens). Local full-size needs multi-GPU; quants reduce that. Previously seen as Ox Alpha.

keep? y/n
```

PROHIBITED (fluffy):
```text
This is getting a lot of attention because people love the mix of…
Worth keeping if you will compare… Easy skip if…
```

PROHIBITED (datasheet / jargon):
```text
- Performance: ...
- Price: ...
320B MoE hybrid sparse+linear…
```

Also PROHIBITED: batch “suggested keep set, reply y to file all.”

**CONSTRAINT:** Enrichment (v1; required on every keep)

- MUST: For each keeper, run WebSearch (and WebFetch/GitHub when a repo URL is found) for homepage, license, local vs API, rough resource needs.
- MUST: Mark uncertain fields as `unknown` or omit the key.
- MUST: Prefer primary project pages over secondary blog/YouTube recaps when both appear.
- MUST NOT: Invent URLs, licenses, or VRAM numbers.
- MUST NOT: Call Perplexity MCP for enrichment in v1.
- MUST NOT: Treat Fabric transcript wording, Polypus summaries, or on-screen frame text as proof of license, GitHub URL, or VRAM.
- MUST NOT: Skip enrichment because media stills already exist.

Enforcement: Every `links.*` value must come from search/fetch results or be omitted; tool history shows WebSearch after each `y`.
Violation: STOP, clear invented fields, re-enrich, re-verify.

CORRECT:
```text
WebSearch "VoiceMem GitHub" → set links.github from result; vram_hint: unknown
optional: download OG image → media.cover; copy Gemma-picked frame → media.gallery
```

PROHIBITED:
```text
Guess github.com/.../voicemem; call perplexity_research for the entry
Write entry with only media.cover and empty/invented links
```

**CONSTRAINT:** Dedup and INDEX

- MUST: Before write, Grep/Glob `library/**/*.md` for matching `name:`.
- MUST: Skip existing `name` unless the user asks to update.
- MUST: After successful writes, run `library-intake index refresh --library <repo>/library` so global `INDEX.md`, each `library/<category>/INDEX.md`, and `library/additions/` week files stay current. Week open/closed status comes from each entry's `added` date (current ISO week is open; older weeks are closed).

Enforcement: Grep `^name: <slug>` across library before Write; INDEX row counts match entry files (excluding `_template.md`); additions week for the keeper's `added` ISO week lists the name.
Violation: STOP, skip or update per user, rebuild indexes.

CORRECT:
```text
Grep name: voicemem → miss → write → index refresh → models/INDEX.md and additions/YYYY-Www.md include voicemem
```

PROHIBITED:
```text
Overwrite existing voicemem.md without asking; leave INDEX and additions stale
```

## Steps

Copy this checklist and track progress:

```
- [ ] Read library/taxonomy.md and library/_template.md
- [ ] If YouTube: tools/library-intake prepare → tmp/library-intake/
- [ ] Load candidates JSON (or agent-parse short paste)
- [ ] For each candidate: ask keep? y/n; wait
- [ ] On y: dedup, WebSearch enrich (required), attach media if any, write entry, `index refresh`
- [ ] On n: skip; next candidate
- [ ] After last candidate: summarize kept vs skipped
```

1. **Load taxonomy** — Read `library/taxonomy.md` and `library/_template.md`. Fail if missing.

2. **Acquire source**
   - **YouTube URL:** From repo root, run:
     ```bash
     cd tools/library-intake && go run ./cmd/library-intake prepare \
       --url "<youtube-url>" --workdir ../../tmp/library-intake
     ```
     Fail closed if transcript or candidates fail. Soft-fail media is OK (`MediaSkipped` in prepare JSON).
   - **Pasted dump:** Use the attached text as-is.

3. **Parse** — Load `tmp/library-intake/<id>.candidates.parsed.json` for YouTube. Skip intro/outro with no tool to file. One candidate = one product/model/system. Do not write yet.
   - Short paste: agent parse is allowed.
   - Set `source` for later frontmatter to the URL or dump label (e.g. `youtube:nZYJdwM-_nI`).

4. **Ask one** — For the next candidate only: name, suggested category, then a **decision cold read**. End with `keep? y/n`.
   - MUST cover, when known: what it is; performance or why it matters; price/cost; how you run it; one caveat.
   - Voice: **direct but readable**. Short complete sentences. Lead with facts. No filler, no hype, no “worth keeping if…” pep talk unless it adds a real decision fork.
   - MUST NOT: Fluffy openers (“people are excited”, “the reason it is getting attention”). MUST NOT: Datasheet bullet stacks as the default. MUST NOT: Paper jargon unless the user asks `explain`.
   - Target length: about 4–8 sentences, or two tight paragraphs.
   - MUST NOT: Ask about another candidate in the same turn.
   - If a Gemma-mapped cover exists for this candidate under `tmp/`, MAY Read that image in the ask turn so the human sees it.
   - If price/performance are missing and the item is more than spectacle, MAY quick-WebSearch those facts only; full enrichment after `y`.
   - On user `explain`, add one plain-language beat, then re-ask `keep? y/n`.

5. **On `n`** — Skip; no entry file. Continue with step 4 for the next candidate (or stop if the user ends intake).

6. **On `y`** — Dedup (`name:` Grep). **WebSearch enrich first** (homepage, GitHub, license, run mode, `vram_hint` or `unknown`). Then attach media: copy Gemma cover/gallery stills and/or a verified project OG/card image into `library/<category>/<slug>/media/`; set optional `media.cover` / `media.gallery`. Write ~320px `thumb-*.jpg` copies and a compact markdown image table under **What it is** for Cursor preview. The Hugo site replaces that table with a GLightbox widget from frontmatter. No separate Gallery section. Write `library/<category>/<slug>.md`. Run:
   ```bash
   cd tools/library-intake && go run ./cmd/library-intake index refresh --library ../../library
   ```
   That refreshes global `INDEX.md`, each category `INDEX.md`, and the weekly additions file (section headings per category; open/closed from `added` dates). Then continue with step 4 for the next candidate.

7. **Done** — After the last candidate (or user stop), summarize kept vs skipped. Do not re-ask batch approval.

## INDEX format

Global `library/INDEX.md` (and matching `library/<category>/INDEX.md`):

```markdown
# Library index

Regenerated by `library-intake index refresh` after writes.

| name | title | category | summary | tags | run | license |
|------|-------|----------|---------|------|-----|---------|
| voicemem | VoiceMem | agents-memory | Local voice memory for companion agents. | memory, open, local, audio | local | open |
```

Weekly additions (`library/additions/YYYY-Www.md`) use calendar titles (for example `Sep 14–20, 2026`) and group keepers under `## <category>` tables (summary / tags / run / license). `library/additions/INDEX.md` lists weeks.
## Pre-completion checklist

- [ ] **Taxonomy loaded:** `library/taxonomy.md` was Read this run
      Method: Confirm Read in tool history before proposal
      Pass: Read occurred
      Fail: STOP, Read taxonomy, re-run proposal if needed
- [ ] **YouTube path held:** If source was a URL, `library-intake` transcript + candidates exist
      Method: Shell/tool history + `tmp/library-intake/` files
      Pass: transcript and candidates.parsed.json present
      Fail: STOP, run prepare
- [ ] **Polypus used for YouTube:** Candidate list came from library-intake / `:1320`
      Method: candidates JSON + CLI history
      Pass: JSON candidates parsed
      Fail: STOP, do not invent list
- [ ] **Enrichment on every keep:** WebSearch/WebFetch ran after each `y` before write
      Method: Tool history for that keeper
      Pass: Search occurred; links traceable or omitted
      Fail: STOP, enrich before write
- [ ] **Stills additive:** Media did not skip enrichment; paths under `library/.../media/` only after `y`
      Method: Diff entry frontmatter + media dir vs enrich history
      Pass: Links from search; stills optional
      Fail: STOP, restore enrichment
- [ ] **Keep gate held:** Each write followed a per-candidate `y`; no batch `y` for many keepers
      Method: Timeline: one ask → one reply → at most one new entry
      Pass: One-at-a-time
      Fail: STOP, revert batch if user requests, resume interactive
- [ ] **Tags legal:** Every tag is in taxonomy allowed list
      Method: Diff tags vs taxonomy
      Pass: All tags listed
      Fail: STOP, fix or ask user to extend taxonomy
- [ ] **Does legal:** Every entry has exactly one `does` from the taxonomy Does table
      Method: Diff `does` vs taxonomy
      Pass: Value listed; INDEX shows does/run/license
      Fail: STOP, fix or ask user to extend taxonomy
- [ ] **No invented links:** Every URL came from WebSearch/WebFetch or was omitted
      Method: Spot-check `links:` block
      Pass: Traceable or omitted
      Fail: STOP, clear inventions
- [ ] **Dedup:** No silent overwrite of existing `name`
      Method: Grep before write
      Pass: Skip or user-approved update
      Fail: STOP, skip duplicate
- [ ] **INDEX fresh:** `index refresh` ran after writes; global + category INDEX rows match entries; additions week lists keepers
      Method: Count rows vs entry files; spot-check additions week
      Pass: Counts match; week file has section for each kept category
      Fail: STOP, run `library-intake index refresh`- [ ] **No Perplexity:** Enrichment used WebSearch/WebFetch only
      Method: Tool history
      Pass: No perplexity_* calls
      Fail: STOP, do not treat Perplexity output as filed metadata

## Additional resources

- Go CLI recipes and schemas: [reference.md](reference.md)
- Go tool: [tools/library-intake/instructions.md](../../tools/library-intake/instructions.md)
- Taxonomy: [library/taxonomy.md](../../library/taxonomy.md)
- Template: [library/_template.md](../../library/_template.md)
