---
name: haru-4b
title: Haru-4B
summary: "Apache-2.0 full-duplex speech model (~4B): listens and speaks in parallel without an ASR→LLM→TTS pipe; ~12 GB NVIDIA VRAM or Apple Silicon via MLX."
category: agents-memory
does: speech
tags:
  - audio
  - open
  - local
  - agentic
status: watch
run: local
license: open
vram_hint: "~12 GB NVIDIA; ~16 GB+ unified memory on Apple Silicon"
links:
  homepage: https://huggingface.co/higurashi-ai/Haru-4B
  github: https://github.com/higurashi-ai/Haru
  weights: https://huggingface.co/higurashi-ai/Haru-4B
source: "youtube:woWo80dPkVk"
added: 2026-11-17
---

# Haru-4B

## What it is

Haru-4B is an open-weight **full-duplex** spoken dialogue model from higurashi-ai: it takes user audio and produces model audio (and aligned text) in one stream, about **7.5 Hz** (~133 ms per frame), without a classic speech-to-text middle step. Implementations ship as **PyTorch** (`haru/`), **MLX** on Apple Silicon (`haru_mlx/`, including a native macOS app), and a **vLLM-Omni** plugin. Code and weights are **Apache-2.0**.

## Why keep

Compare or prototype **local duplex voice agents** when latency and paralinguistics (laughter, backchannels, barge-in) matter more than a thin API wrapper around an LLM plus TTS.

## When to reach for it

You want a self-hosted duplex voice stack on a **single ~12 GB GPU** or a Mac, and you can bring your own character card plus a short reference clip for voice control.

## Caveats

Released as a **pretrained base**, not a polished instruction-tuned assistant; the authors note a gap to human-like conversation. Mac app docs cite **16 GB+** unified memory as practical. The roundup transcript garbled the name as “Hauru for B.”

## Links

- Homepage: https://huggingface.co/higurashi-ai/Haru-4B
- GitHub: https://github.com/higurashi-ai/Haru
- Weights: https://huggingface.co/higurashi-ai/Haru-4B