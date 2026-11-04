---
name: cursor-gpt-5.4-mini-medium
title: Cursor GPT-5.4 Mini (Medium)
kind: service
category: budget
jobs:
  - cheap
  - vision
tags:
  - llm
  - vision
  - closed
  - cursor
  - cheap
  - fast
  - agentic
status: watch
lab: openai
service: cursor
hosts:
  - id: cursor
    runtime_model_id: gpt-5.4-mini
    price_in: 0.75
    price_out: 4.5
    price_note: "Other Models pool; $0.75/$0.075/$4.5 in/cache-read/out per MTok; 90% discount on cached input; Medium is picker reasoning effort; Teams Cursor Token Rate $0.25/MTok; Start has no Other Models pool; hidden by default in some Cursor tables"
    capacity: unknown
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  service_docs: https://cursor.com/docs/models/gpt-5-4-mini
source: "cursor picker seed 2026-09-09; scorecard refresh 2026-11-04"
refreshed: 2026-11-04
---

# Cursor GPT-5.4 Mini (Medium)

## When to reach for it

Use this picker id when you want GPT-5.4 Mini at **Medium** effort for cost-sensitive agent loops with multimodal input. Treat it as a **baseline** next to GPT-5.6 Luna, not a same-generation peer. The card is **kind: service**.

## Signals

Cursor's pricing table lists GPT-5.4 Mini at $0.75 / $4.50 per M input/output and cache read $0.075 with a 90% cached-input discount noted. AA and Arena are not filed.

## My experience

(Operator notes only.)

## Caveats

**Start** has no Other Models pool. Step up to GPT-5.4 or GPT-5.6 Terra/Sol when peak intelligence matters. Prefer Luna when you want the newer cheap GPT-5.6 lane. Confirm the live table before large runs.

## Links

- Service docs: https://cursor.com/docs/models/gpt-5-4-mini
- Models and pricing: https://cursor.com/docs/models-and-pricing
