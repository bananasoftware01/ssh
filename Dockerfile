# Multi-stage Dockerfile for Banana Software SSH Server
FROM golang:alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=1.0.0 -X main.Commit=$(git rev-parse --short HEAD 2>/dev/null || echo 'docker') -X main.Date=$(date -u +%Y-%m-%d)" \
    -o /app/banana-ssh ./cmd/banana-ssh

# Final lightweight runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S banana && adduser -S banana -G banana \
    && mkdir -p /data && chown -R banana:banana /data

WORKDIR /data

COPY --from=builder /app/banana-ssh /usr/local/bin/banana-ssh

USER banana

ENV PORT=2222 \
    HOST=0.0.0.0 \
    SSH_KEY_PATH=/data/keys/banana_ed25519 \
    TZ=Asia/Ho_Chi_Minh

EXPOSE 2222

VOLUME ["/data"]
ENTRYPOINT ["/usr/local/bin/banana-ssh"]
CMD ["serve"]
