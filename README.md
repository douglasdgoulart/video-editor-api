# Video Editor API

## Overview

Welcome to the Video Editor API! This project is designed to provide a scalable and distributed API for processing long-running FFmpeg tasks. The goal is to handle video processing efficiently by leveraging a combination of powerful tools like FFmpeg, Echo, and Kafka.

## Features

- **Scalable API**: Process long-running video tasks efficiently.
- **Distributed Processing**: Utilize Kafka for distributed task management.
- **Automation**: Automated build and testing using Makefile.
- **Configuration Management**: Easy configuration with YAML files.

## Table of Contents

- [Video Editor API](#video-editor-api)
  - [Overview](#overview)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [Setup](#setup)
  - [Usage](#usage)
  - [Testing](#testing)
  - [Contributing](#contributing)
  - [License](#license)

## Setup

Follow these steps to set up the project:

1. **Clone the repository**:
    ```bash
    git clone https://github.com/douglasdgoulart/video-editor-api.git
    cd video-editor-api
    ```

2. **Install dependencies**:
    Ensure you have Go installed. Then, run:
    ```bash
    go mod tidy
    ```

3. **Download FFmpeg**:
    ```bash
    make ffmpeg
    ```

4. **Configuration**:
    Adjust the `config.yaml` file according to your environment and requirements.

5. **Build the project**:
    ```bash
    make build
    ```

## Usage

1. **Run the application**:
    ```bash
    make run
    ```

2. **API Endpoints**:
    - Health Check: `GET /health`
    - Process Video: `POST /process` with form data including the video file and JSON configuration.

## Testing

Run tests to ensure everything is working correctly:
```bash
make test
```

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit them (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a new Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
