---
name: deepseek-v4.1-flash
title: DeepSeek-V4.1-Flash
summary: "DeepSeek multimodal MoE with a **552B** backbone, **8B** active on prefill and **16B** on decode, native image+text input, and up to **1M** context."
category: models
does: multimodal
tags:
  - llm
  - vision
  - open
  - local
  - moe
  - api
  - vram-heavy
status: watch
run: local
license: open
vram_hint: multi-gpu
links:
  homepage: https://api-docs.deepseek.com/news/news260910
  weights: https://huggingface.co/deepseek-ai/DeepSeek-V4.1-Flash
media:
  cover: models/deepseek-v4.1-flash/media/cover.jpg
  gallery:
    - models/deepseek-v4.1-flash/media/g01.jpg
    - models/deepseek-v4.1-flash/media/g02.jpg
    - models/deepseek-v4.1-flash/media/g03.jpg
source: "youtube:nZYJdwM-_nI"
added: 2026-09-16
---

# DeepSeek-V4.1-Flash

## What it is

DeepSeek multimodal MoE with a **552B** backbone, **8B** active on prefill and **16B** on decode, native image+text input, and up to **1M** context. Causal encoder-decoder design with aggressive KV compression. API id `deepseek-flash`; MIT weights on Hugging Face.

| | |
|:---:|:---:|
| ![deepseek-v4.1-flash 1](deepseek-v4.1-flash/media/thumb-cover.jpg) | ![deepseek-v4.1-flash 2](deepseek-v4.1-flash/media/thumb-g01.jpg) |
| ![deepseek-v4.1-flash 3](deepseek-v4.1-flash/media/thumb-g02.jpg) | ![deepseek-v4.1-flash 4](deepseek-v4.1-flash/media/thumb-g03.jpg) |

## Why keep

Compare frontier open multimodal models for agent and coding work, especially when cheap API speed matters or you evaluate long-context / KV-efficient stacks.

## When to reach for it

DeepSeek API bake-offs, or multi-GPU self-host (reference and vLLM/SGLang recipes target roughly **8×** high-memory GPUs; checkpoint ~**510 GB** on disk).

## Caveats

“Flash” and “8B active” do not mean consumer-GPU local. Full weights are already FP8/FP4 mixed; community quants may appear later but official local is multi-GPU class. Confirm current API pricing and routing (legacy V4-Flash / V4-Pro aliases may redirect).

## Links

- Homepage / API news: https://api-docs.deepseek.com/news/news260910
- Weights: https://huggingface.co/deepseek-ai/DeepSeek-V4.1-Flash
