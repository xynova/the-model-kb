---
name: granite-4.2-8b
title: Granite 4.2 8B
kind: sovereign
category: cheap
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
    price_in: 1.0
    price_out: 0.9
    price_note: '{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-670439{display:-webkit-inline-box;display:-webkit-inline-'
    capacity: hobby
    limits_note: '{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-670439{display:-webkit-inline-box;display:-webkit-inline-'
    links:
      deepinfra: https://deepinfra.com/ibm-granite/granite-4.2-8b
aa_intelligence: 5.1
arena_elo: unknown
context_window: 131072
links:
  models_dev: https://models.dev/
  artificial_analysis: https://artificialanalysis.ai/models/granite-4-2-8b
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-06"
aliases:
  - /roster/general/granite-4.2-8b/
---


# Granite 4.2 8B

## When to reach for it

Pick this card when you want an 8B IBM Granite 4.2 reasoning cut on **deepinfra** for cheap chat, coding, or tool use. It is a smaller dense peer in the Granite 4.2 family, not a Gemma 4 26B MoE peer.

## Signals

DeepInfra list prices live in frontmatter. Context is 128k on the model page. Artificial Analysis Intelligence Index is 5.1. The public leaderboard rounds that to 5.

## My experience

(Operator notes only.)

## Caveats

Only **deepinfra** is filed. File industrial capacity only through `capacity-decisions.yaml`.

## Links

- DeepInfra: https://deepinfra.com/ibm-granite/granite-4.2-8b
- models.dev: https://models.dev/
