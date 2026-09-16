---
title: gemma-model
---

Scratch image packaging Gemma 3n E4B Instruct Q6_K GGUF.

| Field | Value |
|-------|-------|
| Image | `xynova/gemma-model` |
| Bake target | `gemma` |
| Directory | `models/gemma-model` |

## Extract

```bash
docker create --name t xynova/gemma-model:latest
docker cp t:/gemma-3n-E4B-it-Q6_K.gguf .
docker rm t
```

## Build

From the repository root:

```bash
make docker-push-target TARGET=gemma
```
