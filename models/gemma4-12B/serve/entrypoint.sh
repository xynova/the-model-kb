#!/bin/sh
set -eu

exec /app/llama-server \
  -m /models/model.gguf \
  --mmproj /models/mmproj.gguf \
  --host 0.0.0.0 \
  --port "${PORT:-8080}" \
  -ngl "${N_GPU_LAYERS:-99}"
