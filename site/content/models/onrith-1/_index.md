---
title: onrith-1
---

Scratch image packaging Ornith 1.0 9B Q4_K_M GGUF.

| Field | Value |
|-------|-------|
| Image | `xynova/onrith-1-model` |
| Bake target | `onrith-1` |
| Directory | `models/onrith-1` |

## Extract

```bash
docker create --name t xynova/onrith-1-model:latest
docker cp t:/ornith-1.0-9b-Q4_K_M.gguf .
docker rm t
```

## Build

```bash
make docker-push-target TARGET=onrith-1
```
