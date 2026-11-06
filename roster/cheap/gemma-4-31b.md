---
name: gemma-4-31b
title: Gemma 4 31B
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
- reasoning
- deepinfra
- cheap
status: watch
lab: google
preferred_host: deepinfra
hosts:
- id: deepinfra
  runtime: deepinfra
  runtime_model_id: google/gemma-4-31B-it
  price_in: 0.2
  price_out: 0.4
  price_note: DeepInfra published list; Priority 1.5x and Flex 0.8x tiers exist
  capacity: hobby
  limits_note: 262,144 context on DeepInfra; fp8; multimodal
  links:
    deepinfra: https://deepinfra.com/google/gemma-4-31B-it
- id: google
  runtime: google
  runtime_model_id: gemma-4-31b-it
  price_note: Polypus gemini_free and gemini_paid routes; Vertex/AI Studio list not
    filed (pricing page fetch timed out)
  capacity: unknown
  limits_note: Confirm model id and region price on Google Cloud before batch jobs
  links:
    provider_pricing: https://cloud.google.com/vertex-ai/generative-ai/pricing
aa_intelligence: 4.4
arena_elo: unknown
context_window: 262144
links:
  models_dev: https://models.dev/
  artificial_analysis: https://artificialanalysis.ai/models/gemma-4-31b
source: polypus catalog diff 2026-11-04
refreshed: '2026-11-05'
aliases:
- /roster/multimodal/gemma-4-31b/
---


# Gemma 4 31B

## When to reach for it

Pick this card when you want the dense Gemma 4 31B Instruct cut with vision. Use **deepinfra** when you have a published list price. Use **google** when you already call Gemma through Gemini/Vertex ids. It is a same-generation **peer** to Gemma 4 26B A4B (different shape: dense 31B vs MoE 26B-A4B).

## Signals

DeepInfra list prices live in frontmatter. Context is 256k on that host. Google list prices are not filed. Artificial Analysis Intelligence Index is 4.4 for Gemma 4 31B (Reasoning). The public leaderboard rounds that to 4.

## My experience

(Operator notes only.)

## Caveats

Do not merge this slug with `gemma-4-26b-a4b`. Confirm the Google runtime id for your project. File industrial capacity only through `capacity-decisions.yaml`.

## Links

- DeepInfra: https://deepinfra.com/google/gemma-4-31B-it
- Vertex pricing: https://cloud.google.com/vertex-ai/generative-ai/pricing
- models.dev: https://models.dev/
