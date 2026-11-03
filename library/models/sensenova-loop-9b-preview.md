---
name: sensenova-loop-9b-preview
title: SenseNova-Loop 9B Preview
summary: "Apache 2.0 looped ~9B multimodal model (~14.6B inference path) for long-context agent workflows."
category: models
does: multimodal
tags:
  - vision
  - llm
  - open
  - local
  - agentic
status: watch
run: local
license: open
vram_hint: vram-heavy
links:
  homepage: https://github.com/OpenSenseNova/SenseNova-Loop
  github: https://github.com/OpenSenseNova/SenseNova-Loop
  weights: https://huggingface.co/sensenova/SenseNova-Loop-9B-Preview
source: "youtube:IGICYFyP_EI"
added: 2026-11-04
---

# SenseNova-Loop 9B Preview

## What it is

**SenseNova-Loop-9B-Preview** is SenseTime’s **image-text-to-text** model for **long-context agent** work: a **looped decoder** reruns a shared layer block with **two-stream manifold hyper-connections (mHC)**, giving ~**14.6B inference-equivalent** compute from ~**9.7B** stored weights. Inputs are **text and images**; training context up to **262K** tokens. **Apache 2.0** on the GitHub repo; HF card tags match.

## Why keep

Try a **looped multimodal backbone** when you compare agent + document UI tasks against flat 9B VLMs, or pair with **[sensenova-6-8-code2video](/library/video/sensenova-6-8-code2video/)** for code-render video vs general looped agents.

## When to reach for it

You can host **custom Transformers code** (`trust_remote_code=True`) and want **English/Chinese** multimodal agents with vendor **LightLLM / vLLM** recipes (often **TP=4** at long context).

## Caveats

**Preview** name; weight repo may lack a LICENSE file on Hub (license text on **GitHub**). Looped inference is slower per token than a flat 9B despite similar storage. Benchmarks in README are **vendor demos** (motion design, PPT, data analysis).

## Links

- Homepage: https://github.com/OpenSenseNova/SenseNova-Loop
- GitHub: https://github.com/OpenSenseNova/SenseNova-Loop
- Weights: https://huggingface.co/sensenova/SenseNova-Loop-9B-Preview
