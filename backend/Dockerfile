# Build stage

FROM --platform=linux/amd64 golang:1.24-alpine AS builder


WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

# Create data directory
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./server"]
