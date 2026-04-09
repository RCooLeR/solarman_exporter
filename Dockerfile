# --- Stage 1: Build the Go binary ---
FROM golang:1.26.2-alpine AS builder

# Install git and CA certs (if needed for Go modules)
RUN apk add --no-cache git ca-certificates

# Set working directory inside the container
WORKDIR /app

# Copy Go source code into the container
COPY . .

# Enable Go Modules (optional if using go.mod)
ENV GO111MODULE=on

# Download deps
RUN go mod download

# Build the Go app for Linux
RUN go build -o solarman-exporter

# --- Stage 2: Create a lightweight image with the binary only ---
FROM alpine:3.23

# Add CA certs for HTTPS support
RUN apk upgrade --no-cache && apk add --no-cache ca-certificates

# Copy binary from builder stage
COPY --from=builder /app/solarman-exporter /usr/bin/solarman-exporter

# Set executable permissions (just in case)
RUN chmod +x /usr/bin/solarman-exporter

# Command to run
ENTRYPOINT ["/usr/bin/solarman-exporter"]