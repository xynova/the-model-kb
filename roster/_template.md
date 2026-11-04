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

State the job and constraint in active voice. Name which **host id** fits (price vs industrial limits vs product bundle). One or two short paragraphs.

## Signals

Restate capability numbers in plain words. Say which source wins if host prices disagree. Omit what is unknown.

## My experience

Personal notes only. Leave empty or `(Operator notes only.)` until the operator has used the model.

## Caveats

Per-host limits, routing quirks, license, neurons, Flex/Priority, proxy notes. Active voice; no hype.

## Links

- Cloudflare pricing: (from host links, if any)
- DeepInfra: (from host links, if any)
- Provider pricing: (from frontmatter, if any)
- Service docs: (from frontmatter, if any)
- Artificial Analysis: (from frontmatter, if any)
- OpenRouter: (from host links, if any)
- models.dev: (from frontmatter, if any)
- Aider: (from frontmatter, if any)
