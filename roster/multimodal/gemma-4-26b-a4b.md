---
name: gemma-4-26b-a4b
title: Gemma 4 26B A4B
kind: sovereign
category: multimodal
jobs:
  - vision
  - reasoning
  - cheap
tags:
  - llm
  - vision
  - open
  - moe
  - reasoning
  - cloudflare
  - cheap
status: reach-for
lab: google
preferred_host: cloudflare
aliases:
  - /roster/multimodal/gemma-4-26b-a4b-cf/
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/google/gemma-4-26b-a4b-it"
    price_in: 0.10
    price_out: 0.30
    price_note: "Workers AI Neurons ($0.011/1k Neurons); 10k Neurons/day free; 9091/27273 neurons per M in/out; confirm CF pricing page"
    capacity: hobby
    limits_note: "10k Neurons/day free tier; verify dashboard after large jobs"
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: openrouter
    runtime: openrouter
    runtime_model_id: google/gemma-4-26b-a4b-it
    price_in: 0.04
    price_out: 0.15
    price_note: "OpenRouter aggregate route; provider mix behind the slug can shift; confirm openrouter.ai model page"
    capacity: hobby
    limits_note: "Marketplace routing; not industrial throughput"
    links:
      openrouter: https://openrouter.ai/google/gemma-4-26b-a4b-it
  - id: google
    runtime: google
    runtime_model_id: gemma-4-26b-a4b-it-maas
    price_in: 0.15
    price_out: 0.60
    price_note: "Vertex AI Gemma open-weights MaaS list; confirm Google Cloud pricing for your region"
    capacity: hobby
    limits_note: "Cloud project quotas apply; not filed as industrial here"
    links:
      provider_pricing: https://cloud.google.com/vertex-ai/generative-ai/pricing
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: google/gemma-4-26B-A4B-it
    price_in: 0.07
    price_out: 0.34
    price_note: "DeepInfra published list on model page; Priority 1.5x and Flex 0.8x tiers exist"
    capacity: industrial
    limits_note: "Operator pick for high-volume batch; DeepInfra supports production throughput vs CF neuron hobby tier"
    links:
      deepinfra: https://deepinfra.com/google/gemma-4-26B-A4B-it
  - id: xinference
    runtime: xinference
    price_note: "Xinference AU Model API; Gemma listed on request on pricing page (Oct 2026); no public list price filed"
    capacity: unknown
    limits_note: "Australian data centres; OpenAI-compatible API; confirm model id with Xinference before use"
    links:
      xinference: https://xinference.co/pricing.html
aa_intelligence: 31
arena_elo: unknown
context_window: 256000
links:
  artificial_analysis: https://artificialanalysis.ai/articles/gemma-4-everything-you-need-to-know
  models_dev: https://models.dev/
source: "roster seed gemma-4 segment 2026-09-08; hosts enriched 2026-11-04; xinference host 2026-11-04"
refreshed: 2026-11-04
---

# Gemma 4 26B A4B

## When to reach for it

Open multimodal MoE when you want portable Gemma 4 weights. Compare **host offers** in the table above: Cloudflare for cheap hobby use, OpenRouter for marketplace routing, Google Vertex for lab-hosted MaaS, DeepInfra when you need **industrial** throughput.

## Signals

CF $0.10 / $0.30 per M (Neurons). OpenRouter aggregate $0.04 / $0.15 per M (confirm on OpenRouter). Google Vertex MaaS about $0.15 / $0.60 per M (confirm region table). DeepInfra $0.07 / $0.34 per M with **industrial** capacity filed for batch production. AA Intelligence Index about 31 for Gemma 4 26B-A4B (Reasoning).

## My experience

(Not used yet.)

## Caveats

Billed in Neurons with a daily free allocation; token-equivalent rates can drift if CF revises the table. Community has reported neuron metering quirks on this model in the past; verify dashboard usage after large jobs. Apache-style open license per Gemma 4 family (confirm current terms on the model card).

## Links

- DeepInfra: https://deepinfra.com/google/gemma-4-26B-A4B-it
- OpenRouter: https://openrouter.ai/google/gemma-4-26b-a4b-it
- Artificial Analysis: https://artificialanalysis.ai/articles/gemma-4-everything-you-need-to-know
- models.dev: https://models.dev/
