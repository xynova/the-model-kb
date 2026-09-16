group "default" {
  targets = ["model"]
}

target "model" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/gemma4-12b-model:latest",
    "docker.io/xynova/gemma4-12b-model:gemma-4-12B-it-QAT-Q4_0"
  ]
  context = "."
  output = ["type=registry"]
}
