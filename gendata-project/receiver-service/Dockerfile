FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем все go.mod файлы
COPY go.mod ./
COPY shared/go.mod ./shared/
COPY receiver-service/go.mod ./receiver-service/

# Копируем shared модуль
COPY shared/ ./shared/

# Копируем receiver-service
COPY receiver-service/ ./receiver-service/

# Переходим в директорию сервиса
WORKDIR /app/receiver-service

# Загружаем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o receiver ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN apk add --no-cache netcat-openbsd

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/receiver-service/receiver .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./receiver"]