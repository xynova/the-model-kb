---
name: unlimited-ocr-2
title: Unlimited OCR 2
summary: "MIT Baidu vision-language model for long-horizon document parsing in one shot (SR-SWA)."
category: models
does: multimodal
tags:
  - vision
  - llm
  - open
  - local
status: watch
run: local
license: open
vram_hint: vram-heavy
links:
  homepage: https://github.com/baidu/Unlimited-OCR
  github: https://github.com/baidu/Unlimited-OCR
  weights: https://huggingface.co/baidu/Unlimited-OCR-2
source: "youtube:IGICYFyP_EI"
added: 2026-11-03
---

# Unlimited OCR 2

## What it is

**Unlimited OCR 2** (~3.3B parameters, Hugging Face `baidu/Unlimited-OCR-2`) is Baidu’s open **image-text-to-text** model for **very long documents**: many pages in one request with **sparse reference sliding-window attention (SR-SWA)** so cross-page structure is not lost. The GitHub repo documents **Transformers**, **vLLM**, and **SGLang** paths; weights are ungated SafeTensors (~6.6 GB on disk).

## Why keep

File or compare a **local long-PDF / book parsing** stack when you need OCR plus layout-aware text in one pass instead of per-page crops.

## When to reach for it

You are building ingestion for **multi-hundred-page scans**, technical manuals, or archives and want an open checkpoint with a published **MarathonDocBench** line rather than a closed document API.

## Caveats

README testing cites **A800-class** hardware; real VRAM depends on page count, resolution, and backend. Hub lists **MIT**; confirm license text on the repo/card before redistribution. Also available via Baidu Cloud and ModelScope if you prefer hosted inference.

## Links

- Homepage: https://github.com/baidu/Unlimited-OCR
- GitHub: https://github.com/baidu/Unlimited-OCR
- Weights: https://huggingface.co/baidu/Unlimited-OCR-2
