---
name: sage-a1-5
title: Sage-A1.5
summary: "InternScience 35B MoE agentic VLM for long-horizon scientific tasks, orchestration, and self-verification."
category: agents-memory
does: research
tags:
  - llm
  - vision
  - open
  - local
  - moe
  - agentic
status: watch
run: local
license: open
vram_hint: multi-gpu
links:
  homepage: https://alpha-innovator.github.io/Sage-A1.5/
  github: https://github.com/InternScience/Agents-A1
  weights: https://huggingface.co/InternScience/Sage-A1.5
media:
  cover: agents-memory/sage-a1-5/media/cover.jpg
source: "youtube:EfyQOAJjx1A"
added: 2026-11-04
---

# Sage-A1.5

## What it is

**Sage-A1.5** (roundup: **Sage A1.5** / **Alpha Innovator**): **~35B MoE** **agentic VLM** from **Shanghai AI Laboratory (InternScience)**, built on **Qwen3.6** and **Agents-A1**. Targets **long-horizon scientific work** with **orchestration**, **verification**, and **open discovery** (literature search, cited reports, figures, molecule/reaction/brain-data demos in the release materials). Weights on Hugging Face (**~72 GB** BF16 card storage; **Q4 GGUF** sibling ~**23 GB**). **Apache 2.0** on the model card; serve via **vLLM** or **SGLang** (project recipes cite **256K** context).

| | |
|:---:|:---:|
| ![sage-a1-5 1](sage-a1-5/media/thumb-cover.jpg) | |

## Why keep

Run or benchmark a **local scientific agent** when comparing MoE agent stacks against general coding models.

## When to reach for it

Self-hosted **OpenAI-compatible** agent endpoints for multi-step research pipelines where you can afford **multi-GPU** or a **quant** build; pair with the project **evaluation** framework when reproducing benchmark claims.

## Caveats

**VRAM-heavy** at full precision; verify quant quality for your domain. Hub org is **InternScience**; marketing also uses **Alpha Innovator** naming. Related code tree is **Agents-A1** (prior A1 line); confirm the exact repo/scripts for A1.5 before production. Roundup benchmarks are **vendor-reported**.

## Links

- Project: https://alpha-innovator.github.io/Sage-A1.5/
- GitHub (Agents-A1): https://github.com/InternScience/Agents-A1
- Weights: https://huggingface.co/InternScience/Sage-A1.5
