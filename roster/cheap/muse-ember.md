---
name: muse-ember
title: Muse Ember
kind: provider
category: cheap
jobs:
- cheap
- coding
- vision
- reasoning
- long-context
tags:
- llm
- vision
- closed
- api
- cheap
- coding
- agentic
- reasoning
status: watch
lab: meta
preferred_host: openrouter
hosts:
- id: openrouter
  runtime: openrouter
  runtime_model_id: meta/muse-ember
  price_in: 0.1
  price_out: 0.5
  price_note: OpenRouter Meta and Google Vertex list $0.10/$0.50 per M; cache read
    $0.004; Meta list above 200K input tokens is $0.30/$1.50 per M
  capacity: hobby
  limits_note: 1,048,576 context; OpenRouter routes Meta and Google Vertex; Standard
    mode only on Meta Model API
  links:
    openrouter: https://openrouter.ai/meta/muse-ember
aa_intelligence: 36
arena_elo: unknown
context_window: 1048576
links:
  models_dev: https://models.dev/
  artificial_analysis: https://artificialanalysis.ai/models/muse-ember
source: youtube:woWo80dPkVk
refreshed: '2026-11-17'
---

# Muse Ember

## When to reach for it

Pick this card when you want Meta Muse Ember on **openrouter** for high-volume agent steps, coding subtasks, and multimodal document work at Haiku-class list price. Pair it with Muse Spark or Muse Wonder when a larger planner should hand off volume work. Access is lab-locked (`kind: provider`).

## Signals

OpenRouter lists $0.10 / $0.50 per M on Meta and Google Vertex rows. Artificial Analysis Intelligence Index is 36 for Muse Ember (Max). Context is 1M.

## My experience

(Operator notes only.)

## Caveats

Meta bills the whole request at the higher long-context rates once input exceeds 200K tokens. OpenRouter can route to Vertex; pin a provider if you need one lab. This is not Muse Spark Open or Muse Glimmer.

## Links

- OpenRouter: https://openrouter.ai/meta/muse-ember
- Artificial Analysis: https://artificialanalysis.ai/models/muse-ember
- models.dev: https://models.dev/