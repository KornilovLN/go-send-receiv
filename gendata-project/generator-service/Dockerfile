FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем все go.mod файлы
COPY go.mod ./
COPY shared/go.mod ./shared/
COPY generator-service/go.mod ./generator-service/

# Копируем shared модуль
COPY shared/ ./shared/

# Копируем generator-service
COPY generator-service/ ./generator-service/

# Переходим в директорию сервиса
WORKDIR /app/generator-service

# Загружаем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o generator ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/generator-service/generator .

# Запускаем приложение
CMD ["./generator"]