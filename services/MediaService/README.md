Media Service
gRPC-микросервис для управления медиафайлами с использованием MinIO в качестве объектного хранилища.


Особенности
gRPC API - высокопроизводительный RPC-интерфейс

MinIO Integration - надежное хранение файлов в S3-совместимом хранилище

Dockerized - полная контейнеризация приложения

Web UI - интерактивный интерфейс для тестирования API через grpcui

Health Checks - автоматические проверки здоровья сервисов


Требования
Docker 20.10+

Docker Compose 2.0+


Быстрый старт
1. Клонирование и настройка
git clone <repository-url>
cd media_service


2. Запуск сервисов
docker-compose up --build

Сервисы будут доступны по следующим адресам:

Media Service gRPC: localhost:50051

MinIO Console: http://localhost:9001 (admin/minioadmin)

gRPC Web UI: http://localhost:8080


3. Проверка работы
Откройте браузер и перейдите на http://localhost:8080 для тестирования методов API.

Структура проекта
media_service/
├── cmd/media_service/     # Точка входа приложения
├── internal/              # Внутренние пакеты
│   ├── core/             # Бизнес-логика и модели
│   ├── grpc_server/      # gRPC сервер и хендлеры
│   ├── storage/          # Слой хранения (MinIO клиент)
│   └── utils/            # Вспомогательные функции
├── pkg/pb/               # Сгенерированные protobuf файлы
├── proto/                # Proto-файлы
├── config/               # Конфигурация приложения
├── tests/                # Тесты
└── scripts/              # Вспомогательные скрипты


API Методы

UploadMedia
Загрузка медиафайла в хранилище.

Request:

protobuf
message UploadMediaRequest {
  string user_id = 1;
  bytes file_bytes = 2;
  string mime_type = 3;
  string file_name = 4;
}
Response:

protobuf
message UploadMediaResponse {
  string media_id = 1;
  string message = 2;
  string url = 3;
}

GetMedia
Получение медиафайла по ID.

DeleteMedia
Удаление медиафайла.

ListMedia
Получение списка файлов пользователя.

GetUrl
Генерация временной ссылки для доступа к файлу.


Конфигурация
Настройки сервиса задаются через переменные окружения:

Переменная	По умолчанию	Описание
GRPC_PORT	50051	Порт gRPC сервера
MINIO_ENDPOINT	minio:9000	Адрес MinIO сервера
MINIO_ACCESS_KEY	minioadmin	Ключ доступа MinIO
MINIO_SECRET_KEY	minioadmin	Секретный ключ MinIO
MINIO_BUCKET	media-service	Бакет для хранения файлов


Тестирование
Запуск тестов
docker-compose exec media-service python -m pytest tests/
Тестовый клиент
python test.py


Генерация protobuf
Для перегенерации protobuf файлов после изменений в .proto:

./scripts/generate_proto.sh
Или вручную:

python -m grpc_tools.protoc \
    --proto_path=proto \
    --python_out=pkg/pb \
    --grpc_python_out=pkg/pb \
    proto/media.proto


Данные
Файлы: хранятся в MinIO бакете media-service

Метаданные: временно хранятся в памяти 


Отладка
Просмотр логов
docker-compose logs media-service
docker-compose logs minio
docker-compose logs grpcui

Проверка здоровья
docker-compose ps

Подключение к контейнеру
docker-compose exec media-service /bin/bash

Лицензия
MIT License