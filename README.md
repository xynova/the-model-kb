# the-model-kb

Scratch-image packages live under `models/`. Place weights in `models/{model-dir}/model-files/` before building.

Aggregate bake file: [`models/docker-bake.hcl`](models/docker-bake.hcl). Run bake commands from the **repository root** (contexts are resolved from the working directory):

```bash
make docker-print
make docker-push
make docker-push-target TARGET=gemma4-12B-serve
```

Or:

```bash
docker buildx bake -f models/docker-bake.hcl --push
```

Per-model builds still work from a model directory with that directory's local `docker-bake.hcl`.

## Public site

Hugo site in [`site/`](site/) publishes three catalogs to GitHub Pages (set Pages source to **GitHub Actions** after the first merge):

| Section | Source | URL path |
|---------|--------|----------|
| Library | [`library/`](library/) | `/library/` |
| Packaged models | [`models/`](models/) docs | `/models/` |
| Roster | [`roster/`](roster/) | `/roster/` |

- Published URL (after first deploy): https://xynova.github.io/the-model-kb/
- Local preview: `make serve`
- Production build: `make build`

See [`site/instructions.md`](site/instructions.md) for theme, mounts, and workflow details.

## Extract

**gemma-model**
```bash
docker create --name t xynova/gemma-model:latest
docker cp t:/gemma-3n-E4B-it-Q6_K.gguf .
docker rm t
```

**gemma4-12b-model**
```bash
docker create --name t xynova/gemma4-12b-model:latest
docker cp t:/gemma-4-12B-it-QAT-Q4_0.gguf .
docker cp t:/mmproj-gemma-4-12B-it-QAT-BF16.gguf .
docker rm t
```

**qwen-model**
```bash
docker create --name t xynova/qwen-model:latest
docker cp t:/Qwen3VL-4B-Instruct-Q4_K_M.gguf .
docker cp t:/mmproj-Qwen3VL-4B-Instruct-F16.gguf .
docker rm t
```

**bge-base-model**
```bash
docker create --name t xynova/bge-base-model:latest
docker cp t:/bge-base-en-v1.5-q4_k_m.gguf .
docker rm t
```

**onrith-1-model**
```bash
docker create --name t xynova/onrith-1-model:latest
docker cp t:/ornith-1.0-9b-Q4_K_M.gguf .
docker rm t
```

## Cloud Run (gemma4-12b-server)

Runnable llama.cpp server image (CUDA). Requires weights in `models/gemma4-12B/static/model-files/` first.

**Build**
```bash
make docker-push-target TARGET=gemma4-12B-serve
```

**Deploy** (1× NVIDIA L4, `asia-southeast1`)

```bash
gcloud run deploy gemma4 \
  --image xynova/gemma4-12b-server:latest \
  --region asia-southeast1 \
  --port 8080 \
  --memory 16Gi \
  --cpu 4 \
  --gpu 1 \
  --gpu-type nvidia-l4 \
  --no-gpu-zonal-redundancy \
  --no-cpu-throttling \
  --cpu-boost \
  --concurrency 4 \
  --timeout 300 \
  --min-instances 1 \
  --max-instances 2 \
  --service-account cloud-run-ai@gen-lang-client-0341029530.iam.gserviceaccount.com \
  --no-allow-unauthenticated \
  --startup-probe tcpSocket.port=8080,timeoutSeconds=240,periodSeconds=240,failureThreshold=1
```

Image is pulled via Artifact Registry Docker Hub proxy in production:
`asia-southeast1-docker.pkg.dev/gen-lang-client-0341029530/docker-hub-proxy/xynova/gemma4-12b-server`.

Health check: `GET /health`. OpenAI-compatible API on `/v1/chat/completions`.

Invocation guide (auth, curl, Python): [models/gemma4-12B/serve/INVOKE.md](models/gemma4-12B/serve/INVOKE.md).
