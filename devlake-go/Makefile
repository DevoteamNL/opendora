NAME=opendora-api

PWD = $(shell pwd)
SHA = $(shell git --no-pager describe --always)
TAG ?= $(shell git tag --points-at HEAD)

HARDWARE=$(shell uname -m)
BUILD_TIME=$(shell date '+%s')
LFLAGS ?= -X main.gitsha=${SHA} -X main.compiled=${BUILD_TIME}

IMG_VERSION ?= latest
IMG ?= devoteam/opendora-backend-api:$(IMG_VERSION)

default: build

.PHONY: dependencies build static

dependencies:
	@go get ./...

build: dependencies
	@echo "--> Compiling the project"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -o bin/${NAME} api/main.go

static:
	@echo "--> Running static compilation"
	@mkdir -p bin
	CGO_ENABLED=0 go build -a -tags netgo -ldflags "-s -w ${LFLAGS}" -o bin/${NAME} api/main.go

.PHONY: clean vet lint format bench cover cover-report
clean:
	@rm -rf ./bin 2>/dev/null
	@rm -rf ./release 2>/dev/null
	@rm -f cover.*

vet:
	@echo "--> Running go vet"
	@go vet ./...

lint:
	@echo "--> Running go lint"
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint golangci-lint run -v

format:
	@echo "--> Running go fmt"
	@gofmt -s -w **/*.go

bench:
	@echo "--> Running go bench"
	@go test ./... -bench=. -benchmem

cover:
	@echo "--> Running go cover"
	@go test ./... --cover

cover-report:
	@echo "--> Running go coverage report"
	@go test ./... -coverprofile cover.out
	@go tool cover -html=cover.out -o cover.html

.PHONY: test all
test:
	@echo "--> Running the tests"
	@go test ./... -v


all: test
	@echo "--> Performing all tests"
	@$(MAKE) gofmt
	@$(MAKE) vet
	@$(MAKE) cover
	@$(MAKE) bench
	@$(MAKE) cover-report

.PHONY: docker-build docker-push
docker-build:
	docker build -t ${IMG} .

docker-push:
	docker push ${IMG}

.PHONY: run-api run-sync
run-api:
	@echo "--> Running the OpenDORA API"
	@go run ./api

run-sync:
	@echo "--> Running the Backstage (Groups) -> DevLake (Teams) Sync"
	@go run ./cmd/group-sync
