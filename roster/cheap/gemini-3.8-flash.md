---
name: gemini-3.8-flash
title: Gemini 3.8 Flash
kind: provider
category: cheap
jobs:
  - cheap
  - vision
  - chat
  - reasoning
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
    runtime_model_id: gemini-3.8-flash
    price_in: 0.75
    price_out: 3.75
    price_note: OpenRouter Google AI Studio headline $0.75/$3.75 per M (50% off noted on fetch); confirm ai.google.dev
    capacity: unknown
    limits_note: Polypus gemini_paid route; 1M context on OpenRouter listing
    links:
      provider_pricing: https://ai.google.dev/gemini-api/docs/pricing
  - id: openrouter
    runtime: openrouter
    runtime_model_id: google/gemini-3.8-flash
    price_in: 0.075
    price_out: 19.0
    price_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    capacity: hobby
    limits_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    links:
      openrouter: https://openrouter.ai/google/gemini-3.8-flash
aa_intelligence: 30.8
arena_elo: unknown
context_window: 1000000
links:
  models_dev: https://models.dev/
  artificial_analysis: https://artificialanalysis.ai/models/gemini-3-8-flash
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-13"
aliases:
  - /roster/budget/gemini-3.8-flash/
---


# Gemini 3.8 Flash

## When to reach for it

Pick this card when you want Google Gemini 3.8 Flash as a **lab-locked API** (`kind: provider`) for cheap multimodal chat. This is not a Cursor picker card.

## Signals

OpenRouter lists $0.75 / $3.75 per M as the headline at fetch, with a 50 percent off note. Context is 1M on that listing. Artificial Analysis Intelligence Index is 30.8 for Gemini 3.8 Flash (High). The public leaderboard rounds that to 31. Confirm Google’s own table before large runs.

## My experience

(Operator notes only.)

## Caveats

Weights are not relocatable. Cursor Gemini Flash cards stay `kind: service`. Promotional OpenRouter headlines can expire.

## Links

- Gemini API pricing: https://ai.google.dev/gemini-api/docs/pricing
- OpenRouter: https://openrouter.ai/google/gemini-3.8-flash
- models.dev: https://models.dev/
