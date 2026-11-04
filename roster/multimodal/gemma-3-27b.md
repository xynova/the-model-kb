---
name: gemma-3-27b
title: Gemma 3 27B
kind: sovereign
category: multimodal
jobs:
  - vision
  - cheap
  - chat
tags:
  - llm
  - vision
  - open
  - deepinfra
  - cheap
status: watch
lab: google
preferred_host: deepinfra
hosts:
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: google/gemma-3-27b-it
    price_in: 0.08
    price_out: 0.16
    price_note: "DeepInfra published list; Priority 1.5x and Flex 0.8x tiers exist"
    capacity: hobby
    limits_note: "131,072 context; multimodal; max 30 images per request on DeepInfra"
    links:
      deepinfra: https://deepinfra.com/google/gemma-3-27b-it
aa_intelligence: unknown
arena_elo: unknown
context_window: 131072
links:
  models_dev: https://models.dev/
source: "polypus catalog diff 2026-11-04"
refreshed: 2026-11-04
---

# Gemma 3 27B

## When to reach for it

Pick this card when you want open Gemma 3 27B Instruct with vision on **deepinfra**. Treat it as a **baseline** next to Gemma 4 26B A4B, not a same-generation peer.

## Signals

DeepInfra list prices live in frontmatter. Context is 128k. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

Older Gemma 3 cut. Only **deepinfra** is filed. Prefer Gemma 4 cards when generation match matters.

## Links

- DeepInfra: https://deepinfra.com/google/gemma-3-27b-it
- models.dev: https://models.dev/
