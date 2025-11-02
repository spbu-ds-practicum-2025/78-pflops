## User Service
### Setup:
1. Install Go dependencies.
In /user_service run:

```
go mod tidy
```
2. Make make sure that the ``protoc-gen-go`` and ``protoc-gen-go-grpc`` are in GOPATH:
To check it, run:
```
which protoc-gen-go
which protoc-gen-go-grpc
```
If the response it empty you need to add  ``protoc-gen-go`` and ``protoc-gen-go-grpc`` to GOPATH:

``export PATH="$PATH:$(go env GOPATH)/bin"``

After that you can generate files
In /user_service run:
```
make deps  // will install protobuf and gRPC plugins
make proto // will generate proto files
```
3. Now you can compile the program:
```
go build ./cmd/user_service
```
