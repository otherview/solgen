# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o solgen ./cmd/solgen

# Final stage
FROM alpine:3.18

# Create non-root user
RUN addgroup -g 1000 solgen && \
    adduser -D -s /bin/sh -u 1000 -G solgen solgen

# Copy the binary from builder stage
COPY --from=builder /build/solgen /usr/local/bin/solgen

# Set up working directory
WORKDIR /sources

# Make sure solgen is executable
RUN chmod +x /usr/local/bin/solgen

# Switch to non-root user
USER solgen

# Set entrypoint to solgen binary
ENTRYPOINT ["solgen"]

# Default command is help (so users can see available commands)
CMD ["--help"]