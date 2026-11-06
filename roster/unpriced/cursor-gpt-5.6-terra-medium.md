---
name: cursor-gpt-5.6-terra-medium
title: Cursor GPT-5.6 Terra (Medium)
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
lab: openai
service: cursor
hosts:
- id: cursor
  runtime_model_id: gpt-5.6-terra
  price_in: 2
  price_out: 12
  price_note: Other Models pool; $2/$2.5/$0.2/$12 in/cache-write/cache-read/out per
    MTok; Fast 2x; long context >272k input 2x and output 1.5x; Teams Cursor Token
    Rate $0.25/MTok; Medium is picker reasoning effort; Max Mode on legacy request
    plans
  capacity: unknown
aa_intelligence: unknown
arena_elo: unknown
context_window: 1000000
links:
  service_docs: https://cursor.com/docs/models/gpt-5-6-terra
source: cursor picker seed 2026-09-09; scorecard refresh 2026-11-04
refreshed: 2026-11-04
aliases:
- /roster/general/cursor-gpt-5.6-terra-medium/
---


# Cursor GPT-5.6 Terra (Medium)

## When to reach for it

Use this picker id when you want the mid GPT-5.6 tier at **Medium** effort for everyday agent work on the **Other Models** pool instead of first-party Grok or Composer. The card is **kind: service**, not a direct OpenAI API pin.

## Signals

Cursor's pricing table lists GPT-5.6 Terra at $2 / $12 per M input/output, cache write $2.50, cache read $0.20. Fast mode is 2x. Long context above 272k uses higher multipliers; context up to 1M. AA and Arena are not filed.

## My experience

(Operator notes only.)

## Caveats

**Start** has no Other Models pool. Teams and Enterprise add Cursor Token Rate. Use Sol for harder long runs and Luna when cost or speed dominates. Confirm the live table before large runs.

## Links

- Service docs: https://cursor.com/docs/models/gpt-5-6-terra
- Models and pricing: https://cursor.com/docs/models-and-pricing
