# Mermaid диаграммы не всегда корректно экспортируются в PDF

## 1. Используйте Mermaid CLI для генерации изображений
* **Установите Mermaid CLI:**
```bash
npm install -g @mermaid-js/mermaid-cli
```
* **Создайте отдельные файлы диаграмм:**
  * **architect.mmd** 
```mermaid
%%{init: {
  'theme': 'base',
  'themeVariables': {
    'primaryColor': '#00ADD8',
    'primaryTextColor': '#ffffff',
    'primaryBorderColor': '#007d9c',
    'lineColor': '#00ADD8',
    'secondaryColor': '#5DC9E2',
    'tertiaryColor': '#CE3262',
    'background': '#1a1a1a',
    'mainBkg': '#2d2d2d'
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
* **Генерируйте изображения:**
```bash
mmdc -i md/modern1/architect.mmd -o md/modern1/architect.png
mmdc -i md/modern1/architect.mmd -o md/modern1/architect.svg
```

## 2. Создайте скрипт для автоматической генерации generate-diagrams.sh
```bash
#!/bin/bash

# generate-diagrams.sh

# Создаем директорию для изображений
mkdir -p md/modern1/images

# Генерируем PNG изображения
mmdc -i md/modern1/architect.mmd -o md/modern1/images/architect.png -w 1200 -H 800
mmdc -i md/modern1/dataflow.mmd -o md/modern1/images/dataflow.png -w 1200 -H 800

# Генерируем SVG (векторные изображения)
mmdc -i md/modern1/architect.mmd -o md/modern1/images/architect.svg
mmdc -i md/modern1/dataflow.mmd -o md/modern1/images/dataflow.svg

echo "Диаграммы сгенерированы в md/modern1/images/"
```
* **Сделайте скрипт исполняемым:**
```bash
chmod +x scripts/generate-diagrams.sh
```

## 3. Обновите Markdown файлы с ссылками на изображения
    architect-with-images.md
### System Architecture

#### Sequence Diagram

![Architecture Sequence Diagram](images/architect.png)

*Диаграмма последовательности взаимодействия компонентов системы*

#### Component Flow

![Component Flow Diagram](images/dataflow.png)

*Схема потоков данных между компонентами*

---

### Mermaid Source Code

<details>
<summary>Исходный код диаграммы последовательности</summary>

```mermaid
%%{init: {'theme':'dark'}}%%
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

</details>
```

## 4. Создайте Makefile для автоматизации

```makefile:Makefile
.PHONY: diagrams pdf clean

# Генерация диаграмм
diagrams:
	@echo "Генерация диаграмм..."
	@mkdir -p md/modern1/images
	@mmdc -i md/modern1/architect.mmd -o md/modern1/images/architect.png -w 1200 -H 800
	@mmdc -i md/modern1/architect.mmd -o md/modern1/images/architect.svg
	@echo "Диаграммы сгенерированы!"

# Генерация PDF с изображениями
pdf: diagrams
	@echo "Генерация PDF..."
	@pandoc md/modern1/architect-with-images.md -o md/modern1/architect.pdf --pdf-engine=wkhtmltopdf
	@echo "PDF сгенерирован!"

# Очистка
clean:
	@rm -rf md/modern1/images/
	@rm -f md/modern1/*.pdf
	@echo "Очистка завершена!"
```

## 5. Запустите генерацию

```bash
make diagrams
```

```bash
make pdf
```


Теперь у вас будут PNG/SVG изображения диаграмм, которые корректно отобразятся в PDF!