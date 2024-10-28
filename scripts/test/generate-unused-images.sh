#!/bin/bash

# generate-unused-images.sh

# This script automates the generation of test Docker images for the Dockclean application.
# It creates multiple Docker images with different version labels, simulating untagged images for testing purposes.
# The script simplifies the process of recreating test images after cleanup operations.
# Ensure that a valid Dockerfile is present in the current directory before running the script.

# Array of versions to build
versions=("1.0" "2.0" "3.0" "4.0")

# Build Docker images one by one
for i in "${!versions[@]}"; do
    version=${versions[$i]}
    
    # Build Docker image with the current version label
    docker build --no-cache --label tag="for-test" --label version="$version" .

    # Display build completion message
    echo -e "\n---------------------------------------------------------"
    echo -e "Build complete for Docker image version: $version"
    echo -e "---------------------------------------------------------\n"

    # Check if there are more versions to build
    if [ $i -lt $((${#versions[@]} - 1)) ]; then
        echo "Starting to build the next Docker image ..."
    fi
done
