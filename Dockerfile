# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.21-alpine AS build

# Security: Install security updates and create non-root user
RUN apk update && apk upgrade && \
    apk add --no-cache \
    gcc \
    musl-dev \
    ca-certificates \
    tzdata && \
    adduser -D -s /bin/sh -u 1000 secretary

ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

WORKDIR /workspace

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build with security flags
RUN CGO_ENABLED=1 go build \
    -ldflags="-w -s -extldflags=-static" \
    -a -installsuffix cgo \
    -o ./secretary ./gateway/main.go

# -----------------------------------------------------------------------------
#  Main Stage
# -----------------------------------------------------------------------------
FROM alpine:3.19

# Security: Install security updates and create non-root user
RUN apk update && apk upgrade && \
    apk add --no-cache \
    ca-certificates \
    tzdata && \
    adduser -D -s /bin/sh -u 1000 secretary && \
    rm -rf /var/cache/apk/*

# Copy timezone data
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

# Create necessary directories with proper permissions
RUN mkdir -p /app/data /app/data/audit && \
    chown -R secretary:secretary /app

# Copy binary
COPY --from=build --chown=secretary:secretary /workspace/secretary /app/secretary
COPY --from=build --chown=secretary:secretary /workspace/banner.txt /app/banner.txt

# Security: Switch to non-root user
USER secretary

WORKDIR /app

# Security: Use non-root port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app/secretary", "--health-check"]

# Set default environment variables for security
ENV SECRETARY_HOST=0.0.0.0 \
    SECRETARY_PORT=8080 \
    SECRETARY_SECURE_COOKIES=true \
    SECRETARY_AUDIT_ENABLED=true \
    SECRETARY_SESSION_MAX_AGE=3600

ENTRYPOINT ["/app/secretary"]
