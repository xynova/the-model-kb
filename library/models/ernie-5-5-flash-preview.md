---
name: ernie-5-5-flash-preview
title: ERNIE 5.5 Flash Preview
summary: "Baidu 280B/14B-active MoE text model (1M context): coding, agentic 3D/game demos via API; Preview weights promised Apache-2.0 (Qianfan / OpenRouter)."
category: models
does: llm
tags:
  - llm
  - api
  - moe
  - open
status: watch
run: api
license: mixed
vram_hint: "Preview weights not on HF yet; Base FP8 ~288 GB class if self-hosting Base"
links:
  homepage: https://ernie.baidu.com/release/ernie-5-5-flash-preview/
  weights: https://huggingface.co/baidu/ERNIE-5.5-Flash-Base
source: "youtube:woWo80dPkVk"
added: 2026-11-17
---

# ERNIE 5.5 Flash Preview

## What it is

**ERNIE 5.5 Flash Preview** is Baidu’s first **post-trained** build on **ERNIE 5.5 Flash Base**: a sparse **MoE** with **280B** total parameters and **14B active** per token, **1M-token** context, **text-only** I/O. Baidu pitches long **agent** jobs: game code (Unreal/Unity/Godot), **FreeCAD/Blender** scripting, knowledge work, and video pipelines driven by tools. **Preview** is available on **Qianfan** and third-party routers (e.g. OpenRouter); Baidu says **Preview weights** will ship under **Apache 2.0** (announced before **25 Nov 2026**). The already-released **Base** checkpoint on Hugging Face uses a **separate Baidu model license**, not Apache.

## Why keep

Track a **large MoE agent/coding** option and compare API cost vs DeepSeek/GLM flash tiers while waiting for downloadable Preview weights.

## When to reach for it

You want **million-token text agents** on a hosted API and can tolerate an **intermediate preview** quality bar; not for native multimodal inputs on the model itself.

## Caveats

Baidu labels Preview an **incomplete post-train checkpoint**. **Self-host Preview** was not public at filing time; local play today means **Base** plus your own SFT, with heavy hardware and custom license terms. Demos in the release are **vendor-run** agent harnesses. **Base** `weights` link is not the Preview checkpoint.

## Links

- Release: https://ernie.baidu.com/release/ernie-5-5-flash-preview/
- Base weights (pre-train): https://huggingface.co/baidu/ERNIE-5.5-Flash-Base