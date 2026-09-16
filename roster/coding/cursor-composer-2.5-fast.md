---
name: cursor-composer-2.5-fast
title: Cursor Composer 2.5 (Fast)
kind: service
category: coding
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
runtime_model_id: composer-2.5
price_in: 3
price_out: 15
price_note: "Cursor Models pool; Fast (default in product) $3/$0.5/$15 in/cache/out per MTok; Standard same intelligence at $0.5/$0.2/$2.5; toggle Fast via picker Edit; Start plan non-fast only; exempt from Cursor Token Rate"
aa_intelligence: unknown
arena_elo: unknown
context_window: 200000
links:
  service_docs: https://cursor.com/docs/models/cursor-composer-2-5
source: "cursor picker seed 2026-09-09"
refreshed: 2026-09-09
---

# Cursor Composer 2.5 (Fast)

## When to reach for it

Cursor's first-party agentic coding model at Fast (product default) for interactive agent-on-code when you want the Cursor Models pool instead of Grok or a third-party model. Access is `kind: service`; you cannot rehost this picker path.

## Signals

Cursor docs (2026-09-09): Fast $3 / $15 per M input/output ($0.50 cached input); Standard $0.50 / $2.50 ($0.20 cached). Same intelligence; Fast is speed/throughput. Context 200k. Draws from Cursor Models pool with Grok 4.6 / 4.5. Built on Moonshot Kimi K2.5 checkpoint per Cursor blog. AA / Arena not filed.

## My experience

(Not used yet.)

## Caveats

Fast burns the pool much faster than Standard; turn Fast off via picker Edit when latency is not the goal. Start (India) locks non-fast. Prefer Grok 4.6 High/xhigh for longer general knowledge work per Cursor help positioning. Confirm current pricing on models-and-pricing.

## Links

- Service docs: https://cursor.com/docs/models/cursor-composer-2-5
- Models and pricing: https://cursor.com/docs/models-and-pricing
- Announcement: https://cursor.com/blog/composer-2-5
