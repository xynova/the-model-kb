---
name: memt-1-0
title: MemT 1.0
summary: "MemTensor memory-native LLMs with fixed-size persistent state across turns; Nano and Flash weight tiers on Hugging Face."
category: agents-memory
does: memory
tags:
  - memory
  - open
  - local
  - moe
  - agentic
status: watch
run: local
license: mixed
vram_hint: vram-heavy
links:
  homepage: https://github.com/MemTensor/MemT-1.0
  github: https://github.com/MemTensor/MemT-1.0
  weights: https://huggingface.co/MemTensor/MemT-1.0-Nano
source: "youtube:EfyQOAJjx1A"
added: 2026-11-04
---

# MemT 1.0

## What it is

**MemT-1.0** from **MemTensor**: **memory-native** causal LMs that keep a **fixed-size persistent memory state** across turns instead of re-reading full chat history (project cites ~**66×** lower storage vs full **128K** KV at their comparison point). Two checkpoints: **MemT-1.0-Nano** (**36B** total, **3B** active, Qwen3.6 MoE backbone) and **MemT-1.0-Flash** (**288B** total, **17B** active). Weights on Hugging Face with **`trust_remote_code`**; **SGLang** integration lives in the GitHub repo.

## Why keep

Benchmark or wire **native long-horizon memory** when comparing against RAG, external memory stores, or stuffing transcripts into context.

## When to reach for it

Agent or companion pipelines where you can host **large downloads** (roundup cited ~**70 GB** Nano and ~**576 GB** Flash) and run **bf16** inference with the published Transformers/SGLang paths.

## Caveats

**License is split:** Nano weights **Apache 2.0**, Flash weights **MIT**, but **MemTensor repo software** (including SGLang integration) is **PolyForm Noncommercial** unless you negotiate commercial terms. Paper is **CC BY-NC-SA**. **Flash** tier is **datacenter-scale** for most users. Roundup kitten demo is illustrative; validate on your memory benchmarks.

## Links

- GitHub: https://github.com/MemTensor/MemT-1.0
- Weights (Nano): https://huggingface.co/MemTensor/MemT-1.0-Nano
- Weights (Flash): https://huggingface.co/MemTensor/MemT-1.0-Flash
- Collection: https://huggingface.co/collections/MemTensor/memt-1
