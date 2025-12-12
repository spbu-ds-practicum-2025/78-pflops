GO_MODULE=78-pflops/services/http_gateway
OUT_DIR=mediapb
PROTO_DIR=../MediaService/proto
SERVICE=http_gateway
GOPATH=$(shell go env GOPATH)

.PHONY: proto deps clean help run build up down logs

proto:
	@echo "Генерация gRPC и protobuf файлов for MediaService proto..."
	@mkdir -p $(OUT_DIR)
	PATH="$(GOPATH)/bin:$(PATH)" protoc --proto_path=$(PROTO_DIR) \
	       --go_out=$(OUT_DIR) --go_opt=module=$(GO_MODULE),Mmedia.proto=$(GO_MODULE)/mediapb \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=module=$(GO_MODULE),Mmedia.proto=$(GO_MODULE)/mediapb \
	       $(PROTO_DIR)/*.proto
	@echo "Генерация завершена"

deps:
	@echo "Установка плагинов protoc..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Плагины установлены в $(GOPATH)/bin"

clean:
	@echo "Очистка сгенерированных файлов..."
	-rm -f $(OUT_DIR)/*.pb.go
	@echo "Очищено"

run:
	@echo "Запуск http_gateway..."
	go run ./main.go

build:
	@echo "Сборка http_gateway..."
	mkdir -p ./bin
	go build -o ./bin/$(SERVICE) ./main.go

up:
	@echo "Запуск Docker окружения (compose at services root)..."
	docker compose up --build -d

down:
	@echo "Остановка Docker окружения..."
	docker compose down

logs:
	@echo "Просмотр логов (http_gateway)..."
	docker compose logs -f http_gateway

help:
	@echo "Доступные команды:"
	@echo "  make proto   - Генерация gRPC и protobuf файлов (MediaService proto)"
	@echo "  make deps    - Установка protoc плагинов"
	@echo "  make clean   - Очистка сгенерированных файлов"
	@echo "  make run     - Запуск (без Docker)"
	@echo "  make build   - Сборка бинарника"
	@echo "  make up      - Запуск Docker окружения"
	@echo "  make down    - Остановка Docker окружения"
	@echo "  make logs    - Логи (http_gateway)"
