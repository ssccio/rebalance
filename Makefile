.PHONY: build test clean install

BINARY_NAME=kubectl-rebalance
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME) cmd/kubectl-rebalance/main.go

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

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