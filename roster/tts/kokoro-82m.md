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
    limits_note: Python package kokoro with espeak-ng and misaki G2P; 54 preset voices; 8 languages
  - id: openrouter
    runtime: openrouter
    runtime_model_id: hexgrad/kokoro-82m
    price_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    capacity: hobby
    limits_note: -4","bg-chart-5","bg-chart-6","bg-chart-7","bg-chart-8","bg-chart-9","bg-chart-10","bg-chart-11","bg-chart-12","bg-chart-13","bg-chart-14","bg-chart-15","bg-chart-16","bg-chart-17","bg-chart-18","bg-chart-19","bg-chart-20"],"triggerClass":"
    links:
      openrouter: https://openrouter.ai/hexgrad/kokoro-82m
    price_in: 1.0
    price_out: 19.0
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: hexgrad/Kokoro-82M
    price_note: ill{from{display:block;}}@keyframes mui-auto-fill{from{display:block;}}@-webkit-keyframes mui-auto-fill-cancel{from{display:block;}}@keyframes mui-auto-fill-cancel{from{display:block;}} .mui-style-d4nlmw{display:-webkit-box;display:-webkit-
    capacity: hobby
    limits_note: ill{from{display:block;}}@keyframes mui-auto-fill{from{display:block;}}@-webkit-keyframes mui-auto-fill-cancel{from{display:block;}}@keyframes mui-auto-fill-cancel{from{display:block;}} .mui-style-d4nlmw{display:-webkit-box;display:-webkit-
    links:
      deepinfra: https://deepinfra.com/hexgrad/Kokoro-82M
    price_in: 1000.0
    price_out: 1000.0
source: operator TTS trial 2026-11-06
refreshed: "2026-11-13"
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
