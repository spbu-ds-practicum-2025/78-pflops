include .env

DB_NAME=${POSTGRES_DB}
DB_USER=${POSTGRES_USER}
GO_MODULE=78-pflops/services/user_service
OUT_DIR=gen
PROTO_DIR=proto
SERVICE=user_service

.PHONY: proto deps clean help

# Local development

proto:
	@echo "Генерация gRPC и protobuf файлов..."
	@mkdir -p $(OUT_DIR)
	protoc --go_out=$(OUT_DIR) --go_opt=module=$(GO_MODULE) \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=module=$(GO_MODULE) \
	       $(PROTO_DIR)/*.proto
	@echo "Генерация завершена"

deps:
	@echo "Установка плагинов protoc..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Плагины установлены"

clean:
	rm -f $(OUT_DIR)/$(PROTO_DIR)/*.pb.go
	@echo "Очищено"

run:
	@echo "Запуск сервиса..."
	go run ./cmd/$(SERVICE)/main.go

build:
	@echo "Сборка сервиса..."
	go build -o ./bin/$(SERVICE) ./cmd/$(SERVICE)/main.go

# Docker commands

up:
	@echo "Запуск Docker контейнера..."
	docker compose up --build -d

down:
	@echo "Остановка Docker контейнера..."
	docker compose down

logs:
	@echo "Просмотр логов Docker контейнера..."
	docker compose logs -f

dbshell:
	@echo "Запуск оболочки базы данных..."
	@echo "Для выхода из оболочки введите '\q'"
	docker exec -it user_service_db psql -U $(DB_USER) -d $(DB_NAME)


# Test

test-register:
	grpcurl -plaintext -d '{"email": "asaadsdffdf@ya.ru", "password": "1234"}' localhost:50051 user.UserService.Register

test-register-docker:
	grpcurl -plaintext -d '{"email": "asaadsdffdf@ya.ru", "password": "1234"}' user_service:50051 user.UserService.Register

test-login:
	grpcurl -plaintext -d '{"email": "test@example.com", "password": "password"}' localhost:50051 user.UserService.Login

test-login-docker:
	grpcurl -plaintext -d '{"email": "test@example.com", "password": "password"}' user_service:50051 user.UserService.Login

help:
	@echo "Доступные команды:"
	@echo "  make proto   - Генерация gRPC и protobuf файлов"
	@echo "  make deps    - Установка необходимых плагинов"
	@echo "  make clean   - Очистка сгенерированных файлов"
	@echo "  make run     - Запуск сервиса (без Docker)"
	@echo "  make build   - Сборка сервиса (без Docker)"
	@echo "  make up      - Запуск Docker контейнера"
	@echo "  make down    - Остановка Docker контейнера"
	@echo "  make logs    - Просмотр логов Docker контейнера"
	@echo "  make dbshell - Запуск оболочки базы данных"
	@echo "  make test-register    - Запуск теста регистрации"
	@echo "  make test-login       - Запуск теста входа"
	@echo "  make test-register-docker - Запуск теста регистрации в Docker"
	@echo "  make test-login-docker - Запуск теста входа в Docker"