---
name: gpt-oss-20b
title: GPT-OSS 20B
kind: sovereign
category: budget
jobs:
  - cheap
  - coding
  - chat
tags:
  - llm
  - open
  - cheap
  - coding
  - cloudflare
status: watch
lab: openai
preferred_host: openrouter
hosts:
  - id: openrouter
    runtime: openrouter
    runtime_model_id: openai/gpt-oss-20b
    price_in: 0.008
    price_out: 0.032
    price_note: "OpenRouter headline aggregate $0.008/$0.032 per M; upstream mix can shift"
    capacity: hobby
    limits_note: "Marketplace routing; not industrial throughput"
    links:
      openrouter: https://openrouter.ai/openai/gpt-oss-20b
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/openai/gpt-oss-20b"
    price_in: 0.20
    price_out: 0.30
    price_note: "Workers AI Neurons at $0.011/1k; 18182/27273 neurons per M in/out; CF table Oct 2026"
    capacity: hobby
    limits_note: "10k Neurons/day free tier; same open weights as the OpenRouter slug"
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  models_dev: https://models.dev/
source: "polypus catalog diff 2026-11-04"
refreshed: 2026-11-04
---

# GPT-OSS 20B

## When to reach for it

Pick this card when you want open GPT-OSS 20B for cheap chat or light coding. Use **openrouter** when the aggregate list is enough. Use **cloudflare** when you already run on Workers AI.

## Signals

OpenRouter headline is far below the Cloudflare neuron-equivalent rate. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

Open weights, not a locked GPT-5.x Cursor service. **openrouter** can shift upstream provider.

## Links

- OpenRouter: https://openrouter.ai/openai/gpt-oss-20b
- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- models.dev: https://models.dev/
