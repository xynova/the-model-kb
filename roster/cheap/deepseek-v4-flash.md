---
name: deepseek-v4-flash
title: DeepSeek V4 Flash
kind: sovereign
category: cheap
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
    runtime_model_id: '@cf/deepseek-ai/deepseek-v4-flash-0731'
    price_in: 0.1
    price_out: 2.51
    price_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    capacity: hobby
    limits_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: deepseek-ai/DeepSeek-V4-Flash-0731
    price_in: 1.0
    price_out: 0.9
    price_note: mation-1p2h4ri{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-d4nlmw{display:-webkit-box;display:-webkit-
    capacity: hobby
    limits_note: mation-1p2h4ri{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-d4nlmw{display:-webkit-box;display:-webkit-
    links:
      deepinfra: https://deepinfra.com/deepseek-ai/DeepSeek-V4-Flash-0731
aa_intelligence: unknown
arena_elo: unknown
context_window: 1048576
links:
  artificial_analysis: https://artificialanalysis.ai/models/deepseek-v4-flash
  models_dev: https://models.dev/
source: polypus catalog diff 2026-11-04
refreshed: "2026-11-13"
aliases:
  - /roster/general/deepseek-v4-flash/
---


# DeepSeek V4 Flash

## When to reach for it

Pick this card when you want open DeepSeek V4 Flash (0731 cut) for cheap reasoning, chat, or agentic coding with a 1M context window, and you do not need vision. Use **deepinfra** when list price matters. Use **cloudflare** when you already run on Workers AI and can pay Workers Paid or AI Gateway credits.

## Signals

Artificial Analysis lists DeepSeek V4 Flash 0731 (Max) with no Intelligence Index on the current pull, so this card does not file a number. The scored DeepSeek row on the public leaderboard is V4.1 Flash, a different cut. Both filed hosts advertise a 1M context window. Host list prices live in frontmatter; **deepinfra** is the cheaper published list.

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
