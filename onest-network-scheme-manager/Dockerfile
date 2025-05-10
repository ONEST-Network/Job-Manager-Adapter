FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scheme-manager-adapter .

# Use a smaller image for the final build
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/scheme-manager-adapter .

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./scheme-manager-adapter"]