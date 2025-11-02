# Путь до proto-файлов
PROTO_DIR=api/proto

# Каталог, куда кладём сгенерированные файлы
OUT_DIR=api/proto

# Команда генерации .pb.go файлов
proto:
	@echo "Генерация gRPC и protobuf файлов..."
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	        --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	        $(PROTO_DIR)/*.proto
	@echo "Генерация завершена"

# Установка необходимых плагинов
deps:
	@echo "Установка плагинов protoc..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Плагины установлены"

# Очистка сгенерированных файлов (опционально)
clean:
	rm -f $(OUT_DIR)/*.pb.go
	@echo "Очищено"
