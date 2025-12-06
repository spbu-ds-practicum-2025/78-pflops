# Stage 1: build Go binary
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user_service ./cmd/user-service

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/user_service .
COPY .env .

EXPOSE ${USER_SERVICE_PORT}

CMD ["./user_service"]
