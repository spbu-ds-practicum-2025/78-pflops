# Stage 1: build Go binary
FROM golang:1.25.1 AS builder

WORKDIR /app

# Install protoc compiler and dependencies
RUN apt-get update && apt-get install -y \
    protobuf-compiler \
    unzip \
    && rm -rf /var/lib/apt/lists/*

COPY . .

RUN go mod tidy
ENV PATH="$PATH:$(go env GOPATH)/bin"

RUN make deps
RUN make proto

RUN go mod download

RUN go build -o user_service ./cmd/user-service

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/user_service .
COPY .env .

EXPOSE ${USER_SERVICE_PORT}

CMD ["./user_service"]
