---
name: drama-3
title: Drama 3
summary: "Fish Audio open-weight TTS for directed performance, voice cloning, and long-form speech."
category: models
does: speech
tags:
  - audio
  - local
  - api
  - mixed
  - moe
status: watch
run: local
license: mixed
vram_hint: consumer
links:
  homepage: https://fish.audio/blog/fish-audio-drama-3/
  github: https://github.com/fishaudio/fish-speech
  weights: https://huggingface.co/fishaudio/drama-3
media:
  cover: models/drama-3/media/cover.jpg
  gallery:
    - models/drama-3/media/g01.jpg
    - models/drama-3/media/g02.jpg
    - models/drama-3/media/g03.jpg
source: "youtube:tDWtcyeNNew"
added: 2026-11-03
---

# Drama 3

## What it is

Fish Audio open-weight TTS for directed performance, voice cloning, and long-form speech. 35B MoE with about 3B active per token. You steer delivery with in-script tags and optional reference audio.

| | |
|:---:|:---:|
| ![drama-3 1](drama-3/media/thumb-cover.jpg) | ![drama-3 2](drama-3/media/thumb-g01.jpg) |
| ![drama-3 3](drama-3/media/thumb-g02.jpg) | ![drama-3 4](drama-3/media/thumb-g03.jpg) |

## Why keep

Compare when choosing a local or API TTS path that takes direction (emotion, pause, emphasis), not only read-aloud.

## When to reach for it

Local via SGLang-Omni or vLLM-Omni from the Hugging Face weights, or the fish.audio API (`model: drama-3`), when a line needs performance control or a cloned voice.

## Caveats

Fish Audio Community License, not OSI-open: free under $500k annual revenue; hosted TTS or cloning as a service needs a separate license; no cloning a real person without consent. Official serving numbers are on H200; the roundup claimed 8 GB VRAM reports. API list price is $50 per million characters (promo pricing expires 12 Nov 2026).

## Links

- Homepage: https://fish.audio/blog/fish-audio-drama-3/
- GitHub: https://github.com/fishaudio/fish-speech
- Weights: https://huggingface.co/fishaudio/drama-3
