---
name: cursor-gemini-3.1-pro
title: Cursor Gemini 3.1 Pro
kind: service
category: multimodal
tags:
  - llm
  - vision
  - closed
  - cursor
  - agentic
status: watch
lab: google
service: cursor
runtime_model_id: gemini-3.1-pro
price_in: 2
price_out: 12
price_note: "Other Models pool; $2/$0.2/$12 in/cache-read/out per MTok (no separate cache-write listed); long context >200k may raise rates (confirm table); Teams Cursor Token Rate $0.25/MTok; Start has no Other Models pool"
aa_intelligence: unknown
arena_elo: unknown
context_window: 1000000
links:
  service_docs: https://cursor.com/docs/models/gemini-3-1-pro
source: "cursor picker seed 2026-09-09"
refreshed: 2026-09-10
---

# Cursor Gemini 3.1 Pro

## When to reach for it

Google Gemini 3.1 Pro in Cursor when you want image-in plus code (UI/UX from mockups, visual understanding) or large-context Gemini, and you accept Other Models pool spend. Still `kind: service` (Cursor picker), not a direct Google API pin.

## Signals

Cursor docs (2026-09-10): $2 / $12 per M input/output; cache read $0.2. Up to 1M context. Positioned for multimodal UI/UX and whole-codebase analysis. AA / Arena not filed.

## My experience

(Not used yet.)

## Caveats

Other Models pool only (absent on Start). Teams/Enterprise add Cursor Token Rate. Prefer Flash when cost/speed dominate and Pro-level vision is not required. Confirm long-context multipliers on models-and-pricing.

## Links

- Service docs: https://cursor.com/docs/models/gemini-3-1-pro
- Models and pricing: https://cursor.com/docs/models-and-pricing
