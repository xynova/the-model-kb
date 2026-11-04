---
name: mercury-decide
title: Mercury Decide
kind: provider
category: budget
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
  - fast
status: watch
lab: inception
preferred_host: openrouter
hosts:
  - id: openrouter
    runtime: openrouter
    runtime_model_id: inception/mercury-decide
    price_in: 0.02
    price_out: 0
    price_note: "OpenRouter headline $0.02/M input, output free; 50% off vs $0.04 list noted on fetch; Inception is the sole provider"
    capacity: hobby
    limits_note: "System One Decisions API, not chat completions; 65,536 context; up to 14 decisions/s on the model page"
    links:
      openrouter: https://openrouter.ai/inception/mercury-decide
aa_intelligence: unknown
arena_elo: unknown
context_window: 65536
links:
  models_dev: https://models.dev/
source: "polypus systemone catalog 2026-11-05"
refreshed: 2026-11-05
---

# Mercury Decide

## When to reach for it

Pick this card when you want Inception Mercury Decide on **openrouter** for System One structured decisions with image-in state. It shares the noul/choice/score contract with Jev. Chat completions SDKs will not work.

## Signals

OpenRouter lists $0.02 per M input and free output (50 percent off the $0.04 list at fetch). Context is 65k. Artificial Analysis chat index is not filed.

## My experience

(Operator notes only.)

## Caveats

This is a lab-locked Decisions API. Promotional OpenRouter headlines can expire. Do not call it as an OpenAI chat model.

## Links

- OpenRouter: https://openrouter.ai/inception/mercury-decide
- models.dev: https://models.dev/
