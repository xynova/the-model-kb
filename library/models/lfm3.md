---
name: lfm3
title: LFM3
summary: "Liquid AI 2B-class multimodal family (text, VL, audio, omni, decision) on one backbone for edge fine-tune."
category: models
does: multimodal
tags:
  - llm
  - vision
  - audio
  - open
  - local
  - quantized
  - hardware
status: watch
run: local
license: mixed
vram_hint: consumer
links:
  homepage: https://www.liquid.ai/blog/lfm3
  weights: https://huggingface.co/collections/LiquidAI/lfm3
media:
  cover: models/lfm3/media/cover.jpg
source: "youtube:IGICYFyP_EI"
added: 2026-11-04
---

# LFM3

## What it is

**LFM3** is Liquid AI’s **~2B** **Liquid Foundation Model** generation: **text**, **vision-language**, **audio**, **omni**, and **decision** checkpoints built from **LFM3-1.7B-Base**, with **GGUF**, **ONNX**, **MLX**, and **FP8/NVFP4** exports for edge devices. The blog cites a tri-hybrid architecture and large-scale pretraining aimed at **specialized fine-tunes on one GPU**. Weights live in the **LFM3 Hugging Face collection**; Liquid also offers playground and API hosting.

| | |
|:---:|:---:|
| ![lfm3 1](lfm3/media/thumb-cover.jpg) | |

## Why keep

Pick a **small multimodal base** when you compare on-device agents, VL, or audio helpers without jumping to 30B stacks.

## When to reach for it

You need **Apple Silicon / quant-friendly** artifacts (~**1 GB** 4-bit text class in Liquid’s marketing) and plan **task-specific fine-tunes** rather than a general cloud-only API.

## Caveats

**LFM Open License v1.0** (not plain Apache 2.0): **free commercial use only below ~$10M annual revenue**; above that, contact Liquid for a paid license. Research and qualified nonprofits have broader terms. Also served via **Liquid AI API** if you do not self-host.

## Links

- Homepage: https://www.liquid.ai/blog/lfm3
- License summary: https://www.liquid.ai/lfm-license
- Weights collection: https://huggingface.co/collections/LiquidAI/lfm3
