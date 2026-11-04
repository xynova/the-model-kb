---
name: glm-4.7-flash
title: GLM-4.7-Flash
kind: sovereign
category: general
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
  - /roster/general/glm-4.7-flash-cf/
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/zai-org/glm-4.7-flash"
    price_in: 0.0605
    price_out: 0.40
    price_note: "Workers AI Neurons at $0.011/1k; 10k Neurons/day free; 5500/36400 neurons per M in/out; confirm CF pricing page"
    capacity: hobby
    limits_note: "10k Neurons/day free tier; 131k context on Workers AI"
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: openrouter
    runtime: openrouter
    runtime_model_id: z-ai/glm-4.7-flash
    price_in: 0.0605
    price_out: 0.40
    price_note: "OpenRouter list; routes Cloudflare or NovitaAI behind the slug; confirm openrouter.ai model page"
    capacity: hobby
    limits_note: "Marketplace routing; 200k context on OpenRouter listing; not industrial throughput"
    links:
      openrouter: https://openrouter.ai/z-ai/glm-4.7-flash
aa_intelligence: 23
arena_elo: unknown
context_window: 131072
links:
  artificial_analysis: https://artificialanalysis.ai/models/releases/glm-4-7-flash
  models_dev: https://models.dev/
source: "roster seed gemma-4 segment peers 2026-09-08; scorecard refresh 2026-11-04"
refreshed: 2026-11-04
---

# GLM-4.7-Flash

## When to reach for it

Pick this card when you want an open Z.ai MoE for chat, cheap reasoning, or agentic tool use, and you do not need vision. Use **cloudflare** when you already run on Workers AI and want the neuron free tier. Use **openrouter** when you want a single slug and can accept Novita or Cloudflare behind the route.

## Signals

Artificial Analysis puts GLM-4.7-Flash (Reasoning) at intelligence index 23 (~3B active / 31B total MoE). **cloudflare** lists $0.0605 / $0.40 per M tokens with a 131k context window. **openrouter** lists the same headline prices with a 200k context window on the model page. Gemma 4 26B-A4B (Reasoning) scores about 26 on the same index if you need a vision peer instead.

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
