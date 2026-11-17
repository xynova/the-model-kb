---
name: zonos-3-mini
title: ZONOS3-Mini
summary: "Zyphra Apache-2.0 TTS: 800M total / ~100M active MoE, zero-shot cloning, steered emotion and delivery; local (GGUF, MLX, WebGPU) or Zyphra Cloud."
category: models
does: speech
tags:
  - audio
  - open
  - local
  - api
  - moe
status: watch
run: local
license: open
vram_hint: "~9 GB CUDA (full PyTorch stack); GGUF / WebGPU paths lower"
links:
  homepage: https://www.zyphra.com/our-work/zonos3-mini
  github: https://github.com/Zyphra/ZONOS3
  weights: https://huggingface.co/Zyphra/ZONOS3-Mini
source: "youtube:woWo80dPkVk"
added: 2026-11-17
---

# ZONOS3-Mini

## What it is

**ZONOS3-Mini** is Zyphra’s compact **open-weight TTS**: **800M** total parameters with about **100M active** per token (latent **MoE**), **zero-shot voice cloning** from a reference clip, and bucketed controls for loudness, bandwidth, and speaking rate plus character shaping (age, accent, texture). Weights are **Apache-2.0**; runtimes cover **SGLang-Omni**, **MLX**, **ggml** (CPU/CUDA/Metal/Vulkan), and **WebGPU** in the browser, with optional **Zyphra Cloud** serving.

## Why keep

Pick or benchmark **steerable local TTS** when you need cloning without a giant codec stack (contrast with ZONOS2-scale models) and want both server and edge paths in one family.

## When to reach for it

You are wiring **narration, assistants, or games** that need multilingual speech, clone-from-sample voices, or fine-grained delivery control, and you can host ~consumer GPU weights or use GGUF/WebGPU builds.

## Caveats

Full **PyTorch/CUDA** serving pulls multiple weight bundles (~**9 GB** in upstream docs). **MIT** runtimes vs **Apache-2.0** weights: read each repo’s license file before commercial ship. The roundup misheard the vendor as “Zyfra.”

## Links

- Homepage: https://www.zyphra.com/our-work/zonos3-mini
- GitHub (runtimes): https://github.com/Zyphra/ZONOS3
- Weights: https://huggingface.co/Zyphra/ZONOS3-Mini