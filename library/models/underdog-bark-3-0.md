---
name: underdog-bark-3-0
title: Underdog Bark 3.0
summary: "Conway Research English ASR in a 163 MB packed download, competitive with much larger Whisper and Parakeet models."
category: models
does: speech
tags:
  - audio
  - open
  - local
  - quantized
status: watch
run: local
license: open
vram_hint: consumer
links:
  homepage: https://huggingface.co/ConwayResearch/Underdog-Bark-3.0
  weights: https://huggingface.co/ConwayResearch/Underdog-Bark-3.0
source: "youtube:EfyQOAJjx1A"
added: 2026-11-03
---

# Underdog Bark 3.0

## What it is

English **automatic speech recognition** from Conway Research (the roundup’s “Bark 3”). A **0.6B** model shipped as **`model.bark`** (~163 MB download). The model card reports **~4.94% average WER** across seven Open ASR Leaderboard English sets, in the same band as NVIDIA **Parakeet TDT 0.6B v3** at a much smaller download. Run locally via **`bark_loader.py`** (PyTorch on CUDA, CPU, or MPS) or **`bark_mlx.py`** on Apple Silicon; Core ML and Windows helper scripts ship on the same Hugging Face repo.

## Why keep

Compare and run when you need **on-device English transcription** without pulling multi-gigabyte Whisper weights, or when sizing a local STT slot on a laptop or phone pipeline.

## When to reach for it

Batch or CLI transcription of English audio on consumer hardware; MLX path on Mac when you want Parakeet-class accuracy from a tiny download. Step up to **Underdog Bark 3.0 Plus** or **Pro** on the same org when you need lower WER; use **Underdog Bark 3.0 Clean** when you want filler stripped for dictation-style text.

## Caveats

**English only.** Derived from NVIDIA Parakeet and Moondream checkpoints; **CC-BY-4.0** with attribution in the repo **NOTICE**. Weights are small but **RAM grows with clip length** on one-pass MLX decode (the card warns about long files on 8–16 GB Macs; split long recordings). Bark rows in the card were scored with the leaderboard code but were **not yet listed on the public leaderboard** at publish time. Not the Suno **Bark** TTS model.

## Links

- Homepage: https://huggingface.co/ConwayResearch/Underdog-Bark-3.0
- Weights: https://huggingface.co/ConwayResearch/Underdog-Bark-3.0
- Org: https://conway.tech
