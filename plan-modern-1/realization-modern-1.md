# Новая структура проекта:
## JSON структура
```json
{
  "project": "gendata-distributed-system",
  "architecture": "microservices",
  "services": {
    "generator": {
      "description": "Генерирует пакеты данных",
      "port": null,
      "dependencies": ["receiver"],
      "outputs": ["raw_data_packets"]
    },
    "receiver": {
      "description": "Принимает данные, парсит и распределяет",
      "ports": [8080, 8081],
      "dependencies": ["rabbitmq", "redis"],
      "inputs": ["raw_data_packets"],
      "outputs": ["parsed_data_to_queue", "web_interface"]
    },
    "filtrator": {
      "description": "Фильтрует и прореживает данные",
      "dependencies": ["rabbitmq", "redis", "opentsdb"],
      "inputs": ["raw_data_from_queue"],
      "outputs": ["filtered_data_to_opentsdb"]
    },
    "influxdb": {
      "description": "Хранит сырые данные",
      "port": 8086,
      "dependencies": ["rabbitmq"],
      "inputs": ["raw_data_from_queue"],
      "outputs": ["stored_raw_data"]
    },
    "questdb": {
      "description": "Альтернатива InfluxDB",
      "port": 9000,
      "dependencies": ["rabbitmq"],
      "inputs": ["raw_data_from_queue"],
      "outputs": ["stored_raw_data"]
    },
    "opentsdb": {
      "description": "Хранит отфильтрованные данные",
      "port": 4242,
      "dependencies": ["rabbitmq"],
      "inputs": ["filtered_data_from_queue"],
      "outputs": ["stored_filtered_data"]
    },
    "observer": {
      "description": "Веб-интерфейс для мониторинга",
      "port": 8082,
      "dependencies": ["rabbitmq", "influxdb", "questdb", "opentsdb"],
      "inputs": ["notifications_from_queue", "db_queries"],
      "outputs": ["web_dashboard"]
    }
  },
  "infrastructure": {
    "rabbitmq": {
      "port": 5672,
      "queues": ["raw_data", "filtered_data", "notifications"],
      "exchanges": ["data_exchange"]
    },
    "redis": {
      "port": 6379,
      "usage": ["caching", "session_storage"]
    }
  },
  "data_flow": {
    "raw_data_path": "generator -> receiver -> rabbitmq -> [influxdb|questdb]",
    "filtered_data_path": "receiver -> rabbitmq -> filtrator -> opentsdb",
    "monitoring_path": "rabbitmq -> observer -> web_dashboard"
  }
}
```

## Диаграммы:
* **architect.puml**
```plantuml
@startuml system_architecture

package "Data Generation" {
  [Generator] as gen
}

package "Data Processing" {
  [Receiver] as recv
  [Filtrator] as filt
}

package "Message Broker" {
  [RabbitMQ] as rmq
  [Redis] as redis
}

package "Data Storage" {
  [InfluxDB] as influx
  [QuestDB] as quest
  [OpenTSDB] as opents
}

package "Monitoring" {
  [Observer] as obs
}

package "Web Interface" {
  [Dashboard] as dash
}

' Data flow connections
gen --> recv : TCP:8080\nraw packets
recv --> rmq : publish\nraw data
recv --> redis : cache
rmq --> influx : consume\nraw data
rmq --> quest : consume\nraw data
rmq --> filt : consume\nraw data
filt --> opents : filtered data
rmq --> obs : notifications
obs --> influx : query data
obs --> quest : query data
obs --> opents : query data
obs --> dash : web interface

' Port annotations
note right of recv : ports: 8080, 8081
note right of influx : port: 8086
note right of quest : port: 9000
note right of opents : port: 4242
note right of obs : port: 8082
note right of rmq : port: 5672
note right of redis : port: 6379

@enduml
```

