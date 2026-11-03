---
name: stream-vibe
title: StreamVibe
summary: "Google/UCSB streaming proactive video agent plus StreamVibe-Bench for implicit, evolving user intent."
category: training-research
does: research
tags:
  - vision
  - agentic
  - api
status: watch
run: api
license: closed
vram_hint: unknown
links:
  homepage: https://stream-vibe.github.io/
  github: https://github.com/Stream-Vibe/Stream-Vibe.github.io
source: "youtube:IGICYFyP_EI"
added: 2026-11-04
---

# StreamVibe

## What it is

**StreamVibe** is a research stack for **proactive video assistants**: a streaming **StreamVibe-VL** model chooses each second whether to stay silent or help, and **StreamVibe-Harness** tracks multiple user requests and task progress in hierarchical memory. **StreamVibe-Bench** (~4.1K eval, ~85K train samples) covers **explicit, implicit, and dynamic** intents (cooking, tool use, UI navigation). Paper: [arXiv:2610.34662](https://arxiv.org/abs/2610.34662).

## Why keep

Track the **when-to-speak** problem for copilots on live video or screen streams, and compare future streaming VLMs against a published bench instead of ad hoc demos.

## When to reach for it

You are designing **intent-aware streaming agents** (implicit goals, mid-task revisions) and need benchmark framing or architecture references before weights ship.

## Caveats

As of filing, the public artifact is the **project site and paper**; the site notes **code links are not yet wired**, and there is **no Hugging Face weight drop** in the release. Treat reported gains as paper metrics until a reproducible checkpoint appears.

## Links

- Homepage: https://stream-vibe.github.io/
- GitHub (project page source): https://github.com/Stream-Vibe/Stream-Vibe.github.io
- Paper: https://arxiv.org/abs/2610.34662
