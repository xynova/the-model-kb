---
name: granite-4.2-8b
title: Granite 4.2 8B
kind: sovereign
category: general
jobs:
  - reasoning
  - cheap
  - chat
  - coding
tags:
  - llm
  - open
  - reasoning
  - deepinfra
  - cheap
  - agentic
status: watch
lab: ibm
preferred_host: deepinfra
hosts:
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: ibm-granite/granite-4.2-8b
    price_in: 0.06
    price_out: 0.25
    price_note: "DeepInfra published list; cached input $0.015/M"
    capacity: hobby
    limits_note: "131,072 native context; Apache 2.0; thinking modes on the model page"
    links:
      deepinfra: https://deepinfra.com/ibm-granite/granite-4.2-8b
aa_intelligence: unknown
arena_elo: unknown
context_window: 131072
links:
  models_dev: https://models.dev/
source: "polypus catalog diff 2026-11-04"
refreshed: 2026-11-04
---

# Granite 4.2 8B

## When to reach for it

Pick this card when you want an 8B IBM Granite 4.2 reasoning cut on **deepinfra** for cheap chat, coding, or tool use. It is a smaller dense peer in the Granite 4.2 family, not a Gemma 4 26B MoE peer.

## Signals

DeepInfra list prices live in frontmatter. Context is 128k on the model page. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

Only **deepinfra** is filed. File industrial capacity only through `capacity-decisions.yaml`.

## Links

- DeepInfra: https://deepinfra.com/ibm-granite/granite-4.2-8b
- models.dev: https://models.dev/
