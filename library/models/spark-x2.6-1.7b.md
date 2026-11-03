---
name: spark-x2.6-1.7b
title: Spark-X2.6-1.7B
summary: "Apache 2.0 ~1.7B on-device agent LLM with thinking levels, tools, and long-context serving paths."
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
  homepage: https://huggingface.co/XHToken/Spark-X2.6-1.7B
  github: https://github.com/XHToken/Spark-X2.6
  weights: https://huggingface.co/XHToken/Spark-X2.6-1.7B
source: "youtube:IGICYFyP_EI"
added: 2026-11-04
---

# Spark-X2.6-1.7B

## What it is

**Spark-X2.6-1.7B** is XHToken’s compact **agentic LLM** for **on-device** use: same architecture as Spark-X2.5-1.7B with stronger post-training on **coding, tools, long context, and multilingual** tasks. The chat template exposes **none / low / high** thinking budgets. Official paths include **Transformers**, **SGLang**, **vLLM** (Spark plugin), **MLX**, **llama.cpp** GGUF, and **Ollama**. **Apache 2.0**.

## Why keep

Compare **~2B-class local agents** against Qwen 3.5-2B, MiniCPM-2B, LFM2.5, or **nanbeige-4.2-3b** when throughput and tool calling on a single GPU or phone-class hardware matter.

## When to reach for it

You want a **small open checkpoint** with documented **spark25 tool-call parsers** and optional **very long context** in server configs (memory permitting), not a cloud-only API.

## Caveats

Benchmark tables are **vendor-published**. Max context in docker examples (**1M tokens**) is not realistic on consumer VRAM. Some integrations need **forked** runtimes (`XHToken/llama.cpp`, Spark plugin for vLLM).

## Links

- Homepage: https://huggingface.co/XHToken/Spark-X2.6-1.7B
- GitHub: https://github.com/XHToken/Spark-X2.6
- Weights: https://huggingface.co/XHToken/Spark-X2.6-1.7B
