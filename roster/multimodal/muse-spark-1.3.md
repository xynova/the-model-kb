---
name: muse-spark-1.3
title: Muse Spark 1.3
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
    runtime_model_id: meta/muse-spark-1.3
    price_in: 1.25
    price_out: 4.25
    price_note: "OpenRouter Meta list $1.25/$4.25 per M; cache read $0.15"
    capacity: hobby
    limits_note: "1,048,576 context; Meta is the sole OpenRouter provider on fetch"
    links:
      openrouter: https://openrouter.ai/meta/muse-spark-1.3
aa_intelligence: 37
arena_elo: unknown
context_window: 1048576
links:
  models_dev: https://models.dev/
source: "polypus catalog diff 2026-11-04"
refreshed: 2026-11-04
---

# Muse Spark 1.3

## When to reach for it

Pick this card when you want Meta Muse Spark 1.3 through **openrouter** for multimodal agent and coding work with a 1M context window. Access is lab-locked (`kind: provider`).

## Signals

OpenRouter lists $1.25 / $4.25 per M. Artificial Analysis Intelligence Index 36.6 (Max) is filed here as 37. Context is 1M.

## My experience

(Operator notes only.)

## Caveats

This is not Muse Spark Open. Do not treat OpenRouter Meta as a relocatable host. Confirm the live table before large runs.

## Links

- OpenRouter: https://openrouter.ai/meta/muse-spark-1.3
- models.dev: https://models.dev/
