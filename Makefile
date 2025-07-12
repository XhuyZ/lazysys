.PHONY: build run clean install uninstall deps test release

# Build the application
build:
	mkdir -p build
	go build -o build/lazysys ./src/...

# Run the application (requires sudo)
run: 
	sudo ./build/lazysys

# Clean build artifacts
clean:
	rm -f lazysys lazysys-linux-*

# Install to system (requires sudo)
install: build
	sudo cp lazysys /usr/local/bin/
	sudo chmod +x /usr/local/bin/lazysys

# Uninstall from system (requires sudo)
uninstall:
	sudo rm -f /usr/local/bin/lazysys

# Get dependencies
deps:
	cd src && go mod tidy

# Run tests
test:
	cd src && go test ./...

# Build for release
release: clean
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/lazysys-linux-amd64 ./src
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o build/lazysys-linux-arm64 ./src

