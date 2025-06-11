package web

import (
	"sync"
	"time"
)

// WebMessage - структура для отображения в веб-интерфейсе
type WebMessage struct {
	Type       string           `json:"type"`
	Timestamp  int64            `json:"timestamp"`
	Header     WebBlockHeader   `json:"header"`
	Parameters []WebSensorValue `json:"parameters"`
	DataSize   int              `json:"data_size"`
	ReceivedAt time.Time        `json:"received_at"`
}

// WebBlockHeader - заголовок блока для веб-интерфейса
type WebBlockHeader struct {
	GID        string    `json:"gid"`
	GIDDecoded string    `json:"gid_decoded"`
	DataType   string    `json:"data_type"`
	DataSymbol string    `json:"data_symbol"`
	ListNum    int       `json:"list_num"`
	ListVer    int       `json:"list_ver"`
	Index      int       `json:"index"`
	NumParams  int       `json:"num_params"`
	Timestamp  time.Time `json:"timestamp"`
}

// WebSensorValue - значение сенсора для веб-интерфейса
type WebSensorValue struct {
	Index       int    `json:"index"`
	Passport    string `json:"passport"`
	DataType    string `json:"data_type"`
	Validation  string `json:"validation"`
	ShiftStatus int    `json:"shift_status"`
	RawValue    string `json:"raw_value"`
	ParsedValue string `json:"parsed_value"`
}

// WebStats - статистика для веб-интерфейса
type WebStats struct {
	mu              sync.RWMutex
	TotalMessages   int            `json:"total_messages"`
	LastUpdate      time.Time      `json:"last_update"`
	MessagesPerType map[string]int `json:"messages_per_type"`
	RecentMessages  []*WebMessage  `json:"recent_messages"`
	IsActive        bool           `json:"is_active"`
	Uptime          string         `json:"uptime"`
	StartTime       time.Time      `json:"-"`
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
