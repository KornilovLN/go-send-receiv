## 1. Dockerfile для каждого сервиса
* **Generator Service Dockerfile**
```Dockerfile 
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod файлы
COPY go.mod go.sum* ./
COPY ../shared ../shared

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o generator ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/generator .

# Запускаем приложение
CMD ["./generator"]
```

* **Receiver Service Dockerfile**
```Dockerfile 
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod файлы
COPY go.mod go.sum* ./
COPY ../shared ../shared

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o receiver ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/receiver .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./receiver"]
```

## 2. Docker Compose в корне проекта
* **docker-compose.yml**
```yml
version: '3.8'

services:
  receiver-service:
    build:
      context: .
      dockerfile: receiver-service/Dockerfile
    container_name: con-receiver
    ports:
      - "8080:8080"
    networks:
      - gendata-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "8080"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  generator-service:
    build:
      context: .
      dockerfile: generator-service/Dockerfile
    container_name: con-gendata
    depends_on:
      receiver-service:
        condition: service_healthy
    networks:
      - gendata-network
    restart: unless-stopped
    environment:
      - RECEIVER_HOST=receiver-service
      - RECEIVER_PORT=8080
      - GENERATION_INTERVAL=5

networks:
  gendata-network:
    driver: bridge
```

## 3. Обновленные Dockerfile с правильным контекстом сборки
    Поскольку используется общий контекст сборки (корень проекта), нужно обновить Dockerfile:
* Обновленный Generator Service Dockerfile
```Dockerfile
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
```

* **Обновленный Receiver Service Dockerfile**
```Dockerfile 
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
```

## 4. Обновленный generator-service/cmd/main.go с переменными окружения
```GO
package main

import (
	"encoding/binary"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"gendata-project/generator-service/internal/generator"
	"gendata-project/shared/protocol"
	"gendata-project/shared/types"
)

func main() {
	log.Println("Запуск генератора данных...")

	// Получаем настройки из переменных окружения
	receiverHost := getEnv("RECEIVER_HOST", "receiver-service")
	receiverPort := getEnv("RECEIVER_PORT", "8080")
	intervalStr := getEnv("GENERATION_INTERVAL", "5")
	
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Printf("Неверный интервал генерации, используется значение по умолчанию: 5 секунд")
		interval = 5
	}

	receiverAddr := receiverHost + ":" + receiverPort
	log.Printf("Адрес получателя: %s", receiverAddr)
	log.Printf("Интервал генерации: %d секунд", interval)

	// Ожидание готовности получателя
	log.Println("Ожидание готовности получателя...")
	time.Sleep(10 * time.Second)

	for {
		// Подключение к получателю
		conn, err := net.Dial("tcp", receiverAddr)
		if err != nil {
			log.Printf("Ошибка подключения к получателю: %v. Повторная попытка через 5 секунд...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Подключение к получателю установлено")
		
		// Генерация данных с заданным интервалом
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		
		dataTypes := []byte{
			types.TypeAnalog,
			types.TypeFloat,
			types.TypeFixed,
			types.TypeInt,
			types.TypeDiscrete,
		}

		connectionActive := true
		for connectionActive {
			select {
			case <-ticker.C:
				// Генерируем данные для каждого типа
				for i, dataType := range dataTypes {
					gid := generator.GIDCreater(1, 2, 1) // Север, NPP 2, Блок 1
					
					header := generator.HeadBlockCreater(
						gid,
						dataType,
						1,                    // ListNum
						1,                    // ListVer
						int16(i),            // Index
						5,                   // NumParams
						time.Now(),
					)

					dataBlock := generator.GenBlock(header)

					// Создаем сообщение
					msg := protocol.Message{
						Type:      "data_block",
						Timestamp: time.Now().Unix(),
						Data:      dataBlock,
					}

					// Отправляем JSON
					jsonData, err := msg.ToJSON()
					if err != nil {
						log.Printf("Ошибка сериализации: %v", err)
						continue
					}

					// Отправляем размер сообщения, затем само сообщение
					size := uint32(len(jsonData))
					if err := binary.Write(conn, binary.LittleEndian, size); err != nil {
						log.Printf("Ошибка отправки размера: %v", err)
						connectionActive = false
						break
					}

					if _, err := conn.Write(jsonData); err != nil {
						log.Printf("Ошибка отправки данных: %v", err)
						connectionActive = false
						break
					}

					log.Printf("Отправлен блок данных типа %c, размер: %d байт", 
						types.DataTypeSymbols[dataType], len(dataBlock))
				}
			}
		}
		
		ticker.Stop()
		conn.Close()
		log.Println("Соединение закрыто, переподключение...")
		time.Sleep(5 * time.Second)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

## 5. .dockerignore файлы
```
.git
.gitignore
README.md
Makefile
scripts/
deploy/
*.md
.env
.env.local
```

## 6. Обновленный Makefile
```Makefile
.PHONY: build deploy logs stop clean status local-run local-stop test build-images

# Сборка локально
build:
	docker-compose build

# Развертывание на удаленном сервере
deploy:
	cd deploy/ansible && chmod +x deploy.sh && ./deploy.sh

# Просмотр логов на удаленном сервере
logs:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose logs -f'

# Остановка сервисов на удаленном сервере
stop:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose down'

# Очистка на удаленном сервере
clean:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose down --rmi all --volumes'

# Проверка статуса сервисов
status:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose ps'

# Локальный запуск для тестирования
local-run:
	docker-compose up --build -d

# Локальная остановка
local-stop:
	docker-compose down

# Локальное тестирование
test:
	chmod +x scripts/test-local.sh && ./scripts/test-local.sh

# Сборка образов
build-images:
	chmod +x scripts/build.sh && ./scripts/build.sh

# Просмотр локальных логов
local-logs:
	docker-compose logs -f

# Перезапуск локальных сервисов
local-restart:
	docker-compose restart

# Показать статус локальных контейнеров
local-status:
	docker-compose ps
```

## 7. Итоговая структура проекта
```
gendata-project/
├── docker-compose.yml          # В корне проекта
├── .dockerignore              # В корне проекта
├── Makefile                   # В корне проекта
├── README.md                  # В корне проекта
├── go.mod                     # В корне проекта
├── generator-service/
│   ├── Dockerfile            # Dockerfile генератора
│   ├── go.mod
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       └── generator/
│           └── generator.go
├── receiver-service/
│   ├── Dockerfile            # Dockerfile получателя
│   ├── go.mod
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── parser/
│       │   └── parser.go
│       └── display/
│           └── display.go
├── shared/
│   ├── go.mod
│   ├── types/
│   │   └── types.go
│   ├── sensorvalue/
│   │   └── sensorvalue.go
│   └── protocol/
│       └── protocol.go
├── deploy/
│   └── ansible/
│       ├── inventory.yml
│       ├── playbook.yml
│       ├── deploy.sh
│       └── ansible.cfg
└── scripts/
    ├── build.sh
    ├── test-local.sh
    └── create-structure.sh
```

## Расположение файлов:
    docker-compose.yml в корне проекта
    Dockerfile в каждом сервисе
    Общий контекст сборки позволяет использовать shared модуль
    Переменные окружения для настройки сервисов