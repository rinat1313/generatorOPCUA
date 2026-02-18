# ЭТАП 1: Сборка (если собираем прямо в Docker)
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# ЭТАП 2: Финальный контейнер
FROM alpine:latest
WORKDIR /root/

# Устанавливаем часовой пояс (опционально, полезно для логов)
RUN apk add --no-cache tzdata

# Копируем бинарник из билдера
COPY --from=builder /app/main .

# КОПИРУЕМ ФРОНТЕНД (твой html файл)
COPY index.html .

# СОЗДАЕМ СТРУКТУРУ ПАПОК
# -p создаст всю цепочку, если чего-то нет
RUN mkdir -p "data" \
    && mkdir -p "Шаблоны/Защит" \
    && mkdir -p "Шаблоны/Команд" \
    && mkdir -p "Шаблоны/Свойства"

COPY ./data ./data

# Открываем порт, который слушает твой Go-сервер
EXPOSE 8080:8080

# Запуск
CMD ["./main"]