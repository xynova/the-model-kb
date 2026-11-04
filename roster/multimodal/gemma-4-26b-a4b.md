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
    price_note: "Workers AI Neurons at $0.011/1k; 10k Neurons/day free; 9091/27273 neurons per M in/out; CF table Oct 2026"
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
source: "roster seed gemma-4 segment 2026-09-08; scorecard refresh 2026-11-04"
refreshed: 2026-11-04
---

# Gemma 4 26B A4B

## When to reach for it

Pick this card when you want portable open Gemma 4 with vision and reasoning on the shelf. Use **cloudflare** for hobby spend and the neuron free tier. Use **openrouter** when you want a cheaper aggregate slug. Use **google** when you run Vertex MaaS in your project. Use **deepinfra** when you filed **industrial** capacity for batch throughput. Use **xinference** only after you confirm model access with the vendor (on-request pricing).

## Signals

Artificial Analysis puts Gemma 4 26B-A4B (Reasoning) at intelligence index about 31. Context is about 256k on the card. Host list prices live in frontmatter; compare hosts by **capacity** and job, not by body prose.

## My experience

(Operator notes only.)

## Caveats

Neurons on **cloudflare** can meter differently than token math suggests; check the dashboard after large jobs. **deepinfra** **industrial** capacity is an operator filing, not an automatic mine result. **xinference** has no public list price on the pricing page for Gemma.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- OpenRouter: https://openrouter.ai/google/gemma-4-26b-a4b-it
- DeepInfra: https://deepinfra.com/google/gemma-4-26B-A4B-it
- Xinference: https://xinference.co/pricing.html
- Artificial Analysis: https://artificialanalysis.ai/articles/gemma-4-everything-you-need-to-know
- models.dev: https://models.dev/
