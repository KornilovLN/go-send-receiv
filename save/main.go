// receiver-service/cmd/main.go
// Это первая версия кода, которая запускает сервер для получения данных от генератора.
// без отображения данных на http-странице.
package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"gendata-project/receiver-service/internal/display"
	"gendata-project/shared/protocol"
)

func main() {
	log.Println("Запуск сервиса получения данных...")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
	defer listener.Close()

	log.Println("Сервер запущен на порту 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка принятия соединения: %v", err)
			continue
		}

		go handleConnection(conn)
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

		// Отображаем данные
		if err := display.DisplayDataBlock(msg.Data); err != nil {
			log.Printf("Ошибка отображения данных: %v", err)
		}
	}
}
