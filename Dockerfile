FROM golang:1.23-alpine

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Copy store master data
COPY store_master.csv .

# Expose port
EXPOSE 8080

# Run the server
CMD ["/app/server"]
