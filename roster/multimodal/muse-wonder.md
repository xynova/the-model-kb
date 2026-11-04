---
name: muse-wonder
title: Muse Wonder
kind: provider
category: multimodal
jobs:
  - vision
  - reasoning
  - coding
  - long-context
tags:
  - llm
  - vision
  - closed
  - api
  - agentic
  - reasoning
status: watch
lab: meta
preferred_host: openrouter
hosts:
  - id: openrouter
    runtime: openrouter
    runtime_model_id: meta/muse-wonder
    price_in: 4
    price_out: 20
    price_note: "OpenRouter Meta list $4/$20 per M; cache read $0.20"
    capacity: hobby
    limits_note: "1,048,576 context; Meta is the sole OpenRouter provider on fetch"
    links:
      openrouter: https://openrouter.ai/meta/muse-wonder
aa_intelligence: 44
arena_elo: unknown
context_window: 1048576
links:
  models_dev: https://models.dev/
source: "polypus catalog diff 2026-11-04"
refreshed: 2026-11-04
---

# Muse Wonder

## When to reach for it

Pick this card when you want Meta Muse Wonder through **openrouter** for long-horizon multimodal agent work. Access is lab-locked (`kind: provider`). Use Muse Spark 1.3 on this shelf when list price matters more than this cut.

## Signals

OpenRouter lists $4 / $20 per M. Artificial Analysis Intelligence Index 43.6 (Max) is filed here as 44. Context is 1M.

## My experience

(Operator notes only.)

## Caveats

This is not an open-weights Muse Glimmer card. Confirm the live table before large runs.

## Links

- OpenRouter: https://openrouter.ai/meta/muse-wonder
- models.dev: https://models.dev/
