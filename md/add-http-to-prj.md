# Добавление HTTP протокола
    Пока сервисы работают по TCP протоколу для передачи данных, а не как HTTP веб-сервер. Порт 8080 используется для TCP соединения между генератором и получателем данных, а не для веб-интерфейса.
    Можно добавить веб-интерфейс для мониторинга:

## Вариант 1: Добавить HTTP endpoint в receiver-service
* **Обновим receiver-service для поддержки HTTP:**
```GO
// receiver-service/main.go
package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"gendata-project/receiver-service/internal/display"
	"gendata-project/shared/protocol"
)

type WebStats struct {
	mu              sync.RWMutex
	TotalMessages   int                    `json:"total_messages"`
	LastMessage     *protocol.Message      `json:"last_message"`
	LastUpdate      time.Time              `json:"last_update"`
	MessagesPerType map[string]int         `json:"messages_per_type"`
	RecentMessages  []*protocol.Message    `json:"recent_messages"`
}

var stats = &WebStats{
	MessagesPerType: make(map[string]int),
	RecentMessages:  make([]*protocol.Message, 0, 10),
}

func main() {
	log.Println("Запуск сервиса получения данных...")

	// Запуск HTTP сервера для веб-интерфейса
	go startWebServer()

	// Запуск TCP сервера для приема данных
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Ошибка запуска TCP сервера: %v", err)
	}
	defer listener.Close()

	log.Println("TCP сервер запущен на порту 8080")
	log.Println("Веб-интерфейс доступен на http://localhost:8081")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка принятия соединения: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func startWebServer() {
	http.HandleFunc("/", webHandler)
	http.HandleFunc("/api/stats", apiHandler)
	http.HandleFunc("/api/messages", messagesHandler)
	
	log.Println("HTTP сервер запущен на порту 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Printf("Ошибка HTTP сервера: %v", err)
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>GenData Monitor</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .card { background: white; padding: 20px; margin: 10px 0; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; }
        .stat-item { text-align: center; padding: 15px; background: #e3f2fd; border-radius: 5px; }
        .stat-value { font-size: 24px; font-weight: bold; color: #1976d2; }
        .stat-label { color: #666; margin-top: 5px; }
        .messages { max-height: 400px; overflow-y: auto; }
        .message { border-left: 4px solid #4caf50; padding: 10px; margin: 5px 0; background: #f9f9f9; }
        .message-header { font-weight: bold; color: #333; }
        .message-time { color: #666; font-size: 12px; }
        .refresh-btn { background: #4caf50; color: white; border: none; padding: 10px 20px; border-radius: 5px; cursor: pointer; }
        .refresh-btn:hover { background: #45a049; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔄 GenData Services Monitor</h1>
        
        <div class="card">
            <h2>📊 Статистика</h2>
            <div class="stats" id="stats">
                <div class="stat-item">
                    <div class="stat-value" id="total-messages">0</div>
                    <div class="stat-label">Всего сообщений</div>
                </div>
                <div class="stat-item">
                    <div class="stat-value" id="last-update">-</div>
                    <div class="stat-label">Последнее обновление</div>
                </div>
            </div>
        </div>

        <div class="card">
            <h2>📨 Последние сообщения</h2>
            <button class="refresh-btn" onclick="refreshData()">🔄 Обновить</button>
            <div class="messages" id="messages">
                <p>Ожидание данных...</p>
            </div>
        </div>
    </div>

    <script>
        function refreshData() {
            fetch('/api/stats')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('total-messages').textContent = data.total_messages;
                    document.getElementById('last-update').textContent = 
                        data.last_update ? new Date(data.last_update).toLocaleString() : '-';
                });

            fetch('/api/messages')
                .then(response => response.json())
                .then(data => {
                    const messagesDiv = document.getElementById('messages');
                    if (data.length === 0) {
                        messagesDiv.innerHTML = '<p>Нет сообщений</p>';
                        return;
                    }
                    
                    messagesDiv.innerHTML = data.map(msg => 
                        '<div class="message">' +
                        '<div class="message-header">Тип: ' + msg.type + '</div>' +
                        '<div class="message-time">Время: ' + new Date(msg.timestamp * 1000).toLocaleString() + '</div>' +
                        '<div>Размер данных: ' + (msg.data ? msg.data.length : 0) + ' байт</div>' +
                        '</div>'
                    ).join('');
                });
        }

        // Автообновление каждые 2 секунды
        setInterval(refreshData, 2000);
        
        // Первоначальная загрузка
        refreshData();
    </script>
</body>
</html>`

	t, _ := template.New("index").Parse(tmpl)
	t.Execute(w, nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	stats.mu.RLock()
	defer stats.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	stats.mu.RLock()
	defer stats.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats.RecentMessages)
}

func updateStats(msg *protocol.Message) {
	stats.mu.Lock()
	defer stats.mu.Unlock()

	stats.TotalMessages++
	stats.LastMessage = msg
	stats.LastUpdate = time.Now()
	stats.MessagesPerType[msg.Type]++

	// Добавляем в список последних сообщений
	stats.RecentMessages = append([]*protocol.Message{msg}, stats.RecentMessages...)
	if len(stats.RecentMessages) > 10 {
		stats.RecentMessages = stats.RecentMessages[:10]
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Новое соединение от %s", conn.RemoteAddr())

	for {
		// Читаем размер сообщения
		var size uint32
		if err := binary.Read(conn, binary.LittleEndian, &size); err != nil {
			if err != io.EOF {
				log.Printf("Ошибка чтения размера: %v", err)
			}
			break
		}

		// Читаем само сообщение
		msgData := make([]byte, size)
		if _, err := io.ReadFull(conn, msgData); err != nil {
			log.Printf("Ошибка чтения сообщения: %v", err)
			break
		}

		// Парсим JSON сообщение
		msg, err := protocol.FromJSON(msgData)
		if err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			continue
		}

		log.Printf("Получено сообщение типа: %s, размер данных: %d байт", 
			msg.Type, len(msg.Data))

		// Обновляем статистику
		updateStats(msg)

		// Отображаем данные в консоли
		if err := display.DisplayDataBlock(msg.Data); err != nil {
			log.Printf("Ошибка отображения данных: %v", err)
		}
	}
}
```

* **Теперь обновим docker-compose.yml для проброса нового порта:**
```yml
version: '3.8'

services:
  receiver-service:
    build:
      context: ./receiver-service
      dockerfile: Dockerfile
    container_name: con-receiver
    ports:
      - "8080:8080"  # TCP порт для приема данных
      - "8081:8081"  # HTTP порт для веб-интерфейса
    networks:
      - gendata-network
    restart: unless-stopped

  generator-service:
    build:
      context: ./generator-service
      dockerfile: Dockerfile
    container_name: con-gendata
    depends_on:
      - receiver-service
    networks:
      - gendata-network
    restart: unless-stopped

networks:
  gendata-network:
    driver: bridge
```

## Вариант 2: Простая проверка статуса через HTTP
    Более простое решение - добавить только endpoint для проверки статуса:
```GO
// status.go
package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type StatusResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		Status:    "running",
		Timestamp: time.Now(),
		Message:   "GenData Receiver Service is running",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
```

## После внесения изменений:
* **Пересоберите и разверните:**
```bash 
make deploy-manual
```

* **Теперь веб-интерфейс будет доступен по адресу:**
```
http://192.168.88.102:8081
```

* **А статус сервиса можно проверить по:**
```
http://192.168.88.102:8081/api/stats
```

* **Порт 8080 остается для TCP соединения между сервисами, а 8081 - для веб-интерфейса мониторинга.**