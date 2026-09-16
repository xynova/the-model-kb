.PHONY: serve build clean site-serve site-build site-clean \
	docker-print docker-push docker-push-target

BAKE_FILE ?= models/docker-bake.hcl

# Library Hugo site (site/)
serve site-serve:
	$(MAKE) -C site serve

build site-build:
	$(MAKE) -C site build

clean site-clean:
	$(MAKE) -C site clean

# Packaged model images (aggregate bake under models/)
# Run these from the repository root. Contexts in the bake file are
# resolved relative to the current working directory.
docker-print:
	docker buildx bake --print -f $(BAKE_FILE)

docker-push:
	docker buildx bake --push -f $(BAKE_FILE)

docker-push-target:
	@test -n "$(TARGET)" || (echo 'Usage: make docker-push-target TARGET=<bake-target>' >&2; exit 1)
	docker buildx bake --push -f $(BAKE_FILE) $(TARGET)
