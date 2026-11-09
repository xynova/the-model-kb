---
name: gemini-omni-1-2-flash
title: Gemini Omni 1.2 Flash
summary: "Google paid-tier API for 3–15s native-audio video generation, editing, extension, and reference-guided clips."
category: video
does: video
tags:
  - video
  - audio
  - api
  - closed
status: watch
run: api
license: closed
links:
  homepage: https://aistudio.google.com/docs/models/gemini-omni-flash
source: "youtube:9RIK3A24ZxA"
added: 2026-11-08
---

# Gemini Omni 1.2 Flash

## What it is

Google's **Omni Flash** video stack (`gemini-omni-1.2-flash` on the Gemini API): **text / image / video in**, **video with native audio** out. Clips are **3–15 seconds** at **360p–4K** (24 fps), with **reference-guided** generation, **NL editing**, **extension**, **upscaling**, and **multi-turn** follow-ups via the **Interactions API**.

## Why keep

Anchor for **comparing paid short-form video APIs** when you choose hosts for motion ads, try-on edits, or conversational re-cut workflows.

## When to reach for it

Use when you already pay for **Gemini API** or **Vertex** and need **GA Omni video** with **15s** output and **~$0.07/s** effective **720p** pricing (token-metered; higher at 1080p/4K).

## Caveats

**No local weights.** **Paid tier** on the official pricing page (no free tier listed for this model). **Not** on GenerateContent; video uses **Interactions API** only. **1.1-flash** routes to **1.2** on **2026-12-03**. Tooling and agent features (grounding, function calling) are largely unsupported on this SKU per Cloud docs.

## Links

- Homepage: https://aistudio.google.com/docs/models/gemini-omni-flash
- Pricing: https://aistudio.google.com/docs/pricing#gemini-omni-1.2-flash
- Product post: https://blog.google/innovation-and-ai/technology/developers-tools/introducing-gemini-omni-1-2-flash/
