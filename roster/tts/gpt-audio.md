---
name: gpt-audio
title: GPT Audio
kind: provider
category: tts
jobs:
  - tts
tags:
  - audio
  - closed
  - api
status: watch
lab: openai
preferred_host: openrouter
hosts:
  - id: openrouter
    runtime: openrouter
    runtime_model_id: openai/gpt-audio
    price_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    capacity: hobby
    limits_note: te()); gtag('config', 'G-R8YZRJS2XN'); html[data-nav-arm='treatment'] [data-nav-tree='control'],html:not([data-nav-arm='treatment']) [data-nav-tree='treatment']{display:none}html[data-nav-arm='treatment'] [data-nav-search]{margin-left:auto;
    links:
      openrouter: https://openrouter.ai/openai/gpt-audio
    price_in: 0.6
    price_out: 2.4
source: operator TTS trial 2026-11-06
refreshed: "2026-11-06"
---

# GPT Audio

## When to reach for it

Do not pick this card for speech reads. Open it only to remember the miss. For local preset-voice TTS, use **Kokoro 82M** on **local**. This is a lab-locked OpenAI audio model (kind **provider**), not an open-weights peer of Kokoro.

## Signals

OpenRouter lists `openai/gpt-audio` with text $2.50 in / $10 out per million tokens and audio $32 in / $64 out per million tokens, matching OpenAI's published `gpt-audio` table. Those audio-token rates are why a short read can cost more than a character-priced TTS. Design Arena audiorealism ranks are not filed.

## My experience

Expensive, and the speech was crap.

## Caveats

OpenRouter forwards to OpenAI only. `openai` is not a minted host id, so first-party API access is not a second row. Do not confuse this card with **gpt-audio-mini** or **gpt-4o-mini-tts**.

## Links

- OpenRouter: https://openrouter.ai/openai/gpt-audio
- OpenAI pricing: https://developers.openai.com/api/docs/pricing
