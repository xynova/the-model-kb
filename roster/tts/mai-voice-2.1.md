---
name: mai-voice-2.1
title: MAI-Voice-2.1
kind: provider
category: tts
jobs:
  - tts
tags:
  - audio
  - closed
  - api
status: watch
lab: microsoft
preferred_host: openrouter
hosts:
  - id: openrouter
    runtime: openrouter
    runtime_model_id: microsoft/mai-voice-2.1
    price_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    capacity: hobby
    limits_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    links:
      openrouter: https://openrouter.ai/microsoft/mai-voice-2.1
    price_in: 1.0
    price_out: 19.0
source: operator TTS trial 2026-11-06
refreshed: "2026-11-06"
---

# MAI-Voice-2.1

## When to reach for it

Pick **openrouter** when you want Microsoft's current speech API and will pay the character rate. This is a lab-locked **provider** card on the Speech shelf, not an open-weights peer of **Kokoro 82M**, and not **GPT Audio**.

## Signals

OpenRouter lists `microsoft/mai-voice-2.1` at $22 per million characters, with Azure as the only upstream. Lab copy on microsoft.ai matches that $22 character price. Design Arena audiorealism ranks (including the older Voice 2 row) are not filed.

## My experience

(Operator notes only.)

## Caveats

Voice ids must include the `:MAI-Voice-2.1` suffix. `microsoft` and `azure` are not minted host ids, so Foundry or Azure Speech is not a second row. The speed cut is the sibling card **mai-voice-2.1-flash**. The older **MAI-Voice-2** slugs are not this card.

## Links

- OpenRouter: https://openrouter.ai/microsoft/mai-voice-2.1
- Microsoft model page: https://microsoft.ai/models/mai-voice-2-1/
