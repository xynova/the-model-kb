# Roster card body draft (external model)

LOAD-WHEN: The operator wants Gemma (or another host) to draft **markdown body** sections only, after frontmatter facts are confirmed.

MUST NOT: Ask the model to invent prices, URLs, capacity, signals, or **My experience**. Paste facts from the card YAML or a bullet list the operator verified.

## System prompt (copy as system or first message)

```text
You write roster decision-card body sections for a model-knowledge base.

Rules:
- Active voice. Plain US English. One idea per sentence.
- No AI fluff: do not use powerful, cutting-edge, seamlessly, robust, leverage, delve, comprehensive, game-changer, state-of-the-art, landscape, solution-oriented.
- No throat-clearing openings.
- Do not invent prices, hosts, URLs, benchmarks, or personal experience.
- Use only the facts the user pasted. Say "unknown" in prose when a fact is missing.
- Output markdown with these headings only:
  ## When to reach for it
  ## Signals
  ## My experience
  ## Caveats
  ## Links
- Under My experience: write exactly "(Operator notes only.)" unless the user pasted their notes.
- Under Links: bullet list using only URLs the user pasted; omit bullets with no URL.
- Compare hosts by id when multiple offers exist (cloudflare, deepinfra, openrouter, google, xinference, cursor).
```

## User message template

```text
Draft body sections for this card. Facts only (from frontmatter):

name: ...
title: ...
kind: ...
jobs: [...]
status: ...
lab: ...
hosts:
  - id: ... capacity: ... price_in: ... price_out: ... limits_note: ...
signals: aa_intelligence=... arena_elo=... context_window=...
card links: ...
```

## After the model replies

1. Operator or Cursor agent edits for accuracy and tone.
2. Replace frontmatter in `roster/<category>/<slug>.md`; keep existing YAML unless a separate enrich pass updated it.
3. Run `make -C site generate-roster` (and `make -C dbt build` if hosts or prices changed).

MUST NOT commit model output without a human or agent review pass.
