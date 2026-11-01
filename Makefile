# Project metadata
NAME    := zplkit
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

# Linker flags
LDFLAGS         := -X main.version=$(VERSION)
RELEASE_LDFLAGS  := $(LDFLAGS) -s -w

# Target platforms for `make release`
PLATFORMS := linux/amd64 linux/arm64 windows/amd64 darwin/amd64 darwin/arm64

.PHONY: build test vet release clean help

build: ## Development build → build/zplkit (debug symbols; version = "dev" if no tags)
	mkdir -p build
	go mod tidy
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o build/$(NAME) ./cmd/zplkit

test: ## Run the test suite
	go test -count=1 ./...

vet: ## Run go vet static analysis
	go vet ./...

release: clean ## Stripped, trimpath, multi-arch release archives → dist/ + checksums.txt
	mkdir -p dist
	go mod tidy
	@set -e; for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d/ -f1); \
		ARCH=$$(echo $$platform | cut -d/ -f2); \
		EXT=$$([ "$$OS" = "windows" ] && echo ".exe" || echo ""); \
		OUT=dist/$(NAME)_$${OS}_$${ARCH}$${EXT}; \
		echo "building $$OS/$$ARCH"; \
		GOOS=$$OS GOARCH=$$ARCH CGO_ENABLED=0 go build -trimpath \
			-ldflags "$(RELEASE_LDFLAGS)" -o $$OUT ./cmd/zplkit; \
	done
	@echo "==> Checksums"
	cd dist && sha256sum * > checksums.txt && cat checksums.txt

clean: ## Remove build/ and dist/
	rm -rf build dist

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\nTargets:\n"} \
	     /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
