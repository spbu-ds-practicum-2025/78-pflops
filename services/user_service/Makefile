PROTO_DIR=proto
OUT_DIR=gen
GO_MODULE=78-pflops/services/user_service

.PHONY: proto deps clean help

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

help:
	@echo "Доступные команды:"
	@echo "  make proto   - Генерация gRPC и protobuf файлов"
	@echo "  make deps    - Установка необходимых плагинов"
	@echo "  make clean   - Очистка сгенерированных файлов"