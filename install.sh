#!/bin/bash

# Function to download a file
download_file() {
    local url=$1
    local output_path=$2
    echo "Downloading $url..."
    curl -L -o "$output_path" "$url"
    echo "Downloaded to $output_path"
}

# Function to determine the platform and architecture
get_platform_arch() {
    local platform=$(uname | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    case $arch in
        x86_64)
            arch="amd64"
            ;;
        arm64)
            arch="arm64"
            ;;
        i686)
            arch="386"
            ;;
        aarch64)
            arch="arm64"
            ;;
        *)
            echo "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    echo "$platform-$arch"
}

# GitHub repository details
REPO_OWNER="shashankx86"
REPO_NAME="att"

# Get the latest release from GitHub API
LATEST_RELEASE_URL="https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest"
echo "Fetching latest release info from $LATEST_RELEASE_URL..."
LATEST_RELEASE_JSON=$(curl -s $LATEST_RELEASE_URL)

# Extract the tag name and assets URL using grep, awk, and sed
TAG_NAME=$(echo "$LATEST_RELEASE_JSON" | grep -oP '"tag_name": "\K(.*?)(?=")')
ASSET_URLS=$(echo "$LATEST_RELEASE_JSON" | grep -oP '"browser_download_url": "\K(.*?)(?=")')

echo "Latest release: $TAG_NAME"

# Get the platform and architecture
PLATFORM_ARCH=$(get_platform_arch)

# Find the appropriate asset for the current platform and architecture
ASSET_URL=$(echo "$ASSET_URLS" | grep "att-$PLATFORM_ARCH")

if [[ -z "$ASSET_URL" ]]; then
    echo "No matching asset found for $PLATFORM_ARCH"
    exit 1
fi

# Download the asset
DOWNLOAD_DIR="/usr/local/bin"
DOWNLOAD_PATH="$DOWNLOAD_DIR/att"

# Create the directory if it doesn't exist
mkdir -p "$DOWNLOAD_DIR"

download_file "$ASSET_URL" "$DOWNLOAD_PATH"

# Make the file executable
chmod +x "$DOWNLOAD_PATH"

echo "att installed successfully in $DOWNLOAD_DIR"
echo "You can run it using the command: att"
