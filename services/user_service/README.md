## User Service
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
5. Теперь вы можете скомпилировать программу:
```
go build ./cmd/user_service
```
### Управление сервисом:
Просмотреть список доступных команд:
```
grpcurl -plaintext localhost:50051 list user.UserService
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
5. Now you can compile the program:
```
go build ./cmd/user_service
```
### Service Management:
To view the list of available commands:
```
grpcurl -plaintext localhost:50051 list user.UserService
```
