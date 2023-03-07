# Define the compiler and compiler flags
GO := go
GOBUILD := $(GO) build
GOFLAGS := -v
FRONTEND_DIR := frontend
BACKEND_DIR := backend

# Define the build targets
all: frontend-build backend-build

# Build the frontend first
frontend-build:
	@test -d $(FRONTEND_DIR)/node_modules || (cd $(FRONTEND_DIR) && npm install --silent)
	@test -d $(FRONTEND_DIR)/static || (cd $(FRONTEND_DIR) && npm run build --silent)
	# @cd $(FRONTEND_DIR) && pnpm run build
	@test -d $(BACKEND_DIR)/src/static || cp -R $(FRONTEND_DIR)/static/* $(BACKEND_DIR)/src/static
	# @cp -R $(FRONTEND_DIR)/static/* $(BACKEND_DIR)/src/static

# Then build the backend
backend-build:
	@cd $(BACKEND_DIR)/src && $(GOBUILD) $(GOFLAGS) -o ssdocker .

# Define the run target
run: all
	@cd $(BACKEND_DIR)/src && ./ssdocker -c config.yml

# Define the clean target
clean:
	@rm -rf $(BACKEND_DIR)/bin/*
	@rm -rf $(BACKEND_DIR)/src/static/*

.PHONY: all frontend-build backend-build run clean
