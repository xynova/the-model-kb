---
name: minicpm5-2b
title: MiniCPM5-2B
summary: "OpenBMB dense ~**2B** on-device LLM (MiniCPM5 series) with long context, tool use, and agent-oriented post-training (RL + OPD)."
category: models
does: llm
tags:
  - llm
  - open
  - local
  - quantized
  - agentic
status: watch
run: local
license: open
vram_hint: consumer
links:
  homepage: https://github.com/OpenBMB/MiniCPM
  github: https://github.com/OpenBMB/MiniCPM
  weights: https://huggingface.co/openbmb/MiniCPM5-2B
media:
  cover: models/minicpm5-2b/media/cover.jpg
  gallery:
    - models/minicpm5-2b/media/g01.jpg
source: "youtube:nZYJdwM-_nI"
added: 2026-09-16
---

# MiniCPM5-2B

## What it is

OpenBMB dense ~**2B** on-device LLM (MiniCPM5 series) with long context, tool use, and agent-oriented post-training (RL + OPD). Apache-2.0 weights plus GGUF, MLX, and GPTQ builds; training datasets also published. Roundup “mini CPM 52B” ASR maps to this **MiniCPM5-2B** release.

| | |
|:---:|:---:|
| ![minicpm5-2b 1](minicpm5-2b/media/thumb-cover.jpg) | ![minicpm5-2b 2](minicpm5-2b/media/thumb-g01.jpg) |

## Why keep

Compare tiny local models when edge footprint and small-size coding/math/agent scores matter more than frontier scale.

## When to reach for it

Phone, laptop, or low-VRAM boxes via transformers / llama.cpp / MLX, or bake-offs against other sub-4B opens.

## Caveats

Vendor harness SOTA claims need independent re-checks; Artificial Analysis sub-4B ranking is a separate public signal. Still a small model: expect limits on hard long-horizon agent work versus larger MoEs.

## Links

- GitHub: https://github.com/OpenBMB/MiniCPM
- Weights: https://huggingface.co/openbmb/MiniCPM5-2B
- GGUF: https://huggingface.co/openbmb/MiniCPM5-2B-GGUF
- Demo: https://huggingface.co/spaces/openbmb/MiniCPM5-2B-Demo
