---
name: sensenova-6-8-code2video
title: SenseNova 6.8 Code2Video
summary: "Apache 2.0 Qwen3.8-27B fine-tune that builds short videos by writing and executing render code."
category: video
does: video
tags:
  - video
  - open
  - local
  - agentic
  - llm
status: watch
run: local
license: open
vram_hint: vram-heavy
links:
  homepage: https://github.com/OpenSenseNova/SenseNova-6.8-Code2Video
  github: https://github.com/OpenSenseNova/SenseNova-6.8-Code2Video
  weights: https://huggingface.co/sensenova/SenseNova-6.8-Code2Video
source: "youtube:IGICYFyP_EI"
added: 2026-11-04
---

# SenseNova 6.8 Code2Video

## What it is

**SenseNova 6.8 Code2Video** is SenseTime’s open **code-driven video** model (~**27B**, fine-tuned from **Qwen3.8-27B**): it plans and writes **executable render code** for short motion-design clips instead of emitting pixels directly. The release includes **agent skill** (`code-video`), prompt templates, example **trajectories**, and deploy notes for **vLLM / SGLang / LightLLM** with very long context (512K target in docs). Weights on Hugging Face are tagged **Apache 2.0**.

## Why keep

Prototype **programmatic video** pipelines where you edit outputs as code, diff frames, or chain other gen tools from one agent session.

## When to reach for it

You already run a **27B-class agent harness** (the project documents **Claude Code**-style workflows) and want open weights for **motion graphics / explainers**, not a closed Runway-style API.

## Caveats

**~51 GB** weight storage; deploy examples assume **multi-GPU tensor parallel**, not a single consumer card. GitHub repo may lack a root LICENSE file even when the **HF card** states Apache 2.0. English and Chinese focus; quality is short-form motion design, not long cinematic diffusion.

## Links

- Homepage: https://github.com/OpenSenseNova/SenseNova-6.8-Code2Video
- GitHub: https://github.com/OpenSenseNova/SenseNova-6.8-Code2Video
- Weights: https://huggingface.co/sensenova/SenseNova-6.8-Code2Video
