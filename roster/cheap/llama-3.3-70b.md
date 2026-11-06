---
name: llama-3.3-70b
title: Llama 3.3 70B
kind: sovereign
category: cheap
jobs:
- chat
tags:
- llm
- open
- cloudflare
status: watch
lab: meta
preferred_host: openrouter
hosts:
- id: cloudflare
  runtime: cloudflare
  runtime_model_id: '@cf/meta/llama-3.3-70b-instruct-fp8-fast'
  price_in: 0.293
  price_out: 2.253
  price_note: Workers AI Neurons at $0.011/1k; 26668/204805 neurons per M in/out;
    CF table Oct 2026
  capacity: hobby
  limits_note: 10k Neurons/day free tier; fp8-fast Workers AI slug
  links:
    cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- id: openrouter
  runtime: openrouter
  runtime_model_id: meta-llama/llama-3.3-70b-instruct
  price_in: 0.1
  price_out: 0.32
  price_note: OpenRouter headline; DeepInfra Turbo often behind the slug; confirm
    openrouter.ai model page
  capacity: hobby
  limits_note: Marketplace routing; not industrial throughput
  links:
    openrouter: https://openrouter.ai/meta-llama/llama-3.3-70b-instruct
aa_intelligence: unknown
arena_elo: unknown
context_window: 128000
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: 2026-11-04
aliases:
- /roster/general/llama-3.3-70b/
---


# Llama 3.3 70B

## When to reach for it

Pick this card when you want a large Llama 3.3 Instruct cut for text chat and you already have it on Polypus. Treat it as a **baseline** next to Llama 4 Scout, not a same-generation vision peer. Use **openrouter** when list price matters. Use **cloudflare** when you already run on Workers AI.

## Signals

Host list prices live in frontmatter. **openrouter** headline is far below the Cloudflare neuron-equivalent rate. Arena and Artificial Analysis are not filed.

## My experience

(Operator notes only.)

## Caveats

This card is text-only. Llama 3.3 is an older cut than Llama 4 Scout on this shelf. **openrouter** can shift upstream provider.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- OpenRouter: https://openrouter.ai/meta-llama/llama-3.3-70b-instruct
- models.dev: https://models.dev/
