# Не идут данные на web страницу в таблицу получаемых данных.
## 1. Проблема в JavaScript файле
* **Создайте правильный dashboard.js:**
```js
// dashboard.js
// Глобальные переменные
let autoScroll = true;
let totalMessages = 0;
let isLoading = false;

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    console.log('Dashboard initialized');
    loadTableData();
    
    // Автообновление каждые 2 секунды
    setInterval(loadTableData, 2000);
    
    // Обновление статистики каждые 5 секунд
    setInterval(updateStats, 5000);
    
    updateStats();
});

// Загрузка данных для таблицы
async function loadTableData() {
    if (isLoading) return;
    isLoading = true;
    
    try {
        console.log('Loading table data...');
        const response = await fetch('/api/messages');
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const messages = await response.json();
        console.log('Received messages:', messages);
        
        updateTable(messages);
        updateConnectionStatus(true);
        
    } catch (error) {
        console.error('Error loading table data:', error);
        updateConnectionStatus(false);
        showError('Ошибка загрузки данных: ' + error.message);
    } finally {
        isLoading = false;
    }
}

// Обновление таблицы
function updateTable(messages) {
    const tbody = document.getElementById('sensor-data');
    if (!tbody) {
        console.error('Table body not found');
        return;
    }
    
    if (!messages || messages.length === 0) {
        tbody.innerHTML = '<tr><td colspan="9" class="no-data">Нет данных</td></tr>';
        return;
    }
    
    let html = '';
    messages.forEach((message, index) => {
        if (message.Parameters && message.Parameters.length > 0) {
            message.Parameters.forEach((param, paramIndex) => {
                html += `
                    <tr>
                        <td>${index + 1}.${paramIndex + 1}</td>
                        <td>${param.Passport || 'N/A'}</td>
                        <td>${param.DataType || 'N/A'}</td>
                        <td>${message.Header?.DataType || 'N/A'}</td>
                        <td>${param.Validation ? '✅' : '❌'}</td>
                        <td>${param.ShiftStatus || 0}</td>
                        <td>${param.ParsedValue || 'N/A'}</td>
                        <td>${param.RawValue || 'N/A'}</td>
                        <td>${formatTimestamp(message.Timestamp)}</td>
                    </tr>
                `;
            });
        } else {
            html += `
                <tr>
                    <td>${index + 1}</td>
                    <td colspan="8">Сообщение без параметров (${message.Type})</td>
                </tr>
            `;
        }
    });
    
    tbody.innerHTML = html;
    
    // Обновляем счетчик
    totalMessages = messages.length;
    const countElement = document.getElementById('total-count');
    if (countElement) {
        countElement.textContent = totalMessages;
    }
    
    // Автопрокрутка
    if (autoScroll) {
        const container = document.getElementById('table-container');
        if (container) {
            container.scrollTop = container.scrollHeight;
        }
    }
}

// Обновление статистики
async function updateStats() {
    try {
        const response = await fetch('/api/stats');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const stats = await response.json();
        console.log('Stats:', stats);
        
        // Обновляем элементы статистики
        updateElement('total-messages', stats.TotalMessages || 0);
        updateElement('last-update', formatTimestamp(stats.LastUpdate));
        updateElement('uptime', stats.Uptime || '0s');
        
        // Обновляем типы сообщений
        updateMessageTypes(stats.MessageTypes || {});
        
    } catch (error) {
        console.error('Error loading stats:', error);
    }
}

// Обновление типов сообщений
function updateMessageTypes(messageTypes) {
    const container = document.getElementById('message-types');
    if (!container) return;
    
    if (Object.keys(messageTypes).length === 0) {
        container.innerHTML = '<div class="no-data">Нет данных</div>';
        return;
    }
    
    let html = '';
    for (const [type, count] of Object.entries(messageTypes)) {
        html += `
            <div class="message-type-item">
                <span class="type-name">${type}</span>
                <span class="type-count">${count}</span>
            </div>
        `;
    }
    
    container.innerHTML = html;
}

// Вспомогательные функции
function updateElement(id, value) {
    const element = document.getElementById(id);
    if (element) {
        element.textContent = value;
    }
}

function formatTimestamp(timestamp) {
    if (!timestamp) return 'N/A';
    
    const date = new Date(timestamp);
    return date.toLocaleString('ru-RU', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    });
}

function updateConnectionStatus(connected) {
    const statusElement = document.getElementById('connection-status');
    if (statusElement) {
        if (connected) {
            statusElement.textContent = '🟢 Подключено';
            statusElement.className = 'status connected';
        } else {
            statusElement.textContent = '🔴 Отключено';
            statusElement.className = 'status disconnected';
        }
    }
}

function showError(message) {
    const tbody = document.getElementById('sensor-data');
    if (tbody) {
        tbody.innerHTML = `<tr><td colspan="9" class="error">${message}</td></tr>`;
    }
}

// Функции управления
function toggleAutoScroll() {
    autoScroll = !autoScroll;
    const btn = document.getElementById('auto-scroll-btn');
    if (btn) {
        btn.textContent = `🔄 Авто-прокрутка: ${autoScroll ? 'ВКЛ' : 'ВЫКЛ'}`;
        btn.className = autoScroll ? 'auto-scroll' : 'auto-scroll disabled';
    }
}

function clearTable() {
    const tbody = document.getElementById('sensor-data');
    if (tbody) {
        tbody.innerHTML = '<tr><td colspan="9" class="loading">Таблица очищена</td></tr>';
    }
    
    const countElement = document.getElementById('total-count');
    if (countElement) {
        countElement.textContent = '0';
    }
}

// Проверка здоровья API
async function checkHealth() {
    try {
        const response = await fetch('/api/health');
        const health = await response.json();
        console.log('Health check:', health);
        return health.status === 'ok';
    } catch (error) {
        console.error('Health check failed:', error);
        return false;
    }
}
```

