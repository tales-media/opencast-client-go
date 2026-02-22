# Copyright 2025 shio solutions GmbH
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Options

OS             ?= $(HOST_OS)
ARCH           ?= $(HOST_ARCH)
VERSION        ?= dev

GIT_COMMIT     ?= $(shell git rev-parse --short HEAD || echo "unknown")
GIT_TREE_STATE ?= $(shell sh -c 'if test -z "$$(git status --porcelain 2>/dev/null)"; then echo clean; else echo dirty; fi')
BUILD_DATE     ?= $(shell date -u +"%Y-%m-%dT%TZ")

# Do not change
HOST_OS     = $(shell which go >/dev/null 2>&1 && go env GOOS)
HOST_ARCH   = $(shell which go >/dev/null 2>&1 && go env GOARCH)
SHELL       = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Targets

.DEFAULT_GOAL:=help

##@ General

.PHONY: help
help: ## Print this help
	@awk 'BEGIN                      { FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n" } \
				/^[a-zA-Z_0-9-]+:.*?##/    { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } \
				/^## [a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", substr($$1, 4), $$2 } \
				/^##@/                     { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' \
				$(MAKEFILE_LIST)

info: ## Print options
	@printf "\n"
	@printf "\033[1m%s\033[0m\n"          "Build"
	@printf "  \033[36m%-15s\033[0m %s\n"   "OS"             "$(OS)"
	@printf "  \033[36m%-15s\033[0m %s\n"   "ARCH"           "$(ARCH)"
	@printf "\n"
	@printf "\033[1m%s\033[0m\n"          "Version Info"
	@printf "  \033[36m%-15s\033[0m %s\n"   "VERSION"        "$(VERSION)"
	@printf "  \033[36m%-15s\033[0m %s\n"   "GIT_COMMIT"     "$(GIT_COMMIT)"
	@printf "  \033[36m%-15s\033[0m %s\n"   "GIT_TREE_STATE" "$(GIT_TREE_STATE)"
	@printf "  \033[36m%-15s\033[0m %s\n"   "BUILD_DATE"     "$(BUILD_DATE)"

##@ Development

.PHONY: generate
generate: generate-modules ## Generate all

.PHONY: generate-modules
generate-modules: ## Generate Go modules files
	@scripts/exec-local generate-modules

.PHONY: fmt
fmt: ## Run go fmt against code
	@scripts/exec-local fmt

.PHONY: vet
vet: ## Run go vet against code
	@scripts/exec-local vet

##@ Build

.PHONY: build
build: generate fmt vet ## Build all binaries