* **dataflow.puml**
```plantuml
@startuml data_flow_sequence
participant Generator as G
participant Receiver as R
participant RabbitMQ as MQ
participant Redis as Cache
participant InfluxDB as IDB
participant Filtrator as F
participant OpenTSDB as OTSDB
participant Observer as O
participant WebDashboard as WD

== Data Generation & Processing ==
G -> R: Send raw data packet (TCP:8080)
R -> R: Parse packet
R -> Cache: Store parsed data
R -> MQ: Publish to raw_data queue

== Parallel Data Storage ==
par Raw Data Storage
  MQ -> IDB: Consume raw data
  IDB -> IDB: Store raw data
else Data Filtering
  MQ -> F: Consume raw data
  F -> F: Apply filters & thresholds
  F -> OTSDB: Send filtered data
  OTSDB -> OTSDB: Store filtered data
end

== Monitoring & Notifications ==
MQ -> O: Send notification
O -> IDB: Query raw data
IDB -> O: Return data
O -> OTSDB: Query filtered data
OTSDB -> O: Return data
O -> WD: Update dashboard
WD -> WD: Display metrics

== Web Interface ==
note over R: Receiver also has\nweb interface (port 8081)
R -> WD: Real-time data display

@enduml
```

## Создание структуры каталогов
```bash
mkdir -p filtrator-service/{cmd,internal/{filtrator,config},pkg}
mkdir -p observer-service/{cmd,internal/{observer,monitoring,config},pkg}
mkdir -p web-dashboard/{cmd,internal/{dashboard,api,websocket},static/{css,js,images},templates}
mkdir -p storage-services/{influxdb,questdb,opentsdb}
mkdir -p message-broker/{rabbitmq,redis}
mkdir -p monitoring/{prometheus,grafana}
mkdir -p docs/{architecture,api}
```

## 1. Создать структуру проекта (restructure-project.sh):
```bash
#!/bin/bash
# restructure-project.sh

echo "Реструктуризация проекта go-send-receiv..."

# Создание новых директорий согласно архитектуре
mkdir -p {filtrator-service,observer-service,web-dashboard}/{cmd,internal,pkg}
mkdir -p storage-services/{influxdb,questdb,opentsdb}/config
mkdir -p message-broker/{rabbitmq,redis}/config
mkdir -p monitoring/{prometheus,grafana}/config
mkdir -p docs/{architecture,api,deployment}
mkdir -p configs/{development,production,testing}

# Создание поддиректорий для сервисов
mkdir -p filtrator-service/internal/{filtrator,algorithms,config}
mkdir -p observer-service/internal/{observer,monitoring,metrics,config}
mkdir -p web-dashboard/internal/{dashboard,api,websocket,handlers}
mkdir -p web-dashboard/{static/{css,js,images},templates/{layouts,partials}}

echo "Структура создана!"
```