## 2. Исправьте обработчики в Go
```GO
// handlers.go
package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// MessagesHandler возвращает последние сообщения в JSON
func (ws *WebServer) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем CORS заголовки
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-Control", "no-cache")

	// Обработка preflight запросов
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("API /api/messages called")

	messages := ws.Stats.GetRecentMessages()
	log.Printf("Returning %d messages", len(messages))

	// Добавляем отладочную информацию
	if len(messages) > 0 {
		log.Printf("First message: Type=%s, ParamsCount=%d", 
			messages[0].Type, len(messages[0].Parameters))
	}

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Error encoding messages: %v", err)
		http.Error(w, "Error encoding messages", http.StatusInternalServerError)
		return
	}
}

// StatsHandler возвращает статистику в JSON
func (ws *WebServer) StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-Control", "no-cache")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("API /api/stats called")

	stats := ws.Stats.GetStats()
	log.Printf("Returning stats: TotalMessages=%d", stats.TotalMessages)

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding stats: %v", err)
		http.Error(w, "Error encoding stats", http.StatusInternalServerError)
		return
	}
}

// HealthHandler проверка здоровья сервиса
func (ws *WebServer) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-Control", "no-cache")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("API /api/health called")

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
		"uptime":    ws.Stats.GetStats().Uptime,
		"service":   "receiver-service",
	}

	json.NewEncoder(w).Encode(response)
}
```

