---
name: gemini-2.5-flash
title: Gemini 2.5 Flash
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
    runtime_model_id: gemini-2.5-flash
    price_in: 0.3
    price_out: 2.5
    price_note: OpenRouter Google AI Studio list $0.30/$2.50 per M; cache read $0.03; confirm ai.google.dev
    capacity: unknown
    limits_note: Polypus gemini_paid/gemini-2.5-flash
    links:
      provider_pricing: https://ai.google.dev/gemini-api/docs/pricing
  - id: openrouter
    runtime: openrouter
    runtime_model_id: google/gemini-2.5-flash
    price_in: 0.1
    price_out: 19.0
    price_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    capacity: hobby
    limits_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    links:
      openrouter: https://openrouter.ai/google/gemini-2.5-flash
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-13"
aliases:
  - /roster/budget/gemini-2.5-flash/
---


# Gemini 2.5 Flash

## When to reach for it

Pick this card when you want Google Gemini 2.5 Flash as a cheap multimodal API. Treat it as a **baseline** next to Gemini 3 Flash on this shelf.

## Signals

OpenRouter lists $0.30 / $2.50 per M on Google AI Studio. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

Older Gemini 2.5 cut. Weights are not relocatable.

## Links

- Gemini API pricing: https://ai.google.dev/gemini-api/docs/pricing
- OpenRouter: https://openrouter.ai/google/gemini-2.5-flash
- models.dev: https://models.dev/
