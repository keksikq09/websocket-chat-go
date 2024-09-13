# Переменные
SERVER_BINARY = server.exe
CLIENT_BINARY = client.exe
SERVER_MAIN = cmd\server\main.go
CLIENT_MAIN = cmd\client\main.go
BUILD_DIR = .build

# Команды для сборки
build-server:
	if not exist $(BUILD_DIR) (mkdir $(BUILD_DIR) && attrib +h $(BUILD_DIR))
	go build -o $(BUILD_DIR)\$(SERVER_BINARY) $(SERVER_MAIN)

build-client:
	if not exist $(BUILD_DIR) (mkdir $(BUILD_DIR) && attrib +h $(BUILD_DIR))
	go build -o $(BUILD_DIR)\$(CLIENT_BINARY) $(CLIENT_MAIN)

# Команда для сборки всего
build: build-server build-client

# Команды для запуска
run-server: build-server
	$(BUILD_DIR)\$(SERVER_BINARY)

run-client: build-client
	$(BUILD_DIR)\$(CLIENT_BINARY)

# Очистка
clean:
	if exist $(BUILD_DIR) (rd /s /q $(BUILD_DIR))

# Обозначаем, что это не файлы
.PHONY: build-server build-client build run-server run-client clean