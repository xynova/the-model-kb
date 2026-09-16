---
title: bge-base-model
---

Scratch image packaging BGE base English v1.5 Q4_K_M GGUF.

| Field | Value |
|-------|-------|
| Image | `xynova/bge-base-model` |
| Bake target | `bge-base` |
| Directory | `models/bge-base-model` |

## Extract

```bash
docker create --name t xynova/bge-base-model:latest
docker cp t:/bge-base-en-v1.5-q4_k_m.gguf .
docker rm t
```

## Build

```bash
make docker-push-target TARGET=bge-base
```
