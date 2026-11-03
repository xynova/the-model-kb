---
name: ideogram-vector-1-0
title: Ideogram Vector 1.0
summary: "Ideogram hosted SVG generator: text-to-vector and image vectorization with editable paths, app and API."
category: creative
does: image
tags:
  - image
  - closed
  - api
status: watch
run: api
license: closed
vram_hint: unknown
links:
  homepage: https://ideogram.ai/models/vector/
source: "youtube:EfyQOAJjx1A"
added: 2026-11-04
---

# Ideogram Vector 1.0

## What it is

**Ideogram Vector 1.0** (`ideogram-vector-1`): **hosted** model for **text-to-SVG** and **image-to-vector** (logos, wordmarks, illustrations, technical drawings). Outputs **SVG** with individually editable shapes for Figma/Illustrator workflows. Live in the **Ideogram app** (Generate vector) and **Ideogram API** (`POST /v2/svg/generate/ideogram-vector-1`, `POST /v2/svg/vectorize/ideogram-vector-1`). API pricing on Ideogram’s model page starts around **$0.13–$0.20 per generated SVG** by quality tier (verify current table).

## Why keep

Compare when you need **production SVG artwork** or **raster-to-vector** in a design pipeline, not only raster image models.

## When to reach for it

Brand and print work where you will **export SVG** and tweak paths in a vector editor; batch or product flows via **API key** and async `generation_id` polling (or webhook).

## Caveats

**Closed SaaS**; no local weights. Independent SVG teardowns note **outlined type** (not live text elements), **flat path lists** (layer grouping **not yet** on vectorize), and **no API edit** pass (regenerate instead). Marketing claims **editable** vectors; validate on your asset types before production lock-in.

## Links

- Model page: https://ideogram.ai/models/vector/
- API: https://developer.ideogram.ai/api-reference/vectors/generate-svg