## 2. Обновленный docker-compose.yml:
```yml
version: '3.8'

services:
  # Message Broker Layer
  rabbitmq:
    image: rabbitmq:3-management
    container_name: go-send-recv-rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: admin123
    volumes:
      - ./message-broker/rabbitmq/config:/etc/rabbitmq
    networks:
      - go-send-recv-net

  redis:
    image: redis:7-alpine
    container_name: go-send-recv-redis
    ports:
      - "6379:6379"
    volumes:
      - ./message-broker/redis/config:/usr/local/etc/redis
    networks:
      - go-send-recv-net

  # Storage Layer
  influxdb:
    image: influxdb:2.7
    container_name: go-send-recv-influxdb
    ports:
      - "8086:8086"
    environment:
      DOCKER_INFLUXDB_INIT_MODE: setup
      DOCKER_INFLUXDB_INIT_USERNAME: admin
      DOCKER_INFLUXDB_INIT_PASSWORD: admin123
      DOCKER_INFLUXDB_INIT_ORG: go-send-recv
      DOCKER_INFLUXDB_INIT_BUCKET: sensor-data
    volumes:
      - influxdb-data:/var/lib/influxdb2
      - ./storage-services/influxdb/config:/etc/influxdb2
    networks:
      - go-send-recv-net

  questdb:
    image: questdb/questdb:7.3.10
    container_name: go-send-recv-questdb
    ports:
      - "9000:9000"
      - "8812:8812"
      - "9009:9009"
    volumes:
      - questdb-data:/var/lib/questdb
      - ./storage-services/questdb/config:/var/lib/questdb/conf
    networks:
      - go-send-recv-net

  opentsdb:
    image: petergrace/opentsdb-docker:latest
    container_name: go-send-recv-opentsdb
    ports:
      - "4242:4242"
    volumes:
      - opentsdb-data:/opt/opentsdb
      - ./storage-services/opentsdb/config:/etc/opentsdb
    networks:
      - go-send-recv-net

  # Data Generation Layer
  generator:
    build:
      context: ./generator-service
      dockerfile: Dockerfile
    container_name: go-send-recv-generator
    environment:
      - RECEIVER_HOST=receiver
      - RECEIVER_PORT=8080
      - GENERATION_INTERVAL=1s
    depends_on:
      - receiver
    networks:
      - go-send-recv-net

  # Data Processing Layer
  receiver:
    build:
      context: ./receiver-service
      dockerfile: Dockerfile
    container_name: go-send-recv-receiver
    ports:
      - "8080:8080"
      - "8081:8081"
    environment:
      - RABBITMQ_URL=amqp://admin:admin123@rabbitmq:5672/
      - REDIS_URL=redis:6379
    depends_on:
      - rabbitmq
      - redis
    networks:
      - go-send-recv-net

  filtrator:
    build:
      context: ./filtrator-service
      dockerfile: Dockerfile
    container_name: go-send-recv-filtrator
    environment:
      - RABBITMQ_URL=amqp://admin:admin123@rabbitmq:5672/
      - OPENTSDB_URL=http://opentsdb:4242
    depends_on:
      - rabbitmq
      - opentsdb
    networks:
      - go-send-recv-net

  # Monitoring Layer
  observer:
    build:
      context: ./observer-service
      dockerfile: Dockerfile
    container_name: go-send-recv-observer
    ports:
      - "8082:8082"
    environment:
      - RABBITMQ_URL=amqp://admin:admin123@rabbitmq:5672/
      - INFLUXDB_URL=http://influxdb:8086
      - QUESTDB_URL=http://questdb:9000
      - OPENTSDB_URL=http://opentsdb:4242
    depends_on:
      - rabbitmq
      - influxdb
      - questdb
      - opentsdb
    networks:
      - go-send-recv-net

  # Web Interface Layer
  web-dashboard:
    build:
      context: ./web-dashboard
      dockerfile: Dockerfile
    container_name: go-send-recv-dashboard
    ports:
      - "8083:8083"
    environment:
      - OBSERVER_URL=http://observer:8082
      - RECEIVER_URL=http://receiver:8081
    depends_on:
      - observer
      - receiver
    networks:
      - go-send-recv-net

  # Monitoring Stack
  prometheus:
    image: prom/prometheus:latest
    container_name: go-send-recv-prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus/config:/etc/prometheus
    networks:
      - go-send-recv-net

  grafana:
    image: grafana/grafana:latest
    container_name: go-send-recv-grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin123
    volumes:
      - grafana-data:/var/lib/grafana
      - ./monitoring/grafana/config:/etc/grafana
    networks:
      - go-send-recv-net

volumes:
  influxdb-data:
  questdb-data:
  opentsdb-data:
  grafana-data:

networks:
  go-send-recv-net:
    driver: bridge
```

## 3. Filtrator Service:
```GO
// main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KornilovLN/go-send-receiv/filtrator-service/internal/filtrator"
)

func main() {
	log.Println("Starting Filtrator Service...")

	// Инициализация фильтратора
	f, err := filtrator.New()
	if err != nil {
		log.Fatalf("Failed to initialize filtrator: %v", err)
	}

	// Запуск фильтратора
	go func() {
		if err := f.Start(); err != nil {
			log.Fatalf("Filtrator failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Filtrator Service...")
	f.Stop()
}
```

