# System Architecture

```mermaid
%%{init: {
  'theme': 'base',
  'themeVariables': {
    'primaryColor': '#00ADD8',
    'primaryTextColor': '#000000',
    'primaryBorderColor': '#007d9c',
    'lineColor': '#333333',
    'secondaryColor': '#5DC9E2',
    'tertiaryColor': '#CE3262',
    'background': '#ffffff',
    'mainBkg': '#f8f9fa',
    'secondBkg': '#e9ecef',
    'tertiaryBkg': '#dee2e6',
    'noteTextColor': '#000000',
    'noteBkgColor': '#e3f2fd',
    'noteBorderColor': '#00ADD8',
    'actorBkg': '#f8f9fa',
    'actorBorder': '#00ADD8',
    'actorTextColor': '#000000',
    'actorLineColor': '#333333',
    'signalColor': '#333333',
    'signalTextColor': '#000000',
    'activationBorderColor': '#007d9c',
    'activationBkgColor': '#e3f2fd',
    'sequenceNumberColor': '#ffffff'
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
```

## Альтернативный вариант с Go-стилем и контрастными цветами:

```markdown:md/modern1/architect-go-contrast.md
# Go Send-Receive Architecture

```mermaid
%%{init: {
  'theme': 'base',
  'themeVariables': {
    'primaryColor': '#00ADD8',
    'primaryTextColor': '#2c3e50',
    'primaryBorderColor': '#007d9c',
    'lineColor': '#2c3e50',
    'secondaryColor': '#5DC9E2',
    'tertiaryColor': '#CE3262',
    'background': '#ffffff',
    'mainBkg': '#ecf0f1',
    'secondBkg': '#bdc3c7',
    'tertiaryBkg': '#95a5a6',
    'noteTextColor': '#2c3e50',
    'noteBkgColor': '#dbeafe',
    'noteBorderColor': '#3b82f6',
    'actorBkg': '#ecf0f1',
    'actorBorder': '#34495e',
    'actorTextColor': '#2c3e50',
    'actorLineColor': '#2c3e50',
    'signalColor': '#2c3e50',
    'signalTextColor': '#2c3e50',
    'activationBorderColor': '#e74c3c',
    'activationBkgColor': '#fadbd8',
    'loopTextColor': '#2c3e50'
  }
}}%%

sequenceDiagram
    participant G as 🚀 Generator
    participant R as 📡 Receiver
    participant MQ as 🐰 RabbitMQ
    participant C as ⚡ Redis Cache
    participant I as 📊 InfluxDB
    participant F as 🔍 Filtrator
    participant O as 📈 OpenTSDB
    participant Obs as 👁️ Observer
    participant W as 🌐 Web UI

    Note over G,W: 📤 Data Generation Phase
    G->>+R: Raw data packets (TCP:8080)
    R->>R: Parse & validate data
    R->>C: Cache parsed data
    R->>MQ: Publish to raw_data exchange
    R-->>-G: ACK
    
    Note over G,W: ⚡ Parallel Processing Phase
    par 💾 Raw Data Storage
        MQ->>+I: Consume from raw_data queue
        I->>-I: Store raw sensor data
    and 🔄 Data Filtering
        MQ->>+F: Consume from raw_data queue
        F->>F: Apply filtering algorithms
        F->>-O: Store filtered data
    and 📊 Monitoring
        MQ->>+Obs: Notification events
        Obs->>I: Query raw data
        Obs->>O: Query filtered data
        Obs->>-W: Update dashboard
    end

    Note over G,W: 🌐 Web Interface Updates
    R-->>W: Direct real-time updates (:8081)
    Obs-->>W: Aggregated monitoring data (:8082)
```
```

## Для dataflow диаграммы:

```markdown:md/modern1/dataflow-contrast.md
## Go Send-Receive - Архитектура компонентов

```mermaid
%%{init: {
  'theme': 'base',
  'themeVariables': {
    'primaryColor': '#00ADD8',
    'primaryTextColor': '#2c3e50',
    'primaryBorderColor': '#007d9c',
    'lineColor': '#2c3e50',
    'secondaryColor': '#5DC9E2',
    'tertiaryColor': '#CE3262',
    'background': '#ffffff',
    'mainBkg': '#ecf0f1',
    'clusterBkg': '#f8f9fa',
    'clusterBorder': '#007d9c',
    'edgeLabelBackground': '#ffffff',
    'nodeTextColor': '#2c3e50',
    'titleColor': '#2c3e50'
  }
}}%%

graph TB
    subgraph DS ["🚀 Data Sources"]
        GEN[🔄 Generator Service]
    end
    
    subgraph PL ["⚡ Processing Layer"]
        REC[📡 Receiver Service<br/>:8080, :8081]
        FILT[🔍 Filtrator Service]
    end
    
    subgraph MB ["💬 Message Broker"]
        RMQ[🐰 RabbitMQ :5672]
        REDIS[⚡ Redis :6379]
    end
    
    subgraph SL ["💾 Storage Layer"]
        INFLUX[📊 InfluxDB :8086<br/>Raw Data]
        QUEST[🗃️ QuestDB :9000<br/>Raw Data Alternative]
        OPENTS[📈 OpenTSDB :4242<br/>Filtered Data]
    end
    
    subgraph ML ["👁️ Monitoring Layer"]
        OBS[🔍 Observer Service :8082]
        WEB[🌐 Web Dashboard]
    end
    
    %% Data Flow
    GEN -->|TCP packets| REC
    REC ==>|parsed data| RMQ
    REC -->|cache| REDIS
    RMQ ==>|raw queue| INFLUX
    RMQ ==>|raw queue| QUEST
    RMQ ==>|raw queue| FILT
    FILT ==>|filtered data| OPENTS
    RMQ -.->|notifications| OBS
    OBS -->|queries| INFLUX
    OBS -->|queries| QUEST
    OBS -->|queries| OPENTS
    OBS ==> WEB
    REC -.->|direct web| WEB

    %% Контрастные цвета для светлого фона
    classDef dataSource fill:#e74c3c,stroke:#c0392b,stroke-width:2px,color:#ffffff
    classDef processing fill:#3498db,stroke:#2980b9,stroke-width:2px,color:#ffffff
    classDef broker fill:#f39c12,stroke:#e67e22,stroke-width:2px,color:#ffffff
    classDef storage fill:#2ecc71,stroke:#27ae60,stroke-width:2px,color:#ffffff
    classDef monitoring fill:#9b59b6,stroke:#8e44ad,stroke-width:2px,color:#ffffff
    
    class GEN dataSource
    class REC,FILT processing
    class RMQ,REDIS broker
    class INFLUX,QUEST,OPENTS storage
    class OBS,WEB monitoring
```
