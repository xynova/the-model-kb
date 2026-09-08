---
name: glm-4.7-flash-cf
title: GLM-4.7-Flash (Cloudflare)
kind: sovereign
category: general
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
runtime: cloudflare
runtime_model_id: "@cf/zai-org/glm-4.7-flash"
price_in: 0.06
price_out: 0.40
price_note: "Workers AI Neurons; 10k Neurons/day free; 5500/36400 neurons per M in/out; confirm CF pricing page"
aa_intelligence: 23
arena_elo: unknown
context_window: 131072
links:
  cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  artificial_analysis: https://artificialanalysis.ai/models/releases/glm-4-7-flash
  models_dev: https://models.dev/
source: "roster seed gemma-4 segment peers 2026-09-08"
refreshed: 2026-09-08
---

# GLM-4.7-Flash (Cloudflare)

## When to reach for it

Open MoE text/reasoning peer to Gemma 4 26B-A4B on the same Cloudflare host when you want slightly cheaper input and strong tool/dialogue use, and you do not need vision.

## Signals

CF list $0.06 / $0.40 per M input/output. Artificial Analysis: Intelligence Index 23 for GLM-4.7-Flash (Reasoning), ~3B active / 31.2B total; Gemma 4 26B-A4B (Reasoning) is about 26 on the same index. Peer for capability band; not a vision peer (Gemma has image input on CF). CF context 131k vs Gemma 256k.

## My experience

(Not used yet.)

## Caveats

Text-only on Workers AI for this card. Neurons billing and free daily allocation apply. Confirm current CF table before large runs.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- Model page: https://developers.cloudflare.com/workers-ai/models/glm-4.7-flash/
- Artificial Analysis: https://artificialanalysis.ai/models/releases/glm-4-7-flash
- models.dev: https://models.dev/
