---
name: mistral-small-3.1-24b
title: Mistral Small 3.1 24B
kind: sovereign
category: cheap
jobs:
  - vision
  - chat
  - cheap
tags:
  - llm
  - vision
  - open
  - cloudflare
  - cheap
status: watch
lab: mistral
preferred_host: cloudflare
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: '@cf/mistralai/mistral-small-3.1-24b-instruct'
    price_in: 0.1
    price_out: 2.51
    price_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    capacity: hobby
    limits_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-06"
aliases:
  - /roster/multimodal/mistral-small-3.1-24b/
---


# Mistral Small 3.1 24B

## When to reach for it

Pick this card when you want open Mistral Small 3.1 Instruct with image-in on **cloudflare** Workers AI. Use it for cheap multimodal chat, not as a Gemma 4 peer unless you only need a smaller vision baseline.

## Signals

Cloudflare lists neuron-equivalent token rates in frontmatter. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

Only **cloudflare** is filed. Confirm vision payloads on the Workers AI model page before you rely on image-in.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- models.dev: https://models.dev/
