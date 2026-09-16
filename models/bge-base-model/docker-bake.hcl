group "default" {
  targets = ["model"]
}

target "model" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/bge-base-model:latest",
    "docker.io/xynova/bge-base-model:bge-base-en-v1.5-q4_k_m"
  ]
  context = "."
  output = ["type=registry"]
}