## 3. Добавьте отладку в WebStats
```GO
// stats.go
package web

import (
	"log"
	"sync"
	"time"
)

// WebStats хранит статистику веб-сервера
type WebStats struct {
	mu              sync.RWMutex
	totalMessages   int
	lastUpdate      time.Time
	startTime       time.Time
	recentMessages  []*WebMessage
	messageTypes    map[string]int
	maxMessages     int
}

// NewWebStats создает новый экземпляр WebStats
func NewWebStats() *WebStats {
	return &WebStats{
		startTime:      time.Now(),
		recentMessages: make([]*WebMessage, 0),
		messageTypes:   make(map[string]int),
		maxMessages:    100, // Храним последние 100 сообщений
	}
}

// UpdateStats обновляет статистику
func (ws *WebStats) UpdateStats(msg *WebMessage) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.totalMessages++
	ws.lastUpdate = time.Now()
	
	// Обновляем типы сообщений
	ws.messageTypes[msg.Type]++
	
	// Добавляем сообщение в список последних
	ws.recentMessages = append(ws.recentMessages, msg)
	
	// Ограничиваем размер списка
	if len(ws.recentMessages) > ws.maxMessages {
		ws.recentMessages = ws.recentMessages[len(ws.recentMessages)-ws.maxMessages:]
	}
	
	log.Printf("Stats updated: Total=%d, Type=%s, RecentCount=%d", 
		ws.totalMessages, msg.Type, len(ws.recentMessages))
}

// GetStats возвращает текущую статистику
func (ws *WebStats) GetStats() *WebStatsResponse {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	uptime := time.Since(ws.startTime)
	
	stats := &WebStatsResponse{
		TotalMessages: ws.totalMessages,
		LastUpdate:    ws.lastUpdate,
		Uptime:        uptime.String(),
		MessageTypes:  make(map[string]int),
	}
	
	// Копируем типы сообщений
	for k, v := range ws.messageTypes {
		stats.MessageTypes[k] = v
	}
	
	return stats
}

// GetRecentMessages возвращает последние сообщения
func (ws *WebStats) GetRecentMessages() []*WebMessage {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	// Создаем копию слайса
	messages := make([]*WebMessage, len(ws.recentMessages))
	copy(messages, ws.recentMessages)
	
	log.Printf("GetRecentMessages: returning %d messages", len(messages))
	
	return messages
}

// WebStatsResponse структура ответа со статистикой
type WebStatsResponse struct {
	TotalMessages int            `json:"TotalMessages"`
	LastUpdate    time.Time      `json:"LastUpdate"`
	Uptime        string         `json:"Uptime"`
	MessageTypes  map[string]int `json:"MessageTypes"`
}
```

## 4. Создайте недостающий template (table.html)
```html
<!DOCTYPE html>
<html>
<head>
    <title>Таблица данных датчиков</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="/static/css/tables.css">
</head>
<body>
    <div class="container">
        <h1>📊 Данные датчиков</h1>
        
        <div class="status" id="connection-status">🟡 Подключение...</div>
        
        <div class="controls">
            <button id="auto-scroll-btn" class="auto-scroll" onclick="toggleAutoScroll()">
                🔄 Авто-прокрутка: ВКЛ
            </button>
            <button onclick="clearTable()">🗑️ Очистить</button>
            <button onclick="loadTableData()">🔄 Обновить</button>
            <span>Всего записей: <span id="total-count">0</span></span>
        </div>
        
        <div class="table-container" id="table-container">
            <table>
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Паспорт</th>
                        <th>Тип данных</th>
                        <th>Код типа</th>
                        <th>Валидность</th>
                        <th>Смещение</th>
                        <th>Значение</th>
                        <th>Сырые байты</th>
                        <th>Время</th>
                    </tr>
                </thead>
                <tbody id="sensor-data">
                    <tr>
                        <td colspan="9" class="loading">🔄 Загрузка данных...</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    
    <script src="/static/js/dashboard.js"></script>
</body>
</html>
```