```GO
// filtrator.go

package filtrator

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
	"github.com/KornilovLN/go-send-receiv/shared/types"
)

type Filtrator struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	opentsdbURL  string
	stopChan     chan bool
}

func New() (*Filtrator, error) {
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	opentsdbURL := getEnv("OPENTSDB_URL", "http://localhost:4242")

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Filtrator{
		conn:        conn,
		channel:     channel,
		opentsdbURL: opentsdbURL,
		stopChan:    make(chan bool),
	}, nil
}

func (f *Filtrator) Start() error {
	// Объявляем очередь для сырых данных
	_, err := f.channel.QueueDeclare(
		"raw_data", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	// Начинаем потреблять сообщения
	msgs, err := f.channel.Consume(
		"raw_data", // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	log.Println("Filtrator started, waiting for messages...")

	for {
		select {
		case msg := <-msgs:
			f.processMessage(msg.Body)
		case <-f.stopChan:
			return nil
		}
	}
}

func (f *Filtrator) processMessage(body []byte) {
	var sensorData types.SensorData
	if err := json.Unmarshal(body, &sensorData); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	// Применяем фильтры
	if f.shouldFilter(sensorData) {
		filteredData := f.applyFilters(sensorData)
		f.sendToOpenTSDB(filteredData)
	}
}

func (f *Filtrator) shouldFilter(data types.SensorData) bool {
	// Логика определения необходимости фильтрации
	return data.Value > 10.0 // Пример: фильтруем значения больше 10
}

func (f *Filtrator) applyFilters(data types.SensorData) types.SensorData {
	// Применяем алгоритмы фильтрации
	filtered := data
	filtered.Value = f.movingAverage(data.Value)
	filtered.Timestamp = time.Now()
	return filtered
}

func (f *Filtrator) movingAverage(value float64) float64 {
	// Простая реализация скользящего среднего
	// В реальности здесь будет более сложная логика
	return value * 0.8
}

func (f *Filtrator) sendToOpenTSDB(data types.SensorData) {
	// Отправка данных в OpenTSDB
	log.Printf("Sending filtered data to OpenTSDB: %+v", data)
	// TODO: Реализовать отправку в OpenTSDB
}

func (f *Filtrator) Stop() {
	close(f.stopChan)
	f.channel.Close()
	f.conn.Close()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

## 4. Observer Service:
```GO
// main.go 
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KornilovLN/go-send-receiv/observer-service/internal/observer"
)

func main() {
	log.Println("Starting Observer Service...")

	obs, err := observer.New()
	if err != nil {
		log.Fatalf("Failed to initialize observer: %v", err)
	}

	// Запуск HTTP сервера
	go func() {
		log.Println("Observer HTTP server starting on :8082")
		if err := http.ListenAndServe(":8082", obs.Router()); err != nil {
			log.Fatalf("Observer HTTP server failed: %v", err)
		}
	}()

	// Запуск мониторинга
	go obs.StartMonitoring()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Observer Service...")
	obs.Stop()
}
```

## 5. Обновленный Makefile:
```Makefile
.PHONY: build run stop clean test docker-build docker-up docker-down logs

# Переменные
SERVICES := generator-service receiver-service filtrator-service observer-service web-dashboard
DOCKER_COMPOSE := docker-compose

# Сборка всех сервисов
build:
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && go mod tidy && go build -o bin/main ./cmd/main.go && cd ..; \
	done

# Запуск через Docker Compose
docker-up:
	@echo "Starting all services with Docker Compose..."
	$(DOCKER_COMPOSE) up -d

# Остановка Docker Compose
docker-down:
	@echo "Stopping all services..."
	$(DOCKER_COMPOSE) down

# Сборка Docker образов
docker-build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

# Полная пересборка и запуск
rebuild: docker-down docker-build docker-up
```

## 6. Web Dashboard Service:
```GO
// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KornilovLN/go-send-receiv/web-dashboard/internal/dashboard"
)

