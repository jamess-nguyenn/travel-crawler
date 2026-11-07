# ============================================================
# Development stage
# ============================================================
FROM golang:1.27 AS dev

RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Install Goose migration tool
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

WORKDIR /app

EXPOSE 8000
EXPOSE 2345


# ============================================================
# Build stage
# ============================================================
FROM golang:1.27 AS builder

WORKDIR /app

# Download dependencies first for better Docker caching
COPY go.mod go.sum ./

RUN go mod download

# Copy application source
COPY . .

# Build accommodation worker
RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin/accommodation-worker \
    ./cmd/worker/accommodation


# ============================================================
# Production accommodation worker stage
# ============================================================
FROM debian:bookworm-slim AS worker

RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy compiled binary
COPY --from=builder /app/bin/accommodation-worker \
    /app/accommodation-worker

# Run as non-root user
RUN useradd --system --no-create-home --uid 10001 appuser \
    && chown appuser:appuser /app/accommodation-worker

USER appuser

ENTRYPOINT ["/app/accommodation-worker"]
