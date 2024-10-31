# Use the official Golang image based on Alpine for the build stage
FROM golang:1.23.2-alpine AS builder

# Install necessary dependencies for building and running the app
RUN apk add --no-cache \
    git \
    chromium \
    chromium-chromedriver \
    nss \
    freetype \
    alsa-lib \
    libx11 \
    libxcomposite \
    libxcursor \
    libxdamage \
    libxi \
    libxtst \
    libxrandr \
    pango \
    ttf-freefont

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application code
COPY . .

# Build the Go app
RUN go build -o app .

# Create a new stage for running the application
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    alsa-lib \
    libx11 \
    libxcomposite \
    libxcursor \
    libxdamage \
    libxi \
    libxtst \
    libxrandr \
    pango \
    ttf-freefont

# Copy the binary from the builder stage
COPY --from=builder /app/app /app/app

# Set the working directory
WORKDIR /app

# Set the environment variable for the port
ENV PORT=10000

# Expose the port
EXPOSE 10000

# Run the app
CMD ["./app"]
