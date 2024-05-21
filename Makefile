# Variables
API_BINARY_NAME=api
JOB_BINARY_NAME=job
BUILD_DIR=bin
API_CMD_DIR=cmd/api
JOB_CMD_DIR=cmd/job
API_MAIN_FILE=$(API_CMD_DIR)/api.go
JOB_MAIN_FILE=$(JOB_CMD_DIR)/job.go

FFMPEG_VERSION=release
FFMPEG_BUILD=amd64
FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-$(FFMPEG_VERSION)-$(FFMPEG_BUILD)-static.tar.xz
FFMPEG_DIR=bin/ffmpeg

# first so "all" becomes default target
all: clean build

# Build the API binary
build-api:
	@echo "Building the API binary"
	@go build -o $(BUILD_DIR)/$(API_BINARY_NAME) $(API_MAIN_FILE)

# Build the Job binary
build-job:
	@echo "Building the Job binary"
	@go build -o $(BUILD_DIR)/$(JOB_BINARY_NAME) $(JOB_MAIN_FILE)

# Build both binaries
build: build-api build-job

# Clean the build directory
clean:
	@echo "Cleaning the build directory"
	@rm -rf $(BUILD_DIR)/*

# Run the API application
run-api: build-api
	@echo "Running the API application"
	@./$(BUILD_DIR)/$(API_BINARY_NAME)

# Run the Job application
run-job: build-job
	@echo "Running the Job application"
	@./$(BUILD_DIR)/$(JOB_BINARY_NAME)

# Run tests
test:
	@echo "Running tests"
	@go test ./...

# Run linting (requires golangci-lint installed)
lint:
	@echo "Linting the code"
	@golangci-lint run

# Download and extract FFmpeg
$(FFMPEG_DIR):
	@echo "Downloading FFmpeg..."
	@mkdir -p $(FFMPEG_DIR)
	@curl -L $(FFMPEG_URL) -o ffmpeg.tar.xz
	@tar -xvf ffmpeg.tar.xz -C $(FFMPEG_DIR) --strip-components=1
	@rm ffmpeg.tar.xz

# Ensure FFmpeg is available
ffmpeg: $(FFMPEG_DIR)

.PHONY: all build-api build-job build clean run-api run-job test lint
