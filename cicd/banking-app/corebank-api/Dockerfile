# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Enable Go module caching
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOMODCACHE=/go/pkg/mod

# Cache dependencies first
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download -x || \
    { sleep 5 && go mod download -x; } || \
    { sleep 10 && go mod download -x; }

# Copy remaining files
COPY . .

# Build with optimizations
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -trimpath -o bank-api .

# Stage 2: Minimal runtime image
FROM alpine:3.21

WORKDIR /app

# Copy only necessary files from builder
COPY --from=builder /app/bank-api .
COPY --from=builder /app/.env .env
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Create non-root user
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080

ENTRYPOINT ["./bank-api"]

# # Stage 1: Build
# FROM golang:1.22-alpine AS builder

# # Install Git (required for go get) and necessary tools
# RUN apk add --no-cache git

# WORKDIR /app

# # Cache dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the rest of the source code
# COPY . .

# # Copy .env if needed
# COPY .env .env

# # Build the binary from the root
# RUN go build -o bank-api .

# # Stage 2: Minimal image to run the app
# FROM alpine:3.19

# # Install ca-certificates for HTTPS and timezone data
# RUN apk add --no-cache ca-certificates tzdata

# WORKDIR /app

# # Copy the binary from the builder
# COPY --from=builder /app/bank-api .

# # Copy .env if needed
# COPY --from=builder /app/.env .env

# # Expose the port used by the Go app
# EXPOSE 8080

# # Run the app
# ENTRYPOINT ["./bank-api"]
