# Xynova Model Docker Reference

## Current models

| Directory | Docker repo | model-files/ contents | Artifacts in image |
|-----------|-------------|----------------------|-------------------|
| gemma-model | xynova/gemma-model | `gemma-3n-E4B-it-Q6_K.gguf` | `/gemma-3n-E4B-it-Q6_K.gguf` |
| gemma4-12B/static | xynova/gemma4-12b-model | `gemma-4-12B-it-QAT-Q4_0.gguf`, `mmproj-...gguf` | `/gemma-4-12B-it-QAT-Q4_0.gguf`, `/mmproj-gemma-4-12B-it-QAT-BF16.gguf` |
| qwen-model | xynova/qwen-model | `Qwen3VL-4B-Instruct-Q4_K_M.gguf`, `mmproj-...gguf` | `/Qwen3VL-4B-Instruct-Q4_K_M.gguf`, `/mmproj-Qwen3VL-4B-Instruct-F16.gguf` |
| bge-base-model | xynova/bge-base-model | `bge-base-en-v1.5-q4_k_m.gguf` | `/bge-base-en-v1.5-q4_k_m.gguf` |
| onrith-1 | xynova/onrith-1-model | `ornith-1.0-9b-Q4_K_M.gguf` | `/ornith-1.0-9b-Q4_K_M.gguf` |

Update this table when adding a model.

## model-files/ convention

- All model weights go in `{model-dir}/model-files/` — never committed to git
- Track `{model-dir}/model-files/.gitkeep` so the directory exists after clone
- Root `.gitignore` uses `**/model-files/*` with `!**/model-files/.gitkeep`

## Dockerfile templates

### Single GGUF

```dockerfile
FROM scratch

COPY model-files/{filename}.gguf /{filename}.gguf

LABEL model="{filename}" \
      format="gguf"
```

### GGUF + mmproj (vision/multimodal)

```dockerfile
FROM scratch

COPY model-files/{main-filename}.gguf /{main-filename}.gguf
COPY model-files/{mmproj-filename}.gguf /{mmproj-filename}.gguf

LABEL model="{main-filename-without-ext}" \
      format="gguf"
```

### HuggingFace / sentence-transformers directory

```dockerfile
FROM scratch

COPY model-files/{model-dir} /{model-dir}

LABEL model="{model-dir}" \
      format="sentence-transformers"
```

Add `.dockerignore` in the model directory if HuggingFace cache may appear:

```
model-files/**/.cache
```

## Root docker-bake.hcl target template

```hcl
target "{bake-target}" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/{repo-name}:latest",
    "docker.io/xynova/{repo-name}:{version-tag}"
  ]
  context = "{model-dir}"
  output = ["type=registry"]
}
```

## README extract snippet templates

### Single GGUF

````markdown
### {Display Name}

1. **Pull the image:**
   ```bash
   docker pull xynova/{repo-name}:latest
   ```

2. **Extract the file:**
   ```bash
   docker create --name temp-{short} xynova/{repo-name}:latest
   docker cp temp-{short}:/{filename}.gguf ./{filename}.gguf
   docker rm temp-{short}
   ```

   One-liner (PowerShell):
   ```powershell
   docker create --name temp-{short} xynova/{repo-name}:latest; docker cp temp-{short}:/{filename}.gguf ./{filename}.gguf; docker rm temp-{short}
   ```
````

### GGUF + mmproj

Add a `docker cp` line per file in the image.

### HuggingFace directory

Use `docker cp temp-{short}:/{model-dir} ./{model-dir}` instead of a single file.
