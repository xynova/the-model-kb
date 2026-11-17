---
name: sophon-realtime
title: Sophon-Realtime
summary: "Adventists AI Apache-2.0 full-duplex voice family: listens while you talk, paralinguistic ears, expressive speech; Tag/Flash/E2E (32B or 4B) via vLLM-Omni or MLX."
category: agents-memory
does: speech
tags:
  - audio
  - open
  - local
  - api
  - agentic
status: watch
run: local
license: mixed
vram_hint: "32B builds ~80 GB GPU; 4B Flash ~24 GB, Tag/E2E ~32 GB (vendor docs)"
links:
  homepage: https://adventists.ai/
  github: https://github.com/adventists-ai/sophon-realtime
  weights: https://huggingface.co/adventists-ai/Sophon-Realtime-v1
source: "youtube:woWo80dPkVk"
added: 2026-11-17
---

# Sophon-Realtime

## What it is

**Sophon-Realtime** (Adventists AI) is an open **full-duplex** speech model family: it listens while the user talks, models **tone and paralinguistics** (DuplexJev), and answers with speech that can laugh, sigh, or backchannel without a separate ASR→LLM→TTS stack. Three builds share one **thinker** brain: **Tag** (recommended), **Flash** (faster, more concurrent users), and **E2E** (preview). Serve with **vLLM-Omni** on Linux GPUs, **MLX** on Apple Silicon, or a hosted **OpenAI Realtime-compatible** trial API. Code and main weights are **Apache-2.0**; **Sophon-Realtime 4B** uses a smaller brain with stricter terms on some bundled components (see repo **NOTICE**).

## Why keep

Compare **open duplex agents** when you care about barge-in, emotional cues, zero-shot cloning, and GPT-Live-style integration, alongside Haru-4B and hosted Fish Live 3.

## When to reach for it

You can host a **large GPU** (32B line) or a **24–32 GB** card (4B Flash/Tag/E2E), or you want to prototype on Adventists’ API before committing hardware.

## Caveats

32B **Tag** weights are **~42 GB** on first download; realistic deployment is **~80 GB** class GPUs for the full family. **4B** trades world knowledge for fit on consumer cards. Third-party and **DuplexJev-4B** licence terms may be **non-commercial** in places; read **NOTICE** before product use. Roundup transcript says “Sofon Realtime.”

## Links

- Homepage: https://adventists.ai/
- GitHub: https://github.com/adventists-ai/sophon-realtime
- Weights hub: https://huggingface.co/adventists-ai/Sophon-Realtime-v1