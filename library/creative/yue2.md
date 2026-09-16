---
name: yue2
title: YuE2
summary: "m-a-p open music model (~3–4B) that plans an editable ABC score (melody/chords) before synthesizing 48 kHz stereo vocals and accompaniment."
category: creative
does: music
tags:
  - audio
  - open
  - local
  - mixed
  - agentic
status: watch
run: local
license: mixed
vram_hint: consumer
links:
  homepage: https://map-yue2.github.io/
  github: https://github.com/multimodal-art-projection/YuE
  weights: https://huggingface.co/m-a-p/YuE2-3B
media:
  cover: creative/yue2/media/cover.jpg
  gallery:
    - creative/yue2/media/g01.jpg
source: "youtube:nZYJdwM-_nI"
added: 2026-09-16
---

# YuE2

## What it is

m-a-p open music model (~3–4B) that plans an editable ABC score (melody/chords) before synthesizing 48 kHz stereo vocals and accompaniment. Supports covers, agentic score edits, and local HF inference. WildSongBench claims competitive or better SongBench averages versus Suno v5/v6 under their selection protocols (including best-of-8).

| | |
|:---:|:---:|
| ![yue2 1](yue2/media/thumb-cover.jpg) | ![yue2 2](yue2/media/thumb-g01.jpg) |

## Why keep

Local alternative to Suno when you want score-level control and open weights instead of a closed SaaS.

## When to reach for it

Consumer **24 GB** NVIDIA setups for lyrics-to-song, covers via SheetSage2 scores, or iterative agent edits on the ABC plan.

## Caveats

Weights are **CC-BY-NC-4.0** (non-commercial). Roundup “UA2” / “7.3 GB” ASR maps to this YuE2 release; peak VRAM on a 4090 is about **11 GiB** for one song, with a 24 GB GPU recommended. Benchmark wins depend on candidate-selection protocol; listen before trusting leaderboard gaps.

## Links

- Homepage: https://map-yue2.github.io/
- GitHub: https://github.com/multimodal-art-projection/YuE
- Weights: https://huggingface.co/m-a-p/YuE2-3B
