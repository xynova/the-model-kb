---
title: gemma4-12B
sidebar:
  open: true
---

Gemma 4 12B packaging: static scratch weights image and a CUDA llama.cpp Cloud Run server image.

| Package | Image | Bake target |
|---------|-------|-------------|
| Static weights | `xynova/gemma4-12b-model` | `gemma4-12B` |
| Server | `xynova/gemma4-12b-server` | `gemma4-12B-serve` |

## Extract static weights

```bash
docker create --name t xynova/gemma4-12b-model:latest
docker cp t:/gemma-4-12B-it-QAT-Q4_0.gguf .
docker cp t:/mmproj-gemma-4-12B-it-QAT-BF16.gguf .
docker rm t
```

## Build

```bash
make docker-push-target TARGET=gemma4-12B
make docker-push-target TARGET=gemma4-12B-serve
```

## Invoke on Cloud Run

See the [invoke guide](/models/gemma4-12b/serve/invoke/).
