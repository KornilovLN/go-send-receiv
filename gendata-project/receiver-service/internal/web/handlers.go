package web

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

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

// TableHandler обрабатывает страницу с табличным отображением данных
func (ws *WebServer) TableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(GetTableTemplate()))
}

// StatsHandler возвращает статистику в JSON
func (ws *WebServer) StatsHandler(w http.ResponseWriter, r *http.Request) {
	//stats := ws.Stats.GetStats()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-Control", "no-cache")

	// Обработка preflight запросов
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	stats := ws.Stats.GetStats()
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Error encoding stats", http.StatusInternalServerError)
		return
	}

}

// MessagesHandler возвращает последние сообщения в JSON
func (ws *WebServer) MessagesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-Control", "no-cache")

	// Обработка preflight запросов
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	messages := ws.Stats.GetRecentMessages()
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Error encoding messages", http.StatusInternalServerError)
		return
	}
}

// HealthHandler проверка здоровья сервиса
func (ws *WebServer) HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Добавляем расширенные CORS заголовки
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-Control", "no-cache")

	// Обработка preflight запросов
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
		"uptime":    ws.Stats.GetStats().Uptime,
	}
	json.NewEncoder(w).Encode(response)
}
