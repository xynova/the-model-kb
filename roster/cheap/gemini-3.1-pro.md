---
name: gemini-3.1-pro
title: Gemini 3.1 Pro
kind: provider
category: cheap
jobs:
  - vision
  - reasoning
  - long-context
tags:
  - llm
  - vision
  - closed
  - api
  - agentic
status: watch
lab: google
preferred_host: google
hosts:
  - id: google
    runtime: google
    runtime_model_id: gemini-3.1-pro-preview
    price_in: 2
    price_out: 12
    price_note: OpenRouter headline $2/$12 per M; confirm ai.google.dev
    capacity: unknown
    limits_note: Polypus gemini_paid/gemini-3.1-pro-preview; 1M context on OpenRouter listing
    links:
      provider_pricing: https://ai.google.dev/gemini-api/docs/pricing
  - id: openrouter
    runtime: openrouter
    runtime_model_id: google/gemini-3.1-pro-preview
    price_in: 0.2
    price_out: 19.0
    price_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    capacity: hobby
    limits_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    links:
      openrouter: https://openrouter.ai/google/gemini-3.1-pro-preview
aa_intelligence: unknown
arena_elo: unknown
context_window: 1000000
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-06"
aliases:
  - /roster/multimodal/gemini-3.1-pro/
---


# Gemini 3.1 Pro

## When to reach for it

Pick this card when you want Google Gemini 3.1 Pro through the **Gemini API or OpenRouter** for vision and long-context reasoning. The Cursor picker lives on `cursor-gemini-3.1-pro` (`kind: service`).

## Signals

OpenRouter lists $2 / $12 per M. Context is 1M on that listing. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

The Polypus id still uses a **preview** suffix. Do not merge this card with the Cursor service card.

## Links

- Gemini API pricing: https://ai.google.dev/gemini-api/docs/pricing
- OpenRouter: https://openrouter.ai/google/gemini-3.1-pro-preview
- models.dev: https://models.dev/
