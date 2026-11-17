---
name: gemini-decision-1
title: Gemini Decision 1
summary: "Google Decisions API model (gemini-decision-1): up to 1,000 parallel noul/choice/rubric questions over shared multimodal input; probabilities, zero output tokens."
category: models
does: multimodal
tags:
  - llm
  - vision
  - video
  - audio
  - api
  - closed
  - agentic
status: watch
run: api
license: closed
vram_hint: unknown
links:
  homepage: https://blog.google/innovation-and-ai/technology/developers-tools/gemini-decisions-api/
source: "youtube:woWo80dPkVk"
added: 2026-11-17
---

# Gemini Decision 1

## What it is

**Gemini Decision 1** (`gemini-decision-1`) powers Google’s **Decisions API**: a **multimodal** model that scores **closed-set questions** in one forward pass over shared **text, image, video, audio, and PDF** input (up to **~262K** tokens). It returns **calibrated probabilities** and **confidence**, not generated prose (**0 output tokens** billed). Question types: **yes/no (`noul`)**, **choice** (2–255 options), and **score** (ordered rubric). One request can include up to **1,000** independent questions. Use cases Google highlights: routing, policy checks, tool selection, and rubric grading inside agent loops.

## Why keep

Compare against **Jev**-style cheap classifiers and full LLM “ask in chat” when you need **batched multimodal decisions** with explicit confidence thresholds.

## When to reach for it

You already use the **Gemini API / AI Studio** and want **parallel policy or routing questions** on the same document, clip, or ticket payload without paying for long completions.

## Caveats

**API-only**; not a general chat model (text generation capability is off for this sku). **noul** probabilities are bounded **0.01–0.99**. Low-confidence paths need your own fallback or human review. Roundup compared it to **Jev**; validate on your domain.

## Links

- Announcement: https://blog.google/innovation-and-ai/technology/developers-tools/gemini-decisions-api/
- Model card: https://aistudio.google.com/docs/models/gemini-decision-1
- Decisions API: https://aistudio.google.com/docs/decisions-api