func main() {
	log.Println("Starting Web Dashboard Service...")

	dash, err := dashboard.New()
	if err != nil {
		log.Fatalf("Failed to initialize dashboard: %v", err)
	}

	// Запуск HTTP сервера
	go func() {
		log.Println("Web Dashboard starting on :8083")
		if err := http.ListenAndServe(":8083", dash.Router()); err != nil {
			log.Fatalf("Web Dashboard failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Web Dashboard...")
	dash.Stop()
}
```

```GO
// dashboard.go
package dashboard

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Dashboard struct {
	observerURL  string
	receiverURL  string
	templates    *template.Template
	upgrader     websocket.Upgrader
}

func New() (*Dashboard, error) {
	observerURL := getEnv("OBSERVER_URL", "http://localhost:8082")
	receiverURL := getEnv("RECEIVER_URL", "http://localhost:8081")

	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Dashboard{
		observerURL: observerURL,
		receiverURL: receiverURL,
		templates:   templates,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}, nil
}

func (d *Dashboard) Router() http.Handler {
	r := mux.NewRouter()
	
	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	
	// Routes
	r.HandleFunc("/", d.handleIndex).Methods("GET")
	r.HandleFunc("/api/metrics", d.handleMetrics).Methods("GET")
	r.HandleFunc("/ws", d.handleWebSocket)
	
	return r
}

func (d *Dashboard) handleIndex(w http.ResponseWriter, r *http.Request) {
	d.templates.ExecuteTemplate(w, "index.html", nil)
}

func (d *Dashboard) handleMetrics(w http.ResponseWriter, r *http.Request) {
	// Получаем метрики от Observer
	// TODO: Реализовать получение метрик
	metrics := map[string]interface{}{
		"total_messages": 1000,
		"filtered_messages": 750,
		"error_rate": 0.05,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (d *Dashboard) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := d.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// WebSocket логика для real-time обновлений
	for {
		// TODO: Отправка real-time данных
		select {
		default:
			// Отправляем данные клиенту
		}
	}
}

func (d *Dashboard) Stop() {
	// Cleanup
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

## 7. Dockerfile для новых сервисов:
```Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

```Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

```Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY templates/ ./templates/
COPY static/ ./static/
CMD ["./main"]
```

## 8. Go модули для новых сервисов:
* **go.mod**
```GO
module github.com/KornilovLN/go-send-receiv/filtrator-service

go 1.21

require (
	github.com/KornilovLN/go-send-receiv/shared v0.0.0
	github.com/streadway/amqp v1.1.0
)

replace github.com/KornilovLN/go-send-receiv/shared => ../shared
```

* **go.mod**
```GO
module github.com/KornilovLN/go-send-receiv/observer-service

go 1.21

require (
	github.com/KornilovLN/go-send-receiv/shared v0.0.0
	github.com/gorilla/mux v1.8.0
	github.com/streadway/amqp v1.1.0
	github.com/influxdata/influxdb-client-go/v2 v2.12.3
)

replace github.com/KornilovLN/go-send-receiv/shared => ../shared
```

* **go.mod**
```GO 
module github.com/KornilovLN/go-send-receiv/web-dashboard

go 1.21

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
)
```

## 9. Конфигурационные файлы:
* **config.yml**
```yml
services:
  generator:
    interval: 1s
    batch_size: 10
  
  receiver:
    tcp_port: 8080
    web_port: 8081
    
  filtrator:
    algorithms:
      - moving_average
      - threshold_filter
    threshold: 10.0
    
  observer:
    port: 8082
    metrics_interval: 30s
    
  dashboard:
    port: 8083
    refresh_interval: 5s

databases:
  influxdb:
    url: "http://localhost:8086"
    token: "your-token"
    org: "go-send-recv"
    bucket: "sensor-data"
    
  questdb:
    url: "http://localhost:9000"
    
  opentsdb:
    url: "http://localhost:4242"

message_broker:
  rabbitmq:
    url: "amqp://admin:admin123@localhost:5672/"
    exchange: "sensor_data"
    
  redis:
    url: "localhost:6379"
    db: 0
```

## 10. Обновленный скрипт развертывания (deploy-full.sh):
```bash
#!/bin/bash

# deploy-full.sh

set -e

echo "🚀 Deploying Go Send-Receive System..."

# Проверка зависимостей
command -v docker >/dev/null 2>&1 || { echo "Docker не установлен!" >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose не установлен!" >&2; exit 1; }

# Остановка существующих контейнеров
echo "🛑 Stopping existing containers..."
docker-compose down

# Сборка образов
echo "🔨 Building Docker images..."
docker-compose build

# Создание сетей и томов
echo "🌐 Creating networks and volumes..."
docker network create go-send-recv-net 2>/dev/null || true

# Запуск инфраструктуры
echo "🏗️ Starting infrastructure services..."
docker-compose up -d rabbitmq redis influxdb questdb opentsdb

# Ожидание готовности инфраструктуры
echo "⏳ Waiting for infrastructure to be ready..."
sleep 30

# Запуск приложений
echo "🚀 Starting application services..."
docker-compose up -d generator receiver filtrator observer web-dashboard

# Запуск мониторинга
echo "📊 Starting monitoring services..."
docker-compose up -d prometheus grafana

echo "✅ Deployment completed!"
echo ""
echo "🌐 Services available at:"
echo "  - Receiver Web UI: http://localhost:8081"
echo "  - Observer API: http://localhost:8082"
echo "  - Web Dashboard: http://localhost:8083"
echo "  - RabbitMQ Management: http://localhost:15672"
echo "  - InfluxDB: http://localhost:8086"
echo "  - QuestDB: http://localhost:9000"
echo "  - OpenTSDB: http://localhost:4242"
echo "  - Grafana: http://localhost:3000"
echo "  - Prometheus: http://localhost:9090"
```

## 11. Обновленная документация:
    # Go Send-Receive System Architecture

    ## Overview

    Система состоит из следующих компонентов:

    ### Data Generation Layer
    - **Generator Service**: Генерирует тестовые данные сенсоров

    ### Data Processing Layer  
    - **Receiver Service**: Принимает TCP пакеты, парсит и публикует в очереди
    - **Filtrator Service**: Применяет алгоритмы фильтрации к сырым данным

    ### Message Broker Layer
    - **RabbitMQ**: Основная очередь сообщений
    - **Redis**: Кэширование и быстрый доступ к данным

    ### Storage Layer
    - **InfluxDB**: Хранение сырых данных временных рядов
    - **QuestDB**: Альтернативное хранилище сырых данных
    - **OpenTSDB**: Хранение отфильтрованных данных

    ### Monitoring Layer
    - **Observer Service**: Мониторинг системы и агрегация метрик
    - **Prometheus**: Сбор метрик
    - **Grafana**: Визуализация метрик

    ### Web Interface Layer
    - **Web Dashboard**: Единый веб-интерфейс для мониторинга
    - **Receiver Web UI**: Встроенный веб-интерфейс receiver'а

    ## Data Flow

    1. Generator → Receiver (TCP:8080)
    2. Receiver → RabbitMQ (publish raw data)
    3. Receiver → Redis (cache)
    4. RabbitMQ → InfluxDB/QuestDB (raw data storage)
    5. RabbitMQ → Filtrator (data processing)
    6. Filtrator → OpenTSDB (filtered data storage)
    7. RabbitMQ → Observer (notifications)
    8. Observer → Web Dashboard (metrics)

    ## Ports

    - Generator: internal
    - Receiver: 8080 (TCP), 8081 (Web)
    - Filtrator: internal
    - Observer: 8082 (API)
    - Web Dashboard: 8083 (Web)
    - RabbitMQ: 5672 (AMQP), 15672 (Management)
    - Redis: 6379
    - InfluxDB: 8086
    - QuestDB: 9000
    - OpenTSDB: 4242
    - Prometheus: 9090
    - Grafana: 3000