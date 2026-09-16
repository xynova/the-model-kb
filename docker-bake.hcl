group "default" {
  targets = ["gemma", "qwen", "bge-base", "gemma4-12B", "gemma4-12B-serve", "onrith-1"]
}

target "gemma" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/gemma-model:latest",
    "docker.io/xynova/gemma-model:gemma-3n-E4B-it-Q6_K"
  ]
  context = "models/gemma-model"
  output = ["type=registry"]
}

target "qwen" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/qwen-model:latest",
    "docker.io/xynova/qwen-model:Qwen3VL-4B-Instruct-Q4_K_M"
  ]
  context = "models/qwen-model"
  output = ["type=registry"]
}

target "bge-base" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/bge-base-model:latest",
    "docker.io/xynova/bge-base-model:bge-base-en-v1.5-q4_k_m"
  ]
  context = "models/bge-base-model"
  output = ["type=registry"]
}

target "gemma4-12B" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/gemma4-12b-model:latest",
    "docker.io/xynova/gemma4-12b-model:gemma-4-12B-it-QAT-Q4_0"
  ]
  context = "models/gemma4-12B/static"
  output = ["type=registry"]
}

target "gemma4-12B-serve" {
  dockerfile = "serve/Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/gemma4-12b-server:latest",
    "docker.io/xynova/gemma4-12b-server:gemma-4-12B-it-QAT-Q4_0"
  ]
  context = "models/gemma4-12B"
  output = ["type=registry"]
}

target "onrith-1" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/onrith-1-model:latest",
    "docker.io/xynova/onrith-1-model:ornith-1.0-9b-Q4_K_M"
  ]
  context = "models/onrith-1"
  output = ["type=registry"]
}
