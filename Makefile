#!/usr/bin/env make
MAKEFLAGS += --silent
SHELL := /usr/bin/env bash
TESTS_IMAGE_NAME := wordcloud-tests
TESTS_IMAGE_DOCKERFILE := tests.Dockerfile
REBUILD ?= # Rebuilds test and app Docker images.
define MAKEFILE_NOTES
- If you've made changes to the wordcloud tests Dockerfile, run 'make bdd' with
  REBUILD=1.
endef
export MAKEFILE_NOTES

usage: ## Prints this help text.
	printf "%s\n\n" "$(PROJECT_DESCRIPTION)"; \
  printf "%s\n\n" "TARGETS"; \
  awk 'BEGIN { FS = "[:]" }; { printf "  %-45s %-45s\n", $$1, $$2 }' <<< \
    $$(fgrep -h '##' $(MAKEFILE_LIST) | fgrep -v '?=' | fgrep -v grep | sed 's/\\$$//' | sed -e 's/##//'); \
  printf "\n\n%s\n\n" "ENVIRONMENT VARIABLES"; \
  awk 'BEGIN { FS = "[;]" }; { printf "  %-45s %-45s\n", $$1, $$2 }' <<< \
    $$(fgrep '?=' $(MAKEFILE_LIST) | grep -v grep | sed 's/\?=.*##//' | sed 's/Makefile://'); \
  printf "\n\n%s\n\n" "NOTES"; \
  echo "$$MAKEFILE_NOTES";

bdd: _build_test_docker_image
bdd: ## Tests word-cloud features.
	docker run -v $(PWD)/features:/sut --rm "$(TESTS_IMAGE_NAME)" /sut

_build_test_docker_image:
	{ test -z "$(REBUILD)" && test -n "$$(docker image ls --quiet "$(TESTS_IMAGE_NAME)")"; } && exit 0; \
	docker image build -t "$(TESTS_IMAGE_NAME)" - < $(TESTS_IMAGE_DOCKERFILE)
