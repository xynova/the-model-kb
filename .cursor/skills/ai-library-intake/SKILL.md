---
name: ai-library-intake
description: >-
  Intake AI news dumps into the curated library/: propose keep/drop, enrich
  keepers via WebSearch, write tagged markdown entries, refresh INDEX.md. Use
  when the user says library intake, catalog roundup, keep into library, ingest
  this into the library, or attaches a weekly AI news dump for filing.
---

# AI library intake

Curate keepers into `library/`. This skill is a **consumer overlay** (not cursor-packs). Catalog root is repo-relative `library/`.

## When to load

Load when the user asks to ingest, catalog, or keep items into the library, or attaches a roundup / chapter notes for filing.

## Core constraints

**CONSTRAINT:** Catalog root and taxonomy

- MUST: Read `library/taxonomy.md` before assigning `category` or `tags`.
- MUST: Write entries under `library/<category>/<slug>.md` using structure from `library/_template.md`.
- MUST NOT: Invent categories or tags absent from `taxonomy.md` (unless the user first updates taxonomy).

Enforcement: Read taxonomy before the proposal table; Grep tags against the allowed list before write.
Violation: STOP, re-Read taxonomy, fix assignments, re-verify.

CORRECT:
```text
category: agents-memory
tags: [memory, open, local, audio]
```

PROHIBITED:
```text
category: misc
tags: [cool, new]
```

**CONSTRAINT:** Keep gate before any write (one candidate at a time)

- MUST: Parse all candidates first (skip intro/outro), then ask about **exactly one** candidate per turn.
- MUST: For each candidate, show name, suggested category, one-line why (future job or drop reason), and ask `keep? y/n` (or Yes/No). Wait for the reply before the next candidate.
- MUST NOT: Present a batch keep list and treat a single `y`/`yes` as approval of every candidate.
- MUST NOT: Auto-keep the whole dump.
- MUST NOT: Write that candidate's entry (or refresh INDEX for it) until the user answers `y` for that candidate.
- MUST: On `n`, skip with no file; move to the next candidate.
- MUST: On `y`, enrich → write → refresh INDEX for that keeper, then ask about the next candidate (unless the user stops the intake).

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

**CONSTRAINT:** Enrichment (v1)

- MUST: For each keeper, run WebSearch (and WebFetch/GitHub when a repo URL is found) for homepage, license, local vs API, rough resource needs.
- MUST: Mark uncertain fields as `unknown` or omit the key.
- MUST NOT: Invent URLs, licenses, or VRAM numbers.
- MUST NOT: Call Perplexity MCP for enrichment in v1.

Enforcement: Every `links.*` value must come from search/fetch results or be omitted.
Violation: STOP, clear invented fields, re-enrich, re-verify.

CORRECT:
```text
WebSearch "VoiceMem GitHub" → set links.github from result; vram_hint: unknown
```

PROHIBITED:
```text
Guess github.com/.../voicemem; call perplexity_research for the entry
```

**CONSTRAINT:** Dedup and INDEX

- MUST: Before write, Grep/Glob `library/**/*.md` for matching `name:`.
- MUST: Skip existing `name` unless the user asks to update.
- MUST: Regenerate `library/INDEX.md` after successful writes (all entry rows).

Enforcement: Grep `^name: <slug>` across library before Write; INDEX row count matches entry files (excluding `_template.md`).
Violation: STOP, skip or update per user, rebuild INDEX.

CORRECT:
```text
Grep name: voicemem → miss → write → INDEX includes voicemem row
```

PROHIBITED:
```text
Overwrite existing voicemem.md without asking; leave INDEX stale
```

## Steps

Copy this checklist and track progress:

```
- [ ] Read library/taxonomy.md and library/_template.md
- [ ] Parse candidates (skip intro/outro fluff)
- [ ] For each candidate: ask keep? y/n; wait
- [ ] On y: dedup, WebSearch enrich, write entry, refresh INDEX
- [ ] On n: skip; next candidate
- [ ] After last candidate: summarize kept vs skipped
```

1. **Load taxonomy** — Read `library/taxonomy.md` and `library/_template.md`. Fail if missing.

2. **Parse** — Extract candidates from the attached dump. Skip intro/outro chapters with no tool to file. One candidate = one product/model/system. Build an ordered list; do not write yet.

3. **Ask one** — For the next candidate only: name, suggested category, then a **decision cold read**. End with `keep? y/n`.
   - MUST cover, when known: what it is; performance or why it matters; price/cost; how you run it; one caveat.
   - Voice: **direct but readable**. Short complete sentences. Lead with facts. No filler, no hype, no “worth keeping if…” pep talk unless it adds a real decision fork.
   - MUST NOT: Fluffy openers (“people are excited”, “the reason it is getting attention”). MUST NOT: Datasheet bullet stacks as the default. MUST NOT: Paper jargon unless the user asks `explain`.
   - Target length: about 4–8 sentences, or two tight paragraphs.
   - MUST NOT: Ask about another candidate in the same turn.
   - If price/performance are missing and the item is more than spectacle, MAY quick-WebSearch those facts only; full enrichment after `y`.
   - On user `explain`, add one plain-language beat, then re-ask `keep? y/n`.

4. **On `n`** — Skip; no entry file. Continue with step 3 for the next candidate (or stop if the user ends intake).

5. **On `y`** — Dedup (`name:` Grep). WebSearch enrich (homepage, GitHub, license, run mode, `vram_hint` or `unknown`). Write `library/<category>/<slug>.md`. Refresh `library/INDEX.md`. Then continue with step 3 for the next candidate.

6. **Done** — After the last candidate (or user stop), summarize kept vs skipped. Do not re-ask batch approval.

## INDEX format

```markdown
# Library index

Regenerated by `ai-library-intake` after writes.

| name | title | category | status | tags | path |
|------|-------|----------|--------|------|------|
| voicemem | VoiceMem | agents-memory | watch | memory, open, local, audio | agents-memory/voicemem.md |
```

## Pre-completion checklist

- [ ] **Taxonomy loaded:** `library/taxonomy.md` was Read this run
      Method: Confirm Read in tool history before proposal
      Pass: Read occurred
      Fail: STOP, Read taxonomy, re-run proposal if needed
- [ ] **Keep gate held:** Each write followed a per-candidate `y`; no batch `y` for many keepers
      Method: Timeline: one ask → one reply → at most one new entry
      Pass: One-at-a-time
      Fail: STOP, revert batch if user requests, resume interactive
- [ ] **Tags legal:** Every tag is in taxonomy allowed list
      Method: Diff tags vs taxonomy
      Pass: All tags listed
      Fail: STOP, fix or ask user to extend taxonomy
- [ ] **No invented links:** Every URL came from WebSearch/WebFetch or was omitted
      Method: Spot-check `links:` block
      Pass: Traceable or omitted
      Fail: STOP, clear inventions
- [ ] **Dedup:** No silent overwrite of existing `name`
      Method: Grep before write
      Pass: Skip or user-approved update
      Fail: STOP, skip duplicate
- [ ] **INDEX fresh:** INDEX lists all category entries after writes
      Method: Count rows vs entry files
      Pass: Counts match
      Fail: STOP, regenerate INDEX
- [ ] **No Perplexity:** Enrichment used WebSearch/WebFetch only
      Method: Tool history
      Pass: No perplexity_* calls
      Fail: STOP, do not treat Perplexity output as filed metadata

## Additional resources

- Frontmatter examples and keep/drop CORRECT/PROHIBITED: [reference.md](reference.md)
- Taxonomy: [library/taxonomy.md](../../library/taxonomy.md)
- Template: [library/_template.md](../../library/_template.md)
