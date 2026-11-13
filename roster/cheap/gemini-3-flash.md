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
    price_out: 1.0
    price_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    capacity: hobby
    limits_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    links:
      openrouter: https://openrouter.ai/google/gemini-3-flash-preview
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-13"
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
