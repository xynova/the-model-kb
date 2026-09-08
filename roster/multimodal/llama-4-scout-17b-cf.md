---
name: llama-4-scout-17b-cf
title: Llama 4 Scout 17B (Cloudflare)
kind: sovereign
category: multimodal
tags:
  - llm
  - vision
  - open
  - moe
  - cloudflare
status: watch
lab: meta
runtime: cloudflare
runtime_model_id: "@cf/meta/llama-4-scout-17b-16e-instruct"
price_in: 0.27
price_out: 0.85
price_note: "Workers AI Neurons; 10k Neurons/day free; 24545/77273 neurons per M in/out; confirm CF pricing page"
aa_intelligence: 10
arena_elo: unknown
context_window: 131000
links:
  cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  artificial_analysis: https://artificialanalysis.ai/models/comparisons/llama-4-scout-vs-gpt-5
  models_dev: https://models.dev/
source: "roster seed gemma-4 segment peers 2026-09-08"
refreshed: 2026-09-08
---

# Llama 4 Scout 17B (Cloudflare)

## When to reach for it

Open multimodal MoE on Cloudflare when you want Llama 4 vision/text on the same host as Gemma 4, and you accept a lower intelligence band and higher CF list price.

## Signals

CF list $0.27 / $0.85 per M (about 3× Gemma 4 input on CF). Artificial Analysis Intelligence Index about 10 for Llama 4 Scout; AA marks reasoning as no. Peer for **modality + open MoE on CF**, not for intelligence/reasoning class next to Gemma 4 26B-A4B (~26–31). CF context 131k (AA lists much larger native windows elsewhere; use the CF figure for this runtime).

## My experience

(Not used yet.)

## Caveats

Not an intelligence peer to Gemma 4 despite multimodal MoE shape. Prefer Gemma 4 or GLM-4.7-Flash on CF when quality-per-dollar matters. Neurons billing applies.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- Model page: https://developers.cloudflare.com/workers-ai/models/llama-4-scout-17b-16e-instruct/
- Artificial Analysis (example compare): https://artificialanalysis.ai/models/comparisons/gpt-5-vs-llama-4-scout
- models.dev: https://models.dev/
