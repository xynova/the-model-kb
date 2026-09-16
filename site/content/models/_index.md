---
title: Packaged models
sidebar:
  open: true
---

Scratch-image Docker packages for transferring GGUF and related weights. Place files under each package `model-files/` directory before building. Build from the repository root with `make docker-push` or `make docker-push-target TARGET=...`.

{{< cards >}}
  {{< card link="/models/gemma-model" title="gemma-model" subtitle="Gemma 3n E4B GGUF scratch image" >}}
  {{< card link="/models/gemma4-12b" title="gemma4-12B" subtitle="Gemma 4 12B static weights and Cloud Run server" >}}
  {{< card link="/models/qwen-model" title="qwen-model" subtitle="Qwen3VL 4B GGUF + mmproj" >}}
  {{< card link="/models/bge-base-model" title="bge-base-model" subtitle="BGE base English embedding GGUF" >}}
  {{< card link="/models/onrith-1" title="onrith-1" subtitle="Ornith 1.0 9B GGUF" >}}
{{< /cards >}}
