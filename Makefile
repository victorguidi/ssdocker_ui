# Define the compiler and compiler flags
GO := go
GOBUILD := $(GO) build
GOFLAGS := -v
FRONTEND_DIR := frontend
BACKEND_DIR := backend

# Define the build targets
all:check frontend-build backend-build

# Test if golang and nodejs are installed
check:
	@which go || (echo "Go is not installed" && exit 1)
	@which npm || (echo "NodeJS is not installed" && exit 1)

# Build the frontend first
frontend-build:
	@cd $(FRONTEND_DIR) && npm install
	@cd $(FRONTEND_DIR) && npm run build
	# @cd $(FRONTEND_DIR) && npm run build
	@test -d $(BACKEND_DIR)/src/static || mkdir $(BACKEND_DIR)/src/static
	@cp -R $(FRONTEND_DIR)/static/* $(BACKEND_DIR)/src/static

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
