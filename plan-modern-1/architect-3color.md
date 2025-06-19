# Go Send-Receive Architecture
## Доступные готовые темы
    default - светлая тема
    dark - темная тема
    forest - зеленая тема
    neutral - нейтральная тема
    base - базовая тема для кастомизации

```mermaid
%%{init: {
  'theme': 'forest',
  'themeVariables': {
    'primaryColor': '#00ADD8',
    'primaryTextColor': '#ffffff',
    'primaryBorderColor': '#007d9c',
    'lineColor': '#00ADD8',
    'secondaryColor': '#5DC9E2',
    'tertiaryColor': '#CE3262',
    'background': '#2e2e2e',
    'mainBkg': '#252526',
    'secondBkg': '#00ADD8',
    'tertiaryBkg': '#5DC9E2',
    'noteTextColor': '#ffffff',
    'noteBkgColor': '#007d9c',
    'noteBorderColor': '#00ADD8',
    'actorBkg': '#252526',
    'actorBorder': '#00ADD8',
    'actorTextColor': '#ffffff',
    'actorLineColor': '#00ADD8',
    'signalColor': '#ffffff',
    'signalTextColor': '#ffffff'
  }
}}%%

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
