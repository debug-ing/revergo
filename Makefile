APP_NAME := revergo
SRC := cmd/main.go
OUTPUT := $(APP_NAME)

all: build

build: $(SRC)
	@echo "Building the project..."
	go build -o $(OUTPUT) $(SRC)

run: build
	@echo "Running the application..."
	./$(OUTPUT)

debug:
	@echo "Running the application with -race"
	go run -race $(SRC)

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm $(OUTPUT)

help:
	@echo "Makefile targets:"
	@echo "  build   - Build the Go project"
	@echo "  run     - Build and run the project"
	@echo "  debug   - Run the application with
	@echo "  test    - Run tests"
	@echo "  clean   - Remove build file"
	@echo "  help    - Show this help message"