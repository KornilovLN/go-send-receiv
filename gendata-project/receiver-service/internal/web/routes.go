package web

import (
	"log"
	"net/http"
)

// SetupRoutes настраивает маршруты для веб-сервера
func (ws *WebServer) SetupRoutes() {
	http.HandleFunc("/", ws.IndexHandler)
	http.HandleFunc("/table", ws.TableHandler) // Новый маршрут для табличного отображения
	http.HandleFunc("/api/stats", ws.StatsHandler)
	http.HandleFunc("/api/messages", ws.MessagesHandler)
	http.HandleFunc("/api/health", ws.HealthHandler)

	// Статические файлы (если понадобятся)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

// Start запускает веб-сервер
func (ws *WebServer) Start(port string) error {
	ws.SetupRoutes()

	log.Printf("HTTP сервер запущен на порту %s", port)
	log.Printf("Веб-интерфейс доступен на http://localhost:%s", port)

	return http.ListenAndServe(":"+port, nil)
}
