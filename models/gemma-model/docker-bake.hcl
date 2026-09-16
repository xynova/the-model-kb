group "default" {
  targets = ["model"]
}

target "model" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/gemma-model:latest",
    "docker.io/xynova/gemma-model:gemma-3n-E4B-it-Q6_K"
  ]
  context = "."
  output = ["type=registry"]
}

