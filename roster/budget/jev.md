---
name: jev
title: Jev
kind: provider
category: budget
jobs:
  - classify
  - cheap
tags:
  - llm
  - closed
  - api
  - cheap
  - fast
status: watch
lab: typesafe
preferred_host: local
hosts:
  - id: local
    runtime: local
    runtime_model_id: typesafe/jev
    price_in: 0.042
    price_out: 0
    price_note: "TypeSafe list $0.042/M input, output free (blog 2026-09-15); Polypus route cf_local/typesafe/jev"
    capacity: unknown
    limits_note: "System One API (noul/choice/score); not a chat model; ~64k context; early access"
aa_intelligence: unknown
arena_elo: unknown
context_window: 64000
source: "polypus systemone catalog 2026-11-05"
refreshed: 2026-11-05
---

# Jev

## When to reach for it

Pick this card when you want TypeSafe Jev for high-frequency **System One** decisions (noul, choice, score) instead of a chat model. Polypus serves it as `cf_local/typesafe/jev`. Use it to classify, route, or guardrail with calibrated probabilities.

## Signals

TypeSafe lists $0.042 per M input tokens and free output. Calls are tens to hundreds of milliseconds on their published note. This is not an Artificial Analysis chat index.

## My experience

(Operator notes only.)

## Caveats

Jev does not generate free-form text. The Polypus `cf_local` prefix is a gateway route, not a Cloudflare Workers AI public slug. There is no `typesafe` host id in taxonomy yet; this row is `local`. Confirm early-access terms on TypeSafe.

## Links

- TypeSafe intro: https://typesafe.ai/blog/introducing-system-one-models-and-jev
- API: https://jevapi.dev/
