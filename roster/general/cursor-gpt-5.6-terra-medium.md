---
name: cursor-gpt-5.6-terra-medium
title: Cursor GPT-5.6 Terra (Medium)
kind: service
category: general
tags:
  - llm
  - closed
  - cursor
  - agentic
  - reasoning
status: watch
lab: openai
service: cursor
runtime_model_id: gpt-5.6-terra
price_in: 2
price_out: 12
price_note: "Other Models pool; $2/$2.5/$0.2/$12 in/cache-write/cache-read/out per MTok; Fast 2x; long context >272k input 2x and output 1.5x; Teams Cursor Token Rate $0.25/MTok; Medium is picker reasoning effort; Max Mode on legacy request plans"
aa_intelligence: unknown
arena_elo: unknown
context_window: 1000000
links:
  service_docs: https://cursor.com/docs/models/gpt-5-6-terra
source: "cursor picker seed 2026-09-09"
refreshed: 2026-09-09
---

# Cursor GPT-5.6 Terra (Medium)

## When to reach for it

OpenAI GPT-5.6 Terra at Medium effort in Cursor when you want the mid GPT-5.6 tier (between Luna and Sol) for everyday agent work and you are willing to spend the Other Models pool instead of first-party Grok/Composer. Still `kind: service` (Cursor picker), not a direct OpenAI API pin.

## Signals

Cursor docs (2026-09-09): $2 / $12 per M input/output; cache write $2.5, cache read $0.2. Fast mode 2x. Long context above 272k: input 2x, output 1.5x; up to 1M tokens. Mid-tier between Sol and Luna per Cursor model page. AA / Arena not filed.

## My experience

(Not used yet.)

## Caveats

Draws from Other Models (Start plan has no Other Models pool). Teams/Enterprise add Cursor Token Rate on third-party. Not the top of the GPT-5.6 family; use Sol for harder long runs, Luna when cost/speed dominate. Confirm current table on models-and-pricing.

## Links

- Service docs: https://cursor.com/docs/models/gpt-5-6-terra
- Models and pricing: https://cursor.com/docs/models-and-pricing
