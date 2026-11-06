---
name: cursor-grok-4.6-high
title: Cursor Grok 4.6 (High)
kind: service
category: unpriced
jobs:
- reasoning
- chat
tags:
- llm
- closed
- cursor
- agentic
- reasoning
status: watch
lab: xai
service: cursor
hosts:
- id: cursor
  runtime_model_id: grok-4.6
  price_in: 2
  price_out: 6
  price_note: Cursor Models pool (with Grok 4.5, Composer 2.5); Standard $2/$0.50/$6
    in/cache/out per MTok; Fast $4/$1/$12; High is the named-model default effort
    (also low/medium/xhigh); Start plan locks medium non-fast
  capacity: unknown
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  service_docs: https://cursor.com/docs/models/grok-4-6
source: cursor picker seed 2026-09-09; scorecard refresh 2026-11-04
refreshed: 2026-11-04
aliases:
- /roster/general/cursor-grok-4.6-high/
---


# Cursor Grok 4.6 (High)

## When to reach for it

Use this picker id when you want Cursor's first-party Grok 4.6 lane at **High** effort for long-horizon agent and knowledge work instead of Composer or a third-party model. You cannot rehost this path; the card is **kind: service**.

## Signals

Cursor's pricing table lists Grok 4.6 Standard at $2 / $6 per M input/output and cache read $0.50; Fast is $4 / $12. Usage draws from the weekly **Cursor Models** pool shared with Composer and other Grok cuts. **High** is the named-model default effort; low, medium, high, and xhigh exist. AA and Arena are not filed.

## My experience

(Operator notes only.)

## Caveats

Plan caps apply; **Start** fixes medium effort without Fast. **Fast** is a speed tier, not the same knob as **High** effort. This card is distinct from Grok Bot. Confirm the live table before large on-demand runs.

## Links

- Service docs: https://cursor.com/docs/models/grok-4-6
- Models and pricing: https://cursor.com/docs/models-and-pricing
