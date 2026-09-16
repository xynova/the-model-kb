---
name: add-xynova-model
description: >-
  Wire up a new ML model directory with Dockerfile, docker-bake.hcl, and README
  updates for Xynova scratch-image packaging. Use when adding a new model to this
  repo, containerizing GGUF or HuggingFace models, or mirroring gemma-model/qwen-model
  patterns.
---

# Add Xynova Model Container

Package a new model under `models/` in a minimal `FROM scratch` image and register it in the root bake file.

## Workflow

Copy this checklist and track progress:

```
- [ ] Inspect model directory and classify type (see reference.md)
- [ ] Choose lowercase Docker Hub repo name (xynova/{name}-model)
- [ ] Create models/{short-dir}/model-files/.gitkeep
- [ ] Place model weights in models/{short-dir}/model-files/
- [ ] Create models/{short-dir}/Dockerfile (COPY from model-files/)
- [ ] Create models/{short-dir}/docker-bake.hcl
- [ ] Add .dockerignore if HuggingFace .cache may exist under model-files/
- [ ] Add target to root docker-bake.hcl default group
- [ ] Update README.md (structure, build list, extract section)
- [ ] Validate: docker buildx bake --print
- [ ] Optional local test: docker buildx build --load .
```

## Step 1: Classify the model

| Type | Example | model-files/ layout |
|------|---------|---------------------|
| Single GGUF | models/gemma-model | `model-files/{name}.gguf` |
| GGUF + mmproj | models/qwen-model, models/gemma4-12B/static | `model-files/{main}.gguf` + `model-files/{mmproj}.gguf` |
| HuggingFace dir | (example) | `model-files/{model-dir}/` |

All weights live in `{model-dir}/model-files/` (gitignored). Dockerfiles COPY from `model-files/` and place artifacts at image root (`/filename.gguf` or `/model-dir/`).

## Step 2: Naming rules

- **Docker Hub repo**: lowercase only — `xynova/gemma4-12b-model`, not `gemma4-12B-model`
- **Version tag**: may keep model casing — `:gemma-4-12B-it-QAT-Q4_0`
- **Bake target**: HCL identifier, e.g. `gemma4-12B` (not a Docker tag)
- **Filesystem dir**: under `models/`; may use mixed case — `models/gemma4-12B/`

## Step 3: Create model Dockerfile

Always `FROM scratch`. Set `LABEL model` and `LABEL format` (`gguf` or `sentence-transformers`).

See [reference.md](reference.md) for copy-paste templates.

## Step 4: Create model docker-bake.hcl

```hcl
group "default" {
  targets = ["model"]
}

target "model" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/{repo-name}:latest",
    "docker.io/xynova/{repo-name}:{version-tag}"
  ]
  context = "."
  output = ["type=registry"]
}
```

## Step 5: Register in root docker-bake.hcl

1. Add target name to `group "default" { targets = [...] }`
2. Add matching `target "{name}"` block with same tags and `context = "models/{short-dir}"`

## Step 6: Update README.md

Add to: directory structure, build-all tags list, individual build section, extract section.

For multimodal GGUF, document **both** main and mmproj extraction. For HuggingFace dirs, extract the directory path.

## Step 7: Validate

```bash
cd {model-dir}
docker buildx bake --print
```

If error `repository name must be lowercase`, fix repo segment in all bake files and README.

For a local smoke test (no push):

```bash
docker buildx build --platform linux/amd64 -t {repo-name}:test --load .
```

## Additional resources

- Templates and current model registry: [reference.md](reference.md)
- Repo overview: [README.md](../../README.md)
