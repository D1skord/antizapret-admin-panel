# Makefile for AntiZapret Admin Panel
.PHONY: help run dev-backend dev-frontend build install-deps create-frontend-dist

# --- Tooling Setup ---
# Define a local bin directory for our tools. This makes the setup self-contained.
LOCAL_BIN := $(CURDIR)/bin
# Define the full path to the air executable we want to use.
AIR_CMD := $(LOCAL_BIN)/air
# --- End Tooling Setup ---

BINARY_NAME=antizapret-admin
FRONTEND_DIR=./frontend

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  help           Show this help message."
	@echo "  run            Shows instructions on how to run development servers."
	@echo "  dev-backend    Run the backend development server with hot-reload."
	@echo "  dev-frontend   Run the frontend development server."
	@echo "  build          Build the production-ready binary."
	@echo "  install-deps   Install all Go and Node.js dependencies."

# Development
run:
	@echo "To start development, please run 'make dev-backend' and 'make dev-frontend' in separate terminals."

# dev-backend now depends on the 'air' executable file itself.
# If the file doesn't exist, 'make' will run the rule below to create it.
dev-backend: $(AIR_CMD) create-frontend-dist
	@echo "Starting backend server with hot-reload..."
	@$(AIR_CMD)

# This rule tells 'make' how to build the '$(AIR_CMD)' file if it's missing.
$(AIR_CMD):
	@echo "---"
	@echo "Setup: 'air' tool not found. Installing it locally into $(LOCAL_BIN)..."
	@mkdir -p $(LOCAL_BIN)
	@# By setting GOBIN, we tell 'go install' exactly where to put the binary.
	@GOBIN=$(LOCAL_BIN) go install github.com/air-verse/air@latest
	@echo "Setup complete. 'air' is now at $(AIR_CMD)"
	@echo "---"

dev-frontend:
	@echo "Starting frontend development server..."
	@cd $(FRONTEND_DIR) && npm run dev

# Target to ensure frontend/dist directory exists for Go embed
create-frontend-dist:
	@echo "Ensuring $(FRONTEND_DIR)/dist directory exists..."
	@mkdir -p $(FRONTEND_DIR)/dist

# Build
build:
	@echo "Building frontend assets..."
	@cd $(FRONTEND_DIR) && npm install && npm run build
	@echo "Building Go binary..."
	@go build -o $(BINARY_NAME)
	@echo "Build complete: ./${BINARY_NAME}"

# Dependencies
install-deps:
	@echo "Installing Go dependencies..."
	@go mod tidy
	@echo "Installing Node.js dependencies..."
	@cd $(FRONTEND_DIR) && npm install