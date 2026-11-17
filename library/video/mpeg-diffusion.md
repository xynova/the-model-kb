---
name: mpeg-diffusion
title: MPEG-Diffusion
summary: "Mixed-resolution video DiT method: self-selected HR subject region, dual-resolution attention, blended decode; ~4.8× Wan2.1 14B speed at ~30% tokens."
category: video
does: video
tags:
  - video
status: watch
run: local
license: unknown
vram_hint: "Benchmarks on Wan2.1 1.3B/14B; ~30% full-res tokens"
links:
  homepage: https://dennismenn.github.io/mpeg-diffusion/
source: "youtube:woWo80dPkVk"
added: 2026-11-17
---

# MPEG-Diffusion

## What it is

**MPEG-Diffusion** (UT Austin + Roblox) is a **mixed-resolution** recipe for **video DiTs**: the model **picks the high-res region** from intermediate DiT features (no preset mask), runs **dual-resolution attention** (separate K/V per resolution), and **merges HR/LR latents** before VAE decode to avoid seams. On **Wan2.1**, authors report up to **~4.8×** denoising speed on **14B** (~**3.8×** on **1.3B**) while using about **30%** of full-resolution tokens, with user-study quality near full res. Also demonstrated on **FLUX.2** images. **Paper:** [arXiv:2611.12069](https://arxiv.org/abs/2611.12069).

## Why keep

Reach for it when you are **accelerating Wan-class video DiTs** and want **automatic saliency** for HR tokens instead of hand-drawn foveation masks (contrast with Phase-Aligned RoPE / Foveated Diffusion baselines in the paper).

## When to reach for it

You run or fine-tune **Wan2.1** (or similar DiTs) and need a cited **mixed-res** stack with self-selected regions once **official code** lands or you port from the paper.

## Caveats

Authors said **code was under internal review** (Nov 2026); **no official public repo** yet. Third-party ports exist but are **not** the paper release. Fast-moving backgrounds can show **HR/LR boundaries**. Roundup name matches **MPEG Diffusion**.

## Links

- Project page: https://dennismenn.github.io/mpeg-diffusion/
- Paper: https://arxiv.org/abs/2611.12069