# Задача
    В этом проекте сервис generator создает пакеты данных и выбрасывает их в сеть.
    Сервис receiver создает сервер для работы своего web представителя как go-рутину на порту 8081, также создает слушателя на порту 8080 и принимает пакеты данных, сам их препарирует и сам выводит в консоль посредством логов и в веб интерфейс своей подсистемой, заключенной в каталоги internal/web и static.
    Надо, чтобы данные приходили в receiver, но тот в свою очередь (будучи подписанным на брокера Rabbitmq+Redis в канал и очередь гнал бы данные подписчикам сервисам fitrator и сервису InfluxDB (или QuestDB)), а те в свою очередь должны быть подписчиками в соответствующих каналах и очередях.
    Сервис Filtrator прореживает данные опираясь на свои алгоритмы в пределах соответствия порогам сверху и снизу и тренду и также имея подписку как отправитель - отправляет обработанные данные  сервису хранения, построенному на базе OpenTSDB.
    Сервис OpenTSDB также подписан у брокера как получатель данных после прореживания от сервиса Filtrator.
    Сервис же InfluxDB (или QuestDB) подписаны как получатели сырых данных.
    Возможно построения отдельного сервиса observer, который будет обращаться к сервису БД InfluxDB (или QuestDB) по пришествии ему сообщения от rabbitmq.
    Observer  будет содержать в себе весь функционал веб интерфейса.
    Требуется: структура проекта как json, plantuml и как Diagramms формат.
    Диаграмма всей структуры проекта с отношениями.
    Диаграмма движения данных, запросов, ответов, логов и пр. во времени.
    Возможно, еще какие-то диаграммы.

## JSON структура проекта
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

## PlantUML - Архитектура системы
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

## PlantUML - Диаграмма последовательности (Data Flow)
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

## Mermaid - Архитектура компонентов
```mermaid
graph TB
    subgraph "Data Sources"
        GEN[Generator Service]
    end
    
    subgraph "Processing Layer"
        REC[Receiver Service<br/>:8080, :8081]
        FILT[Filtrator Service]
    end
    
    subgraph "Message Broker"
        RMQ[RabbitMQ :5672]
        REDIS[Redis :6379]
    end
    
    subgraph "Storage Layer"
        INFLUX[InfluxDB :8086<br/>Raw Data]
        QUEST[QuestDB :9000<br/>Raw Data Alternative]
        OPENTS[OpenTSDB :4242<br/>Filtered Data]
    end
    
    subgraph "Monitoring Layer"
        OBS[Observer Service :8082]
        WEB[Web Dashboard]
    end
    
    %% Data Flow
    GEN -->|TCP packets| REC
    REC -->|parsed data| RMQ
    REC -->|cache| REDIS
    RMQ -->|raw queue| INFLUX
    RMQ -->|raw queue| QUEST
    RMQ -->|raw queue| FILT
    FILT -->|filtered data| OPENTS
    RMQ -->|notifications| OBS
    OBS -->|queries| INFLUX
    OBS -->|queries| QUEST
    OBS -->|queries| OPENTS
    OBS --> WEB
    REC -->|direct web| WEB
```

## Mermaid - Диаграмма потоков данных
```mermaid
sequenceDiagram
    participant G as Generator
    participant R as Receiver
    participant MQ as RabbitMQ
    participant C as Redis Cache
    participant I as InfluxDB
    participant F as Filtrator
    participant O as OpenTSDB
    participant Obs as Observer
    participant W as Web UI

    Note over G,W: Data Generation Phase
    G->>R: Raw data packets (TCP:8080)
    R->>R: Parse & validate data
    R->>C: Cache parsed data
    R->>MQ: Publish to raw_data exchange

    Note over G,W: Parallel Processing Phase
    par Raw Data Storage
        MQ->>I: Consume from raw_data queue
        I->>I: Store raw sensor data
    and Data Filtering
        MQ->>F: Consume from raw_data queue
        F->>F: Apply filtering algorithms
        F->>O: Store filtered data
    and Monitoring
        MQ->>Obs: Notification events
        Obs->>I: Query raw data
        Obs->>O: Query filtered data
        Obs->>W: Update dashboard
    end

    Note over G,W: Web Interface Updates
    R->>W: Direct real-time updates (:8081)
    Obs->>W: Aggregated monitoring data (:8082)
```