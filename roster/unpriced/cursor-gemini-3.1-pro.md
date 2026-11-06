---
name: cursor-gemini-3.1-pro
title: Cursor Gemini 3.1 Pro
kind: service
category: unpriced
jobs:
- vision
- reasoning
tags:
- llm
- vision
- closed
- cursor
- agentic
status: watch
lab: google
service: cursor
hosts:
- id: cursor
  runtime_model_id: gemini-3.1-pro
  price_in: 2
  price_out: 12
  price_note: Other Models pool; $2/$0.2/$12 in/cache-read/out per MTok (no separate
    cache-write listed); long context >200k may raise rates (confirm table); Teams
    Cursor Token Rate $0.25/MTok; Start has no Other Models pool
  capacity: unknown
aa_intelligence: unknown
arena_elo: unknown
context_window: 1000000
links:
  service_docs: https://cursor.com/docs/models/gemini-3-1-pro
source: cursor picker seed 2026-09-09; scorecard refresh 2026-11-04
refreshed: 2026-11-04
aliases:
- /roster/multimodal/cursor-gemini-3.1-pro/
---


# Cursor Gemini 3.1 Pro

## When to reach for it

Use this picker id when you want Gemini 3.1 Pro in Cursor for image-in plus code or large-context Gemini work on the **Other Models** pool. The card is **kind: service**, not a direct Google API pin.

## Signals

Cursor's pricing table lists Gemini 3.1 Pro at $2 / $12 per M input/output and cache read $0.20. Context up to 1M. AA and Arena are not filed.

## My experience

(Operator notes only.)

## Caveats

**Start** has no Other Models pool. Teams and Enterprise add Cursor Token Rate. Prefer Gemini 3 Flash when cost and speed dominate and Pro-level vision is not required. Confirm long-context multipliers on the live table.

## Links

- Service docs: https://cursor.com/docs/models/gemini-3-1-pro
- Models and pricing: https://cursor.com/docs/models-and-pricing
