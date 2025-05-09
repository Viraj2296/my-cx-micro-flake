#!/bin/bash

# Detect the OS if not provided as an argument
if [ -z "$1" ]; then
  echo "No OS specified. Detecting OS..."
  case "$(uname -s)" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    CYGWIN*|MINGW*|MSYS*) OS="windows";;
    *)          echo "Unsupported OS detected. Exiting."; exit 1;;
  esac
  echo "Detected OS: $OS"
else
  OS=$1
fi

# Set the architecture variable
ARCH="amd64"  # Default architecture, modify as needed

# Set the output binary name and directory
OUTPUT_DIR="./bin"
OUTPUT_NAME="gateway_entry"

# Modify output name for Windows OS
if [ "$OS" == "windows" ]; then
  OUTPUT_NAME="$OUTPUT_NAME.exe"
fi

# Create the output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# Build the Go project
echo "Building for OS: $OS, Arch: $ARCH"
GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT_DIR/$OUTPUT_NAME -ldflags="-X main.version=1.0.4" ./cmd/gateway_entry.go

# Check if the build was successful
if [ $? -eq 0 ]; then
  echo "Build successful! Binary created at $OUTPUT_DIR/$OUTPUT_NAME"
else
  echo "Build failed!"
  exit 1
fi
