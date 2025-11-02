
  

## User Service

  

### Настройка:

1. Настройте переменные окружения.

Находясь в ``/user_service`` создайте файл ``.env`` с содержимым аналогично файлу ``.env.example``

  

2. Установите необходимые зависимости для Go.

Находясь в  ``/user_service`` выполните команду:

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

В /user_service запустите:

```

make deps // установит плагины protobuf и gRPC

make proto // сгенерирует файлы proto

```

4. Теперь вы можете скомпилировать программу:

```

go build ./cmd/user_service

```

  

### (English) Setup:

1. Configure environment variables
Create ``.env`` file (see ``.env.example``) in ``/user_service``

3. Install the Go dependencies.

In the ``/user_service`` directory, run:

  

```

go mod tidy

```

2. Make sure that the``protoc-gen-go`` and ``protoc-gen-go-grpc`` are in GOPATH.

To check, run:

```

which protoc-gen-go

which protoc-gen-go-grpc

```

If the response is empty, add protoc-gen-go and protoc-gen-go-grpc to GOPATH.

  

``export PATH="$PATH:$(go env GOPATH)/bin"``

  

After that, you can generate the files.

In /user_service, run:

```

make deps // Installs the protobuf and gRPC plugins.

make proto // will generate proto files

```

3. Now, you can compile the program.

```

build ./cmd/user_service.

```