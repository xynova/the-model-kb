---
name: llama-4-scout-17b
title: Llama 4 Scout 17B
kind: sovereign
category: multimodal
jobs:
  - vision
tags:
  - llm
  - vision
  - open
  - moe
  - cloudflare
status: watch
lab: meta
preferred_host: cloudflare
aliases:
  - /roster/multimodal/llama-4-scout-17b-cf/
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/meta/llama-4-scout-17b-16e-instruct"
    price_in: 0.27
    price_out: 0.85
    price_note: "Workers AI Neurons at $0.011/1k; 10k Neurons/day free; 24545/77273 neurons per M in/out; confirm CF pricing page"
    capacity: hobby
    limits_note: "10k Neurons/day free tier; 131k context on Workers AI"
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: openrouter
    runtime: openrouter
    runtime_model_id: meta-llama/llama-4-scout
    price_in: 0.10
    price_out: 0.30
    price_note: "OpenRouter list; routes DeepInfra or Vertex behind the slug; confirm openrouter.ai model page"
    capacity: hobby
    limits_note: "Marketplace routing; 1.3M context on OpenRouter listing; not industrial throughput"
    links:
      openrouter: https://openrouter.ai/meta-llama/llama-4-scout-17b-16e-instruct
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: meta-llama/Llama-4-Scout-17B-16E-Instruct
    price_in: 0.10
    price_out: 0.30
    price_note: "DeepInfra published list; Priority 1.5x and Flex 0.8x tiers exist"
    capacity: hobby
    limits_note: "Lower list price than CF; confirm DeepInfra model page before batch jobs"
    links:
      deepinfra: https://deepinfra.com/meta-llama/Llama-4-Scout-17B-16E-Instruct
  - id: xinference
    runtime: xinference
    price_note: "Xinference AU Model API; Llama listed on request on pricing page (Oct 2026); no public list price filed"
    capacity: unknown
    limits_note: "Australian data centres; OpenAI-compatible API; confirm model id with Xinference before use"
    links:
      xinference: https://xinference.co/pricing.html
aa_intelligence: 10
arena_elo: unknown
context_window: 131000
links:
  artificial_analysis: https://artificialanalysis.ai/models/comparisons/llama-4-scout-vs-gpt-5
  models_dev: https://models.dev/
source: "roster seed gemma-4 segment peers 2026-09-08; scorecard refresh 2026-11-04"
refreshed: 2026-11-04
---

# Llama 4 Scout 17B

## When to reach for it

Use this card when you need open Llama 4 vision on a hosted API and you accept a low intelligence index next to Gemma 4 on the same shelf. Pick **cloudflare** if you already meter Neurons on Workers AI. Pick **openrouter** or **deepinfra** when you want a lower list rate ($0.10 / $0.30 per M) and can handle routing or a second vendor.

## Signals

Artificial Analysis puts Llama 4 Scout at intelligence index about 10 with reasoning marked no. **cloudflare** lists $0.27 / $0.85 per M with a 131k context window. **openrouter** and **deepinfra** list $0.10 / $0.30 per M on their public pages (OpenRouter shows a much larger context window on the aggregate slug). This is a modality peer on CF, not a reasoning peer to Gemma 4 26B-A4B (~26–31).

## My experience

(Operator notes only.)

## Caveats

Do not treat Scout as a quality substitute for Gemma 4 or GLM-5.3-Flash when reasoning density matters. **xinference** stays on-request with no filed list price. Neurons and free-tier rules apply on **cloudflare**.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- Cloudflare model: https://developers.cloudflare.com/workers-ai/models/llama-4-scout-17b-16e-instruct/
- OpenRouter: https://openrouter.ai/meta-llama/llama-4-scout-17b-16e-instruct
- DeepInfra: https://deepinfra.com/meta-llama/Llama-4-Scout-17B-16E-Instruct
- Artificial Analysis: https://artificialanalysis.ai/models/comparisons/llama-4-scout-vs-gpt-5
- models.dev: https://models.dev/
