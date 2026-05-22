# Define variables
GO      := go
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean

# Define the name of the output binary
BINARY_NAME := dnslookup

.PHONY: all build clean

all: build ## Default target, builds the project

build: ## Compiles the Go program
	$(GOBUILD) -o $(BINARY_NAME) .

clean: ## Removes the generated binary and cache
	$(GOCLEAN)
	rm -f $(BINARY_NAME)