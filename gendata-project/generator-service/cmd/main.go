package main

import (
	"encoding/binary"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"gendata-project/generator-service/internal/generator"
	"gendata-project/shared/protocol"
	"gendata-project/shared/types"
)

const (
	LISTNUM   = 1 // ListNum
	LISTVER   = 1 // ListVer
	INDEX     = 1 // Index
	NUMPARAMS = 5 // NumParams
)

func main() {
	log.Println("Запуск генератора данных...")

	// Получаем настройки из переменных окружения
	receiverHost := getEnv("RECEIVER_HOST", "receiver-service")
	receiverPort := getEnv("RECEIVER_PORT", "8080")
	intervalStr := getEnv("GENERATION_INTERVAL", "1")

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Printf("Неверный интервал генерации, используется значение по умолчанию: 1 секунд")
		interval = 1
	}

	receiverAddr := receiverHost + ":" + receiverPort
	log.Printf("Адрес получателя: %s", receiverAddr)
	log.Printf("Интервал генерации: %d секунд", interval)

	// Ожидание готовности получателя
	log.Println("Ожидание готовности получателя...")
	time.Sleep(10 * time.Second)

	for {
		// Подключение к получателю
		conn, err := net.Dial("tcp", receiverAddr)
		if err != nil {
			log.Printf("Ошибка подключения к получателю: %v. Повторная попытка через 1 секунд...", err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Println("Подключение к получателю установлено")

		// Генерация данных с заданным интервалом
		ticker := time.NewTicker(time.Duration(interval) * time.Second)

		dataTypes := []byte{
			types.TypeAnalog,
			types.TypeFloat,
			types.TypeFixed,
			types.TypeInt,
			types.TypeDiscrete,
		}

		connectionActive := true
		for connectionActive {
			select {
			case <-ticker.C:
				// Генерируем данные для каждого типа
				for _, dataType := range dataTypes {
					gid := generator.GIDCreater(1, 2, 1) // Север, NPP 2, Блок 1

					header := generator.HeadBlockCreater(
						gid,
						dataType,
						byte(LISTNUM),
						byte(LISTVER),
						int16(INDEX),
						int16(NUMPARAMS),
						//1,  // ListNum
						//1,  // ListVer
						//1,  //int16(i), // Index
						//5,  // NumParams
						time.Now(),
					)

					dataBlock := generator.GenBlock(header)

					// Создаем сообщение
					msg := protocol.Message{
						Type:      "data_block",
						Timestamp: time.Now().Unix(),
						Data:      dataBlock,
					}

					// Отправляем JSON
					jsonData, err := msg.ToJSON()
					if err != nil {
						log.Printf("Ошибка сериализации: %v", err)
						continue
					}

					// Отправляем размер сообщения, затем само сообщение
					size := uint32(len(jsonData))
					if err := binary.Write(conn, binary.LittleEndian, size); err != nil {
						log.Printf("Ошибка отправки размера: %v", err)
						connectionActive = false
						break
					}

					if _, err := conn.Write(jsonData); err != nil {
						log.Printf("Ошибка отправки данных: %v", err)
						connectionActive = false
						break
					}

					log.Printf("Отправлен блок данных типа %c, размер: %d байт",
						types.DataTypeSymbols[dataType], len(dataBlock))
				}
			}
		}

		ticker.Stop()
		conn.Close()
		log.Println("Соединение закрыто, переподключение...")
		time.Sleep(1 * time.Second)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
