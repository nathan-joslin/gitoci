all: build

REV=$(shell git describe --long --tags --match='v*' --dirty 2>/dev/null || git rev-list -n1 HEAD)

# This is the default. It can be overridden in the main Makefile after
# including build.make.
REGISTRY_NAME=reg.lion.act3-ace.ai
IMAGE_REPO ?= $(REGISTRY_NAME)/gitoci

.PHONY: generate
generate: tool/controller-gen tool/crd-ref-docs
	go generate ./...

.PHONY: build
build: generate
	@mkdir -p bin
	go build -o bin/gitoci ./cmd/gitoci

.PHONY: build
build-linux: generate
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/gitoci--linux--amd64 ./cmd/gitoci

.PHONY: install
install: generate
	go install ./cmd/gitoci

.PHONY: build-docker
build-docker: generate
	docker build -t reg.github.com/act3-ai/gitoci .

.PHONY: docker-run
docker-run: generate
	docker run -it reg.github.com/act3-ai/gitoci:latest

.PHONY: ko
ko: tool/ko
	VERSION=$(REV) KO_DOCKER_REPO=$(IMAGE_REPO) tool/ko build -B --platform=all --image-label version=$(REV) ./cmd/gitoci

.PHONY: test
test: generate
	go test ./...

.PHONY: lint
lint: tool/golangci-lint
	tool/golangci-lint run

############################################################
# External tools
############################################################

# renovate: datasource=go depName=sigs.k8s.io/controller-tools
CONTROLLER_GEN_VERSION?=v0.16.5
# renovate: datasource=go depName=k8s.io/code-generator
CONVERSION_GEN_VERSION?=v0.31.2
# renovate: datasource=go depName=github.com/elastic/crd-ref-docs
CRD_REF_DOCS_VERSION?=v0.1.0
# renovate: datasource=go depName=github.com/google/ko
KO_VERSION?=v0.17.1
# renovate: datasource=go depName=github.com/golangci/golangci-lint
GOLANGCILINT_VERSION?=latest

# Installs all tools
.PHONY: tool
tool: tool/controller-gen tool/crd-ref-docs tool/ko tool/golangci-lint tool/go-md2man

# controller-gen: generates copy functions for CRDs
tool/controller-gen: tool/.controller-gen.$(CONTROLLER_GEN_VERSION)
	GOBIN=$(PWD)/tool go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

tool/.controller-gen.$(CONTROLLER_GEN_VERSION):
	@rm -f tool/.controller-gen.*
	@mkdir -p tool
	touch $@

# conversion-gen: generates conversion functions for CRDs
tool/conversion-gen: tool/.conversion-gen.$(CONVERSION_GEN_VERSION)
	GOBIN=$(PWD)/tool go install k8s.io/code-generator/cmd/conversion-gen@$(CONVERSION_GEN_VERSION)

tool/.conversion-gen.$(CONVERSION_GEN_VERSION):
	@rm -f tool/.conversion-gen.*
	@mkdir -p tool
	touch $@

# crd-ref-docs: Generates markdown documentation for CRDs
tool/crd-ref-docs: tool/.crd-ref-docs.$(CRD_REF_DOCS_VERSION)
	GOBIN=$(PWD)/tool go install github.com/elastic/crd-ref-docs@$(CRD_REF_DOCS_VERSION)

tool/.crd-ref-docs.$(CRD_REF_DOCS_VERSION):
	@rm -f tool/.crd-ref-docs.*
	@mkdir -p tool
	touch $@

# ko: builds application images for Go projects
tool/ko: tool/.ko.$(KO_VERSION)
	GOBIN=$(PWD)/tool go install github.com/google/ko@$(KO_VERSION)

tool/.ko.$(KO_VERSION):
	@rm -f tool/.ko.*
	@mkdir -p tool
	touch $@

# golangci-lint: lints Go code
tool/golangci-lint: tool/.golangci-lint.$(GOLANGCILINT_VERSION)
	GOBIN=$(PWD)/tool go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

tool/.golangci-lint.$(GOLANGCILINT_VERSION):
	@rm -f tool/.golangci-lint.*
	@mkdir -p tool
	touch $@
