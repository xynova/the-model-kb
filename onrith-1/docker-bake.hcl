group "default" {
  targets = ["model"]
}

target "model" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/onrith-1-model:latest",
    "docker.io/xynova/onrith-1-model:ornith-1.0-9b-Q4_K_M"
  ]
  context = "."
  output = ["type=registry"]
}
