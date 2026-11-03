---
name: nanbeige-4.3-33b
title: Nanbeige4.3-33B
summary: "Apache 2.0 ~33B agentic coding and office cowork LLM with 512K context and XML tool calling."
category: models
does: llm
tags:
  - llm
  - open
  - local
  - agentic
status: watch
run: local
license: open
vram_hint: vram-heavy
links:
  homepage: https://huggingface.co/Nanbeige/Nanbeige4.3-33B
  weights: https://huggingface.co/Nanbeige/Nanbeige4.3-33B
source: "youtube:IGICYFyP_EI"
added: 2026-11-04
---

# Nanbeige4.3-33B

## What it is

**Nanbeige4.3-33B** is Nanbeige LLM Lab’s (~BOSS Zhipin) dense **~33B** chat model for **agentic coding and office cowork**: **512K** context, **thinking** and **tool-call** chat templates (XML tools recommended), post-trained from **Nanbeige4.3-33B-Base**. Weights are ungated on Hugging Face (~**66 GB** BF16). Card and metadata tag **Apache 2.0** (no standalone LICENSE file in the weight repo at filing time).

## Why keep

Host or benchmark a **mid-size open agent teacher** when you compare terminal, SWE, and planning harnesses against Qwen 3.8 / Gemma 4 class models, or distill down to smaller Nanbeige checkpoints later.

## When to reach for it

You can serve **tensor-parallel ~4 GPU** (official **SGLang** / **vLLM** recipes) and want **`preserve_thinking=True`** for multi-turn tool and code-agent runs. Pair with **[nanbeige-4.2-3b](/library/models/nanbeige-4.2-3b/)** when the job is **local small agent** rather than flagship teacher quality.

## Caveats

Benchmark tables on the card are **publisher-run** with mixed harnesses; treat rivals’ scores as directional. **English and Chinese** focus. Confirm Apache grant from the **HF card** if you need a LICENSE file on disk.

## Links

- Homepage: https://huggingface.co/Nanbeige/Nanbeige4.3-33B
- Weights: https://huggingface.co/Nanbeige/Nanbeige4.3-33B
- Base checkpoint: https://huggingface.co/Nanbeige/Nanbeige4.3-33B-Base
