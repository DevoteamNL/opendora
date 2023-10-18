NAME=opendora-api

PWD = $(shell pwd)
SHA = $(shell git --no-pager describe --always)
TAG ?= $(shell git tag --points-at HEAD)

HARDWARE = $(shell uname -m)
BUILD_TIME = $(shell date '+%s')
LFLAGS ?= -X main.gitsha=${SHA} -X main.compiled=${BUILD_TIME}

TAG := $(if $(TAG),$(TAG),$(shell git symbolic-ref --short HEAD)-$(SHA))

IMAGE_VERSION ?= $(TAG)
IMAGE ?= ghcr.io/devoteamnl/opendora/opendora-api:$(IMAGE_VERSION)

export DOCKER_CLI_EXPERIMENTAL=enabled

default: build

.PHONY: clean # Clean up the project (binaries and reports)
clean:
	@rm -rf ./bin 2>/dev/null
	@rm -rf ./release 2>/dev/null
	@rm -f cover.*

.PHONY: dependencies # Add dependencies to current module and install them
dependencies:
	@go get ./...

.PHONY: build # Build the OpenDORA API
build: dependencies clean
	@echo "--> Compiling the project"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -o bin/${NAME} api/main.go

.PHONY: build-static # Build the OpenDORA API to a static binary
build-static: clean
	@echo "--> Running static compilation"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -tags netgo -ldflags "-s -w ${LFLAGS}" -o bin/${NAME} api/main.go

.PHONY: vet # Run go vet looking for suspicious constructs
vet:
	@echo "--> Running go vet"
	@go vet ./...

.PHONY: lint # Runs golangci-lint
lint:
	@echo "--> Running go lint"
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint golangci-lint run -v

.PHONY: check-format # Check if there are go files in need of formatting
check-format:
	@echo "--> Checking go fmt"
	@gofmt -e -l **/*.go

.PHONY: format # Format go files
format:
	@echo "--> Running go fmt"
	@gofmt -s -w **/*.go

.PHONY: test # Run the tests
test:
	@echo "--> Running the tests"
	@go test ./... -v

.PHONY: bench # Run the tests generating their benchmarks
bench:
	@echo "--> Running the tests with benchmarking"
	@go test ./... -bench=. -benchmem

.PHONY: cover # Run the tests generating coverage profiles
cover:
	@echo "--> Running go cover"
	@go test ./... --cover

.PHONY: cover-report # Generate coverage report for the default profile
cover-report: cover
	@echo "--> Running go coverage report"
	@go test ./... -coverprofile cover.out
	@go tool cover -html=cover.out -o cover.html

.PHONY: docker-build # Build the multi arch Docker Image
docker-build:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
				--label "org.opencontainers.image.source=https://github.com/devoteamnl/opendora" \
           		--output "type=docker,push=false" \
           		--tag $(IMAGE) \
           		.

.PHONY: docker-push # Build and push the multi arch Docker Image
docker-push:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
				--label "org.opencontainers.image.source=https://github.com/devoteamnl/opendora" \
           		--platform linux/amd64,linux/arm64 \
           		--output "type=image,push=true" \
           		--tag $(IMAGE) \
           		.

.PHONY: run-api # Run the OpenDORA API
run-api: build
	@echo "--> Running the OpenDORA API"
	@go run ./api

.PHONY: run-sync # Run Backstage (Groups) -> DevLake (Teams) Sync
run-sync: build
	@echo "--> Running the Backstage (Groups) -> DevLake (Teams) Sync"
	@go run ./cmd/group-sync

.PHONY: help # Show this list of commands
help:
	@echo "OpenDORA Backend"
	@echo "Usage: make [target]"
	@echo ""
	@echo "target	description" | expand -t20
	@echo "-----------------------------------"
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1	\2/' | expand -t20
