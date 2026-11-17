---
name: gemini-decision-1
title: Gemini Decision 1
kind: provider
category: classify
jobs:
- classify
- cheap
- vision
tags:
- llm
- vision
- video
- audio
- closed
- api
- cheap
- fast
- agentic
status: watch
lab: google
preferred_host: google
hosts:
- id: google
  runtime: google
  runtime_model_id: gemini-decision-1
  price_in: 0.1
  price_out: 0
  price_note: Gemini API paid list $0.10/M input (text, image, video, audio, PDF);
    output not billed (0 tokens); no context cache; Batch API not supported
  capacity: unknown
  limits_note: Decisions API only (noul/choice/score); up to 1,000 questions and 262,144
    input tokens per request; not a chat model
  links:
    provider_pricing: https://ai.google.dev/gemini-api/docs/pricing
- id: openrouter
  runtime: openrouter
  runtime_model_id: google/gemini-decision-1
  price_in: 0.1
  price_out: 0
  price_note: OpenRouter Google AI Studio list $0.10/M input, output free
  capacity: hobby
  limits_note: OpenRouter Decisions API, not chat completions; Google AI Studio is
    the sole provider on fetch
  links:
    openrouter: https://openrouter.ai/google/gemini-decision-1
aa_intelligence: unknown
arena_elo: unknown
context_window: 262144
links:
  models_dev: https://models.dev/
source: youtube:woWo80dPkVk
refreshed: '2026-11-17'
---

# Gemini Decision 1

## When to reach for it

Pick this card when you want Google Gemini Decision 1 for **System One-style** routing, policy checks, or rubric scores over shared text or media. Call the **google** Gemini Decisions API (`gemini-decision-1`) or the same sku on **openrouter**. Use Jev or Mercury Decide when you need a cheaper text-first decider.

## Signals

Google lists $0.10 per M input tokens and does not bill output. OpenRouter matches that Google AI Studio row. Context is 262k. Artificial Analysis chat Intelligence Index is not filed; Google and AA publish separate decision-board numbers that this card does not store as `aa_intelligence`.

## My experience

(Operator notes only.)

## Caveats

Chat completions SDKs will not work. Google does not cache this sku or support Batch. Vendor JevBench rows are launch self-report, not an independent typed-decision leaf.

## Links

- Gemini API pricing: https://ai.google.dev/gemini-api/docs/pricing
- OpenRouter: https://openrouter.ai/google/gemini-decision-1
- models.dev: https://models.dev/