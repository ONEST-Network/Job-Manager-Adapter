FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/server

# Default configuration
ENV PORT=8080 \
    LOG_LEVEL=info \
    LOG_PRETTY_PRINT=false \
    DB_MAX_OPEN_CONNS=25 \
    DB_MAX_IDLE_CONNS=25

EXPOSE 8080

CMD ["./main"] 