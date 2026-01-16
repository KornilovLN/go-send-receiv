FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod файлы
COPY go.mod ./
COPY shared/go.mod ./shared/
COPY receiver-service/go.mod ./receiver-service/

# Копируем исходный код
COPY shared/ ./shared/
COPY receiver-service/ ./receiver-service/

# Переходим в директорию receiver-service и работаем с модулями
WORKDIR /app/receiver-service

# Инициализируем модули
RUN go mod tidy

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o receiver ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/receiver-service/receiver .
COPY --from=builder /app/receiver-service/static ./static

EXPOSE 8080 8081
CMD ["./receiver"]

