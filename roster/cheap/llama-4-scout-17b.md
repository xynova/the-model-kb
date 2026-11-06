---
name: llama-4-scout-17b
title: Llama 4 Scout 17B
kind: sovereign
category: cheap
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
  - /roster/multimodal/llama-4-scout-17b/
  - /roster/multimodal/llama-4-scout-17b-cf/
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: '@cf/meta/llama-4-scout-17b-16e-instruct'
    price_in: 0.1
    price_out: 2.51
    price_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    capacity: hobby
    limits_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: openrouter
    runtime: openrouter
    runtime_model_id: meta-llama/llama-4-scout
    price_in: 0.3
    price_out: 19.0
    price_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    capacity: hobby
    limits_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    links:
      openrouter: https://openrouter.ai/meta-llama/llama-4-scout-17b-16e-instruct
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: meta-llama/Llama-4-Scout-17B-16E-Instruct
    price_in: 1.0
    price_out: 0.9
    price_note: '{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-670439{display:-webkit-inline-box;display:-webkit-inline-'
    capacity: hobby
    limits_note: '{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-670439{display:-webkit-inline-box;display:-webkit-inline-'
    links:
      deepinfra: https://deepinfra.com/meta-llama/Llama-4-Scout-17B-16E-Instruct
  - id: xinference
    runtime: xinference
    price_note: Xinference AU Model API; Llama listed on request on pricing page (Oct 2026); no public list price filed
    capacity: unknown
    limits_note: ine-tune models unless you agree in writing. They’re processed in Australian data centres or on your own hardware."}},{"@type":"Question","name":"What does Enterprise add over Community?","acceptedAnswer":{"@type":"Answer","text":"Up to 2
    links:
      xinference: https://xinference.co/pricing.html
aa_intelligence: unknown
arena_elo: unknown
context_window: 131000
links:
  artificial_analysis: https://artificialanalysis.ai/models/comparisons/llama-4-scout-vs-gpt-5
  models_dev: https://models.dev/
source: roster seed gemma-4 segment peers 2026-09-08; scorecard refresh 2026-11-04
refreshed: "2026-11-06"
---


# Llama 4 Scout 17B

## When to reach for it

Use this card when you need open Llama 4 vision on a hosted API and you accept a low intelligence index next to Gemma 4 on the same shelf. Pick **cloudflare** if you already meter Neurons on Workers AI. Pick **openrouter** or **deepinfra** when you want a lower list rate ($0.10 / $0.30 per M) and can handle routing or a second vendor.

## Signals

Artificial Analysis lists Llama 4 Scout with no Intelligence Index on the current pull. **cloudflare** lists $0.27 / $0.85 per M with a 131k context window. **openrouter** and **deepinfra** list $0.10 / $0.30 per M on their public pages (OpenRouter shows a much larger context window on the aggregate slug).

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
