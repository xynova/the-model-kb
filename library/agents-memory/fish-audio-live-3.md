---
name: fish-audio-live-3
title: Fish Audio Live 3
summary: "Fish Audio full-duplex voice engine for agents: listens while speaking, short acks in ~186 ms; Fish Agent Platform, API, or enterprise on-prem (English)."
category: agents-memory
does: speech
tags:
  - audio
  - api
  - agentic
  - closed
status: watch
run: api
license: closed
vram_hint: unknown
links:
  homepage: https://fish.audio/live3/
source: "youtube:woWo80dPkVk"
added: 2026-11-17
---

# Fish Audio Live 3

## What it is

**Live 3** is Fish Audio’s **full-duplex** speech model for **voice agents**: it keeps listening while it talks, handles backchannels and barge-in, and can emit short acknowledgments in about **186 ms** while a longer reply is prepared. You wire it through the **Fish Agent Platform** or **API** (`voice.engine: live3`); **enterprise on-prem** is offered. **English only** at launch. Weights are **not** open-sourced as of the Nov 2026 announcement.

## Why keep

Compare hosted **duplex agent** timing and interruption behavior against open stacks (Haru-4B, Sophon, and similar) when the job is product voice UX, not training your own speech LM.

## When to reach for it

You are building or evaluating a **Fish Audio agent** and need conversational rhythm (acknowledgments, yielding on real interruptions) without stitching ASR, LLM, and TTS turn-taking by hand.

## Caveats

**API / platform** first; no public weight drop yet (the roundup notes Fish has open-sourced other speech models, but not Live 3). Pricing is per-minute on the product pages (promotional rates may apply). Latency claims cover **small reactions**, not full reply time-to-first-audio.

## Links

- Homepage: https://fish.audio/live3/
- Announcement: https://fish.audio/blog/introducing-live-3/
- Agent API (engine `live3`): https://docs.fish.audio/api-reference/endpoint/agent/create-agent