---
name: gemma-4-26b-a4b
title: Gemma 4 26B A4B
kind: sovereign
category: cheap
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
  - /roster/multimodal/gemma-4-26b-a4b/
  - /roster/multimodal/gemma-4-26b-a4b-cf/
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: '@cf/google/gemma-4-26b-a4b-it'
    price_in: 0.1
    price_out: 2.51
    price_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    capacity: hobby
    limits_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: openrouter
    runtime: openrouter
    runtime_model_id: google/gemma-4-26b-a4b-it
    price_in: 0.051
    price_out: 1.0
    price_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    capacity: hobby
    limits_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    links:
      openrouter: https://openrouter.ai/google/gemma-4-26b-a4b-it
  - id: google
    runtime: google
    runtime_model_id: gemma-4-26b-a4b-it-maas
    price_in: 300.0
    price_out: 5.77
    price_note: ver{-webkit-box-shadow:0 1px 0 rgba(0,0,0,0.15),0 1px 2px rgba(0,0,0,0.2);box-shadow:0 1px 0 rgba(0,0,0,0.15),0 1px 2px rgba(0,0,0,0.2)}.gb_Eb .gb_z,.gb_Fb .gb_z{position:absolute;right:1px}.gb_Hb.gb_aa,.gb_bb.gb_aa,.gb_z.gb_aa{-webkit-box-
    capacity: hobby
    limits_note: ver{-webkit-box-shadow:0 1px 0 rgba(0,0,0,0.15),0 1px 2px rgba(0,0,0,0.2);box-shadow:0 1px 0 rgba(0,0,0,0.15),0 1px 2px rgba(0,0,0,0.2)}.gb_Eb .gb_z,.gb_Fb .gb_z{position:absolute;right:1px}.gb_Hb.gb_aa,.gb_bb.gb_aa,.gb_z.gb_aa{-webkit-box-
    links:
      provider_pricing: https://cloud.google.com/vertex-ai/generative-ai/pricing
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: google/gemma-4-26B-A4B-it
    price_in: 1.0
    price_out: 0.9
    price_note: mation-1p2h4ri{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-d4nlmw{display:-webkit-box;display:-webkit-
    capacity: industrial
    limits_note: mation-1p2h4ri{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-d4nlmw{display:-webkit-box;display:-webkit-
    links:
      deepinfra: https://deepinfra.com/google/gemma-4-26B-A4B-it
  - id: xinference
    runtime: xinference
    price_note: Xinference AU Model API; Gemma listed on request on pricing page (Oct 2026); no public list price filed
    capacity: unknown
    limits_note: ine-tune models unless you agree in writing. They’re processed in Australian data centres or on your own hardware."}},{"@type":"Question","name":"What does Enterprise add over Community?","acceptedAnswer":{"@type":"Answer","text":"Up to 2
    links:
      xinference: https://xinference.co/pricing.html
aa_intelligence: unknown
arena_elo: unknown
context_window: 256000
links:
  artificial_analysis: https://artificialanalysis.ai/articles/gemma-4-everything-you-need-to-know
  models_dev: https://models.dev/
source: roster seed gemma-4 segment 2026-09-08; scorecard refresh 2026-11-04
refreshed: "2026-11-13"
---


# Gemma 4 26B A4B

## When to reach for it

Pick this card when you want portable open Gemma 4 with vision and reasoning on the shelf. Use **cloudflare** for hobby spend and the neuron free tier. Use **openrouter** when you want a cheaper aggregate slug. Use **google** when you run Vertex MaaS in your project. Use **deepinfra** when you filed **industrial** capacity for batch throughput. Use **xinference** only after you confirm model access with the vendor (on-request pricing).

## Signals

Artificial Analysis lists Gemma 4 26B-A4B (Reasoning) with no Intelligence Index on the current pull. Context is about 256k on the card. Host list prices live in frontmatter; compare hosts by **capacity** and job, not by body prose.

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
