---
name: cursor-composer-2.5-fast
title: Cursor Composer 2.5 (Fast)
kind: service
category: coding
jobs:
  - coding
tags:
  - llm
  - closed
  - coding
  - cursor
  - agentic
  - fast
status: watch
lab: cursor
service: cursor
hosts:
  - id: cursor
    runtime_model_id: composer-2.5
    price_in: 3
    price_out: 15
    price_note: "Cursor Models pool; Fast (default in product) $3/$0.5/$15 in/cache/out per MTok; Standard same intelligence at $0.5/$0.2/$2.5; toggle Fast via picker Edit; Start plan non-fast only; exempt from Cursor Token Rate"
    limits_note: "Cursor Models pool; Fast $3/$15 per M in/out; Standard $0.5/$2.5; weekly pool shared with Grok"
    capacity: unknown
aa_intelligence: unknown
arena_elo: unknown
context_window: 200000
links:
  service_docs: https://cursor.com/docs/models/cursor-composer-2-5
source: "cursor picker seed 2026-09-09; scorecard refresh 2026-11-04"
refreshed: 2026-11-04
---

# Cursor Composer 2.5 (Fast)

## When to reach for it

Use this picker id when you want Cursor's first-party coding agent at **Fast** speed inside the Cursor Models pool. You cannot rehost this path; the card is **kind: service**.

## Signals

Cursor's pricing table lists Composer 2.5 (Fast) at $3 / $15 per M input/output and Standard at $0.50 / $2.50 with the same model weights. Context is 200k. AA and Arena are not filed on this card.

## My experience

(Operator notes only.)

## Caveats

**Fast** spends the weekly Cursor Models pool faster than Standard; turn Fast off in the picker when latency is not worth the cost. Plus plan uses fixed medium effort without Fast. Confirm the live table on models-and-pricing before large runs.

## Links

- Service docs: https://cursor.com/docs/models/cursor-composer-2-5
- Models and pricing: https://cursor.com/docs/models-and-pricing
- Announcement: https://cursor.com/blog/composer-2-5
