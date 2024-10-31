# Use the official Golang image as a build stage
FROM golang:1.23.2-bullseye AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application code
COPY . .

# Build the Go app
RUN go build -o app .

# Use a lightweight image to run the app
FROM debian:bullseye-slim

# Install Chrome dependencies
RUN apt-get update && apt-get install -y \
    curl \
    libnss3 \
    libxss1 \
    libasound2 \
    libx11-xcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxi6 \
    libxtst6 \
    libxrandr2 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libgtk-3-0 \
    chromium && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary from the build stage
COPY --from=builder /app/app /app/app

# Set the working directory and specify the default port
WORKDIR /app
ENV PORT=10000

# Expose the port
EXPOSE 10000

# Run the app
CMD ["./app"]
