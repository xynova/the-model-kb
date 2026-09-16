# Invoke Gemma 4 on Cloud Run

**Base URL:** `https://gemma4-rfcsyvwaua-as.a.run.app`

The service runs [llama.cpp server](https://github.com/ggml-org/llama.cpp) with an OpenAI-compatible API. Cloud Run requires IAM authentication on every request.

## Prerequisites

1. **Google Cloud account** with permission to invoke the service (`roles/run.invoker`).
2. **gcloud CLI** installed and logged in.

**For shell/curl testing:**
```bash
gcloud auth login
gcloud config set project gen-lang-client-0341029530
```

**For Go, Python, or other SDKs (recommended):**
```bash
gcloud auth application-default login
gcloud config set project gen-lang-client-0341029530
```

`gcloud auth login` alone does not set up Application Default Credentials. SDKs use ADC to mint audience-bound ID tokens automatically.

If you get `403 Forbidden`, ask a project admin to grant invoker access:

```bash
gcloud run services add-iam-policy-binding gemma4 \
  --region=asia-southeast1 \
  --member="user:YOU@example.com" \
  --role="roles/run.invoker"
```

Replace `user:YOU@example.com` with your Google account, or use a service account email for machine-to-machine calls.

## Authentication

### Quick testing (curl)

```bash
export GEMMA4_URL="https://gemma4-rfcsyvwaua-as.a.run.app"
export TOKEN="$(gcloud auth print-identity-token)"
```

Tokens expire after ~1 hour. Fine for manual testing; use an SDK for application code ([Google guidance](https://cloud.google.com/docs/authentication/get-id-token)).

### Application code (Go / Python / production)

Mint an **ID token** for the Cloud Run service URL (the audience). On GCP workloads the metadata server provides credentials; locally use `gcloud auth application-default login`.

| Environment | Credential source |
|-------------|-------------------|
| Mac / local dev | `gcloud auth application-default login` |
| Cloud Run / GCE / GKE | Service account attached to the workload |

The caller needs `roles/run.invoker`. Grant to a service account for prod:

```bash
gcloud run services add-iam-policy-binding gemma4 \
  --region=asia-southeast1 \
  --member="serviceAccount:YOUR-SA@gen-lang-client-0341029530.iam.gserviceaccount.com" \
  --role="roles/run.invoker"
```

## Quick checks

**Health**

```bash
curl -sS "$GEMMA4_URL/health" \
  -H "Authorization: Bearer $TOKEN"
```

**List models** (the `model` field in requests is ignored; one model is loaded at startup)

```bash
curl -sS "$GEMMA4_URL/v1/models" \
  -H "Authorization: Bearer $TOKEN" | jq
```

## Chat completion (curl)

```bash
curl -sS "$GEMMA4_URL/v1/chat/completions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemma-4-12b",
    "messages": [
      {"role": "system", "content": "You are a helpful assistant."},
      {"role": "user", "content": "Explain what a GGUF file is in two sentences."}
    ],
    "temperature": 0.7,
    "max_tokens": 256
  }' | jq -r '.choices[0].message.content'
```

## Streaming (curl)

```bash
curl -sS "$GEMMA4_URL/v1/chat/completions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemma-4-12b",
    "stream": true,
    "messages": [
      {"role": "user", "content": "Write a haiku about cloud inference."}
    ],
    "max_tokens": 128
  }'
```

## Go

Install the Google auth package:

```bash
go get google.golang.org/api/idtoken
```

### Reusable helper

Use the **root service URL as the audience**, even when calling a subpath like `/v1/chat/completions`:

```go
package gemma4

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/api/idtoken"
)

const DefaultURL = "https://gemma4-rfcsyvwaua-as.a.run.app"

// NewHTTPClient returns a client that attaches Cloud Run ID tokens automatically.
// audience must be the service root URL (no path).
func NewHTTPClient(ctx context.Context, audience string, timeout time.Duration) (*http.Client, error) {
	client, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return nil, err
	}
	if timeout > 0 {
		client.Timeout = timeout
	}
	return client, nil
}
```

### Chat completion

```go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"google.golang.org/api/idtoken"
)

const serviceURL = "https://gemma4-rfcsyvwaua-as.a.run.app"

func main() {
	ctx := context.Background()

	client, err := idtoken.NewClient(ctx, serviceURL)
	if err != nil {
		panic(err)
	}
	client.Timeout = 15 * time.Minute

	body, _ := json.Marshal(map[string]any{
		"model": "gemma-4-12b",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello from Go"},
		},
		"max_tokens": 128,
	})

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost,
		serviceURL+"/v1/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, _ := io.ReadAll(resp.Body)
	fmt.Println("status:", resp.Status)
	fmt.Println(string(out))
}
```

Run locally after `gcloud auth application-default login`:

```bash
go run .
```

### Manual token (custom client stack)

```go
ts, err := idtoken.NewTokenSource(ctx, serviceURL)
tok, err := ts.Token()

req, _ := http.NewRequestWithContext(ctx, "GET", serviceURL+"/health", nil)
req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
resp, err := http.DefaultClient.Do(req)
```

`TokenSource` refreshes tokens automatically; prefer `idtoken.NewClient` when possible.

### dspy-go

Use the OpenAI-compatible provider with an authenticated HTTP client:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/darwishdev/dspy-go/pkg/core"
	"github.com/darwishdev/dspy-go/pkg/llms"
	"google.golang.org/api/idtoken"
)

const serviceURL = "https://gemma4-rfcsyvwaua-as.a.run.app"

func newGemma4LLM(ctx context.Context) (*llms.OpenAILLM, error) {
	client, err := idtoken.NewClient(ctx, serviceURL)
	if err != nil {
		return nil, err
	}
	client.Timeout = 15 * time.Minute

	return llms.NewOpenAILLM(
		core.ModelID("gemma-4-12b"), // any string; server ignores it
		llms.WithOpenAIBaseURL(serviceURL),
		llms.WithOpenAIPath("/v1/chat/completions"),
		llms.WithOpenAITimeout(15*time.Minute),
		llms.WithHTTPClient(client),
	)
}

func main() {
	ctx := context.Background()
	llm, err := newGemma4LLM(ctx)
	if err != nil {
		panic(err)
	}

	resp, err := llm.Generate(ctx, "What is a GGUF file?")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Content)
}
```

Same code works in production: attach a service account with `roles/run.invoker` to your Cloud Run job — no `gcloud login` required.

## Python (OpenAI SDK)

```bash
pip install openai google-auth
```

```python
import subprocess
from openai import OpenAI

URL = "https://gemma4-rfcsyvwaua-as.a.run.app"

def cloud_run_token() -> str:
    return subprocess.check_output(
        ["gcloud", "auth", "print-identity-token"],
        text=True,
    ).strip()

client = OpenAI(
    base_url=f"{URL}/v1",
    api_key="unused",  # required by SDK; Cloud Run auth uses the header below
    default_headers={"Authorization": f"Bearer {cloud_run_token()}"},
)

response = client.chat.completions.create(
    model="gemma-4-12b",
    messages=[
        {"role": "user", "content": "What is Gemma 4?"},
    ],
    max_tokens=256,
)

print(response.choices[0].message.content)
```

### Python from another GCP service (no gcloud)

Use the metadata server or Application Default Credentials:

```python
import google.auth.transport.requests
import google.oauth2.id_token
from openai import OpenAI

URL = "https://gemma4-rfcsyvwaua-as.a.run.app"

def cloud_run_token() -> str:
    auth_req = google.auth.transport.requests.Request()
    return google.oauth2.id_token.fetch_id_token(auth_req, URL)

client = OpenAI(
    base_url=f"{URL}/v1",
    api_key="unused",
    default_headers={"Authorization": f"Bearer {cloud_run_token()}"},
)
```

Requires `gcloud auth application-default login` locally, or a service account with `roles/run.invoker` on GCP.

## Vision (optional)

The server image includes an mmproj file for multimodal input. Pass a base64-encoded image via `image_data` and reference it in the message with `[img-ID]`:

```bash
IMG_B64="$(base64 -i photo.jpg | tr -d '\n')"

curl -sS "$GEMMA4_URL/v1/chat/completions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"model\": \"gemma-4-12b\",
    \"messages\": [
      {\"role\": \"user\", \"content\": \"Describe this image: [img-1]\"}
    ],
    \"image_data\": [{\"id\": 1, \"data\": \"$IMG_B64\"}],
    \"max_tokens\": 256
  }" | jq -r '.choices[0].message.content'
```

## API reference

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Liveness check |
| `/v1/models` | GET | List loaded model |
| `/v1/chat/completions` | POST | Chat (OpenAI-compatible) |
| `/v1/completions` | POST | Text completion |
| `/v1/embeddings` | POST | Embeddings |

Full llama.cpp server docs: https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md

## Operational notes

- **Hardware:** 1× NVIDIA L4 (24 GiB VRAM), 4 vCPU, 16 GiB RAM. `min-instances: 1` avoids cold starts; model load still takes ~1–2 minutes on revision deploy.
- **GPU layers:** The server image uses `ghcr.io/ggml-org/llama.cpp:server-cuda` and offloads layers with `-ngl` (default 99). Override via `N_GPU_LAYERS` env on the Cloud Run service if needed.
- **Timeouts:** Request timeout is 300s in production. Increase if long generations are cut off:
  ```bash
  gcloud run services update gemma4 --region=asia-southeast1 --timeout=900
  ```
- **Concurrency:** Set to 4. Keep low for large models to avoid OOM on the L4.
- **Scaling:** `max-instances: 2`, GPU zonal redundancy disabled.
