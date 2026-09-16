---
name: xynova-model-docker-ops
description: >-
  Build, push, and extract Xynova ML model Docker images on Docker Hub. Use when
  running docker buildx bake, pushing model images, pulling and extracting GGUF or
  HuggingFace models from containers, or troubleshooting Docker tag/registry errors.
---

# Xynova Model Docker Operations

Operational tasks for the scratch-image model transfer workflow in this repo.

## Build and push

### Prerequisites

```bash
docker login
```

### Build all models (from repo root)

```bash
make docker-push
# or: docker buildx bake -f models/docker-bake.hcl --push
```

### Build one model (aggregate target)

```bash
make docker-push-target TARGET=gemma
# or: docker buildx bake -f models/docker-bake.hcl --push gemma
```

### Build one model (local bake file)

Place weights in `{model-dir}/model-files/` first, then:

```bash
cd {model-dir}
docker buildx bake --push
```

Examples: `models/gemma-model`, `models/gemma4-12B`, `models/qwen-model`, `models/bge-base-model`

### Validate bake config (no build)

Aggregate (from repo root):

```bash
make docker-print
# or: docker buildx bake --print -f models/docker-bake.hcl
```

Local bake file:

```bash
cd {model-dir}
docker buildx bake --print
```

### Local build without push

```bash
cd {model-dir}
docker buildx build --platform linux/amd64 -t {repo-name}:test --load .
```

## Extract models from images

Pattern: `pull` → `create` → `cp` → `rm`

### Single GGUF

```bash
docker pull xynova/{repo-name}:latest
docker create --name temp-{short} xynova/{repo-name}:latest
docker cp temp-{short}:/{filename}.gguf ./{filename}.gguf
docker rm temp-{short}
```

### HuggingFace directory

```bash
docker cp temp-{short}:/{model-dir} ./{model-dir}
```

### PowerShell one-liner

Chain with `;` — see [README.md](../../README.md) for per-model commands.

## Current images

| Image | Directory |
|-------|-----------|
| xynova/gemma-model | models/gemma-model |
| xynova/gemma4-12b-model | models/gemma4-12B/static |
| xynova/qwen-model | models/qwen-model |
| xynova/bge-base-model | models/bge-base-model |

Full artifact paths: [add-xynova-model/reference.md](../add-xynova-model/reference.md)

## Troubleshooting

### `repository name must be lowercase`

Docker Hub repo names must be all lowercase. Fix `docker.io/xynova/...` in:

- `{model-dir}/docker-bake.hcl`
- `models/docker-bake.hcl`
- `README.md` pull/extract examples

Version tags may keep mixed case.

### `permission denied` connecting to Docker daemon

Ensure Docker Desktop is running. Re-run with host Docker access if in a sandboxed environment.

### Large build context / slow transfer

- Weights must be in `{model-dir}/model-files/` — do not scatter files elsewhere
- Exclude HuggingFace cache: `model-files/**/.cache` in `.dockerignore`

### Scratch image has no shell

Cannot `docker run ... ls`. Verify contents with `docker create` + `docker cp`, or inspect build output.

## Adding a new model

Use the **add-xynova-model** skill — do not hand-roll Dockerfile/bake patterns without it.
