# ==============================================================================
# Makefile helper functions for docker image
#

DOCKER := docker

REGISTRY_PREFIX ?=registry-qa-test.vecps.com/qa
SERVICE_PREFIX ?=opentelemetry

# Determine images names by stripping out the dir names
IMAGES ?= hertz-client hertz-server kitex-server multiple-server

VERSION=latest

ifeq (${IMAGES},)
  $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

.PHONY: image.push
image.push: $(addprefix image.push.,  $(IMAGES))

.PHONY: image.push.%
image.push.%:
	@echo "===========> Pushing image $* $(VERSION) to $(REGISTRY_PREFIX)"
	$(DOCKER) tag $(SERVICE_PREFIX)-$*:$(VERSION) $(REGISTRY_PREFIX)/$*:$(VERSION)
	$(DOCKER) push $(REGISTRY_PREFIX)/$*:$(VERSION)
