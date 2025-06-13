package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"gendata-project/receiver-service/internal/display"
	//"gendata-project/receiver-service/internal/parser"
	"gendata-project/receiver-service/internal/web"
	"gendata-project/shared/protocol"
)

func main() {
	log.Println("Запуск сервиса получения данных...")

	// Создаем веб-сервер
	webServer := web.NewWebServer()

	// Запуск HTTP сервера для веб-интерфейса
	go func() {
		if err := webServer.Start("8081"); err != nil {
			log.Printf("Ошибка HTTP сервера: %v", err)
		}
	}()

	// Запуск TCP сервера для приема данных
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Ошибка запуска TCP сервера: %v", err)
	}
	defer listener.Close()

	log.Println("TCP сервер запущен на порту 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка принятия соединения: %v", err)
			continue
		}

		go handleConnection(conn, webServer)
	}
}

func handleConnection(conn net.Conn, webServer *web.WebServer) {
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

		// Отображаем данные в консоли (используем display)
		if err := display.DisplayDataBlock(msg.Data); err != nil {
			log.Printf("Ошибка отображения данных: %v", err)
		}

		// Обрабатываем для веб-интерфейса (используем web)
		webServer.ProcessMessage(msg)
	}
}
