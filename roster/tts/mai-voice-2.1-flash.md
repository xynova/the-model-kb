---
name: mai-voice-2.1-flash
title: MAI-Voice-2.1-Flash
kind: provider
category: tts
jobs:
  - tts
tags:
  - audio
  - closed
  - api
  - fast
status: watch
lab: microsoft
preferred_host: openrouter
hosts:
  - id: openrouter
    runtime: openrouter
    runtime_model_id: microsoft/mai-voice-2.1-flash
    price_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    capacity: hobby
    limits_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    links:
      openrouter: https://openrouter.ai/microsoft/mai-voice-2.1-flash
    price_in: 1.0
    price_out: 1.0
source: operator TTS trial 2026-11-06
refreshed: "2026-11-13"
---

# MAI-Voice-2.1-Flash

## When to reach for it

Pick **openrouter** when you want the 2.1 Microsoft speech line at the lower character rate and lower listed latency. This is the speed cut of **MAI-Voice-2.1**, not a second host row on that card, not **Kokoro 82M**, and not the retiring **MAI-Voice-2-Flash** slug.

## Signals

OpenRouter lists `microsoft/mai-voice-2.1-flash` at $15 per million characters versus $22 on **MAI-Voice-2.1**, with Azure as the only upstream. Design Arena audiorealism ranks are not filed.

## My experience

(Operator notes only.)

## Caveats

Voice ids must include the `:MAI-Voice-2.1-Flash` suffix. `microsoft` and `azure` are not minted host ids. Do not call this model with a Voice 2.1 voice id.

## Links

- OpenRouter: https://openrouter.ai/microsoft/mai-voice-2.1-flash
