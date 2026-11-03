---
name: ming-3-0-tts
title: Ming 3.0 TTS
summary: "Ant inclusionAI MIT TTS with voice design, cloning, paralinguistic tags, and dense or MoE checkpoints."
category: models
does: speech
tags:
  - audio
  - open
  - local
  - moe
status: watch
run: local
license: open
vram_hint: consumer
links:
  homepage: https://github.com/inclusionAI/Ming-3.0-tts
  github: https://github.com/inclusionAI/Ming-3.0-tts
  weights: https://huggingface.co/inclusionAI/Ming-3.0-tts-0.5B
source: "youtube:EfyQOAJjx1A"
added: 2026-11-03
---

# Ming 3.0 TTS

## What it is

Ant Group **inclusionAI** open **text-to-speech** family (**MIT**). Two checkpoints: **Ming-3.0-tts-0.5B** (dense **0.5B** backbone, full bundle ~**3 GB** on HF) and **Ming-3.0-tts-tiny** (**7.9B** MoE, **~1.3B** active, sharded ~**18 GB**). Supports **multilingual** speech, **zero-shot cloning**, **natural-language voice design**, inline **paralinguistic tags** (emotion, breath, SFX), and **multi-speaker podcast** mode on the tiny variant. **Inference code lives on GitHub**; Hugging Face repos distribute weights only.

## Why keep

Compare when choosing a **local MIT-licensed TTS** with directed performance and dialect coverage alongside **Drama 3** or cloud voices.

## When to reach for it

Clone the GitHub repo, download **0.5B** for lighter GPU runs or **tiny** for podcast-style multi-speaker output. Use when you need **tagged expressive speech** or **prompt-defined voices** without a proprietary API.

## Caveats

**HF weights alone are not runnable**; follow the GitHub README for pinned deps (tiny uses a separate `requirements-tiny.txt` and Python **3.10.15**). Official **0.5B** tests cite **H800/H20**; consumer success varies. **Tiny** docs recommend **H800-class** hardware. Ant Group asks you **not to impersonate real people without consent** and to **label synthetic audio**. Roundup quality claims sit near **Fish Drama 3**, not a clear leap past it.

## Links

- GitHub: https://github.com/inclusionAI/Ming-3.0-tts
- Weights (0.5B): https://huggingface.co/inclusionAI/Ming-3.0-tts-0.5B
- Weights (tiny): https://huggingface.co/inclusionAI/Ming-3.0-tts-tiny
