---
title: qwen-model
---

Scratch image packaging Qwen3VL 4B Instruct GGUF plus mmproj.

| Field | Value |
|-------|-------|
| Image | `xynova/qwen-model` |
| Bake target | `qwen` |
| Directory | `models/qwen-model` |

## Extract

```bash
docker create --name t xynova/qwen-model:latest
docker cp t:/Qwen3VL-4B-Instruct-Q4_K_M.gguf .
docker cp t:/mmproj-Qwen3VL-4B-Instruct-F16.gguf .
docker rm t
```

## Build

```bash
make docker-push-target TARGET=qwen
```
