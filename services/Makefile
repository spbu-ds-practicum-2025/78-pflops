.PHONY: up down logs ps build test proto help

# Compose-level commands for all services

up:
	@echo "Запуск всех сервисов (user, ad, media)..."
	docker compose up --build -d

down:
	@echo "Остановка всех сервисов..."
	docker compose down

logs:
	@echo "Просмотр логов всех сервисов..."
	docker compose logs -f

ps:
	@echo "Список контейнеров..."
	docker compose ps

build:
	@echo "Сборка всех Go-сервисов..."
	$(MAKE) -C user_service build
	$(MAKE) -C ad_service build

proto:
	@echo "Генерация protobuf/gRPC для всех Go-сервисов..."
	$(MAKE) -C user_service proto
	$(MAKE) -C ad_service proto

test:
	@echo "Запуск тестов user_service и ad_service..."
	cd user_service && go test ./...
	cd ad_service && go test ./...

help:
	@echo "Доступные команды (общий Makefile в services):"
	@echo "  make up       - Запуск всех сервисов (docker-compose)"
	@echo "  make down     - Остановка всех сервисов"
	@echo "  make logs     - Логи всех сервисов"
	@echo "  make ps       - Список контейнеров"
	@echo "  make build    - Сборка бинарников user_service и ad_service"
	@echo "  make proto    - Генерация protobuf для user_service и ad_service"
	@echo "  make test     - Запуск всех Go-тестов в user_service и ad_service"