---
name: auk
title: AuK
summary: "Tencent Hunyuan 1.5B speech foundation model that does TTS, zero-shot voice cloning, content micro-edits, enhancement, and related tasks through natural-language instructions."
category: models
does: speech
tags:
  - audio
  - open
  - local
status: watch
run: local
license: open
vram_hint: consumer
links:
  homepage: https://auk-project.github.io/
  github: https://github.com/Tencent-Hunyuan/AuK
  weights: https://huggingface.co/tencent/AuK
media:
  cover: models/auk/media/cover.jpg
  gallery:
    - models/auk/media/g01.jpg
    - models/auk/media/g02.jpg
source: "youtube:nZYJdwM-_nI"
added: 2026-09-16
---

# AuK

## What it is

Tencent Hunyuan 1.5B speech foundation model that does TTS, zero-shot voice cloning, content micro-edits, enhancement, and related tasks through natural-language instructions. AuK-Flash is the distilled four-step variant for faster inference.

| | |
|:---:|:---:|
| ![auk 1](auk/media/thumb-cover.jpg) | ![auk 2](auk/media/thumb-g01.jpg) |
| ![auk 3](auk/media/thumb-g02.jpg) |   |

## Why keep

Local speech generate-and-edit in one stack when you need clone, word-level rewrite, denoise, or emotion/timbre changes without a closed API.

## When to reach for it

Content pipelines that edit existing speech or generate new takes on a consumer GPU; compare AuK vs AuK-Flash for quality versus speed.

## Caveats

Official VRAM guidance is thin; weight download is roughly 6–7 GB before runtime overhead. MIT license on the released code and weights; still verify provenance and consent for any cloned voice.

## Links

- Homepage: https://auk-project.github.io/
- GitHub: https://github.com/Tencent-Hunyuan/AuK
- Weights: https://huggingface.co/tencent/AuK
- Flash weights: https://huggingface.co/tencent/AuK-Flash
