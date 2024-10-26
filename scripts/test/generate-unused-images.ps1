# generate-unused-images.ps1

# This script automates the generation of test Docker images for the Dockclean application. 
# It creates multiple Docker images with different version labels, simulating untagged images for testing purposes.
# The script simplifies the process of recreating test images after cleanup operations.
# Ensure that a valid Dockerfile is present in the current directory before running the script.

# Array of versions to build
$versions = @("1.0", "2.0", "3.0", "4.0")

# Initialize the index for tracking progress
$index = 0

# Build Docker images one by one
foreach ($version in $versions) {
    # Build Docker image with the current version label
    docker build --no-cache --label tag="for-test" --label version="$version" .

    # Display build completion message
    Write-Host "`n---------------------------------------------------------`n"
    Write-Host "Build complete for Docker image version: $version" -ForegroundColor Green
    Write-Host "`n---------------------------------------------------------`n"

    # Check if there are more versions to build
    if ($index -lt ($versions.Length - 1)) {
        Write-Host "Starting to build the next Docker image ...`n" -ForegroundColor Green
    }

    # Increment the index
    $index++
}
