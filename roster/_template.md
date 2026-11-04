---
name: example-slug
title: Example Title
kind: sovereign
category: general
jobs:
  - chat
  - cheap
tags:
  - llm
  - open
status: watch
lab: meta
preferred_host: cloudflare
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/example/model"
    price_in: unknown
    price_out: unknown
    price_note: unknown
    capacity: hobby
    limits_note: unknown
    links:
      cloudflare_pricing: https://developers.cloudflare.com/workers-ai/platform/pricing/
aa_intelligence: unknown
arena_elo: unknown
context_window: unknown
links:
  artificial_analysis: https://artificialanalysis.ai/leaderboards/models
source: "manual YYYY-MM-DD"
refreshed: YYYY-MM-DD
---

# Example Title

## When to reach for it

Concrete job trigger: situation or constraint that makes this the right pick. State `kind` implications (portable open weights vs lab-locked API vs product-bundled). Name which host row fits which constraint (price vs industrial limits).

## Signals

Short prose restating capability numbers and which source won if prices conflict. Summarize host rows (price and capacity). Omit or say unknown when enrichment found nothing.

## My experience

Personal notes only. Leave empty or a placeholder until you have used the model; the skill MUST NOT invent this section.

## Caveats

Limits, routing quirks, license, neuron free tier, Flex/Priority, proxy-via-gateway notes, or plan caps. Per-host limits belong here or in host `limits_note`.

## Links

- Cloudflare pricing: (from host links, if any)
- DeepInfra: (from host links, if any)
- Provider pricing: (from frontmatter, if any)
- Service docs: (from frontmatter, if any)
- Artificial Analysis: (from frontmatter, if any)
- OpenRouter: (from host links, if any)
- models.dev: (from frontmatter, if any)
- Aider: (from frontmatter, if any)
