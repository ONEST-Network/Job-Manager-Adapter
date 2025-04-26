# syntax=docker/dockerfile:1
FROM golang:1.24.2

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o scheme-manager .

EXPOSE 8080

CMD ["./scheme-manager"]