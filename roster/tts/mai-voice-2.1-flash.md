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
  price_note: OpenRouter $15 per 1M characters (page 2026-11-06); Azure listed as the
    upstream at the same character rate; not token $/M
  capacity: hobby
  limits_note: OpenRouter proxies Azure; not a sovereign host; not EU in-region routing;
    voice ids include the model suffix (for example en-US-Harper:MAI-Voice-2.1-Flash);
    mp3 or pcm 24 kHz mono; OpenRouter P50 latency about 0.38s vs about 0.95s on MAI-Voice-2.1
  links:
    openrouter: https://openrouter.ai/microsoft/mai-voice-2.1-flash
source: operator TTS trial 2026-11-06
refreshed: 2026-11-06
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
