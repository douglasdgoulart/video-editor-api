FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /video-editor-app /app/cmd/app/app.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY bin/ /bin
COPY --from=build-stage /video-editor-app /video-editor-app
COPY config.yaml /config.yaml

EXPOSE 8080

ENTRYPOINT ["/video-editor-app"]
