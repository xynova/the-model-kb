---
name: claude-haiku-5.5
title: Claude Haiku 5.5
kind: provider
category: budget
jobs:
  - cheap
  - coding
  - chat
tags:
  - llm
  - closed
  - api
  - cheap
  - coding
status: watch
lab: anthropic
preferred_host: openrouter
hosts:
  - id: openrouter
    runtime: openrouter
    runtime_model_id: anthropic/claude-haiku-5.5
    price_in: 0.10
    price_out: 0.50
    price_note: "OpenRouter headline; Anthropic and Vertex rows $0.10/$0.50 per M; cache read $0.01"
    capacity: hobby
    limits_note: "OpenRouter proxies Anthropic; not a sovereign host"
    links:
      openrouter: https://openrouter.ai/anthropic/claude-haiku-5.5
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  models_dev: https://models.dev/
source: "polypus catalog diff 2026-11-04"
refreshed: 2026-11-04
---

# Claude Haiku 5.5

## When to reach for it

Pick this card when you want Anthropic Claude Haiku 5.5 through **openrouter** as a cheap coding and chat API. Weights are lab-locked (`kind: provider`).

## Signals

OpenRouter lists $0.10 / $0.50 per M on Anthropic and Vertex rows. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

There is no Anthropic host id in this catalog. File other proxies only after `providers-mint` if you need a new company.

## Links

- OpenRouter: https://openrouter.ai/anthropic/claude-haiku-5.5
- models.dev: https://models.dev/
