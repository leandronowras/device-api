# =============================================================================
# Builder stage: Compile Go application with CGO enabled for DuckDB
# =============================================================================
FROM golang:1.25-bookworm AS builder

# Install build dependencies for CGO and DuckDB
RUN apt-get update && apt-get install -y \
    ca-certificates \
    git \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
# - CGO_ENABLED=1: Required for go-duckdb driver
# - ldflags="-s -w": Strip debug info to reduce binary size
RUN CGO_ENABLED=1 \
    go build -ldflags="-s -w" -o /out/device-api ./cmd/app

# =============================================================================
# Runtime stage: Minimal Debian image with runtime dependencies only
# =============================================================================
FROM debian:bookworm-slim

# Install runtime dependencies
# - ca-certificates: for HTTPS connections
# - libstdc++6: required by DuckDB
# - curl: for health checks and scripts
# - jq: for JSON processing in scripts
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libstdc++6 \
    curl \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user and group for security
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /out/device-api /app/device-api

# Copy scripts folder
COPY scripts /app/scripts

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
