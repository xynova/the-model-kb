---
name: nanbeige-4.2-3b
title: Nanbeige4.2-3B
summary: "3B non-embedding agentic LLM with a two-pass Looped Transformer; strong small-model tool and SWE scores."
category: models
does: llm
tags:
  - llm
  - open
  - local
  - agentic
  - quantized
status: watch
run: local
license: open
vram_hint: consumer
links:
  homepage: https://arxiv.org/abs/2607.22083
  github: https://github.com/Nanbeige/llama.cpp
  weights: https://huggingface.co/Nanbeige/Nanbeige4.2-3B
source: "youtube:fMaPaK4gcNE"
added: 2026-11-04
---

# Nanbeige4.2-3B

## What it is

**Nanbeige4.2-3B** is a compact **agentic** chat model (~**3B** non-embedding params, ~**4B** total with embeddings) built on a **Looped Transformer**: the same decoder stack runs **two passes** for extra effective depth without a second weight set. Apache-2.0 weights on Hugging Face, with official forks for **vLLM**, **SGLang**, **llama.cpp**, and **Ollama** under the Nanbeige org.

## Why keep

Catalog a **consumer-VRAM** local agent when you want to compare small open models on **tool use**, **office**, and **SWE** benches instead of defaulting to Qwen3.5-9B-class sizes.

## When to reach for it

Local OpenClaw-style assistants, MCP-heavy workflows, or Terminal-Bench / SWE-Bench spot checks on one GPU; use **thinking** mode and `preserve_thinking=true` for multi-turn tool runs per the model card.

## Caveats

Needs **`trust_remote_code`** and a pinned **transformers** version on the HF path; looped KV behavior differs from plain Llama stacks. Published tables are author-run with specific scaffolds (OpenHands, Terminus 2, in-house office suites). Long-context (256K) is supported but not free at runtime.

## Links

- Paper: https://arxiv.org/abs/2607.22083
- Weights: https://huggingface.co/Nanbeige/Nanbeige4.2-3B (plus FP8/GPTQ variants in the Nanbeige4.2-3B collection)
- Reference inference forks: https://github.com/Nanbeige/vllm (branch `nanbeige42`), https://github.com/Nanbeige/sglang, https://github.com/Nanbeige/llama.cpp
