---
name: glm-4.7-flash
title: GLM-4.7-Flash
kind: sovereign
category: cheap
jobs:
  - reasoning
  - cheap
  - chat
tags:
  - llm
  - open
  - moe
  - reasoning
  - cloudflare
  - cheap
  - agentic
status: watch
lab: zai
preferred_host: cloudflare
aliases:
  - /roster/general/glm-4.7-flash/
  - /roster/general/glm-4.7-flash-cf/
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: '@cf/zai-org/glm-4.7-flash'
    price_in: 0.1
    price_out: 2.51
    price_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    capacity: hobby
    limits_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: openrouter
    runtime: openrouter
    runtime_model_id: z-ai/glm-4.7-flash
    price_in: 0.4
    price_out: 19.0
    price_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    capacity: hobby
    limits_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    links:
      openrouter: https://openrouter.ai/z-ai/glm-4.7-flash
aa_intelligence: unknown
arena_elo: unknown
context_window: 131072
links:
  artificial_analysis: https://artificialanalysis.ai/models/releases/glm-4-7-flash
  models_dev: https://models.dev/
source: roster seed gemma-4 segment peers 2026-09-08; scorecard refresh 2026-11-04
refreshed: "2026-11-06"
---


# GLM-4.7-Flash

## When to reach for it

Pick this card when you want an open Z.ai MoE for chat, cheap reasoning, or agentic tool use, and you do not need vision. Use **cloudflare** when you already run on Workers AI and want the neuron free tier. Use **openrouter** when you want a single slug and can accept Novita or Cloudflare behind the route.

## Signals

Artificial Analysis lists GLM-4.7-Flash (Reasoning) with no Intelligence Index on the current pull (~3B active / 31B total MoE). **cloudflare** lists $0.0605 / $0.40 per M tokens with a 131k context window. **openrouter** lists the same headline prices with a 200k context window on the model page.

## My experience

(Operator notes only.)

## Caveats

This Workers AI card is text-only. Neurons billing and the 10k Neurons/day free allocation apply on **cloudflare**. **openrouter** can shift upstream provider; check routing if latency or uptime matters. DeepInfra is not filed here until a public model page confirms list pricing.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- Cloudflare model: https://developers.cloudflare.com/workers-ai/models/glm-4.7-flash/
- OpenRouter: https://openrouter.ai/z-ai/glm-4.7-flash
- Artificial Analysis: https://artificialanalysis.ai/models/releases/glm-4-7-flash
- models.dev: https://models.dev/