## 5. Создайте CSS стили
```css
/* tables.css */
/* Основные стили */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    min-height: 100vh;
    color: #333;
}

.container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 20px;
}

h1 {
    text-align: center;
    color: white;
    margin-bottom: 20px;
    text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

/* Статус подключения */
.status {
    text-align: center;
    padding: 10px;
    border-radius: 8px;
    margin-bottom: 20px;
    font-weight: bold;
    font-size: 16px;
}

.status.connected {
    background: rgba(76, 175, 80, 0.2);
    color: #4CAF50;
    border: 2px solid #4CAF50;
}

.status.disconnected {
    background: rgba(244, 67, 54, 0.2);
    color: #f44336;
    border: 2px solid #f44336;
}

/* Панель управления */
.controls {
    display: flex;
    gap: 15px;
    align-items: center;
    margin-bottom: 20px;
    padding: 15px;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    backdrop-filter: blur(10px);
    flex-wrap: wrap;
}

.controls button {
    padding: 10px 20px;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.3s ease;
    font-size: 14px;
}

.auto-scroll {
    background: #4CAF50;
    color: white;
}

.auto-scroll:hover {
    background: #45a049;
    transform: translateY(-2px);
}

.auto-scroll.disabled {
    background: #757575;
}

.controls button:not(.auto-scroll) {
    background: #2196F3;
    color: white;
}

.controls button:not(.auto-scroll):hover {
    background: #1976D2;
    transform: translateY(-2px);
}

.controls span {
    color: white;
    font-weight: bold;
    margin-left: auto;
}

#total-count {
    color: #FFD700;
    font-size: 18px;
}

/* Контейнер таблицы */
.table-container {
    background: rgba(255, 255, 255, 0.95);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    max-height: 70vh;
    overflow-y: auto;
}

/* Стили таблицы */
table {
    width: 100%;
    border-collapse: collapse;
    font-size: 14px;
}

thead {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    position: sticky;
    top: 0;
    z-index: 10;
}

thead th {
    padding: 15px 10px;
    text-align: left;
    font-weight: bold;
    border-bottom: 2px solid rgba(255, 255, 255, 0.2);
}

tbody tr {
    transition: background-color 0.2s ease;
}

tbody tr:nth-child(even) {
    background-color: rgba(0, 0, 0, 0.02);
}

tbody tr:hover {
    background-color: rgba(102, 126, 234, 0.1);
}

tbody td {
    padding: 12px 10px;
    border-bottom: 1px solid #e0e0e0;
    vertical-align: top;
}

/* Специальные ячейки */
.loading, .no-data, .error {
    text-align: center;
    padding: 40px;
    font-style: italic;
    color: #666;
}

.loading {
    background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
    background-size: 200% 100%;
    animation: loading 1.5s infinite;
}

.error {
    color: #f44336;
    background-color: rgba(244, 67, 54, 0.1);
}

@keyframes loading {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
}

/* Адаптивность */
@media (max-width: 1200px) {
    .container {
        padding: 15px;
    }
    
    table {
        font-size: 12px;
    }
    
    thead th, tbody td {
        padding: 8px 6px;
    }
}

@media (max-width: 768px) {
    .controls {
        flex-direction: column;
        align-items: stretch;
    }
    
    .controls button {
        margin-bottom: 10px;
    }
    
    .table-container {
        overflow-x: auto;
    }
    
    table {
        min-width: 800px;
    }
}

/* Скролбар */
.table-container::-webkit-scrollbar {
    width: 8px;
    height: 8px;
}

.table-container::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 4px;
}

.table-container::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 4px;
}

.table-container::-webkit-scrollbar-thumb:hover {
    background: #555;
}

/* Анимации */
@keyframes fadeIn {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
}

tbody tr {
    animation: fadeIn 0.3s ease;
}
```

