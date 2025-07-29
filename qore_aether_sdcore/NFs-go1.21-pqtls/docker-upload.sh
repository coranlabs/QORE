#!/bin/bash

# Set Docker registry details
DOCKER_REGISTRY="lakshyachopra"
DOCKER_TAG="v1"
DIRECTORY="/NFs-go-1.21"

# Log in to Docker (assuming you have environment variables DOCKER_USERNAME and DOCKER_PASSWORD set)
echo "Logging into Docker..."
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

# Loop through each folder in the directory and build, tag, and push the image
for folder in "$DIRECTORY"/*; do
    if [ -d "$folder" ]; then
        BASE_NAME=$(basename "$folder")
        IMAGE_NAME="${BASE_NAME}-tls"
        FULL_IMAGE_NAME="$DOCKER_REGISTRY/$IMAGE_NAME:$DOCKER_TAG"
        
        echo "Building Docker image for $IMAGE_NAME..."
        docker build -t "$FULL_IMAGE_NAME" "$folder"
        
        echo "Pushing Docker image $FULL_IMAGE_NAME..."
        docker push "$FULL_IMAGE_NAME"
    fi
done

echo "All images have been built and pushed successfully!"
