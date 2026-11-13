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
    price_in: 0.1
    price_out: 2.51
    price_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    capacity: hobby
    limits_note: )(); Pricing · Cloudflare Workers AI docs {"@context":"https://schema.org","@type":"TechArticle","@id":"https://developers.cloudflare.com/workers-ai/platform/pricing/#page","headline":"Pricing","description":"Workers AI pricing is based on
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
  - id: openrouter
    runtime: openrouter
    runtime_model_id: z-ai/glm-5.3-flash
    price_in: 0.47
    price_out: 19.0
    price_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    capacity: hobby
    limits_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    links:
      openrouter: https://openrouter.ai/z-ai/glm-5.3-flash
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: zai-org/GLM-5.3-Flash
    price_in: 1.0
    price_out: 0.9
    price_note: mation-1p2h4ri{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-d4nlmw{display:-webkit-box;display:-webkit-
    capacity: unknown
    limits_note: mation-1p2h4ri{0%{stroke-dasharray:1px,200px;stroke-dashoffset:0;}50%{stroke-dasharray:100px,200px;stroke-dashoffset:-15px;}100%{stroke-dasharray:100px,200px;stroke-dashoffset:-125px;}} .mui-style-d4nlmw{display:-webkit-box;display:-webkit-
    links:
      deepinfra: https://deepinfra.com/zai-org/GLM-5.3-Flash
  - id: xinference
    runtime: xinference
    price_note: Xinference AU Model API; GLM 5.3 on request on pricing page (Oct 2026); no public list price filed
    capacity: unknown
    limits_note: ine-tune models unless you agree in writing. They’re processed in Australian data centres or on your own hardware."}},{"@type":"Question","name":"What does Enterprise add over Community?","acceptedAnswer":{"@type":"Answer","text":"Up to 2
    links:
      xinference: https://xinference.co/pricing.html
aa_intelligence: 27.8
arena_elo: unknown
context_window: 1048576
links:
  artificial_analysis: https://artificialanalysis.ai/models/glm-5-3-flash
  models_dev: https://models.dev/
source: roster seed gemma-4 segment peers 2026-09-08; scorecard refresh 2026-11-04
refreshed: "2026-11-13"
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
