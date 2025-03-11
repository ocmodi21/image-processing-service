# Dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Use a small alpine image for the final container
FROM alpine:3.15

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Copy store master data
COPY store_master.csv .

# Expose port
EXPOSE 8080

# Run the server
CMD ["/app/server", "-addr=:8080", "-store-path=./store_master.csv", "-workers=4"]
