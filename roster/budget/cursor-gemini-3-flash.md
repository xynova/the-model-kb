---
name: cursor-gemini-3-flash
title: Cursor Gemini 3 Flash
kind: service
category: budget
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
runtime_model_id: gemini-3-flash
price_in: 0.5
price_out: 3
price_note: "Other Models pool; $0.5/$0.05/$3 in/cache-read/out per MTok; Teams Cursor Token Rate $0.25/MTok; Start has no Other Models pool; hidden by default in some Cursor tables"
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  service_docs: https://cursor.com/docs/models/gemini-3-flash
source: "cursor picker seed 2026-09-09"
refreshed: 2026-09-10
---

# Cursor Gemini 3 Flash

## When to reach for it

Google Gemini 3 Flash in Cursor for simple coding, quick edits, and executing plans from a stronger model when speed and cost dominate. Still `kind: service`. Treat as a **baseline** next to Gemini 3.1 Pro (not the same capability class).

## Signals

Cursor docs (2026-09-10): $0.50 / $3.00 per M input/output; cache read $0.05. Positioned as speed-optimized and among the cheapest Gemini options in Cursor. AA / Arena not filed.

## My experience

(Not used yet.)

## Caveats

Not for peak multimodal or hard reasoning; use Gemini 3.1 Pro (or newer Flash cuts if you enable them) when that matters. Other Models pool only. Confirm current table on models-and-pricing.

## Links

- Service docs: https://cursor.com/docs/models/gemini-3-flash
- Models and pricing: https://cursor.com/docs/models-and-pricing
