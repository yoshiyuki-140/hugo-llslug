#!/bin/bash

# If no version is specified as a command line argument, fetch the latest version.
if [ -z "$1" ]; then
    VERSION=$(curl -s https://api.github.com/repos/yoshiyuki-140/hugo-llslug/releases/latest | grep -o '"tag_name": *"[^"]*"' | sed 's/"tag_name": *"//' | sed 's/"//')
    if [ -z "$VERSION" ]; then
        echo "Failed to fetch the latest version"
        exit 1
    fi
else
    VERSION=$1
fi

OS=$(uname -s)
ARCH=$(uname -m)
URL="https://github.com/yoshiyuki-140/hugo-llslug/releases/download/${VERSION}/hugo-llslug_${OS}_${ARCH}.tar.gz"

echo "Start to install."
echo "VERSION=$VERSION, OS=$OS, ARCH=$ARCH"
echo "URL=$URL"

TMP_DIR=$(mktemp -d)
curl -L $URL -o $TMP_DIR/hugo-llslug.tar.gz
tar -xzvf $TMP_DIR/hugo-llslug.tar.gz -C $TMP_DIR
sudo mv $TMP_DIR/hugo-llslug /usr/local/bin/hugo-llslug
sudo chmod +x /usr/local/bin/hugo-llslug

rm -rf $TMP_DIR

if [ -f "/usr/local/bin/hugo-llslug" ]; then
  echo "[SUCCESS] hugo-llslug $VERSION has been installed to /usr/local/bin"
else
  echo "[FAIL] hugo-llslug $VERSION is not installed."
fi