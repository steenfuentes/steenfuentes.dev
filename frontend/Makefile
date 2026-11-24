# Frontend Makefile

# Variables
IMAGE_NAME ?= steenfuentes-frontend
TAG ?= latest
VITE_API_URL ?= /api
PORT ?= 3000

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make dev          - Run development server"
	@echo "  make build        - Build production bundle"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container locally"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test-build   - Test production build locally"

# Development
.PHONY: dev
dev:
	npm run dev

# Build production bundle
.PHONY: build
build:
	npm run build

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image: $(IMAGE_NAME):$(TAG)"
	@echo "API URL: $(VITE_API_URL)"
	docker build \
		--build-arg VITE_API_URL=$(VITE_API_URL) \
		-t $(IMAGE_NAME):$(TAG) \
		.

# Run Docker container locally
.PHONY: docker-run
docker-run:
	docker run \
		-p $(PORT):3000 \
		-e VITE_API_URL=http://localhost:8080 \
		$(IMAGE_NAME):$(TAG)

# Test production build locally (without Docker)
.PHONY: test-build
test-build: build
	VITE_API_URL=http://localhost:8080 \
	PORT=$(PORT) \
	node build

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf build .svelte-kit node_modules/.vite

# Install dependencies
.PHONY: install
install:
	npm ci

# Format code
.PHONY: format
format:
	npm run format

# Lint code
.PHONY: lint
lint:
	npm run lint

# Type check
.PHONY: check
check:
	npm run check