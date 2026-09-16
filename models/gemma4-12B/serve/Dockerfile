FROM ghcr.io/ggml-org/llama.cpp:server-cuda

COPY static/model-files/gemma-4-12B-it-QAT-Q4_0.gguf /models/model.gguf
COPY static/model-files/mmproj-gemma-4-12B-it-QAT-BF16.gguf /models/mmproj.gguf
COPY serve/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

LABEL model="gemma-4-12B-it-QAT-Q4_0" \
      format="gguf"

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
