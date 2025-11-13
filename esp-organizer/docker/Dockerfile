# build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.buildVersion=$(git rev-parse --short HEAD 2>/dev/null || echo 'dev')" \
    -o /server cmd/server/main.go

# final stage
FROM alpine:latest

# Install required packages in a single layer
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    aws-cli \
    curl \
    wget \
    && update-ca-certificates

# Create non-root user for security
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /server .

# Create logs directory with proper permissions
RUN mkdir -p /app/logs && \
    chown -R appuser:appuser /app/logs

# Switch to non-root user
USER appuser

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Define the command to run your app
CMD ["./server"]