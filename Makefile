.PHONY: build test clean install lint fmt release-dry-run

BINARY_NAME=kubectl-rebalance
INSTALL_PATH=/usr/local/bin
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/kubectl-rebalance

test:
	go test -v -race ./...

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -rf dist/

install: build
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/

uninstall:
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run go fmt on all files
fmt:
	go fmt ./...

# Run linters
lint:
	go vet ./...
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

# Test goreleaser without publishing
release-dry-run:
	goreleaser release --snapshot --clean --skip-publish

# Create a new version tag
tag:
	@echo "Current tags:"
	@git tag -l | tail -5
	@read -p "Enter new version (e.g., v0.2.0): " version; \
	git tag -a $$version -m "Release $$version" && \
	echo "Created tag $$version. Push with: git push origin $$version"