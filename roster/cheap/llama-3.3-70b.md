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
    price_in: 0.1
    price_out: 2.51
    price_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    capacity: hobby
    limits_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: openrouter
    runtime: openrouter
    runtime_model_id: meta-llama/llama-3.3-70b-instruct
    price_in: 0.32
    price_out: 1.0
    price_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    capacity: hobby
    limits_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    links:
      openrouter: https://openrouter.ai/meta-llama/llama-3.3-70b-instruct
aa_intelligence: unknown
arena_elo: unknown
context_window: 128000
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-13"
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
