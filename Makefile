# Variables
BINARY_NAME=app
BUILD_DIR=bin
CMD_DIR=cmd/app
MAIN_FILE=$(CMD_DIR)/app.go

FFMPEG_VERSION=release
FFMPEG_BUILD=amd64
FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-$(FFMPEG_VERSION)-$(FFMPEG_BUILD)-static.tar.xz
FFMPEG_DIR=bin/ffmpeg

API_NAMESPACE=default
JOB_NAMESPACE=default
KAKFA_NAMESPACE=kafka

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
test: ffmpeg
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

deploy-api:
	@helm upgrade --install video-editor-api --namespace $(API_NAMESPACE) deploy/helm/app/chart/ --values deploy/helm/app/api.yaml --wait --timeout 2m

deploy-job:
	@helm upgrade --install video-editor-job --namespace $(JOB_NAMESPACE) deploy/helm/app/chart/ --values deploy/helm/app/job.yaml --wait --timeout 2m

install-kafka:
	@kubectl create namespace $(KAKFA_NAMESPACE) --dry-run -o yaml | kubectl apply -f -
	@helm upgrade --install $(KAKFA_NAMESPACE) oci://registry-1.docker.io/bitnamicharts/kafka \
		--namespace kafka --set listeners.client.protocol=PLAINTEXT \
		--set listeners.controller.protocol=PLAINTEXT \
		--set listeners.external.protocol=PLAINTEXT \
		--wait --timeout 2m
	@kubectl run kafka-create-topic --rm -i --tty --namespace $(KAKFA_NAMESPACE) --image=bitnami/kafka:latest -- \
		kafka-topics.sh --create --if-not-exists --bootstrap-server kafka.kafka.svc.cluster.local:9092 \
		--replication-factor 3 --partitions 100 --topic event

uninstall-kafka:
	@helm uninstall kafka --namespace kafka

.PHONY: all build clean run test lint deploy-api install-kafka uninstall-kafka
