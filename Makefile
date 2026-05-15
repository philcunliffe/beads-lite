# Makefile for beads-lite (SQLite-only fork)

.PHONY: all build test install help

BUILD_DIR := .
GIT_BUILD := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
INSTALL_DIR := $(HOME)/.local/bin
BUILD_TAGS := sqlite_lite

# Pure-Go SQLite (modernc.org/sqlite) — no CGO required.
export CGO_ENABLED := 0

all: build

# Build the bd-lite binary.
build:
	@echo "Building bd-lite..."
	@go build -tags "$(BUILD_TAGS)" -ldflags="-X main.Build=$(GIT_BUILD)" -o $(BUILD_DIR)/bd-lite ./cmd/bd

# Run tests.
test:
	@echo "Running tests..."
	@go test -tags "$(BUILD_TAGS)" ./...

# Install bd-lite into $(INSTALL_DIR); the bd shim at $(INSTALL_DIR)/bd
# is expected to forward to bd-lite.
install: build
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/bd-lite $(INSTALL_DIR)/bd-lite
	@echo "Installed bd-lite -> $(INSTALL_DIR)/bd-lite"

help:
	@echo "beads-lite Makefile targets:"
	@echo "  build    Build ./bd-lite (sqlite_lite tag)"
	@echo "  test     Run go test against the sqlite_lite build"
	@echo "  install  Copy bd-lite into $(INSTALL_DIR)"
