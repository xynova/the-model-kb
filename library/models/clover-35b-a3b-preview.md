---
name: clover-35b-a3b-preview
title: Clover 35B-A3B Preview
summary: "Apache 2.0 35B MoE diffusion language model (~3B active) with block diffusion and 128K context."
category: models
does: llm
tags:
  - llm
  - open
  - local
  - moe
  - agentic
status: watch
run: local
license: open
vram_hint: vram-heavy
links:
  homepage: https://diffuspace.ai/blog/clover-35b/
  github: https://github.com/DiffuSpace/Clover-35B-A3B-Preview
  weights: https://huggingface.co/Diffuspace/Clover-35B-A3B-Preview
source: "youtube:IGICYFyP_EI"
added: 2026-11-04
---

# Clover 35B-A3B Preview

## What it is

**Clover-35B-A3B-Preview** is DiffuSpace’s first **diffusion language model (dLLM)**: a **35B MoE** text stack (~**3B active** per token) derived from **Qwen3.6-35B-A3B**, using **block diffusion** with bidirectional attention inside each block and up to **128K** context. Inference paths include **Transformers**, **SGLang**, and **llama.cpp** per the release docs. **Apache 2.0** on weights and GitHub.

## Why keep

Experiment with **non-autoregressive text generation** (parallel token blocks, higher throughput) while staying in open weights comparable to Qwen3.6-35B-A3B-class agents.

## When to reach for it

You already run **MoE agent harnesses** and want to compare **dLLM decoding** on coding and tool benchmarks, or fine-tune on Clover’s text-only checkpoint layout.

## Caveats

**Preview** release; diffusion LLM tooling is newer than causal vLLM defaults. Derived from Qwen assets with documented **ATTRIBUTIONS** changes (vision/MTP stripped). VRAM still **datacenter-class** despite low active params. Launch blog metrics are **vendor-reported**.

## Links

- Homepage: https://diffuspace.ai/blog/clover-35b/
- GitHub: https://github.com/DiffuSpace/Clover-35B-A3B-Preview
- Weights: https://huggingface.co/Diffuspace/Clover-35B-A3B-Preview
