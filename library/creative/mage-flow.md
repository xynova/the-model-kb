---
name: mage-flow
title: Mage-Flow
summary: "Microsoft 4B native-resolution stack for text-to-image and instruction editing, with Base, RL, and 4-step Turbo checkpoints."
category: creative
does: image
tags:
  - image
  - open
  - local
status: watch
run: local
license: open
vram_hint: consumer
links:
  homepage: https://microsoft.github.io/Mage/flow/
  github: https://github.com/microsoft/Mage
  weights: https://huggingface.co/microsoft/Mage-Flow
source: "youtube:fMaPaK4gcNE"
added: 2026-11-04
---

# Mage-Flow

## What it is

Microsoft’s **Mage-Flow** is a ~4B rectified-flow image stack: **Mage-VAE** plus a native-resolution **MMDiT** backbone. Checkpoints cover text-to-image (**Mage-Flow**) and instruction editing (**Mage-Flow-Edit**), each in Base, RL-aligned, and **4-step Turbo** variants on Hugging Face.

## Why keep

File it when you need an **MIT open** gen-plus-edit pair that targets mid-GPU VRAM instead of 20B+ FLUX-class weights, and you want both generation and editing from one family.

## When to reach for it

Local ComfyUI or diffusers runs at 1024² when API image tools are off limits, or when comparing open editors against Qwen-Image / FLUX on your own prompts.

## Caveats

Roundup and Microsoft benches are self-reported; independent rankings may differ. Full checkpoints are multi-component repos (not a single tiny GGUF); plan disk and VRAM for ~17GB-class installs before quants land.

## Links

- Homepage: https://microsoft.github.io/Mage/flow/
- GitHub: https://github.com/microsoft/Mage (see `mage_flow/`)
- Weights: https://huggingface.co/microsoft/Mage-Flow (plus `-Edit`, `-Turbo`, `-Base` siblings)
