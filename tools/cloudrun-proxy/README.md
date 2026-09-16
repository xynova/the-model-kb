# cloudrun-proxy

Local OpenAI-compatible proxy for private Cloud Run services. Handles Google ID-token auth (ADC), token refresh, retries, and streaming pass-through so tools like RooCode can use a stable `http://127.0.0.1:4318` endpoint.

## Prerequisites (local dev)

Impersonate the project invoker service account (SDK auto-refresh, no `gcloud auth login`):

```bash
gcloud auth application-default login
gcloud auth application-default set-quota-project gen-lang-client-0341029530
```

Configured in `config.yaml` under `auth.impersonate_service_account`.

Your Google user needs `roles/iam.serviceAccountTokenCreator` on that service account:

```bash
gcloud iam service-accounts add-iam-policy-binding \
  cloud-run-invoker@gen-lang-client-0341029530.iam.gserviceaccount.com \
  --member="user:YOU@example.com" \
  --role="roles/iam.serviceAccountTokenCreator"
```

Each model's Cloud Run service must grant `roles/run.invoker` to the invoker SA (one-time, if not already granted):

```bash
gcloud run services add-iam-policy-binding gemma4 \
  --region=asia-southeast1 \
  --member="serviceAccount:cloud-run-invoker@gen-lang-client-0341029530.iam.gserviceaccount.com" \
  --role="roles/run.invoker"
```

**Fallback (no impersonation):** clear `auth.impersonate_service_account` in `config.yaml`.

### `invalid_grant` / `invalid_rapt` / reauth errors

Impersonation uses your **Application Default Credentials** refresh token. Google periodically requires interactive re-login:

```bash
gcloud auth application-default login
gcloud auth application-default set-quota-project gen-lang-client-0341029530
```

Restart the proxy after re-auth. Verify ADC works:

```bash
gcloud auth application-default print-access-token
```

If your org enforces advanced protection / session re-auth, you may need to complete the browser flow when prompted.

## Configuration

Edit `config.yaml`. Precedence: **defaults → `config.yaml` → env vars** (env vars only override `server` and `auth`).

```yaml
server:
  listen_addr: 127.0.0.1:4318
  max_concurrent: 1
  reject_when_busy: true
  rate_limit_cooldown: 60s

auth:
  impersonate_service_account: cloud-run-invoker@gen-lang-client-0341029530.iam.gserviceaccount.com

defaults:
  chat_path: /v1/chat/completions
  request_timeout: 900s
  upstream_model: gemma-4-12b
  enable_thinking: false
  min_max_tokens: 512
  min_max_tokens_thinking: 2048

models:
  gemma-4-cloudrun:
    audience: https://gemma4-rfcsyvwaua-as.a.run.app
    url: https://gemma4-rfcsyvwaua-as.a.run.app

  # another-service:
  #   audience: https://other-xyz.a.run.app
  #   url: https://other-xyz.a.run.app
  #   upstream_model: other-model-name
```

The **map key** under `models` is the model ID clients send (e.g. RooCode Model ID). Each entry can override any `defaults` field. Models sharing the same `audience` + impersonation SA reuse one auth client.

| Env var | Overrides |
|---------|-----------|
| `CONFIG_FILE` | Config file path (default `config.yaml`) |
| `LISTEN_ADDR` | `server.listen_addr` |
| `LOCAL_SECRET` | `server.local_secret` |
| `DEBUG` | `server.debug` |
| `GEMMA_MAX_CONCURRENT` | `server.max_concurrent` |
| `IMPERSONATE_SERVICE_ACCOUNT` | `auth.impersonate_service_account` |

## Run

```bash
cd tools/cloudrun-proxy
make run
```

Other targets: `make build`, `make start`, `make smoke`.

Smoke test:

```bash
curl -sS http://127.0.0.1:4318/v1/models | jq
curl -sS http://127.0.0.1:4318/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"gemma-4-cloudrun","messages":[{"role":"user","content":"Hello"}],"max_tokens":32}'
```

## RooCode

| Field | Value |
|-------|-------|
| API Provider | OpenAI Compatible |
| Base URL | `http://127.0.0.1:4318/v1` (no trailing slash) |
| API Key | `local` (any non-empty string) |
| Model ID | A key from `models` in `config.yaml` (e.g. `gemma-4-cloudrun`) |

**Request timeout:** Gemma on Cloud Run can take **30–90+ seconds** before the first token. In RooCode set `roo-cline.apiRequestTimeout` to **900** or higher (**seconds**, not ms). See [RooCode timeout](#roocode-timeout) below.

`server.max_concurrent: 1` with `reject_when_busy: true` (default): a second local request gets **429** immediately instead of queueing, so RooCode can't pile up parallel calls that trigger Google **Rate exceeded**. After Cloud Run returns 429, `rate_limit_cooldown` (default 60s) blocks further upstream attempts.

**Thinking mode (Gemma 4):** `defaults.enable_thinking: false` by default because RooCode only reads `message.content`. Override per model or in `defaults`. The proxy copies `reasoning_content` into `content` when needed.

| Mode | Config | Typical `max_tokens` |
|------|--------|----------------------|
| Normal chat (default) | `enable_thinking: false` | 512+ (`defaults.min_max_tokens`) |
| Reasoning workloads | `enable_thinking: true` | 2048+ (`defaults.min_max_tokens_thinking`) |

See also [models/gemma4-12B/serve/INVOKE.md](../../models/gemma4-12B/serve/INVOKE.md) for Cloud Run service details.

## Logging

Every request is logged at **info** (`http` line with method, path, status, duration).

Chat completions also log: route selection, upstream status, timings, and auth errors with a **hint** when ADC expires.

For full detail (request summaries, auth client creation, upstream body sizes), enable debug:

```yaml
server:
  debug: true
```

Or: `DEBUG=true make run`

## RooCode timeout

RooCode uses a **global VS Code/Cursor setting** (not in the provider profile UI). Value is in **seconds**.

**Option A — Settings UI**

1. `Cmd + ,` (Mac) or `Ctrl + ,` (Windows) to open Settings
2. Search: `roo api request timeout` or `roo-cline.apiRequestTimeout`
3. Set **Roo Code: Api Request Timeout** to `900` (15 min) or `1200` (20 min)
4. Reload the window: Command Palette → `Developer: Reload Window`

**Option B — settings.json**

1. Command Palette → `Preferences: Open User Settings (JSON)`
2. Add:

```json
"roo-cline.apiRequestTimeout": 900
```

3. Save and reload the window

Default is `600` (10 minutes). Range `0`–`3600`; `0` means no timeout on recent RooCode versions.

If requests still die around **5 minutes**, update RooCode to the latest version (older builds had a separate Undici 300s limit).
