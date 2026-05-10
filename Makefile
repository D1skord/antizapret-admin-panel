# Makefile for AntiZapret Admin Panel
.PHONY: help dev build build-release docker-build

BINARY_NAME=antizapret-admin-panel

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  help           Show this help message."
	@echo "  dev            Start full dev environment via Docker (backend :8080, frontend :5173)."
	@echo "  docker-build   Rebuild Docker images (after Dockerfile changes)."
	@echo "  build          Build the production-ready binary (Linux amd64)."

dev:
	@docker compose up

docker-build:
	@docker compose build

build:
	@echo "Building frontend..."
	@docker run --rm -v $(CURDIR)/frontend:/app -w /app node:22-alpine sh -c "npm ci && npm run build"
	@echo "Building Go binary for Linux..."
	@docker run --rm -v $(CURDIR):/app -w /app golang:1.25-alpine \
		sh -c "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME)"
	@echo "Build complete: ./$(BINARY_NAME)"
