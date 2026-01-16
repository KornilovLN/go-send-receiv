Вы правы, давайте сделаем правильную структуру с разделением на отдельные файлы. Вот полная организация:

Структура файлов
receiver-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── display/
│   │   └── display.go
│   ├── parser/
│   │   └── parser.go
│   └── web/
│       ├── handlers.go
│       ├── routes.go
│       ├── templates.go
│       └── types.go
├── go.mod
└── Dockerfile

Copy

Apply

1. Web типы данных
package web

import (
	"sync"
	"time"
)

// WebMessage - структура для отображения в веб-интерфейсе
type WebMessage struct {
	Type         string                     `json:"type"`
	Timestamp    int64                      `json:"timestamp"`
	Header       WebBlockHeader             `json:"header"`
	Parameters   []WebSensorValue           `json:"parameters"`
	DataSize     int                        `json:"data_size"`
	ReceivedAt   time.Time                  `json:"received_at"`
}

// WebBlockHeader - заголовок блока для веб-интерфейса
type WebBlockHeader struct {
	GID         string    `json:"gid"`
	GIDDecoded  string    `json:"gid_decoded"`
	DataType    string    `json:"data_type"`
	DataSymbol  string    `json:"data_symbol"`
	ListNum     int       `json:"list_num"`
	ListVer     int       `json:"list_ver"`
	Index       int       `json:"index"`
	NumParams   int       `json:"num_params"`
	Timestamp   time.Time `json:"timestamp"`
}

// WebSensorValue - значение сенсора для веб-интерфейса
type WebSensorValue struct {
	Index       int     `json:"index"`
	Passport    string  `json:"passport"`
	DataType    string  `json:"data_type"`
	Validation  string  `json:"validation"`
	ShiftStatus int     `json:"shift_status"`
	RawValue    string  `json:"raw_value"`
	ParsedValue string  `json:"parsed_value"`
}

// WebStats - статистика для веб-интерфейса
type WebStats struct {
	mu              sync.RWMutex
	TotalMessages   int                    `json:"total_messages"`
	LastUpdate      time.Time              `json:"last_update"`
	MessagesPerType map[string]int         `json:"messages_per_type"`
	RecentMessages  []*WebMessage          `json:"recent_messages"`
	IsActive        bool                   `json:"is_active"`
	Uptime          string                 `json:"uptime"`
	StartTime       time.Time              `json:"-"`
}

// NewWebStats создает новый экземпляр статистики
func NewWebStats() *WebStats {
	return &WebStats{
		MessagesPerType: make(map[string]int),
		RecentMessages:  make([]*WebMessage, 0, 10),
		IsActive:        true,
		StartTime:       time.Now(),
	}
}

// UpdateStats обновляет статистику новым сообщением
func (ws *WebStats) UpdateStats(webMsg *WebMessage) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.TotalMessages++
	ws.LastUpdate = time.Now()
	ws.MessagesPerType[webMsg.Type]++
	ws.Uptime = time.Since(ws.StartTime).Round(time.Second).String()

	// Добавляем в список последних сообщений
	ws.RecentMessages = append([]*WebMessage{webMsg}, ws.RecentMessages...)
	if len(ws.RecentMessages) > 10 {
		ws.RecentMessages = ws.RecentMessages[:10]
	}
}

// GetStats возвращает копию статистики
func (ws *WebStats) GetStats() WebStats {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	// Создаем копию для безопасного возврата
	statsCopy := WebStats{
		TotalMessages:   ws.TotalMessages,
		LastUpdate:      ws.LastUpdate,
		MessagesPerType: make(map[string]int),
		RecentMessages:  make([]*WebMessage, len(ws.RecentMessages)),
		IsActive:        ws.IsActive,
		Uptime:          ws.Uptime,
	}

	// Копируем карту типов сообщений
	for k, v := range ws.MessagesPerType {
		statsCopy.MessagesPerType[k] = v
	}

	// Копируем слайс сообщений
	copy(statsCopy.RecentMessages, ws.RecentMessages)

	return statsCopy
}

// GetRecentMessages возвращает копию последних сообщений
func (ws *WebStats) GetRecentMessages() []*WebMessage {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	messages := make([]*WebMessage, len(ws.RecentMessages))
	copy(messages, ws.RecentMessages)
	return messages
}

