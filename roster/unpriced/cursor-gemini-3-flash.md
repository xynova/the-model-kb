---
name: cursor-gemini-3-flash
title: Cursor Gemini 3 Flash
kind: service
category: unpriced
jobs:
- cheap
- chat
tags:
- llm
- closed
- cursor
- cheap
- fast
- agentic
status: watch
lab: google
service: cursor
hosts:
- id: cursor
  runtime_model_id: gemini-3-flash
  price_in: 0.5
  price_out: 3
  price_note: Other Models pool; $0.5/$0.05/$3 in/cache-read/out per MTok; Teams Cursor
    Token Rate $0.25/MTok; Start has no Other Models pool; hidden by default in some
    Cursor tables
  capacity: unknown
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  service_docs: https://cursor.com/docs/models/gemini-3-flash
source: cursor picker seed 2026-09-09; scorecard refresh 2026-11-04
refreshed: 2026-11-04
aliases:
- /roster/budget/cursor-gemini-3-flash/
---


# Cursor Gemini 3 Flash

## When to reach for it

Use this picker id when you want a fast, cheap Gemini lane in Cursor for simple coding and quick edits. Treat it as a **baseline** next to Gemini 3.1 Pro, not the same capability class. The card is **kind: service**.

## Signals

Cursor's pricing table lists Gemini 3 Flash at $0.50 / $3 per M input/output and cache read $0.05. AA and Arena are not filed.

## My experience

(Operator notes only.)

## Caveats

**Start** has no Other Models pool. Step up to Gemini 3.1 Pro when vision or reasoning depth matters. Confirm the live table before large runs.

## Links

- Service docs: https://cursor.com/docs/models/gemini-3-flash
- Models and pricing: https://cursor.com/docs/models-and-pricing
