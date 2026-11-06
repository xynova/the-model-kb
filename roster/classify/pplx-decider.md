---
name: pplx-decider
title: Perplexity Decider V1.2
kind: provider
category: classify
jobs:
- classify
- cheap
- vision
tags:
- llm
- vision
- closed
- api
- cheap
status: watch
lab: perplexity
preferred_host: openrouter
hosts:
- id: openrouter
  runtime: openrouter
  runtime_model_id: perplexity/pplx-decider-v1.2
  price_in: 0.02
  price_out: 0
  price_note: OpenRouter Perplexity list $0.02/M input, output free
  capacity: hobby
  limits_note: System One Decisions API, not chat completions; 262,144 context; 27B
    cut
  links:
    openrouter: https://openrouter.ai/perplexity/pplx-decider-v1.2
aa_intelligence: unknown
arena_elo: unknown
context_window: 262144
links:
  models_dev: https://models.dev/
source: polypus systemone catalog 2026-11-05
refreshed: 2026-11-05
aliases:
- /roster/budget/pplx-decider/
---


# Perplexity Decider V1.2

## When to reach for it

Pick this card when you want Perplexity Decider V1.2 on **openrouter** for System One classification, routing, or rubric scores with optional image state. Use Decider Flash V1.2 when the cheaper 4B cut is enough.

## Signals

OpenRouter lists $0.02 per M input and free output. Context is 256k. Artificial Analysis chat index is not filed.

## My experience

(Operator notes only.)

## Caveats

Chat completions SDKs will not work. This is the 27B V1.2 cut, not Flash.

## Links

- OpenRouter: https://openrouter.ai/perplexity/pplx-decider-v1.2
- models.dev: https://models.dev/
