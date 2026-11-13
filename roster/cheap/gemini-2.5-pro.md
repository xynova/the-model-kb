---
name: gemini-2.5-pro
title: Gemini 2.5 Pro
kind: provider
category: cheap
jobs:
  - reasoning
  - vision
  - chat
tags:
  - llm
  - vision
  - closed
  - api
status: watch
lab: google
preferred_host: google
hosts:
  - id: google
    runtime: google
    runtime_model_id: gemini-2.5-pro
    price_in: 1.25
    price_out: 10
    price_note: OpenRouter headline $1.25/$10 per M; confirm ai.google.dev
    capacity: unknown
    limits_note: Polypus gemini_paid/gemini-2.5-pro; 1M context on OpenRouter listing
    links:
      provider_pricing: https://ai.google.dev/gemini-api/docs/pricing
  - id: openrouter
    runtime: openrouter
    runtime_model_id: google/gemini-2.5-pro
    price_in: 0.125
    price_out: 19.0
    price_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    capacity: hobby
    limits_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    links:
      openrouter: https://openrouter.ai/google/gemini-2.5-pro
aa_intelligence: unknown
arena_elo: unknown
context_window: 1000000
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-13"
aliases:
  - /roster/general/gemini-2.5-pro/
---


# Gemini 2.5 Pro

## When to reach for it

Pick this card when you want Google Gemini 2.5 Pro as a lab-locked API. Treat it as a **baseline** next to Gemini 3.1 Pro on this shelf.

## Signals

OpenRouter lists $1.25 / $10 per M. Context is 1M on that listing. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

Older Gemini 2.5 cut. Weights are not relocatable.

## Links

- Gemini API pricing: https://ai.google.dev/gemini-api/docs/pricing
- OpenRouter: https://openrouter.ai/google/gemini-2.5-pro
- models.dev: https://models.dev/
