---
name: deepseek-v4-flash
title: DeepSeek V4 Flash
kind: sovereign
category: general
jobs:
  - reasoning
  - cheap
  - chat
  - coding
  - long-context
tags:
  - llm
  - open
  - moe
  - reasoning
  - cloudflare
  - deepinfra
  - cheap
  - agentic
status: watch
lab: deepseek
preferred_host: deepinfra
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/deepseek-ai/deepseek-v4-flash-0731"
    price_in: 0.44
    price_out: 1.32
    price_note: "Workers AI Neurons at $0.011/1k; $0.014/M cached input; 40000/120000 neurons per M in/out; CF table Oct 2026; Workers Paid or AI Gateway credits required"
    capacity: hobby
    limits_note: "Not on the 10k Neurons/day free list; 1,048,576 context on Workers AI; reasoning_effort none/low/high/max"
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: deepseek-ai/DeepSeek-V4-Flash-0731
    price_in: 0.06
    price_out: 0.18
    price_note: "DeepInfra published list; cached input $0.015/M; Priority 1.5x and Flex 0.8x tiers exist"
    capacity: hobby
    limits_note: "1,048,576 context; fp8; function calling; lower list price than CF"
    links:
      deepinfra: https://deepinfra.com/deepseek-ai/DeepSeek-V4-Flash-0731
aa_intelligence: 52
arena_elo: unknown
context_window: 1048576
links:
  artificial_analysis: https://artificialanalysis.ai/models/deepseek-v4-flash
  models_dev: https://models.dev/
source: "polypus catalog diff 2026-11-04"
refreshed: 2026-11-04
---

# DeepSeek V4 Flash

## When to reach for it

Pick this card when you want open DeepSeek V4 Flash (0731 cut) for cheap reasoning, chat, or agentic coding with a 1M context window, and you do not need vision. Use **deepinfra** when list price matters. Use **cloudflare** when you already run on Workers AI and can pay Workers Paid or AI Gateway credits.

## Signals

Artificial Analysis puts DeepSeek V4 Flash 0731 (Reasoning, Max Effort) at intelligence index 52. That sits well above GLM-4.7-Flash on the same index. Both filed hosts advertise a 1M context window. Host list prices live in frontmatter; **deepinfra** is the cheaper published list.

## My experience

(Operator notes only.)

## Caveats

This card is text-only. **cloudflare** does not include this model in the daily Neurons free allocation. File industrial capacity only through `capacity-decisions.yaml` after you confirm limits. DeepSeek's first-party API is not a host row here.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- Cloudflare model: https://developers.cloudflare.com/workers-ai/models/deepseek-v4-flash-0731/
- DeepInfra: https://deepinfra.com/deepseek-ai/DeepSeek-V4-Flash-0731
- Artificial Analysis: https://artificialanalysis.ai/models/deepseek-v4-flash
- models.dev: https://models.dev/
