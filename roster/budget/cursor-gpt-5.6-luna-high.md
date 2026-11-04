---
name: cursor-gpt-5.6-luna-high
title: Cursor GPT-5.6 Luna (High)
kind: service
category: budget
jobs:
  - cheap
  - reasoning
tags:
  - llm
  - closed
  - cursor
  - cheap
  - fast
  - agentic
  - reasoning
status: watch
lab: openai
service: cursor
hosts:
  - id: cursor
    runtime_model_id: gpt-5.6-luna
    price_in: 0.2
    price_out: 1.2
    price_note: "Other Models pool; $0.2/$0.25/$0.02/$1.2 in/cache-write/cache-read/out per MTok; Fast 2x; long context >272k input 2x and output 1.5x; Teams Cursor Token Rate $0.25/MTok; High is picker reasoning effort on the cheap base; Max Mode on legacy request plans"
    capacity: unknown
aa_intelligence: unknown
arena_elo: unknown
context_window: 1000000
links:
  service_docs: https://cursor.com/docs/models/gpt-5-6-luna
source: "cursor picker seed 2026-09-09; scorecard refresh 2026-11-04"
refreshed: 2026-11-04
---

# Cursor GPT-5.6 Luna (High)

## When to reach for it

Use this picker id when you want the smallest GPT-5.6 tier at **High** effort for high-volume, cost-first agent loops on the **Other Models** pool. The card is **kind: service**.

## Signals

Cursor's pricing table lists GPT-5.6 Luna at $0.20 / $1.20 per M input/output, cache write $0.25, cache read $0.02. Fast mode is 2x. Long context above 272k uses higher multipliers; context up to 1M. AA and Arena are not filed.

## My experience

(Operator notes only.)

## Caveats

**Start** has no Other Models pool. Teams and Enterprise add Cursor Token Rate. **High** raises compute on a cheap base; **Fast** is a separate multiplier. Step up to Terra or Sol when peak intelligence matters. Confirm the live table before large runs.

## Links

- Service docs: https://cursor.com/docs/models/gpt-5-6-luna
- Models and pricing: https://cursor.com/docs/models-and-pricing
