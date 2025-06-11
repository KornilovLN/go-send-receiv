FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod файлы
COPY go.mod ./
COPY shared/go.mod ./shared/
COPY generator-service/go.mod ./generator-service/

# Копируем исходный код
COPY shared/ ./shared/
COPY generator-service/ ./generator-service/

# Переходим в директорию generator-service
WORKDIR /app/generator-service

# Инициализируем модули
RUN go mod tidy

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o generator ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/generator-service/generator .
CMD ["./generator"]
