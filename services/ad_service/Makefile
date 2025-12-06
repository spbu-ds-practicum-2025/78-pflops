include .env

GO_MODULE=78-pflops/services/ad_service
OUT_DIR=pb
PROTO_DIR=proto
SERVICE=ad_service
GOPATH=$(shell go env GOPATH)

.PHONY: proto deps clean help run build up down logs dbshell test-create-ad test-get-ad test-list-ads

proto:
	@echo "Генерация gRPC и protobuf файлов..."
	mkdir -p "$(OUT_DIR)"
	PATH="$(GOPATH)/bin:$(PATH)" protoc --go_out=. --go_opt=module=$(GO_MODULE) \
	       --go-grpc_out=. --go-grpc_opt=module=$(GO_MODULE) \
	       -I $(PROTO_DIR) \
	       $(PROTO_DIR)/*.proto
	@echo "Генерация завершена"

deps:
	@echo "Установка плагинов protoc..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Плагины установлены в $(GOPATH)/bin"

clean:
	@echo "Очистка сгенерированных файлов..."
	rm -rf "$(OUT_DIR)/ad_service/pb/"*.pb.go
	@echo "Очищено"

run:
	@echo "Запуск сервиса..."
	go run ./cmd/ad-service/main.go

build:
	@echo "Сборка сервиса..."
	mkdir -p ./bin
	go build -o ./bin/$(SERVICE) ./cmd/ad-service/main.go

up:
	@echo "Запуск Docker контейнеров..."
	docker-compose up --build -d

down:
	@echo "Остановка Docker контейнеров..."
	docker-compose down

logs:
	@echo "Просмотр логов..."
	docker-compose logs -f

dbshell:
	@echo "Подключение к PostgreSQL..."
	docker exec -it ad_service_db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

# Тестовые вызовы через grpcurl

test-create-ad:
	grpcurl -plaintext -d '{"user_id":"550e8400-e29b-41d4-a716-446655440000","title":"Test Ad","description":"Test Description","price":12345}' localhost:50052 ad.AdService.CreateAd

test-get-ad:
	grpcurl -plaintext -d '{"id":"REPLACE_ID"}' localhost:50052 ad.AdService.GetAd

test-list-ads:
	grpcurl -plaintext -d '{"text":"","page":1,"page_size":10}' localhost:50052 ad.AdService.ListAds

help:
	@echo "Доступные команды:"
	@echo "  make proto          - Генерация gRPC и protobuf файлов"
	@echo "  make deps           - Установка плагинов protoc"
	@echo "  make clean          - Очистка сгенерированных файлов"
	@echo "  make run            - Запуск (без Docker)"
	@echo "  make build          - Сборка бинарника"
	@echo "  make up             - Запуск Docker окружения"
	@echo "  make down           - Остановка Docker окружения"
	@echo "  make logs           - Логи контейнеров"
	@echo "  make dbshell        - Консоль PostgreSQL"
	@echo "  make test-create-ad - Тест создания объявления"
	@echo "  make test-get-ad    - Тест получения объявления (заменить REPLACE_ID)"
	@echo "  make test-list-ads  - Тест поиска объявлений"
