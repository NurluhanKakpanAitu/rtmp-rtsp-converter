// filepath: /Users/nurlykhankakpan/Desktop/rtmp-rtsp-converter/README.md
# RTMP to RTSP Converter

Сервис для конвертации видеопотоков из формата RTMP в формат RTSP с поддержкой множественных потоков.

## 🚀 Возможности

- ✅ Конвертация RTMP → RTSP в реальном времени
- ✅ Поддержка множественных потоков одновременно
- ✅ REST API для управления потоками
- ✅ Подробное логирование
- ✅ Docker поддержка
- ✅ Простое тестирование и демонстрация

## 📋 Требования

- Go 1.21+
- FFmpeg
- Docker (опционально)

## ⚡ Быстрый старт

### Локальная установка

1. **Установка зависимостей:**
```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt update && sudo apt install -y ffmpeg

# Или используйте Makefile
make install-ffmpeg-mac  # для macOS
make install-ffmpeg-ubuntu  # для Ubuntu
```

2. **Сборка и запуск:**
```bash
git clone <repository>
cd rtmp-rtsp-converter
make deps
make run
```

### Docker запуск

```bash
# Простой запуск
make docker-build
make docker-run

# Или с Docker Compose (включает тестовый RTMP сервер)
make docker-compose-up
```

## 🔧 Конфигурация

Настройки находятся в `configs/config.yaml`:

```yaml
server:
  port: 8080          # HTTP API порт
  host: "0.0.0.0"

rtsp:
  port: 8554          # RTSP порт
  host: "0.0.0.0"

logging:
  level: "info"       # debug, info, warn, error
  format: "json"      # json или text

streams:
  max_concurrent: 10  # максимум одновременных потоков
  timeout: 30s
```

## 📡 API Endpoints

### Создание потока
```bash
POST /api/v1/streams
Content-Type: application/json

{
  "rtmp_url": "rtmp://source.com/live/stream",
  "stream_id": "my-stream"  # опционально
}
```

### Получение информации о потоке
```bash
GET /api/v1/streams/{stream_id}
```

### Список всех потоков
```bash
GET /api/v1/streams
```

### Остановка потока
```bash
DELETE /api/v1/streams/{stream_id}
```

### Проверка здоровья сервиса
```bash
GET /api/v1/health
```

## 🧪 Тестирование

### 1. Запуск тестового окружения

```bash
# Запуск с Docker Compose (включает RTMP сервер)
make docker-compose-up

# Или локально
make run
```

### 2. Создание тестового RTMP потока

```bash
# В отдельном терминале создаем тестовый поток
make create-test-rtmp
```

### 3. Создание конвертации

```bash
# Создаем поток через API
make test-stream

# Или вручную
curl -X POST http://localhost:8080/api/v1/streams \
  -H "Content-Type: application/json" \
  -d '{"rtmp_url": "rtmp://localhost:1935/live/test", "stream_id": "test-stream"}'
```

### 4. Просмотр RTSP потока

```bash
# С помощью FFplay
make view-rtsp

# Или VLC
vlc rtsp://localhost:8554/test-stream

# Или в браузере (через веб-плеер)
open http://localhost:8080/viewer/test-stream
```

### 5. Проверка статуса

```bash
# Список потоков
make test-list

# Остановка потока
make test-stop
```

## 🛠 Команды Makefile

```bash
make build          # Сборка бинарника
make run             # Сборка и запуск
make dev             # Запуск в режиме разработки
make test            # Запуск тестов
make clean           # Очистка
make deps            # Установка зависимостей

# Docker
make docker-build    # Сборка Docker образа
make docker-run      # Запуск в Docker
make docker-compose-up    # Запуск с Docker Compose
make docker-compose-down  # Остановка Docker Compose

# Тестирование
make test-stream     # Создание тестового потока
make test-list       # Список потоков
make test-stop       # Остановка тестового потока
make create-test-rtmp # Создание тестового RTMP
make view-rtsp       # Просмотр RTSP потока
```

## 📊 Мониторинг и логи

### Логи приложения
```bash
# Docker Compose
make docker-compose-logs

# Локально - логи выводятся в консоль
```

### Статистика RTMP сервера
При использовании Docker Compose доступна статистика RTMP:
- http://localhost:8081/stat

## 🎯 Примеры использования

### Конвертация потока с YouTube Live
```bash
curl -X POST http://localhost:8080/api/v1/streams \
  -H "Content-Type: application/json" \
  -d '{
    "rtmp_url": "rtmp://a.rtmp.youtube.com/live2/YOUR_STREAM_KEY",
    "stream_id": "youtube-stream"
  }'

# Просмотр: rtsp://localhost:8554/youtube-stream
```

### Конвертация потока с дрона
```bash
curl -X POST http://localhost:8080/api/v1/streams \
  -H "Content-Type: application/json" \
  -d '{
    "rtmp_url": "rtmp://drone-ip:1935/live/video",
    "stream_id": "drone-feed"
  }'
```

### Множественные потоки
```bash
# Поток 1
curl -X POST http://localhost:8080/api/v1/streams \
  -d '{"rtmp_url": "rtmp://source1.com/live/stream1", "stream_id": "cam1"}'

# Поток 2  
curl -X POST http://localhost:8080/api/v1/streams \
  -d '{"rtmp_url": "rtmp://source2.com/live/stream2", "stream_id": "cam2"}'

# Просмотр
vlc rtsp://localhost:8554/cam1
vlc rtsp://localhost:8554/cam2
```

## 🚨 Решение проблем

### FFmpeg не найден
```bash
# Проверка установки
ffmpeg -version

# Установка
make install-ffmpeg-mac    # macOS
make install-ffmpeg-ubuntu # Ubuntu
```

### Порты заняты
Измените порты в `configs/config.yaml` или используйте переменные окружения:
```bash
SERVER_PORT=8081 RTSP_PORT=8555 ./rtmp-rtsp-converter
```

### Проблемы с RTMP источником
Проверьте доступность источника:
```bash
ffprobe rtmp://your-source/path
```

## 🏗 Архитектура

```
┌─────────────────┐    ┌───────────────────┐    ┌─────────────────┐
│   RTMP Source   │───▶│  RTMP→RTSP Conv.  │───▶│   RTSP Client   │
│  (Drone, etc.)  │    │    (FFmpeg)       │    │  (VLC, etc.)    │
└─────────────────┘    └───────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────┐
                       │  REST API   │
                       │ (Management)│
                       └─────────────┘
```

## 📝 Лицензия

MIT License