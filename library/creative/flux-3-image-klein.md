---
name: flux-3-image-klein
title: FLUX 3 Image Klein
summary: "BFL 5.6B open-weight image model for instruction edit, T2I, references, and bbox layout."
category: creative
does: image
tags:
  - image
  - open
  - local
  - api
  - mixed
status: watch
run: local
license: mixed
vram_hint: consumer
links:
  homepage: https://bfl.ai/models/flux-3-image-klein
  weights: https://huggingface.co/black-forest-labs/flux-3-image-klein-base
media:
  cover: creative/flux-3-image-klein/media/cover.jpg
  gallery:
    - creative/flux-3-image-klein/media/g01.jpg
    - creative/flux-3-image-klein/media/g02.jpg
    - creative/flux-3-image-klein/media/g03.jpg
source: "youtube:tDWtcyeNNew"
added: 2026-11-04
---

# FLUX 3 Image Klein

## What it is

BFL 5.6B open-weight image model for instruction edit, T2I, references, and bbox layout. Base and distilled checkpoints on Hugging Face; optional BFL API to 4K. Uses a separate prompt upsampler (`flux-3-image-klein-pup`) in local Diffusers or ComfyUI pipelines.

| | |
|:---:|:---:|
| ![flux-3-image-klein 1](flux-3-image-klein/media/thumb-cover.jpg) | ![flux-3-image-klein 2](flux-3-image-klein/media/thumb-g01.jpg) |
| ![flux-3-image-klein 3](flux-3-image-klein/media/thumb-g02.jpg) | ![flux-3-image-klein 4](flux-3-image-klein/media/thumb-g03.jpg) |

## Why keep

Compare or fine-tune a local open editor when you need bbox layout, multi-reference compose, or 4K-capable weights on one consumer GPU.

## When to reach for it

Self-host `black-forest-labs/flux-3-image-klein-base` (50-step, LoRA training) or the distilled Klein for 16-step edits; or call the BFL API with the same prompt and bbox format.

## Caveats

FLUX Kommunity License: non-commercial use always; commercial **output** only for Qualifying Users under $5M annual revenue; commercial fine-tuning or weight redistribution needs a BFL commercial license. Roundup ComfyUI tests reported weak faces and anatomy on the base checkpoint; community LoRAs may help. Weights are gated on Hugging Face.

## Links

- Homepage: https://bfl.ai/models/flux-3-image-klein
- Weights: https://huggingface.co/black-forest-labs/flux-3-image-klein-base
- Docs: https://docs.bfl.ml/flux_3/flux3_image_klein
