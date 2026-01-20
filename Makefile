.PHONY: build run test clean install help

# Default target
help:
	@echo "LibreOffice-as-a-Service (Go version)"
	@echo ""
	@echo "Available targets:"
	@echo "  make build    - Build the Go binary"
	@echo "  make run      - Run the application"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make install  - Install dependencies and build"
	@echo "  make test     - Run tests"
	@echo ""
	@echo "For Node.js version, use: npm run start"

build:
	go build -o libreoffice-as-a-service .

run: build
	./libreoffice-as-a-service

clean:
	rm -f libreoffice-as-a-service

install:
	go mod download
	go build -o libreoffice-as-a-service .

test:
	go test ./...
