# SuggestIT backend. The binary is built natively and copied into a scratch
# image, so Docker never compiles Go.

BIN    := bin/suggestit
GOARCH ?= $(shell go env GOARCH)
IMAGE  ?= suggestit-backend:dev

.PHONY: all build test vet fmt fmt-check check run docker clean

all: check build

## build: static linux binary for the Docker image (GOARCH=amd64 in CI)
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -trimpath -ldflags="-s -w" -o $(BIN) .

## test: unit tests (no database or network needed)
test:
	go test ./...

## vet: go vet
vet:
	go vet ./...

## fmt: format all Go code
fmt:
	gofmt -w .

## fmt-check: fail if any Go file is not gofmt'ed
fmt-check:
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then echo "not gofmt'ed:"; echo "$$unformatted"; exit 1; fi

## check: everything CI checks
check: fmt-check vet test

## run: run natively (export the env vars from the README first)
run:
	go run .

## docker: build the image from a freshly built binary
docker: build
	docker build -t $(IMAGE) .

## clean: remove build output
clean:
	rm -rf bin
