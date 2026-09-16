---
name: code-world-model
title: Code World Model
summary: "Keeps an AI world consistent by storing rules and state in code (a coding agent as the “brain”), then uses a video model (MiniMax-H3 + CWM LoRA) only for what you see."
category: agents-memory
does: memory
tags:
  - agentic
  - video
  - open
  - local
  - memory
status: watch
run: local
license: open
vram_hint: multi-gpu
links:
  homepage: https://buaacyw.github.io/cwm/
  github: https://github.com/buaacyw/code-world-model
  weights: https://huggingface.co/NTU-yiwen/awm-minimax-h3-new1344-lora-checkpoints
media:
  cover: agents-memory/code-world-model/media/cover.jpg
  gallery:
    - agents-memory/code-world-model/media/g01.jpg
    - agents-memory/code-world-model/media/g02.jpg
    - agents-memory/code-world-model/media/g03.jpg
source: "youtube:4wjHNgMLeyY"
added: 2026-08-31
---

# Code World Model

## What it is

Keeps an AI world consistent by storing rules and state in code (a coding agent as the “brain”), then uses a video model (MiniMax-H3 + CWM LoRA) only for what you see.

| | |
|:---:|:---:|
| ![code-world-model 1](code-world-model/media/thumb-cover.jpg) | ![code-world-model 2](code-world-model/media/thumb-g01.jpg) |
| ![code-world-model 3](code-world-model/media/thumb-g02.jpg) | ![code-world-model 4](code-world-model/media/thumb-g03.jpg) |

## Why keep

Useful when building long interactive scenes that must remember events and rules instead of relying on video memory alone.

## When to reach for it

Agent + video / game-like worlds where state has to stay coherent across a long session.

## Caveats

Repo is Apache-2.0; the MiniMax-H3 video backbone has its own community license. Needs a serious GPU setup.

## Links

- Homepage: https://buaacyw.github.io/cwm/
- GitHub: https://github.com/buaacyw/code-world-model
- Weights: https://huggingface.co/NTU-yiwen/awm-minimax-h3-new1344-lora-checkpoints
