group "default" {
  targets = ["model"]
}

target "model" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/qwen-model:latest",
    "docker.io/xynova/qwen-model:Qwen3VL-4B-Instruct-Q4_K_M"
  ]
  context = "."
  output = ["type=registry"]
}

