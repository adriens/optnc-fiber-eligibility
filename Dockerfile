# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY docs/ ./docs/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o opt-eligibility ./cmd/api

# Final stage - Minimal Chrome
FROM alpine:3.19

# Install only necessary packages
RUN apk add --no-cache \
    chromium \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# Create app user
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/opt-eligibility .

# Set ownership
RUN chown app:app /app/opt-eligibility

# Set chromium environment
ENV CHROMIUM_PATH=/usr/bin/chromium-browser
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROMEDP_DISABLE_SANDBOX=true

USER app

EXPOSE 8080

ENTRYPOINT ["/app/opt-eligibility"]
CMD ["api", "8080"]
