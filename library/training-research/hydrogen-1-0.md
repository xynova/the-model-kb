---
name: hydrogen-1-0
title: Hydrogen-1.0
summary: "MIT-licensed 753B GLM-5.3 science post-train for chemistry, materials, and life-science research assistance."
category: training-research
does: research
tags:
  - llm
  - open
  - local
  - vram-heavy
status: watch
run: local
license: open
vram_hint: multi-gpu
links:
  homepage: https://blog.deepprinciple.com/introducing-hydrogen-1.0/
  weights: https://huggingface.co/deep-principle/Hydrogen-1.0
source: "youtube:9RIK3A24ZxA"
added: 2026-11-09
---

# Hydrogen-1.0

## What it is

**Hydrogen-1.0** (Deep Principle) is a **753B-parameter** scientific model **post-trained from GLM-5.3** on domain knowledge, **electronic lab records**, and **discovery trajectories** (including failed runs) across **chemistry, materials science, and life sciences**. Open **BF16** weights are **MIT-licensed** on Hugging Face (gated access); **NVFP4** and **MXFP4** builds target **NVIDIA Blackwell** and **AMD MI350** class hosts. **1M-token** context and up to **128K** output; the card documents **SGLang on 8 GPUs** (`--tp 8`).

## Why keep

Compare when you need an **open-weight science frontier** (experiment planning, reaction reasoning, structured lab assistance) rather than a general coding LLM.

## When to reach for it

You operate **multi-GPU inference** (or accept a hosted API such as **NeuroSpark** / **Mira**) and want to benchmark against other **science-tuned** open models on hard benchmarks like **FrontierScience-Research**.

## Caveats

Outputs are **not verified experimental results**; the model card requires expert review. **Full datasets and training recipes are not released.** Quantized checkpoints may behave differently from BF16. Gated HF download and **datacenter-scale** memory only; not a consumer run.

## Links

- Homepage: https://blog.deepprinciple.com/introducing-hydrogen-1.0/
- Weights: https://huggingface.co/deep-principle/Hydrogen-1.0
