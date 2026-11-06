---
name: motion-bricks
title: MotionBricks
summary: "Real-time neural motion generation for SOMA human characters and Unitree G1, with smart-primitive control."
category: video
does: video
tags:
  - video
  - open
  - local
  - robotics
status: watch
run: local
license: mixed
vram_hint: unknown
links:
  homepage: https://nvlabs.github.io/motionbricks/
  github: https://github.com/NVlabs/motionbricks
  weights: https://huggingface.co/nvidia/MotionBricks
media:
  cover: video/motion-bricks/media/cover.jpg
source: "youtube:IGICYFyP_EI"
added: 2026-11-03
---

# MotionBricks

## What it is

NVIDIA’s open **real-time motion generation** stack: one modular latent backbone targets navigation and object interaction via “smart primitives,” with shipped **SOMA** human and **Unitree G1** models (bundled ONNX plus Hugging Face SafeTensors). The project page claims very high FPS and low latency at large skill counts; the repo includes Warp/Newton/TensorRT runtimes, UE-style demos in Python, and editors for primitives and scenes.

| | |
|:---:|:---:|
| ![motion-bricks 1](motion-bricks/media/thumb-cover.jpg) | |

## Why keep

Compare or wire **local character and humanoid motion** when you need controllable locomotion and interaction beyond hand-authored clips, especially alongside GR00T whole-body control work.

## When to reach for it

You are prototyping **game or robotics animation** from velocity, heading, style, or proxy keyframes, and want open weights with a documented runtime rather than a closed mocap or video API.

## Caveats

**License is split:** Apache 2.0 code, **NVIDIA Open Model Agreement** for checkpoints/ONNX, and **evaluation-only** terms on bundled sample motions (not for shipped products). Training data and full UE showcase assets are not the same as the slim GitHub release; VRAM and FPS depend on backend and scene.

## Links

- Homepage: https://nvlabs.github.io/motionbricks/
- GitHub: https://github.com/NVlabs/motionbricks
- Weights: https://huggingface.co/nvidia/MotionBricks
