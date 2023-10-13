#!/bin/sh

# Enable the 'exit on error' option
set -e

# Install helm-diff plugin
helm plugin install https://github.com/databus23/helm-diff --version v3.8.1

# Set URL for latest release of helmfile
URL=https://github.com/helmfile/helmfile/releases/download/v0.157.0/helmfile_0.157.0_linux_amd64.tar.gz

# Download helmfile binary
mkdir helmfile
wget -c $URL -O - | tar -xz -C helmfile

# Check if current user has write permission to /usr/local/bin
if [ -w /usr/local/bin ]
then
  # If write permission is present, set SUDO to empty string
  SUDO=""
else
  # If write permission is not present, set SUDO to "sudo"
  SUDO=sudo
fi

# Make helmfile binary executable
$SUDO chmod +x ./helmfile/helmfile

# Move helmfile binary to /usr/local/bin
$SUDO mv ./helmfile/helmfile /usr/local/bin

# Remove helmfile folder
rm -rf ./helmfile
