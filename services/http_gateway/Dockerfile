FROM golang:1.25-alpine AS builder

WORKDIR /app

# Сначала копируем только файлы модулей для кешируемого go mod download
COPY http_gateway/go.mod http_gateway/

# В http_gateway/go.mod есть replace на ../ad_service и ../user_service,
# поэтому в контексте сборки должны существовать эти директории.
COPY ad_service ad_service
COPY user_service user_service
COPY http_gateway http_gateway

WORKDIR /app/http_gateway
RUN go mod tidy

RUN go build -o /bin/http-gateway main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /bin/http-gateway /usr/local/bin/http-gateway

ENV HTTP_GATEWAY_PORT=8081
CMD ["http-gateway"]
