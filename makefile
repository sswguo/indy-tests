# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=./build
# indy build-tests package
PACKAGE_NAME_BUILD=buildtests
BINARY_NAME_BUILD=$(BUILD_DIR)/indy-build-tests
# indy promote-test package
PACKAGE_NAME_PROMOTE=promotetest
BINARY_NAME_PROMOTE=$(BUILD_DIR)/indy-promote-test

# other indy package

build: 
# build indy-build-tests command
		$(GOBUILD) -trimpath -o $(BINARY_NAME_BUILD) -v ./cmd/$(PACKAGE_NAME_BUILD)
# build indy-promote-test command
		$(GOBUILD) -trimpath -o $(BINARY_NAME_PROMOTE) -v ./cmd/$(PACKAGE_NAME_PROMOTE)
# build other commands 
#    $(GOBUILD) -trimpath -o $(BINARY_NAME_OTHER) -v ./cmd/$(PACKAGE_NAME_OTHER)
.PHONY: build

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
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -o $(BINARY_NAME_PROMOTE)-linux -v ./cmd/$(PACKAGE_NAME_PROMOTE)
.PHONY: build-linux
