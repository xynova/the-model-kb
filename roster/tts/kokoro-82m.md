---
name: kokoro-82m
title: Kokoro 82M
kind: sovereign
category: tts
jobs:
- tts
- cheap
tags:
- audio
- open
- local
- cheap
- deepinfra
status: watch
lab: hexgrad
preferred_host: local
hosts:
- id: local
  runtime: local
  runtime_model_id: hexgrad/Kokoro-82M
  price_note: Local inference; Apache 2.0 weights; no list character price
  capacity: local
  limits_note: Python package kokoro with espeak-ng and misaki G2P; 54 preset voices;
    8 languages
- id: openrouter
  runtime: openrouter
  runtime_model_id: hexgrad/kokoro-82m
  price_note: OpenRouter $0.62 per 1M characters (page 2026-11-06); not token $/M
  capacity: hobby
  limits_note: POST /api/v1/audio/speech; not EU in-region routing; routes DeepInfra
    or Together
  links:
    openrouter: https://openrouter.ai/hexgrad/kokoro-82m
- id: deepinfra
  runtime: deepinfra
  runtime_model_id: hexgrad/Kokoro-82M
  price_note: DeepInfra $0.99 per 1M characters (page 2026-11-06); Priority and Flex
    tiers exist
  capacity: hobby
  limits_note: Native inference hexgrad/Kokoro-82M; OpenRouter lists this upstream
    cheaper than DeepInfra's own page
  links:
    deepinfra: https://deepinfra.com/hexgrad/Kokoro-82M
source: operator TTS trial 2026-11-06
refreshed: 2026-11-06
---

# Kokoro 82M

## When to reach for it

Pick **local** when you want Apache 2.0 speech on your own machine and can live with preset voices. Use **openrouter** when you want the same weights as an audio API at the OpenRouter character rate. Use **deepinfra** when you want DeepInfra's own endpoint and accept that host's character list price.

This card is the local **baseline** for the Speech shelf, not a Cheap LLM peer and not a cloning or instruction-following TTS flagship.

## Signals

hexgrad publishes 82M StyleTTS 2 weights (v1.0, 8 languages, 54 voices) on Hugging Face. OpenRouter lists $0.62 per million characters; DeepInfra lists $0.99 per million characters. Those are character rates, not token `$/M`. Artificial Analysis TTS Elo and Design Arena audiorealism ranks are not filed.

## My experience

Ran locally for free.

## Caveats

Local install needs `espeak-ng`. Voices are presets, not cloned speakers. Unrelated sites with "kokoro" in the domain are not this project. Together appears on OpenRouter as an upstream; it is not a minted host id here.

## Links

- Hugging Face weights: https://huggingface.co/hexgrad/Kokoro-82M
- GitHub: https://github.com/hexgrad/kokoro
- OpenRouter: https://openrouter.ai/hexgrad/kokoro-82m
- DeepInfra: https://deepinfra.com/hexgrad/Kokoro-82M
