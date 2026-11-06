---
name: granite-4.0-h-micro
title: Granite 4.0 H Micro
kind: sovereign
category: cheap
jobs:
- cheap
- chat
tags:
- llm
- open
- cloudflare
- cheap
- fast
status: watch
lab: ibm
preferred_host: cloudflare
hosts:
- id: cloudflare
  runtime: cloudflare
  runtime_model_id: '@cf/ibm-granite/granite-4.0-h-micro'
  price_in: 0.017
  price_out: 0.112
  price_note: Workers AI Neurons at $0.011/1k; 1542/10158 neurons per M in/out; CF
    table Oct 2026
  capacity: hobby
  limits_note: 10k Neurons/day free tier; tiny IBM Granite cut
  links:
    cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: 2026-11-04
aliases:
- /roster/budget/granite-4.0-h-micro/
---


# Granite 4.0 H Micro

## When to reach for it

Pick this card when you want a very cheap IBM Granite text cut on **cloudflare**. Treat it as a **baseline** for tiny chat, not a Granite 4.2 8B or Gemma 4 peer.

## Signals

Cloudflare list rates in frontmatter are among the lowest LLM rows on Workers AI. Artificial Analysis is not filed.

## My experience

(Operator notes only.)

## Caveats

This is a micro cut. Do not use it for hard reasoning. Only **cloudflare** is filed.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- models.dev: https://models.dev/
