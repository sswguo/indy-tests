# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
BUILD_DIR=./build
# indy build-tests package
PACKAGE_NAME_BUILD=root
BINARY_NAME_BUILD=$(BUILD_DIR)/indy-test
# other indy package

build: 
# build indy-build-tests command
		$(GOBUILD) -trimpath -o $(BINARY_NAME_BUILD) -v ./cmd/$(PACKAGE_NAME_BUILD)
# build other commands 
#    $(GOBUILD) -trimpath -o $(BINARY_NAME_OTHER) -v ./cmd/$(PACKAGE_NAME_OTHER)
.PHONY: build

format:
		$(GOFMT) ./...
.PHONY: format

test: 
		$(GOTEST) -v ./...
.PHONY: test

clean: 
		$(GOCLEAN) ./...
		rm -rf $(BUILD_DIR)
.PHONY: clean


# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -o $(BINARY_NAME_BUILD)-linux -v ./cmd/$(PACKAGE_NAME_BUILD)
.PHONY: build-linux