Copy

Apply

types.go
2. Web обработчики
package web

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"gendata-project/receiver-service/internal/parser"
	"gendata-project/shared/protocol"
	"gendata-project/shared/sensorvalue"
	"gendata-project/shared/types"
)

// WebServer представляет веб-сервер
type WebServer struct {
	Stats *WebStats
}

// NewWebServer создает новый веб-сервер
func NewWebServer() *WebServer {
	return &WebServer{
		Stats: NewWebStats(),
	}
}

// ConvertToWebMessage конвертирует protocol.Message в WebMessage
func (ws *WebServer) ConvertToWebMessage(msg *protocol.Message) *WebMessage {
	// Парсим заголовок из бинарных данных
	header, err := parser.ParseHeader(msg.Data)
	if err != nil {
		return &WebMessage{
			Type:       msg.Type,
			Timestamp:  msg.Timestamp,
			DataSize:   len(msg.Data),
			ReceivedAt: time.Now(),
			Header: WebBlockHeader{
				GID:        "Error",
				GIDDecoded: "Parse Error",
				DataType:   "Unknown",
				DataSymbol: "?",
			},
			Parameters: []WebSensorValue{},
		}
	}

	// Парсим параметры
	params, err := parser.ParseParameters(msg.Data, header.NumParams)
	if err != nil {
		params = []sensorvalue.SensorValue{}
	}

	// Конвертируем заголовок
	webHeader := WebBlockHeader{
		GID:        fmt.Sprintf("0x%02X", header.GID),
		GIDDecoded: parser.GetGIDString(header.GID),
		DataType:   fmt.Sprintf("0x%02X", header.DType),
		DataSymbol: string(types.DataTypeSymbols[header.DType]),
		ListNum:    int(header.Nlist),
		ListVer:    int(header.Vlist),
		Index:      int(header.Index),
		NumParams:  int(header.NumParams),
		Timestamp:  header.Timestamp,
	}

	// Конвертируем параметры
	webParams := make([]WebSensorValue, len(params))
	for i, param := range params {
		webParams[i] = WebSensorValue{
			Index:       i + 1,
			Passport:    fmt.Sprintf("0x%02X", param.Pasport),
			DataType:    string(types.DataTypeSymbols[param.DataType()]),
			Validation:  param.Validation(),
			ShiftStatus: int(param.ShiftStatus()),
			RawValue:    fmt.Sprintf("% X", param.RawValue),
			ParsedValue: parseValueForWeb(param),
		}
	}

	return &WebMessage{
		Type:       msg.Type,
		Timestamp:  msg.Timestamp,
		Header:     webHeader,
		Parameters: webParams,
		DataSize:   len(msg.Data),
		ReceivedAt: time.Now(),
	}
}

// parseValueForWeb парсит значение сенсора для отображения в веб
func parseValueForWeb(param sensorvalue.SensorValue) string {
	switch param.DataType() {
	case types.TypeAnalog, types.TypeFloat:
		value := binary.LittleEndian.Uint32(param.RawValue[:])
		floatVal := math.Float32frombits(value)
		return fmt.Sprintf("%.4f", floatVal)

	case types.TypeFixed:
		shift := param.Pasport & 0x07
		intVal := int32(binary.LittleEndian.Uint32(param.RawValue[:]))
		fixedVal := float32(intVal) / float32(math.Pow10(int(shift)))
		return fmt.Sprintf("%.2f (shift=%d)", fixedVal, shift)

	case types.TypeInt:
		intVal := int32(binary.LittleEndian.Uint32(param.RawValue[:]))
		return fmt.Sprintf("%d", intVal)

	case types.TypeDiscrete:
		return fmt.Sprintf("[% X]", param.RawValue)

	default:
		return "Unknown"
	}
}

// ProcessMessage обрабатывает новое сообщение
func (ws *WebServer) ProcessMessage(msg *protocol.Message) {
	webMsg := ws.ConvertToWebMessage(msg)
	ws.Stats.UpdateStats(webMsg)
}

// IndexHandler обрабатывает главную страницу
func (ws *WebServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(GetIndexTemplate()))
}

// StatsHandler возвращает статистику в JSON
func (ws *WebServer) StatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := ws.Stats.GetStats()
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Error encoding stats", http.StatusInternalServerError)
		return
	}
}

