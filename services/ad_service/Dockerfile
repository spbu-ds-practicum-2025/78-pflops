# Stage 1: build Go binary
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ad_service ./cmd/ad-service

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/ad_service .
COPY .env .

EXPOSE ${AD_SERVICE_PORT}

CMD ["./ad_service"]
