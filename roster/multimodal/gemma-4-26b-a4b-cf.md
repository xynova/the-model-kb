---
name: gemma-4-26b-a4b-cf
title: Gemma 4 26B A4B (Cloudflare)
kind: sovereign
category: multimodal
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
runtime: cloudflare
runtime_model_id: "@cf/google/gemma-4-26b-a4b-it"
price_in: 0.10
price_out: 0.30
price_note: "Workers AI Neurons ($0.011/1k Neurons); 10k Neurons/day free; 9091/27273 neurons per M in/out; confirm CF pricing page"
aa_intelligence: 31
arena_elo: unknown
context_window: 256000
links:
  cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  artificial_analysis: https://artificialanalysis.ai/articles/gemma-4-everything-you-need-to-know
  models_dev: https://models.dev/
source: "roster seed gemma-4 segment 2026-09-08"
refreshed: 2026-09-08
---

# Gemma 4 26B A4B (Cloudflare)

## When to reach for it

Open multimodal MoE on Cloudflare when you want Gemma 4 quality with portable weights and a cheap CF list price. Prefer this over lab Gemini when sovereignty (host switching) matters.

## Signals

CF list price $0.10 / $0.30 per M input/output tokens (Neurons behind the scenes). Artificial Analysis article cites Intelligence Index about 31 for Gemma 4 26B A4B (Reasoning); confirm on AA before relying on the number. 256k context; tools, reasoning, and vision per CF model page.

## My experience

(Not used yet.)

## Caveats

Billed in Neurons with a daily free allocation; token-equivalent rates can drift if CF revises the table. Community has reported neuron metering quirks on this model in the past; verify dashboard usage after large jobs. Apache-style open license per Gemma 4 family (confirm current terms on the model card).

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- Model page: https://developers.cloudflare.com/workers-ai/models/gemma-4-26b-a4b-it/
- Artificial Analysis: https://artificialanalysis.ai/articles/gemma-4-everything-you-need-to-know
- models.dev: https://models.dev/
