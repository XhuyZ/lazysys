.PHONY: build run clean install

# Build the application
build:
	go build -o lazysys .

# Run the application (requires sudo)
run: build
	sudo ./lazysys

# Clean build artifacts
clean:
	rm -f lazysys

# Install to system (requires sudo)
install: build
	sudo cp lazysys /usr/local/bin/
	sudo chmod +x /usr/local/bin/lazysys

# Uninstall from system (requires sudo)
uninstall:
	sudo rm -f /usr/local/bin/lazysys

# Get dependencies
deps:
	go mod tidy

# Run tests
test:
	go test ./...

# Build for release
release: clean
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o lazysys-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o lazysys-linux-arm64 . 