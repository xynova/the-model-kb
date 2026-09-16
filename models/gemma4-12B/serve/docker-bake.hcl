group "default" {
  targets = ["serve"]
}

target "serve" {
  dockerfile = "serve/Dockerfile"
  platforms = ["linux/amd64"]
  tags = [
    "docker.io/xynova/gemma4-12b-server:latest",
    "docker.io/xynova/gemma4-12b-server:gemma-4-12B-it-QAT-Q4_0"
  ]
  context = ".."
  output = ["type=registry"]
}
