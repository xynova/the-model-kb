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
  price_note: OpenRouter Google AI Studio list $0.30/$2.50 per M; cache read $0.03;
    confirm ai.google.dev
  capacity: unknown
  limits_note: Polypus gemini_paid/gemini-2.5-flash
  links:
    provider_pricing: https://ai.google.dev/gemini-api/docs/pricing
- id: openrouter
  runtime: openrouter
  runtime_model_id: google/gemini-2.5-flash
  price_in: 0.3
  price_out: 2.5
  price_note: OpenRouter headline; Google AI Studio row $0.30/$2.50
  capacity: hobby
  limits_note: Marketplace routing
  links:
    openrouter: https://openrouter.ai/google/gemini-2.5-flash
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: 2026-11-04
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
