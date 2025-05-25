.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Go параметры
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Имена
BINARY_NAME=rtmp-rtsp-converter
BINARY_UNIX=$(BINARY_NAME)_unix

# Сборка
build:
    $(GOBUILD) -o $(BINARY_NAME) -v cmd/server/main.go

# Запуск
run:
    $(GOBUILD) -o $(BINARY_NAME) -v cmd/server/main.go
    ./$(BINARY_NAME)

# Запуск в режиме разработки
dev:
    $(GOCMD) run cmd/server/main.go

# Тесты
test:
    $(GOTEST) -v ./...

# Очистка
clean:
    $(GOCLEAN)
    rm -f $(BINARY_NAME)
    rm -f $(BINARY_UNIX)

# Загрузка зависимостей
deps:
    $(GOMOD) download
    $(GOMOD) tidy

# Docker сборка
docker-build:
    docker build -t $(BINARY_NAME) .

# Docker запуск
docker-run:
    docker run -p 8080:8080 -p 8554:8554 $(BINARY_NAME)

# Docker Compose
docker-compose-up:
    docker-compose up -d

docker-compose-down:
    docker-compose down

docker-compose-logs:
    docker-compose logs -f

# Тестовые команды
test-stream:
    curl -X POST http://localhost:8080/api/v1/streams \
        -H "Content-Type: application/json" \
        -d '{"rtmp_url": "rtmp://localhost:1935/live/test", "stream_id": "test-stream"}'

test-list:
    curl http://localhost:8080/api/v1/streams

test-stop:
    curl -X DELETE http://localhost:8080/api/v1/streams/test-stream

# Установка FFmpeg (для macOS)
install-ffmpeg-mac:
    brew install ffmpeg

# Установка FFmpeg (для Ubuntu/Debian)
install-ffmpeg-ubuntu:
    sudo apt update && sudo apt install -y ffmpeg

# Создание тестового RTMP потока
create-test-rtmp:
    ffmpeg -f lavfi -i testsrc2=duration=3600:size=1280x720:rate=30 \
        -f lavfi -i sine=frequency=1000:duration=3600 \
        -c:v libx264 -preset fast -c:a aac \
        -f flv rtmp://localhost:1935/live/test

# Просмотр RTSP потока
view-rtsp:
    ffplay rtsp://localhost:8554/test-stream