## 6. Исправьте main.go для лучшей отладки
```GO
// main.go
package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"gendata-project/receiver-service/internal/display"
	"gendata-project/receiver-service/internal/web"
	"gendata-project/shared/protocol"
)

func main() {
	log.Println("🚀 Запуск сервиса получения данных...")

	// Создаем веб-сервер
	webServer := web.NewWebServer()

	// Запуск HTTP сервера для веб-интерфейса
	go func() {
		log.Println("🌐 Запуск HTTP сервера на порту 8081...")
		if err := webServer.Start("8081"); err != nil {
			log.Printf("❌ Ошибка HTTP сервера: %v", err)
		}
	}()

	// Запуск TCP сервера для приема данных
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("❌ Ошибка запуска TCP сервера: %v", err)
	}
	defer listener.Close()

	log.Println("🔌 TCP сервер запущен на порту 8080")
	log.Println("🌐 Веб-интерфейс доступен на:")
	log.Println("   - Главная страница: http://localhost:8081")
	log.Println("   - Таблица данных: http://localhost:8081/table")
	log.Println("   - API статистики: http://localhost:8081/api/stats")
	log.Println("   - API сообщений: http://localhost:8081/api/messages")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("🛑 Получен сигнал завершения...")
		listener.Close()
		os.Exit(0)
	}()

	// Основной цикл приема соединений
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("⚠️ Ошибка принятия соединения: %v", err)
			continue
		}
		go handleConnection(conn, webServer)
	}
}

func handleConnection(conn net.Conn, webServer *web.WebServer) {
	defer conn.Close()
	log.Printf("🔗 Новое соединение от %s", conn.RemoteAddr())

	messageCount := 0

	for {
		// Читаем размер сообщения
		var size uint32
		if err := binary.Read(conn, binary.LittleEndian, &size); err != nil {
			if err != io.EOF {
				log.Printf("❌ Ошибка чтения размера: %v", err)
			}
			break
		}

		// Читаем само сообщение
		msgData := make([]byte, size)
		if _, err := io.ReadFull(conn, msgData); err != nil {
			log.Printf("❌ Ошибка чтения сообщения: %v", err)
			break
		}

		// Парсим JSON сообщение
		msg, err := protocol.FromJSON(msgData)
		if err != nil {
			log.Printf("❌ Ошибка парсинга JSON: %v", err)
			continue
		}

		messageCount++
		log.Printf("📨 [%d] Получено сообщение типа: %s, размер данных: %d байт",
			messageCount, msg.Type, len(msg.Data))

		// Отображаем данные в консоли
		if err := display.DisplayDataBlock(msg.Data); err != nil {
			log.Printf("⚠️ Ошибка отображения данных: %v", err)
		}

		// Обрабатываем для веб-интерфейса
		webServer.ProcessMessage(msg)
		log.Printf("✅ Сообщение обработано для веб-интерфейса")
	}

	log.Printf("🔌 Соединение с %s закрыто. Обработано сообщений: %d", 
		conn.RemoteAddr(), messageCount)
}
```

## 7. Создайте скрипт для тестирования API
```bash
#!/bin/bash
# test-api.sh

echo "🧪 Тестирование API receiver-service..."

BASE_URL="http://localhost:8081"

echo ""
echo "1. Проверка health endpoint:"
curl -s "$BASE_URL/api/health" | jq '.' || echo "❌ Health endpoint недоступен"

echo ""
echo "2. Проверка stats endpoint:"
curl -s "$BASE_URL/api/stats" | jq '.' || echo "❌ Stats endpoint недоступен"

echo ""
echo "3. Проверка messages endpoint:"
curl -s "$BASE_URL/api/messages" | jq '.' || echo "❌ Messages endpoint недоступен"

echo ""
echo "4. Проверка главной страницы:"
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" "$BASE_URL/"

echo ""
echo "5. Проверка страницы таблицы:"
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" "$BASE_URL/table"

echo ""
echo "✅ Тестирование завершено"
```
* **Сделайте скрипт исполняемым:**
```bash
chmod +x scripts/test-api.sh
```

## 8. Запустите и проверьте
```bash
# Перезапустите receiver-service
docker-compose restart receiver

# Проверьте логи
docker-compose logs -f receiver
```

### Тестируйте API
```bash
./scripts/test-api.sh
```

### Теперь откройте в браузере:
```
http://localhost:8081/table - для таблицы данных
http://localhost:8081 - для главной страницы
```

