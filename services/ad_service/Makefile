include .env

SHELL := cmd.exe
.SHELLFLAGS := /C

GO_MODULE=78-pflops/services/ad_service
OUT_DIR=pb
PROTO_DIR=proto
SERVICE=ad_service

.PHONY: proto deps clean help run build up down logs dbshell test-create-ad test-get-ad test-search-ad

proto:
	@echo Генерация gRPC и protobuf файлов...
	if not exist "$(OUT_DIR)" mkdir "$(OUT_DIR)"
	protoc --go_out=. --go_opt=module=$(GO_MODULE) \
	       --go-grpc_out=. --go-grpc_opt=module=$(GO_MODULE) \
	       $(PROTO_DIR)\*.proto
	@echo Генерация завершена

deps:
	@echo Установка плагинов protoc...
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo Плагины установлены

clean:
	if exist "$(OUT_DIR)\ad_service\pb\*.pb.go" del /Q "$(OUT_DIR)\ad_service\pb\*.pb.go"
	@echo Очищено

run:
	@echo Запуск сервиса...
	go run .\cmd\ad-service\main.go

build:
	@echo Сборка сервиса...
	if not exist ".\bin" mkdir ".\bin"
	go build -o .\bin\$(SERVICE).exe .\cmd\ad-service\main.go

up:
	@echo Запуск Docker контейнеров...
	docker compose up --build -d

down:
	@echo Остановка Docker контейнеров...
	docker compose down

logs:
	@echo Просмотр логов...
	docker compose logs -f

dbshell:
	@echo Подключение к PostgreSQL...
	docker exec -it ad_service_db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

# Тестовые вызовы через grpcurl (порт локальный или контейнерный)

test-create-ad:
	powershell -Command "grpcurl -plaintext -d '{\"title\":\"Test Ad\",\"description\":\"Desc\",\"price\":12345,\"category_id\":\"00000000-0000-0000-0000-000000000000\",\"condition\":\"NEW\"}' localhost:50052 ad.AdService.CreateAd"

test-get-ad:
	powershell -Command "grpcurl -plaintext -d '{\"id\":\"REPLACE_ID\"}' localhost:50052 ad.AdService.GetAd"

test-search-ad:
	powershell -Command "grpcurl -plaintext -d '{\"text\":\"Test\",\"page\":1,\"page_size\":10}' localhost:50052 ad.AdService.SearchAds"

help:
	@echo Доступные команды:
	@echo   make proto          - Генерация gRPC и protobuf файлов
	@echo   make deps           - Установка плагинов protoc
	@echo   make clean          - Очистка сгенерированных файлов
	@echo   make run            - Запуск (без Docker)
	@echo   make build          - Сборка бинарника
	@echo   make up             - Запуск Docker окружения
	@echo   make down           - Остановка Docker окружения
	@echo   make logs           - Логи контейнеров
	@echo   make dbshell        - Консоль PostgreSQL
	@echo   make test-create-ad - Тест создания объявления
	@echo   make test-get-ad    - Тест получения объявления (заменить REPLACE_ID)
	@echo   make test-search-ad - Тест поиска объявлений
