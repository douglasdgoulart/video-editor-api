# Variables
BINARY_NAME=app
BUILD_DIR=bin
CMD_DIR=cmd/app
MAIN_FILE=$(CMD_DIR)/app.go

FFMPEG_VERSION=release
FFMPEG_BUILD=amd64
FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-$(FFMPEG_VERSION)-$(FFMPEG_BUILD)-static.tar.xz
FFMPEG_DIR=bin/ffmpeg

# first so "all" becomes default target
all: clean ffmpeg build

# Build the binary
build: ffmpeg
	@echo "Building the binary"
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Clean the build directory
clean:
	@echo "Cleaning the build directory"
	@rm -rf $(BUILD_DIR)/*

# Run the application
run: build
	@echo "Running the application"
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests"
	@go test -v -p 1 ./...

coverage-report:
	@echo "Running tests with coverage report"
	@go test -v -p 1 -coverprofile=coverage.out.tmp  ./... 
	@cat coverage.out.tmp | grep -v "_mock.go" > coverage.out
	@go tool cover -html=coverage.out

# Run linting (requires golangci-lint installed)
lint:
	@echo "Linting the code"
	@go fmt ./...
	@golangci-lint run

# Download and extract FFmpeg
$(FFMPEG_DIR):
	@echo "Downloading FFmpeg..."
	@mkdir -p $(FFMPEG_DIR)
	@curl -L $(FFMPEG_URL) -o ffmpeg-$(FFMPEG_VERSION)-$(FFMPEG_BUILD)-static.tar.xz
	@curl -L $(FFMPEG_URL).md5 -o ffmpeg-$(FFMPEG_VERSION)-$(FFMPEG_BUILD)-static.tar.xz.md5
	@echo "Verifying MD5 checksum..."
	@md5sum --quiet -c ffmpeg-$(FFMPEG_VERSION)-$(FFMPEG_BUILD)-static.tar.xz.md5
	@tar -xf ffmpeg-$(FFMPEG_VERSION)-$(FFMPEG_BUILD)-static.tar.xz -C $(FFMPEG_DIR) --strip-components=1
	@rm ffmpeg-$(FFMPEG_VERSION)-$(FFMPEG_BUILD)-static.tar.xz ffmpeg-$(FFMPEG_VERSION)-$(FFMPEG_BUILD)-static.tar.xz.md5

# Ensure FFmpeg is available
ffmpeg: $(FFMPEG_DIR)

.PHONY: all build clean run test lint