## 9. Проверьте структуру файлов
### Убедитесь, что у вас правильная структура:
```
receiver-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── web/
│   │   ├── templates/
│   │   │   ├── index.html
│   │   │   └── table.html
│   │   ├── handlers.go
│   │   ├── stats.go
│   │   ├── types.go
│   │   ├── server.go
│   │   └── templates.go
│   ├── display/
│   └── parser/
├── static/
│   ├── css/
│   │   ├── main.css
│   │   └── tables.css
│   └── js/
│       └── dashboard.js
└── Dockerfile
```

## 10. Исправьте types.go если его нет
```GO
package web

import "time"

// WebMessage представляет сообщение для веб-интерфейса
type WebMessage struct {
	Type        string             `json:"Type"`
	Timestamp   time.Time          `json:"Timestamp"`
	Header      WebBlockHeader     `json:"Header"`
	Parameters  []WebSensorValue   `json:"Parameters"`
	DataSize    int                `json:"DataSize"`
	ReceivedAt  time.Time          `json:"ReceivedAt"`
}

// WebBlockHeader представляет заголовок блока данных для веб-интерфейса
type WebBlockHeader struct {
	GID         string    `json:"GID"`
	GIDDecoded  string    `json:"GIDDecoded"`
	DataType    string    `json:"DataType"`
	DataSymbol  string    `json:"DataSymbol"`
	ListNum     int       `json:"ListNum"`
	ListVer     int       `json:"ListVer"`
	Index       int       `json:"Index"`
	NumParams   int       `json:"NumParams"`
	Timestamp   time.Time `json:"Timestamp"`
}

// WebSensorValue представляет значение сенсора для веб-интерфейса
type WebSensorValue struct {
	Index       int    `json:"Index"`
	Passport    string `json:"Passport"`
	DataType    string `json:"DataType"`
	Validation  bool   `json:"Validation"`
	ShiftStatus int    `json:"ShiftStatus"`
	RawValue    string `json:"RawValue"`
	ParsedValue string `json:"ParsedValue"`
}

// WebStatsResponse структура ответа со статистикой
type WebStatsResponse struct {
	TotalMessages int            `json:"TotalMessages"`
	LastUpdate    time.Time      `json:"LastUpdate"`
	Uptime        string         `json:"Uptime"`
	MessageTypes  map[string]int `json:"MessageTypes"`
}
```

## 11. Обновите server.go
```GO
// server.go
package web

import (
	"log"
	"net/http"
	"path/filepath"
)

// SetupRoutes настраивает маршруты для веб-сервера
func (ws *WebServer) SetupRoutes() {
	// API маршруты
	http.HandleFunc("/api/stats", ws.StatsHandler)
	http.HandleFunc("/api/messages", ws.MessagesHandler)
	http.HandleFunc("/api/health", ws.HealthHandler)
	
	// Веб-страницы
	http.HandleFunc("/", ws.IndexHandler)
	http.HandleFunc("/table", ws.TableHandler)
	
	// Статические файлы
	staticDir := "./static/"
	if !fileExists(staticDir) {
		// Попробуем найти статические файлы относительно исполняемого файла
		staticDir = filepath.Join("receiver-service", "static")
		if !fileExists(staticDir) {
			log.Printf("⚠️ Статические файлы не найдены в %s", staticDir)
			staticDir = "./static/" // fallback
		}
	}
	
	log.Printf("📁 Статические файлы: %s", staticDir)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
}

// Start запускает веб-сервер
func (ws *WebServer) Start(port string) error {
	ws.SetupRoutes()
	
	log.Printf("🌐 HTTP сервер запущен на порту %s", port)
	log.Printf("📊 Веб-интерфейс доступен на http://localhost:%s", port)
	log.Printf("📋 API endpoints:")
	log.Printf("   - GET /api/health")
	log.Printf("   - GET /api/stats") 
	log.Printf("   - GET /api/messages")
	log.Printf("   - GET /table")
	
	return http.ListenAndServe(":"+port, nil)
}

// fileExists проверяет существование файла или директории
func fileExists(path string) bool {
	if _, err := http.Dir(path).Open("/"); err != nil {
		return false
	}
	return true
}
```

