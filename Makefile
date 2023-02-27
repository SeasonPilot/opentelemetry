include scripts/make-rules/image.mk

all: run push

run:
	docker-compose up --build

cache-run:
	docker-compose up

## push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@$(MAKE) image.push