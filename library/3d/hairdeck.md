---
name: hairdeck
title: HairDeck
summary: "Single-image to animation-ready stylized hair cards with spring-bone chains; GLB and VRM export."
category: 3d
does: 3d
tags:
  - 3d
  - image
  - open
  - local
  - vision
status: watch
run: local
license: mixed
vram_hint: consumer
links:
  homepage: https://huggingface.co/hyz317/HairDeck
  github: https://github.com/hyz317/HairDeck
  weights: https://huggingface.co/hyz317/HairDeck
source: "youtube:9RIK3A24ZxA"
added: 2026-11-08
---

# HairDeck

## What it is

HairDeck generates **stylized 3D hair cards** and **spring-bone animation chains** from one character image. Output includes per-card UVs and colors, GLB hair meshes, and VRM avatar export with head-turn preview video.

## Why keep

Use it when you need **anime or stylized avatar hair** from a reference still without hand-modeling cards and rigging bones, or when comparing single-image hair pipelines.

## When to reach for it

Reach for it for **VTuber/VRM workflows**, game avatar hair blocks, or local batch tests where card-based hair plus simple physics motion is enough.

## Caveats

HairDeck checkpoints are **Apache-2.0**, but the pipeline also pulls **SAM 3, DINOv3, BiRefNet, and Dora-VAE** with their own licenses and Hugging Face acceptance steps. Local setup targets **Python 3.12, PyTorch 2.8, and an NVIDIA CUDA GPU** (native Windows supported); **Apple Silicon runs on MPS** per upstream docs. Weights download is a multi-model bundle, not a single tiny file.

## Links

- Homepage: https://huggingface.co/hyz317/HairDeck
- GitHub: https://github.com/hyz317/HairDeck
- Weights: https://huggingface.co/hyz317/HairDeck
