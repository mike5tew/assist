FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go module files
# Changed: Removed 'esp-organizer/' prefix since we are now inside that folder
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
# Changed: Removed 'esp-organizer/' prefix
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies (including poppler-utils for pdftotext)
RUN apk --no-cache add ca-certificates poppler-utils

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Create non-root user for security
RUN addgroup -S appuser && adduser -S appuser -G appuser
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

EXPOSE 8080

CMD ["./server"]