// MessagesHandler возвращает последние сообщения в JSON
func (ws *WebServer) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages := ws.Stats.GetRecentMessages()
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Error encoding messages", http.StatusInternalServerError)
		return
	}
}

// HealthHandler проверка здоровья сервиса
func (ws *WebServer) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
		"uptime":    ws.Stats.GetStats().Uptime,
	}
	json.NewEncoder(w).Encode(response)
}

Copy

Apply

handlers.go
3. Маршруты
package web

import (
	"log"
	"net/http"
)

// SetupRoutes настраивает маршруты для веб-сервера
func (ws *WebServer) SetupRoutes() {
	http.HandleFunc("/", ws.IndexHandler)
	http.HandleFunc("/api/stats", ws.StatsHandler)
	http.HandleFunc("/api/messages", ws.MessagesHandler)
	http.HandleFunc("/api/health", ws.HealthHandler)
	
	// Статические файлы (если понадобятся)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

// Start запускает веб-сервер
func (ws *WebServer) Start(port string) error {
	ws.SetupRoutes()
	
	log.Printf("HTTP сервер запущен на порту %s", port)
	log.Printf("Веб-интерфейс доступен на http://localhost:%s", port)
	
	return http.ListenAndServe(":"+port, nil)
}

Copy

Apply

routes.go
4. HTML шаблоны
package web

// GetIndexTemplate возвращает HTML шаблон главной страницы
func GetIndexTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>GenData Monitor</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        * { box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            margin: 0; 
            padding: 20px; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container { 
            max-width: 1400px; 
            margin: 0 auto; 
            background: rgba(255,255,255,0.95);
            border-radius: 15px;
            padding: 30px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
        }
        .header { 
            text-align: center; 
            margin-bottom: 30px; 
            padding-bottom: 20px;
            border-bottom: 2px solid #e9ecef;
        }
        .header h1 { 
            color: #2c3e50; 
            margin: 0; 
            font-size: 2.5em;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.1);
        }
        .status { 
            display: inline-block; 
            padding: 8px 20px; 
            border-radius: 25px; 
            background: linear-gradient(45deg, #27ae60, #2ecc71);
            color: white; 
            font-size: 14px; 
            font-weight: bold;
            margin-top: 10px;
            box-shadow: 0 4px 8px rgba(39,174,96,0.3);
        }
        
        .cards { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));

Copy

Apply

templates.go

Продолжи



Context
— none
4. HTML шаблоны (продолжение)
package web

// GetIndexTemplate возвращает HTML шаблон главной страницы
func GetIndexTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>GenData Monitor</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        * { box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            margin: 0; 
            padding: 20px; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container { 
            max-width: 1400px; 
            margin: 0 auto; 
            background: rgba(255,255,255,0.95);
            border-radius: 15px;
            padding: 30px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
        }
        .header { 
            text-align: center; 
            margin-bottom: 30px; 
            padding-bottom: 20px;
            border-bottom: 2px solid #e9ecef;
        }
        .header h1 { 
            color: #2c3e50; 
            margin: 0; 
            font-size: 2.5em;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.1);
        }
        .status { 
            display: inline-block; 
            padding: 8px 20px; 
            border-radius: 25px; 
            background: linear-gradient(45deg, #27ae60, #2ecc71);
            color: white; 
            font-size: 14px; 
            font-weight: bold;
            margin-top: 10px;
            box-shadow: 0 4px 8px rgba(39,174,96,0.3);
        }
        
        .cards { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 25px; 
            margin-bottom: 40px; 
        }
        .card { 
            background: linear-gradient(145deg, #ffffff, #f8f9fa);
            padding: 25px; 
            border-radius: 15px; 
            box-shadow: 0 8px 25px rgba(0,0,0,0.1);
            border: 1px solid rgba(255,255,255,0.2);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }
        .card:hover {
            transform: translateY(-5px);
            box-shadow: 0 15px 35px rgba(0,0,0,0.15);
        }
        .card h3 { 
            margin: 0 0 15px 0; 
            color: #34495e; 
            font-size: 1.2em;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .stat-value { 
            font-size: 2.5em; 
            font-weight: bold; 
            color: #3498db; 
            margin: 15px 0; 
            text-shadow: 1px 1px 2px rgba(0,0,0,0.1);
        }
        .stat-label { 
            color: #7f8c8d; 
            font-size: 14px; 
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        .messages-container { 
            background: linear-gradient(145deg, #ffffff, #f8f9fa);
            border-radius: 15px; 
            box-shadow: 0 8px 25px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .messages-header { 
            padding: 25px; 
            border-bottom: 2px solid #e9ecef;
            background: linear-gradient(90deg, #f8f9fa, #ffffff);
            display: flex; 
            justify-content: space-between; 
            align-items: center; 
        }
        .messages-header h3 { 
            margin: 0; 
            color: #2c3e50; 
            font-size: 1.4em;
        }
        .refresh-btn { 
            background: linear-gradient(45deg, #3498db, #2980b9);
            color: white; 
            border: none; 
            padding: 12px 24px; 
            border-radius: 25px; 
            cursor: pointer; 
            transition: all 0.3s ease;
            font-weight: bold;
            box-shadow: 0 4px 12px rgba(52,152,219,0.3);
        }
        .refresh-btn:hover { 
            background: linear-gradient(45deg, #2980b9, #1f4e79);
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(52,152,219,0.4);
        }
        
        .messages-list { 
            max-height: 700px; 
            overflow-y: auto; 
        }
        .message { 
            border-bottom: 1px solid #e9ecef;
            padding: 25px; 
            transition: background-color 0.3s ease;
        }
        .message:hover {
            background-color: #f8f9fa;
        }
        .message:last-child { 
            border-bottom: none; 
        }
        .message-header { 
            display: flex; 
            justify-content: space-between; 
            align-items: center; 
            margin-bottom: 20px; 
        }
        .message-type { 
            background: linear-gradient(45deg, #e74c3c, #c0392b);
            color: white; 
            padding: 6px 16px; 
            border-radius: 20px; 
            font-size: 12px; 
            font-weight: bold;
            box-shadow: 0 2px 8px rgba(231,76,60,0.3);
        }
        .message-time { 
            color: #7f8c8d; 
            font-size: 14px; 
            font-family: monospace;
        }

Copy

Apply

templates.go

Продолжи



Context
— none
4. HTML шаблоны (продолжение)
        
        .block-header { 
            background: linear-gradient(145deg, #ecf0f1, #d5dbdb);
            padding: 20px; 
            border-radius: 10px; 
            margin-bottom: 20px;
            border-left: 4px solid #3498db;
        }
        .block-header-grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 15px; 
        }
        .block-header-item { 
            display: flex; 
            justify-content: space-between;
            align-items: center;
            padding: 8px 0;
        }
        .block-header-label { 
            font-weight: bold; 
            color: #2c3e50; 
            font-size: 14px;
        }
        .block-header-value { 
            color: #34495e; 
            font-family: monospace;
            background: rgba(255,255,255,0.7);
            padding: 4px 8px;
            border-radius: 4px;
        }
        
        .parameters { 
            margin-top: 20px; 
        }
        .parameters-header { 
            font-weight: bold; 
            color: #2c3e50; 
            margin-bottom: 15px;
            font-size: 16px;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .parameter { 
            background: linear-gradient(145deg, #f8f9fa, #e9ecef);
            padding: 15px; 
            margin: 8px 0; 
            border-radius: 8px; 
            border-left: 4px solid #3498db;
            transition: all 0.3s ease;
        }
        .parameter:hover {
            background: linear-gradient(145deg, #e9ecef, #dee2e6);
            transform: translateX(5px);
        }
        .parameter-grid { 
            display: grid; 
            grid-template-columns: 60px 100px 120px 1fr;
            gap: 15px; 
            align-items: center; 
        }
        .parameter-index { 
            font-weight: bold; 
            color: #e74c3c; 
            font-size: 16px;
            text-align: center;
        }
        .parameter-type { 
            background: linear-gradient(45deg, #3498db, #2980b9);
            color: white; 
            padding: 4px 12px; 
            border-radius: 15px; 
            font-size: 12px; 
            text-align: center;
            font-weight: bold;
            box-shadow: 0 2px 6px rgba(52,152,219,0.3);
        }
        .parameter-validation {
            font-size: 12px;
            padding: 2px 8px;
            border-radius: 10px;
            text-align: center;
            font-weight: bold;
        }
        .parameter-validation.valid {
            background: #d4edda;
            color: #155724;
        }
        .parameter-validation.invalid {
            background: #f8d7da;
            color: #721c24;
        }
        .parameter-value { 
            font-family: 'Courier New', monospace;
            background: rgba(255,255,255,0.8);
            padding: 8px 12px; 
            border-radius: 6px;
            border: 1px solid #dee2e6;
            font-size: 14px;
        }
        
        .loading { 
            text-align: center; 
            padding: 60px; 
            color: #7f8c8d;
            font-size: 18px;
        }
        .no-data { 
            text-align: center; 
            padding: 60px; 
            color: #95a5a6;
            font-size: 18px;
        }
        
        .spinner {
            border: 4px solid #f3f3f3;
            border-top: 4px solid #3498db;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 20px auto;
        }
        
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .message-types-list {
            display: flex;
            flex-direction: column;
            gap: 8px;
        }
        .message-type-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 8px 12px;
            background: rgba(52,152,219,0.1);
            border-radius: 6px;
            border-left: 3px solid #3498db;
        }
        .message-type-name {
            font-weight: bold;
            color: #2c3e50;
        }
        .message-type-count {
            background: #3498db;
            color: white;
            padding: 2px 8px;
            border-radius: 10px;
            font-size: 12px;
            font-weight: bold;
        }
        
        @media (max-width: 768px) {
            .container { padding: 15px; }
            .cards { grid-template-columns: 1fr; }
            .parameter-grid { 
                grid-template-columns: 1fr;
                gap: 8px;
            }
            .block-header-grid {
                grid-template-columns: 1fr;
            }
            .messages-header {
                flex-direction: column;
                gap: 15px;
                text-align: center;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔄 GenData Services Monitor</h1>
            <div class="status" id="status">🟢 Активен</div>
        </div>
        
        <div class="cards">
            <div class="card">
                <h3>📊 Всего сообщений</h3>
                <div class="stat-value" id="total-messages">0</div>
                <div class="stat-label">Получено пакетов</div>
            </div>
            <div class="card">
                <h3>⏰ Последнее обновление</h3>
                <div class="stat-value" id="last-update" style="font-size: 1.2em;">-</div>
                <div class="stat-label">Время получения</div>
            </div>
            <div class="card">
                <h3>🕐 Время работы</h3>
                <div class="stat-value" id="uptime" style="font-size: 1.5em;">-</div>
                <div class="stat-label">Uptime сервиса</div>
            </div>
            <div class="card">
                <h3>📈 Типы данных</h3>
                <div id="message-types" class="message-types-list">
                    <div class="loading">Загрузка...</div>
                </div>
            </div>
        </div>
        
        <div class="messages-container">
            <div class="messages-header">
                <h3>📨 Последние сообщения</h3>
                <button class="refresh-btn" onclick="loadData()">🔄 Обновить</button>
            </div>
            <div class="messages-list" id="messages-list">
                <div class="loading">
                    <div class="spinner"></div>
                    Загрузка данных...
                </div>
            </div>
        </div>
    </div>

    <script>
        let isLoading = false;
        
        function formatTime(timestamp) {
            return new Date(timestamp * 1000).toLocaleString('ru-RU');
        }
        
        function formatDateTime(dateStr) {
            return new Date(dateStr).toLocaleString('ru-RU');
        }
        
        function updateStats(stats) {
            document.getElementById('total-messages').textContent = stats.total_messages || 0;
            document.getElementById('last-update').textContent = 
                stats.last_update ? formatDateTime(stats.last_update) : '-';
            document.getElementById('uptime').textContent = stats.uptime || '-';
            
            // Обновляем типы сообщений
            const typesContainer = document.getElementById('message-types');
            if (stats.messages_per_type && Object.keys(stats.messages_per_type).length > 0) {
                let typesHtml = '';
                for (const [type, count] of Object.entries(stats.messages_per_type)) {
                    typesHtml += ` + "`" + `
                        <div class="message-type-item">
                            <span class="message-type-name">${type}</span>
                            <span class="message-type-count">${count}</span>
                        </div>
                    ` + "`" + `;
                }
                typesContainer.innerHTML = typesHtml;
            } else {
                typesContainer.innerHTML = '<div class="no-data">Нет данных</div>';
            }
        }
        
        function renderMessages(messages) {
            const container = document.getElementById('messages-list');
            
            if (!messages || messages.length === 0) {
                container.innerHTML = '<div class="no-data">📭 Нет сообщений</div>';
                return;
            }
            
            let html = '';
            messages.forEach(msg => {
                html += ` + "`" + `
                    <div class="message">
                        <div class="message-header">
                            <span class="message-type">${msg.type}</span>
                            <span class="message-time">${formatDateTime(msg.received_at)}</span>
                        </div>
                        
                        <div class="block-header">
                            <div class="block-header-grid">
                                <div class="block-header-item">
                                    <span class="block-header-label">GID:</span>
                                    <span class="block-header-value">${msg.header.gid}</span>
                                </div>
                                <div class="block-header-item">
                                    <span class="block-header-label">Регион:</span>
                                    <span class="block-header-value">${msg.header.gid_decoded}</span>
                                </div>
                                <div class="block-header-item">
                                    <span class="block-header-label">Тип данных:</span>
                                    <span class="block-header-value">${msg.header.data_symbol} (${msg.header.data_type})</span>
                                </div>
                                <div class="block-header-item">
                                    <span class="block-header-label">Список:</span>
                                    <span class="block-header-value">${msg.header.list_num}.${msg.header.list_ver}</span>
                                </div>
                                <div class="block-header-item">
                                    <span class="block-header-label">Индекс:</span>
                                    <span class="block-header-value">${msg.header.index}</span>
                                </div>
                                <div class="block-header-item">
                                    <span class="block-header-label">Параметров:</span>
                                    <span class="block-header-value">${msg.header.num_params}</span>
                                </div>
                                <div class="block-header-item">
                                    <span class="block-header-label">Время:</span>
                                    <span class="block-header-value">${formatDateTime(msg.header.timestamp)}</span>
                                </div>
                                <div class="block-header-item">
                                    <span class="block-header-label">Размер:</span>
                                    <span class="block-header-value">${msg.data_size} байт</span>
                                </div>
                            </div>
                        </div>
                        
                        <div class="parameters">
                            <div class="parameters-header">🔧 Параметры (${msg.parameters.length})</div>
                ` + "`" + `;
                
                msg.parameters.forEach(param => {
                    const validationClass = param.validation.toLowerCase() === 'valid' ? 'valid' : 'invalid';
                    html += ` + "`" + `
                        <div class="parameter">
                            <div class="parameter-grid">
                                <div class="parameter-index">[${param.index}]</div>
                                <div class="parameter-type">${param.data_type}</div>
                                <div class="parameter-validation ${validationClass}">${param.validation}</div>
                                <div class="parameter-value">
                                    <strong>Значение:</strong> ${param.parsed_value}<br>
                                    <strong>Raw:</strong> ${param.raw_value}<br>
                                    <strong>Passport:</strong> ${param.passport}
                                </div>
                            </div>
                        </div>
                    ` + "`" + `;
                });
                
                html += ` + "`" + `
                        </div>
                    </div>
                ` + "`" + `;
            });
            
            container.innerHTML = html;
        }
        
        async function loadData() {
            if (isLoading) return;
            isLoading = true;
            
            try {
                // Загружаем статистику
                const statsResponse = await fetch('/api/stats');
                if (statsResponse.ok) {
                    const stats = await statsResponse.json();
                    updateStats(stats);
                }
                
                // Загружаем сообщения
                const messagesResponse = await fetch('/api/messages');
                if (messagesResponse.ok) {
                    const messages = await messagesResponse.json();
                    renderMessages(messages);
                }
                
                // Обновляем статус
                document.getElementById('status').innerHTML = '🟢 Активен';
                document.getElementById('status').style.background = 'linear-gradient(45deg, #27ae60, #2ecc71)';
                
            } catch (error) {
                console.error('Ошибка загрузки данных:', error);
                document.getElementById('status').innerHTML = '🔴 Ошибка';
                document.getElementById('status').style.background = 'linear-gradient(45deg, #e74c3c, #c0392b)';
            } finally {
                isLoading = false;
            }
        }
        
        // Загружаем данные при загрузке страницы
        document.addEventListener('DOMContentLoaded', loadData);
        
        // Автоматическое обновление каждые 3 секунды
        setInterval(loadData, 3000);
    </script>
</body>