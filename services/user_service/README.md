
## User Service

Краткое описание

- Назначение: сервис управления пользователями (регистрация, авторизация, базовая валидация токена и профиль).
- Транспорт: gRPC.
- Хранилище: PostgreSQL.
- Аутентификация: JWT (HS256), хеширование паролей — bcrypt.

API (gRPC)

- Register(email, password, name) -> { id, token, email, name }
- Login(email, password) -> { token }
- ValidateToken(token) -> { user_id, valid }
- GetProfile(user_id) -> { user_id, name }
- UpdateProfile(user_id, name) -> { success, message }
- DeleteUser(user_id, password) -> { success, message }

Где посмотреть контракты: `proto/user.proto` (сгенерированные файлы: `gen/proto`).

Архитектура каталогов

- `internal/db` — подключение к PostgreSQL и инициализация схемы.
- `internal/repository` — доступ к данным (pgx/pool).
- `internal/service` — бизнес-логика (валидации, хеши, токены).
- `cmd/user-service` — gRPC-сервер и регистрация хендлеров.
- `proto`, `gen/proto` — protobuf-описания и сгенерированный код.

Тесты

- Юнит‑тесты: `go test ./internal/service`
- Покрываются кейсы регистрации/логина, проверки валидаций и генерации токена.

### Настройка:

1. Настройте переменные окружения.
Находясь в ``/user_service`` создайте файл ``.env`` с содержимым аналогично файлу ``.env.example``

2. Установите необходимые зависимости для Go.
Находясь в ``/user_service`` выполните команду:
```
go mod tidy
```

3. Убедитесь, что ``protoc-gen-go`` и ``protoc-gen-go-grpc`` находятся в GOPATH:
Чтобы проверить это, выполните:
```
which protoc-gen-go
which protoc-gen-go-grpc
```
Если ответ пустой, вам нужно добавить ``protoc-gen-go`` и ``protoc-gen-go-grpc`` в GOPATH:
``export PATH="$PATH:$(go env GOPATH)/bin"``
После этого вы можете сгенерировать файлы.

4. Генерация proto файлов
В /user_service запустите:
```
make deps // установит плагины protobuf и gRPC
make proto // сгенерирует файлы proto
```

5. Теперь вы можете запустить контейнер:
```
make up
```

### Управление сервисом:
Просмотреть список доступных команд:
```
make help
```

## User Service

### Setup:

1. Configure environment variables.
While in the ``/user_service`` directory, create a ``.env`` file with content similar to the ``.env.example`` file.

2. Install the necessary dependencies for Go.
While in the ``/user_service`` directory, run the command:
```
go mod tidy
```

3. Ensure that ``protoc-gen-go`` and ``protoc-gen-go-grpc`` are in your GOPATH:
To check this, run:
```
which protoc-gen-go
which protoc-gen-go-grpc
```
If the output is empty, you need to add ``protoc-gen-go`` and ``protoc-gen-go-grpc`` to your GOPATH:
``export PATH="$PATH:$(go env GOPATH)/bin"``
After this, you can generate the files.

4. Generate proto files
In the /user_service directory, run:
```
make deps // will install protobuf and gRPC plugins
make proto // will generate the proto files
```

5. Now you can build the container:
```
make up
```

### Service Management:
To view the list of available commands:
```
make help
```
