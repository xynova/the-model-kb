---
name: glm-5.3-flash-cf
title: GLM-5.3-Flash (Cloudflare)
kind: sovereign
category: multimodal
tags:
  - llm
  - vision
  - open
  - moe
  - reasoning
  - cloudflare
  - agentic
status: watch
lab: zai
runtime: cloudflare
runtime_model_id: "@cf/zai-org/glm-5.3-flash"
price_in: 0.15
price_out: 0.50
price_note: "Workers Paid or AI Gateway credits required; Neurons behind the scenes; $0.03/M cached input; 13636/45455 neurons per M in/out; confirm CF pricing page"
aa_intelligence: 57
arena_elo: unknown
context_window: 1048576
links:
  cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  artificial_analysis: https://artificialanalysis.ai/models/glm-5-3-flash
  models_dev: https://models.dev/
source: "roster seed gemma-4 segment peers 2026-09-08"
refreshed: 2026-09-08
---

# GLM-5.3-Flash (Cloudflare)

## When to reach for it

Open multimodal MoE on Cloudflare when you want a clear step-up from Gemma 4 26B-A4B (stronger intelligence, vision, long context) and you accept Paid/credits billing and a larger active-parameter class.

## Signals

CF list $0.15 / $0.50 per M input/output ($0.03/M cached input). Artificial Analysis Intelligence Index about 57 for GLM-5.3-Flash vs roughly 26–31 for Gemma 4 26B-A4B. MoE with ~18B active / 320B total (step-up vs Gemma ~4B active). Vision, reasoning, tools; CF context about 1M tokens. Same-host compare; not the same size band.

## My experience

(Not used yet.)

## Caveats

Requires Workers Paid or prepaid AI Gateway credits (not free-tier Workers AI). Step-up class, not an apples-to-apples active-param peer to Gemma 4. Neurons metering still applies.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- Model page: https://developers.cloudflare.com/workers-ai/models/glm-5.3-flash/
- Artificial Analysis: https://artificialanalysis.ai/models/glm-5-3-flash
- models.dev: https://models.dev/
