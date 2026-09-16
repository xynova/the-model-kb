---
name: unifolm-wla
title: UnifoLM-WLA-1.0
summary: "Unitree World-Language-Action humanoid policy (~**6B** total stack) that combines embodied reasoning, future dynamic-region prediction, and action generation for **64** real-robot tasks (54 tabletop, 10 whole-body)."
category: robotics
does: robotics
tags:
  - robotics
  - vision
  - open
  - local
status: watch
run: local
license: open
vram_hint: consumer
links:
  homepage: https://unigen-x.github.io/unifolm-wla.github.io/
  github: https://github.com/unitreerobotics/unifolm-wla
  weights: https://huggingface.co/collections/unitreerobotics/unifolm-wla-10
media:
  cover: robotics/unifolm-wla/media/cover.jpg
source: "youtube:nZYJdwM-_nI"
added: 2026-09-16
---

# UnifoLM-WLA-1.0

## What it is

Unitree World-Language-Action humanoid policy (~**6B** total stack) that combines embodied reasoning, future dynamic-region prediction, and action generation for **64** real-robot tasks (54 tabletop, 10 whole-body). Supports two-finger grippers and multiple five-finger hands. Released pieces include UnifoLM-ER-1 and UnifoLM-ER-Flow on Hugging Face (Apache-2.0) plus the `unifolm-wla` GitHub repo.

| | |
|:---:|:---:|
| ![unifolm-wla 1](unifolm-wla/media/thumb-cover.jpg) |   |

## Why keep

Compare compact onboard humanoid VLAs when you care about whole-body plus tabletop in one Unitree-oriented stack.

## When to reach for it

Unitree G1 / similar humanoid experiments, or embodied-reasoning bake-offs against other open ER/VLA checkpoints in the UnifoLM family.

## Caveats

Roundup “Uni LM” maps here. The project page advertises a full WLA-1.0; verify which HF artifacts (ER-1, ER-Flow, action expert) you actually need for closed-loop control. Needs a compatible robot or sim. Earlier UnifoLM-VLA-0 family used different licenses on some cards; confirm each card before commercial use.

## Links

- Homepage: https://unigen-x.github.io/unifolm-wla.github.io/
- GitHub: https://github.com/unitreerobotics/unifolm-wla
- Weights collection: https://huggingface.co/collections/unitreerobotics/unifolm-wla-10
- ER-1: https://huggingface.co/unitreerobotics/UnifoLM-ER-1
- ER-Flow: https://huggingface.co/unitreerobotics/UnifoLM-ER-Flow
