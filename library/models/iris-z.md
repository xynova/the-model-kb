---
name: iris-z
title: Iris-Z
summary: "Speridlabs unified RGB+depth generator (9.1B or 32.2B): text→image/RGBD, image→depth, depth→image; Apache-2.0 SDK, FLUX.2 autoencoder gated."
category: models
does: image
tags:
  - image
  - vision
  - open
  - local
  - vram-heavy
status: watch
run: local
license: mixed
vram_hint: "~33 GB (small) or ~112 GB (big) for full tasks; depth-only can omit text encoder"
links:
  homepage: https://speridlabs.com/research/iris-z
  github: https://github.com/speridlabs/iris-z-sdk
  weights: https://huggingface.co/speridlabs/iris-z
source: "youtube:woWo80dPkVk"
added: 2026-11-17
---

# Iris-Z

## What it is

**Iris-Z** (Speridlabs) is one **diffusion transformer** for **images and relative depth**: same weights do **text→image**, **text→RGBD**, **image→depth**, and **depth-conditioned image** generation. Depth is **relative log distance** (not metric meters), encoded in the shared token sequence with RGB, not a separate depth head. Releases: **Iris-Z-small** (9.1B) and **Iris-Z-big** (32.2B). SDK and sampling code are **Apache-2.0**; weights live on Hugging Face. Inference pulls the **FLUX.2** autoencoder (**gated** on Hugging Face).

## Why keep

Use when you want **one local stack** for both **T2I** and **monocular depth** (or depth-guided repainting) without chaining a specialist depth net beside a generator.

## When to reach for it

You can host **~33 GB+** VRAM (small) or a **100 GB+** box (big), accept **FLUX.2** license gating, and need depth benchmarks competitive with **MoGe-3**-class specialists per the authors’ tables.

## Caveats

**Relative** depth only; multi-view / video consistency not the focus. **CUDA 13** stack in upstream SDK docs. Roundup VRAM figures (~18/64 GB) understates official **~33 / ~112 GB** “all tasks” estimates. Name on air: **Iris Z**.

## Links

- Research: https://speridlabs.com/research/iris-z
- Paper: https://arxiv.org/abs/2611.16010
- SDK: https://github.com/speridlabs/iris-z-sdk
- Weights: https://huggingface.co/speridlabs/iris-z