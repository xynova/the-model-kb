---
name: muse-ember
title: Muse Ember
summary: "Meta Muse family execution model (muse-ember-1.0): multimodal inputs, 1M context, tool/routing/sub-agent volume work via Model API and Muse Code."
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
  homepage: https://research.meta.ai/blog/muse-ember
source: "youtube:woWo80dPkVk"
added: 2026-11-17
---

# Muse Ember

## What it is

**Muse Ember** (`muse-ember-1.0`) is Meta’s **small, cost-focused** model in the **Muse** line for **high-volume agent work**: tool calls, coding subtasks, routing, extraction, and batched analysis over many files or images. It accepts **text, image, video, audio, and PDF** and returns text, with up to **~1M** input tokens and large max output on **Meta Model API** and **Muse Code**. Meta positions it beside **Muse Wonder** / **Muse Spark**: a larger model plans, Ember runs the repeated steps. **No open weights** (contrast **Muse Spark Open** / **Muse Glimmer**).

## Why keep

Compare **API cost and latency** for agent harnesses and multimodal doc workflows against Haiku-class models when the job is volume, not frontier reasoning.

## When to reach for it

You already use **Meta Model API** or **Muse Code** and need a **cheaper multimodal executor** for sub-agents, ticket triage, or parallel file/image passes.

## Caveats

**Proprietary**; pricing and long-context tiers apply above **200K** input tokens. Roundup intro briefly said **Muse Amber** (transcript error). Benchmark claims are vendor-reported; verify on your workloads.

## Links

- Announcement: https://research.meta.ai/blog/muse-ember
- API docs: https://dev.meta.ai/docs/muse-ember
- Model catalog: https://dev.meta.ai/docs/models