---
name: glm-5.3-flash
title: GLM-5.3-Flash
kind: sovereign
category: cheap
jobs:
- vision
- reasoning
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
preferred_host: cloudflare
aliases:
- /roster/multimodal/glm-5.3-flash/
- /roster/multimodal/glm-5.3-flash-cf/
hosts:
- id: cloudflare
  runtime: cloudflare
  runtime_model_id: '@cf/zai-org/glm-5.3-flash'
  price_in: 0.15
  price_out: 0.5
  price_note: Workers Paid or AI Gateway credits required; $0.03/M cached input; 13636/45455
    neurons per M in/out; confirm CF pricing page
  capacity: hobby
  limits_note: Paid Workers or AI Gateway credits required on CF; ~1M context on CF
    listing
  links:
    cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- id: openrouter
  runtime: openrouter
  runtime_model_id: z-ai/glm-5.3-flash
  price_in: 0.0269
  price_out: 0.0996
  price_note: OpenRouter aggregate list; upstream provider mix can shift; confirm
    openrouter.ai model page
  capacity: hobby
  limits_note: Marketplace routing; headline OR prices differ from CF/DeepInfra list;
    verify before budgeting
  links:
    openrouter: https://openrouter.ai/z-ai/glm-5.3-flash
- id: deepinfra
  runtime: deepinfra
  runtime_model_id: zai-org/GLM-5.3-Flash
  price_in: 0.15
  price_out: 0.5
  price_note: DeepInfra published list; $0.03/M cached input on model page
  capacity: unknown
  limits_note: Strong agentic/coding signals on public benchmarks; operator has not
    filed industrial capacity
  links:
    deepinfra: https://deepinfra.com/zai-org/GLM-5.3-Flash
- id: xinference
  runtime: xinference
  price_note: Xinference AU Model API; GLM 5.3 on request on pricing page (Oct 2026);
    no public list price filed
  capacity: unknown
  limits_note: Australian data centres; OpenAI-compatible API; confirm model id with
    Xinference before use
  links:
    xinference: https://xinference.co/pricing.html
aa_intelligence: 27.8
arena_elo: unknown
context_window: 1048576
links:
  artificial_analysis: https://artificialanalysis.ai/models/glm-5-3-flash
  models_dev: https://models.dev/
source: roster seed gemma-4 segment peers 2026-09-08; scorecard refresh 2026-11-04
refreshed: '2026-11-05'
---


# GLM-5.3-Flash

## When to reach for it

Reach for this card when you want a multimodal Z.ai model with vision, tools, and about 1M context on **cloudflare**. You need Workers Paid or AI Gateway credits on CF. Use **openrouter** when the aggregate route is cheaper on paper but you accept provider churn. Use **deepinfra** when you want the same $0.15 / $0.50 list as CF with cached-input pricing on the model page.

## Signals

Artificial Analysis Intelligence Index is 27.8 for GLM-5.3-Flash. The public leaderboard rounds that to 28. Gemma 4 26B-A4B has no index on the same pull. The model is a large MoE (~18B active / 320B total), so treat it as a different shape from Gemma. **cloudflare** and **deepinfra** both list $0.15 / $0.50 per M on public pages. **openrouter** lists a much lower headline rate; treat it as a route quote, not a guarantee of CF parity.

## My experience

(Operator notes only.)

## Caveats

**cloudflare** is not on the free Neurons tier for this model. **openrouter** prices can diverge from CF and DeepInfra. **xinference** is on-request only. File industrial capacity only through `capacity-decisions.yaml` after you confirm limits.

## Links

- Cloudflare pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
- Cloudflare model: https://developers.cloudflare.com/workers-ai/models/glm-5.3-flash/
- OpenRouter: https://openrouter.ai/z-ai/glm-5.3-flash
- DeepInfra: https://deepinfra.com/zai-org/GLM-5.3-Flash
- Artificial Analysis: https://artificialanalysis.ai/models/glm-5-3-flash
- models.dev: https://models.dev/
