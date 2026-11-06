---
name: gemini-3-flash
title: Gemini 3 Flash
kind: provider
category: cheap
jobs:
  - cheap
  - vision
  - chat
tags:
  - llm
  - vision
  - closed
  - api
  - cheap
  - fast
status: watch
lab: google
preferred_host: google
hosts:
  - id: google
    runtime: google
    runtime_model_id: gemini-3-flash-preview
    price_in: 0.5
    price_out: 3
    price_note: OpenRouter Google AI Studio list $0.50/$3.00 per M; cache read $0.05; confirm ai.google.dev
    capacity: unknown
    limits_note: Polypus gemini_paid/gemini-3-flash-preview; preview id may change
    links:
      provider_pricing: https://ai.google.dev/gemini-api/docs/pricing
  - id: openrouter
    runtime: openrouter
    runtime_model_id: google/gemini-3-flash-preview
    price_in: 0.1
    price_out: 19.0
    price_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    capacity: hobby
    limits_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    links:
      openrouter: https://openrouter.ai/google/gemini-3-flash-preview
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-06"
aliases:
  - /roster/budget/gemini-3-flash/
---


# Gemini 3 Flash

## When to reach for it

Pick this card when you want Google Gemini 3 Flash through the **Gemini API or OpenRouter**, not through Cursor. The Cursor picker lives on `cursor-gemini-3-flash` (`kind: service`).

## Signals

OpenRouter lists $0.50 / $3.00 per M on Google AI Studio and Vertex rows. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

The Polypus id still uses a **preview** suffix. Do not merge this card with the Cursor service card.

## Links

- Gemini API pricing: https://ai.google.dev/gemini-api/docs/pricing
- OpenRouter: https://openrouter.ai/google/gemini-3-flash-preview
- models.dev: https://models.dev/
