# Ad Service (Moonshine Marketplace)

## Назначение
Сервис отвечает за управление объявлениями: создание, обновление, удаление, поиск, изображения, избранное и отзывы.

## Стек
- Go 1.25
- PostgreSQL (отдельная база)
- gRPC

## Порты
- gRPC: `50052`

## Быстрый старт (Dev)
```bash
# Зайдите в каталог
cd services/ad_service
# Сгенерируйте бинарник
go build ./cmd/ad-service
# Запустите (предварительно должен быть Postgres)
./ad-service
```

## Тесты
- Запуск всех тестов:
```powershell
Set-Location -Path services/ad_service; go test ./...
```
- Запуск только тестов доменного сервиса:
```powershell
Set-Location -Path services/ad_service; go test -count=1 ./internal/service
```

## API сервиса (Go)
Минимальный интерфейс бизнес-логики в `internal/service/ad_service.go`:

```go
CreateAd(ctx context.Context, userID, title, description string, price int64) (*model.Ad, error)
GetAd(ctx context.Context, adID string) (*model.Ad, error)
ListAds(ctx context.Context, f Filters) ([]model.Ad, int, error)
UpdateAd(ctx context.Context, adID, userID string, title, description *string, price *int64) error
DeleteAd(ctx context.Context, adID, userID string) error
AttachMedia(ctx context.Context, adID, mediaID string) error
```

Примечания:
- `Filters` — простой фильтр (текст, категория, границы цены, состояние, пагинация).
- Для простоты `CreateAd` выставляет дефолты: `Condition=NEW`, `CategoryID=00000000-0000-0000-0000-000000000000`.
- `AttachMedia` сохраняет `mediaID` как URL в таблицу `ad_images`.

## Структура
```
internal/
  model/        # Доменные сущности
  db/           # Подключение и миграции
  repository/   # Доступ к данным (позже)
  service/      # Бизнес-логика (позже)
  ...
cmd/
  ad-service/   # Точка входа
```

## Миграции
SQL файлы в `internal/db/migrations`. Использовать golang-migrate или аналог.

## Следующие шаги
1. Добавить `ad.proto` и сгенерировать gRPC код
2. Реализовать репозиторий и сервисный слой
3. Встроить проверку JWT и интеграцию с User Service
4. Добавить загрузку изображений через Media Service
5. Кэширование и события

## Лицензия
См. корневой `LICENSE`.

## Обновление gRPC API
Методы в gRPC UI берутся из скомпилированного protobuf. После изменения `proto/ad.proto` нужно пересгенерировать gRPC код и пересобрать сервис.

1) Установить генераторы (один раз):
```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
$env:PATH = ("$(go env GOPATH)\bin;" + $env:PATH)
```

2) Сгенерировать код (запустить в `services/ad_service`):
```powershell
protoc -I proto --go_out=. --go-grpc_out=. proto/ad.proto
```

3) Собрать и запустить:
```powershell
go build ./cmd/ad-service
./ad-service
```

В API сейчас доступны RPC:
- CreateAd(user_id, title, description, price)
- GetAd(ad_id)
- ListAds(filters)
- UpdateAd(ad_id, user_id, title?, description?, price?)
- DeleteAd(ad_id, user_id)
- AttachMedia(ad_id, media_id)