## 12. Создайте отладочный скрипт
```bash
#!/bin/bash
# debug-receiver.sh

echo "🔍 Отладка receiver-service..."

# Проверяем статус контейнера
echo "📦 Статус контейнера:"
docker ps | grep receiver || echo "❌ Контейнер receiver не запущен"

echo ""
echo "🌐 Проверка портов:"
netstat -tlnp | grep :808 || echo "❌ Порты 8080/8081 не прослушиваются"

echo ""
echo "📋 Последние логи receiver:"
docker logs --tail 20 go-send-recv-receiver

echo ""
echo "🧪 Тестирование API endpoints:"

# Health check
echo "1. Health check:"
curl -s -w "Status: %{http_code}\n" http://localhost:8081/api/health

echo ""
echo "2. Stats check:"
curl -s -w "Status: %{http_code}\n" http://localhost:8081/api/stats

echo ""
echo "3. Messages check:"
curl -s -w "Status: %{http_code}\n" http://localhost:8081/api/messages

echo ""
echo "4. Static files check:"
curl -s -w "Status: %{http_code}\n" http://localhost:8081/static/js/dashboard.js

echo ""
echo "📁 Проверка файловой структуры в контейнере:"
docker exec go-send-recv-receiver ls -la /root/static/ 2>/dev/null || echo "❌ Статические файлы не найдены в контейнере"

echo ""
echo "🔍 Отладка завершена"
```
```bash
chmod +x scripts/debug-receiver.sh
```

## 13. Обновите Dockerfile для receiver-service
```Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем исполняемый файл
COPY --from=builder /app/main .

# Копируем статические файлы
COPY --from=builder /app/static ./static/

# Создаем директории если их нет
RUN mkdir -p ./static/css ./static/js

# Проверяем что файлы скопировались
RUN ls -la ./static/

# Открываем порты
EXPOSE 8080 8081

# Запускаем приложение
CMD ["./main"]
```

## 14. Пересоберите и перезапустите
```bash
# Остановите контейнеры
docker-compose down

# Пересоберите образ receiver
docker-compose build --no-cache receiver

# Запустите заново
docker-compose up -d

# Проверьте логи
docker-compose logs -f receiver
```

## 15. Проверьте в браузере
* **Откройте в браузере:**
```
http://localhost:8081/api/health - должен вернуть JSON со статусом
http://localhost:8081/api/stats - должен вернуть статистику
http://localhost:8081/api/messages - должен вернуть массив сообщений
http://localhost:8081/table - должна открыться таблица с данными
```

## 16. Если данные все еще не идут, добавьте принудительную генерацию тестовых данных
```GO
// test_data.go
package web

import (
	"time"
	"gendata-project/shared/protocol"
)

// GenerateTestData создает тестовые данные для отладки
func (ws *WebServer) GenerateTestData() {
	// Создаем тестовое сообщение
	testData := []byte{
		0x01, 0x02, 0x03, 0x04, // GID и другие поля заголовка
		0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C,
		0x0D, 0x0E, 0x0F, 0x10,
		// Данные параметров
		0x11, 0x12, 0x13, 0x14,
		0x15, 0x16, 0x17, 0x18,
	}

	testMsg := &protocol.Message{
		Type:      "TEST_DATA",
		Data:      testData,
		Timestamp: time.Now(),
	}

	// Обрабатываем тестовое сообщение
	ws.ProcessMessage(testMsg)
}
```
* **Добавьте в main.go генерацию тестовых данных:**
```GO
// main.go
//...
// Добавьте после создания webServer
go func() {
	time.Sleep(5 * time.Second) // Ждем 5 секунд после запуска
	for {
		webServer.GenerateTestData()
		time.Sleep(10 * time.Second) // Генерируем тестовые данные каждые 10 секунд
	}
}()
//...
```

### Запустите отладочный скрипт:
```bash
./scripts/debug-receiver.sh
```
* **И проверьте таблицу в браузере: http://localhost:8081/table**