# Build stage
FROM golang:1.21-alpine AS builder

# Install git dan ca-certificates untuk dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/api

# Production stage
FROM alpine:3.19

# Install ca-certificates untuk HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Jakarta

# Create non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/main .

# Copy docs untuk Swagger
COPY --from=builder /app/docs ./docs

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run binary
CMD ["./main"]
