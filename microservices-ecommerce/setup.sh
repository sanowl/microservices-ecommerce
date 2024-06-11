#!/bin/bash

# List of services
services=("user-service" "product-service" "order-service" "api-gateway")

# Iterate over each service
for service in "${services[@]}"; do
    echo "Setting up $service..."

    # Navigate to the service directory
    cd "$service"

    # Initialize the go module (if not already initialized)
    if [ ! -f go.mod ]; then
        go mod init "$service"
    fi

    # Install Gorilla Mux
    go get -u github.com/gorilla/mux

    # Ensure go.sum exists to prevent errors
    if [ ! -f go.sum ]; then
        go mod tidy
    fi

    # Add Gorilla Mux to go.mod (if not already present)
    if ! grep -q "github.com/gorilla/mux" go.mod; then
        echo "require github.com/gorilla/mux v1.8.0" >> go.mod
    fi

    # Ensure Dockerfile contains necessary commands
    if ! grep -q "RUN go mod download" Dockerfile; then
        echo "Updating Dockerfile for $service..."
        sed -i '' '/WORKDIR \/app/a\
COPY go.mod ./ \
COPY go.sum ./ \
RUN go mod download
' Dockerfile
    fi

    # Return to the root directory
    cd ..

    echo "$service setup completed."
done

echo "All services have been set up."
