# =============================================================================
# Builder stage: Compile Go application with CGO enabled for DuckDB
# =============================================================================
FROM golang:1.25.1-alpine AS builder

# Install build dependencies for CGO and DuckDB
RUN apk add --no-cache \
    ca-certificates \
    git \
    build-base \
    musl-dev

WORKDIR /src

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
# - CGO_ENABLED=1: Required for go-duckdb driver
# - ldflags="-s -w": Strip debug info to reduce binary size
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/device-api ./cmd/app

# =============================================================================
# Runtime stage: Minimal Alpine image with runtime dependencies only
# =============================================================================
FROM alpine:latest

# Install runtime dependencies
# - ca-certificates: for HTTPS connections
# - libstdc++/libgcc: required by DuckDB
# - curl: for health checks
RUN apk add --no-cache \
    ca-certificates \
    libstdc++ \
    libgcc \
    curl

# Create non-root user and group for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /out/device-api /app/device-api

# Create data directory for persistent storage and symlink for backward compatibility
# The app expects ./devices.db, but we store it in /data for volume mounting
RUN mkdir -p /data && \
    ln -s /data/devices.db /app/devices.db && \
    chown -R appuser:appgroup /app /data

# Environment variables (can be overridden at runtime)
ENV DEVICE_DB_PATH=/data/devices.db
ENV PORT=8080

# Expose application port
EXPOSE 8080

# Switch to non-root user
USER appuser

# Health check (uses the list devices endpoint)
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -fsS http://localhost:8080/v1/devices || exit 1

# Start the application
ENTRYPOINT ["/app/device-api"]
