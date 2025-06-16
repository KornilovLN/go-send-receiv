# Тест-испытание поставщика и приемника данных

## 1. Структура проекта: Генерация папок проекта
```bash
mkdir -p gendata-project/{generator-service,receiver-service,shared,deploy/ansible}
cd gendata-project
```

## 2. Общие компоненты (shared) для обоих сервисов
```GO
// types.go
package types

import "time"

const (
	HeaderLength = 12
	PayloadSize  = 5

	// Типы данных
	TypeAnalog   = 0x10
	TypeFloat    = 0x30
	TypeFixed    = 0x20
	TypeInt      = 0x40
	TypeDiscrete = 0x50

	// Маски и сдвиги для GID
	RegionMask   = 0xE0
	RegionShift  = 5
	RegionCheck  = 0x07
	NppMask      = 0x18
	NppShift     = 3
	NppCheck     = 0x03
	EblockMask   = 0x07
	EblockShift  = 0
	EblockCheck  = 0x07
)

// Region один из 5 возможных для размещения станции
const (
	RegionWest  = 0x00 // Западный регион
	RegionCentr = 0x01 // Центральный регион
	RegionEast  = 0x02 // Восточный регион
	RegionSouth = 0x03 // Южный регион
    RegionNorth = 0x04 // Северный регион
)

var DataTypeSymbols = map[byte]rune{
	TypeAnalog:   'A', // Аналоговый тип данных
	TypeFloat:    'F', // Тип данных с плавающей точкой
	TypeFixed:    'I', // Фиксированный тип данных
	TypeInt:      'B', // Целочисленный тип данных
	TypeDiscrete: 'D', // Дискретный тип данных
}

// RegionString - строковые представления регионов.
var RegionString = map[byte]string{
	RegionWest:  "Region-West",
	RegionCentr: "Region-Centr",
	RegionEast:  "Region-East",
	RegionSouth: "Region-South",
    RegionNorth: "Region-North",
}

// BlockData представляет block of data в СПД.
// type BlockData struct {
type BlockHeader struct {
	GID byte // Group ID (мл.3 бита - номер энергоблока, ст.5 бит - группа)
	// [7-6 bits - регион, 5-3 bits - номер АЭС, 2-0 bits - Энергоблок]
	DType     byte      // Тип данных (0x00-аналог,0x10-float,0x20-fixed,0x30-int,0x40-discrete)
	Nlist     byte      // Номер списка (0-255)  		(но меньше 255, а 255 - это error)
	Vlist     byte      // Версия списка (0-255) 		(но гораздо меньше 255, а 255 - это error)
	Index     int16     // Индекс в списке (0-65535) 	(но меньше 65535, а 65535 - это error)
	NumParams int16     // Колич.дан. в блоке (0-65535) (но меньше 65535, а 65535 - это error)
	Timestamp time.Time // Время в формате Unix 		(секунды с 1970-01-01T00:00:00Z)
}

// DataPacket представляет блок данных передаваемый по сети
type DataPacket struct {
	Header BlockHeader
	Data   []byte
}
```

```GO
// sensor.go
package sensorvalue

import "gendata-project/shared/types"

// Структура SensorValue для хранения данных сенсоров
type SensorValue struct {
	Pasport  byte
	RawValue [4]byte
}

// Возвращает тип данных из Pasport
func (sv SensorValue) DataType() byte {
	return sv.Pasport & 0x70
}

// Получение символа типа данных по значению кода типа из Pasport
func (sv *SensorValue) TypeSymbol() byte {
	return types.DataTypeSymbols[sv.DataType()]
}

func (sv SensorValue) Validation() string {
	switch (sv.Pasport >> 3) & 0x01 {
	case 0:
		return "Valid"
	case 1:
		return "Invalid"
	default:
		return "Unknown"
	}
}

func (sv SensorValue) ShiftStatus() byte {
	return sv.Pasport & 0x07
}
```

```GO
//protocol.go
package protocol

import (
	"encoding/json"
	"gendata-project/shared/types"
)

type Message struct {
	Type      string            `json:"type"`
	Timestamp int64             `json:"timestamp"`
	Header    types.BlockHeader `json:"header"`
	Data      []byte            `json:"data"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
```


## 3. Генератор данных (generator-service)
```GO
// generator.go
package generator

import (
	"encoding/binary"
	"math"
	"time"

	"gendata-project/shared/sensorvalue"
	"gendata-project/shared/types"
)

func GenAnalog(value float32, shift byte) []byte {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], math.Float32bits(value))
	return append([]byte{types.TypeAnalog | shift}, bytes[:]...)
}

func GenFloat(value float32, shift byte) []byte {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], math.Float32bits(value))
	return append([]byte{types.TypeFloat | shift}, bytes[:]...)
}

func GenInt(value int32, shift byte) []byte {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], uint32(value))
	return append([]byte{types.TypeInt | shift}, bytes[:]...)
}

func GenFixed(value float32, shift byte) []byte {
	scaled := int32(math.Round(float64(value) * math.Pow10(int(shift))))
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], uint32(scaled))
	return append([]byte{types.TypeFixed | shift}, bytes[:]...)
}

func GenDiscrete(b1, b2, b3, b4 byte) []byte {
	return []byte{
		types.TypeDiscrete | 0x00,
		b1, b2, b3, b4,
	}
}

func GIDCreater(Region byte, Nnpp byte, Eblock byte) byte {
	var GID byte = 0x00
	GID = (Region & types.RegionCheck) << types.RegionShift
	GID |= (Nnpp & types.NppCheck) << types.NppShift
	GID |= (Eblock & types.EblockCheck) << types.EblockShift
	return GID
}

func HeadBlockCreater(
	GID byte,
	DataType byte,
	ListNum byte,
	ListVer byte,
	Index int16,
	NumParams int16,
	Timestamp time.Time,
) []byte {
	header := make([]byte, types.HeaderLength)

	header[0] = GID
	header[1] = DataType
	header[2] = ListNum
	header[3] = ListVer

	binary.LittleEndian.PutUint16(header[4:6], uint16(Index))
	binary.LittleEndian.PutUint16(header[6:8], uint16(NumParams))

	unixTime := uint32(Timestamp.Unix())
	binary.LittleEndian.PutUint32(header[8:12], unixTime)

	return header
}

func GenBlock(header []byte) []byte {
	if len(header) != types.HeaderLength {
		panic("Header должен быть 12 байт")
	}

	numParams := int16(binary.LittleEndian.Uint16(header[6:8]))
	packet := make([]byte, 0, types.HeaderLength+int(numParams)*types.PayloadSize)
	packet = append(packet, header...)

	for i := int16(0); i < numParams; i++ {
		var param []byte

		switch header[1] {
		case types.TypeAnalog:
			param = GenAnalog(float32(i)-486.57, 0x05)
		case types.TypeFloat:
			param = GenFloat(float32(i)+3.1415927, 0x00)
		case types.TypeFixed:
			param = GenFixed(float32(i)+59.27, 0x02)
		case types.TypeInt:
			param = GenInt(int32(i)+54321, 0x00)
		case types.TypeDiscrete:
			val := int16(10)
			b1 := byte((val + i) % 256)
			b2 := byte(((val + i) + 1) % 256)
			b3 := byte(((val + i) + 2) % 256)
			b4 := byte(((val + i) + 3) % 256)
			param = GenDiscrete(b1, b2, b3, b4)
		default:
			param = make([]byte, types.PayloadSize)
		}

		packet = append(packet, param...)
	}

	return packet
}
```

```GO
// main.go
package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"gendata-project/generator-service/internal/generator"
	"gendata-project/shared/protocol"
	"gendata-project/shared/types"
)

func main() {
	log.Println("Запуск генератора данных...")

	// Подключение к получателю
	conn, err := net.Dial("tcp", "receiver-service:8080")
	if err != nil {
		log.Fatalf("Ошибка подключения к получателю: %v", err)
	}
	defer conn.Close()

	// Генерация данных каждые 5 секунд
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	dataTypes := []byte{
		types.TypeAnalog,
		types.TypeFloat,
		types.TypeFixed,
		types.TypeInt,
		types.TypeDiscrete,
	}

	for {
		select {
		case <-ticker.C:
			// Генерируем данные для каждого типа
			for i, dataType := range dataTypes {
				gid := generator.GIDCreater(1, 2, 1) // Север, NPP 2, Блок 1
				
				header := generator.HeadBlockCreater(
					gid,
					dataType,
					1,                    // ListNum
					1,                    // ListVer
					int16(i),            // Index
					5,                   // NumParams
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
					continue
				}

				if _, err := conn.Write(jsonData); err != nil {
					log.Printf("Ошибка отправки данных: %v", err)
					continue
				}

				log.Printf("Отправлен блок данных типа %c, размер: %d байт", 
					types.DataTypeSymbols[dataType], len(dataBlock))
			}
		}
	}
}
```

```Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o generator ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/generator .
CMD ["./generator"]
```

## 4. Получатель данных (receiver-service)
```GO
// parser.go
package parser

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"gendata-project/shared/sensorvalue"
	"gendata-project/shared/types"
)

func ParseHeader(data []byte) (types.BlockHeader, error) {
	if len(data) < types.HeaderLength {
		return types.BlockHeader{}, fmt.Errorf("данные слишком короткие: %d байт", len(data))
	}

	return types.BlockHeader{
		GID:       data[0],
		DType:     data[1],
		Nlist:     data[2],
		Vlist:     data[3],
		Index:     int16(binary.LittleEndian.Uint16(data[4:6])),
		NumParams: int16(binary.LittleEndian.Uint16(data[6:8])),
		Timestamp: time.Unix(int64(binary.LittleEndian.Uint32(data[8:12])), 0),
	}, nil
}

func ParseParameters(data []byte, numParams int16) ([]sensorvalue.SensorValue, error) {
	params := make([]sensorvalue.SensorValue, 0, numParams)

	for i := 0; i < int(numParams); i++ {
		offset := types.HeaderLength + i*types.PayloadSize
		if offset+types.PayloadSize > len(data) {
			return nil, fmt.Errorf("недостаточно данных для параметра %d", i)
		}

		paramData := data[offset : offset+types.PayloadSize]
		params = append(params, sensorvalue.SensorValue{
			Pasport:  paramData[0],
			RawValue: [4]byte{paramData[1], paramData[2], paramData[3], paramData[4]},
		})
	}

	return params, nil
}

func GetRegionGID(GID byte) byte {
	return (GID & types.RegionMask) >> types.RegionShift
}

func GetNPPGID(GID byte) byte {
	return (GID & types.NppMask) >> types.NppShift
}

func GetEblockGID(GID byte) byte {
	return (GID & types.EblockMask) >> types.EblockShift
}

func GetGIDString(GID byte) string {
	region := GetRegionGID(GID)
	npp := GetNPPGID(GID)
	eblock := GetEblockGID(GID)

	return fmt.Sprintf("{%s, NPP: %d, Eblock: %d}",
		types.RegionString[region], npp, eblock)
}
```

```GO
// display.go
package display

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"gendata-project/receiver-service/internal/parser"
	"gendata-project/shared/sensorvalue"
	"gendata-project/shared/types"
)

const LINELEN = 80

func PrintHeaderLine(title string, char rune, length int) string {
	if len(title) >= length-4 {
		return strings.Repeat(string(char), length)
	}
	
	padding := (length - len(title) - 2) / 2
	leftPad := strings.Repeat(string(char), padding)
	rightPad := strings.Repeat(string(char), length-len(title)-2-padding)
	
	return leftPad + " " + title + " " + rightPad
}

func PrintHeader(h types.BlockHeader) {
	fmt.Println(PrintHeaderLine("Заголовок блока", '=', LINELEN))

	fmt.Printf("GID: 0x%02X => ", h.GID)
	fmt.Println(parser.GetGIDString(h.GID))

	fmt.Printf("Dtype: %c\n", types.DataTypeSymbols[h.DType])
	fmt.Printf("Nlist: %d\n", h.Nlist)
	fmt.Printf("Vlist: %d\n", h.Vlist)
	fmt.Printf("Index: %d\n", h.Index)
	fmt.Printf("N prm: %d\n", h.NumParams)
	fmt.Printf("%s\n", h.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Println()
}

func PrintParameter(p sensorvalue.SensorValue, index int) {
	fmt.Printf("[%3d]  ", index+1)
	fmt.Printf("Psp 0x%02X = {", p.Pasport)
	fmt.Printf("%c ", types.DataTypeSymbols[p.DataType()])
	fmt.Printf("%v ", p.Validation())
	fmt.Printf("%d}  ", p.ShiftStatus())
	fmt.Printf("raw [% X]  ", p.RawValue)

	switch p.DataType() {
	case types.TypeAnalog, types.TypeFloat:
		value := binary.LittleEndian.Uint32(p.RawValue[:])
		floatVal := math.Float32frombits(value)
		fmt.Printf("float32:(0x%08X) %.4f\n", value, floatVal)

	case types.TypeFixed:
		shift := p.Pasport & 0x07
		intVal := int32(binary.LittleEndian.Uint32(p.RawValue[:]))
		fixedVal := float32(intVal) / float32(math.Pow10(int(shift)))
		fmt.Printf("fixed:(shift=%d) %.2f\n", shift, fixedVal)

	case types.TypeInt:
		intVal := int32(binary.LittleEndian.Uint32(p.RawValue[:]))
		fmt.Printf("int32:(0x%08X) %d\n", uint32(intVal), intVal)

	case types.TypeDiscrete:
		fmt.Printf("discrete:[% X]\n", p.RawValue)

	default:
		fmt.Println("  Неизвестный тип данных")
	}
}

func DisplayDataBlock(data []byte) error {
	header, err := parser.ParseHeader(data)
	if err != nil {
		return fmt.Errorf("ошибка парсинга заголовка: %v", err)
	}
	
	PrintHeader(header)

	params, err := parser.ParseParameters(data, header.NumParams)
	if err != nil {
		return fmt.Errorf("ошибка парсинга параметров: %v", err)
	}

	stroka := fmt.Sprintf("Параметры (%d)", len(params))
	fmt.Println(PrintHeaderLine(stroka, '=', LINELEN))
	for i, param := range params {
		PrintParameter(param, i)
	}
	
	fmt.Println(strings.Repeat("-", LINELEN))
	return nil
}
```

```GO
// main.go
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
```

```Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o receiver ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/receiver .
EXPOSE 8080
CMD ["./receiver"]
```

## 5. Go модули
* **go.mod gendata-project**
```GO
module gendata-project

go 1.21

require ()
```

* **go.mod gendata-project/generator-service**
```GO
module gendata-project/generator-service

go 1.21

replace gendata-project/shared => ../shared

require (
    gendata-project/shared v0.0.0-00010101000000-000000000000
)   
```


* **go.mod для gendata-project/receiver-service**
```GO 
module gendata-project/receiver-service

go 1.21

replace gendata-project/shared => ../shared

require (
    gendata-project/shared v0.0.0-00010101000000-000000000000
)
```

* **go.mod для module gendata-project/shared**
```GO
module gendata-project/shared

go 1.21
```

## 6. Docker Compose
* **docker-compose.yml**
```yml
version: '3.8'

services:
  receiver-service:
    build:
      context: ./receiver-service
      dockerfile: Dockerfile
    container_name: con-receiver
    ports:
      - "8080:8080"
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

## 7. Ansible для развертывания
* **inventory.yml**
```yml 
all:
  hosts:
    srv102:
      ansible_host: 192.168.88.102
      ansible_user: starmark
      ansible_ssh_private_key_file: ~/.ssh/id_rsa
```

* **playbook.yml**
```yml
---
- name: Deploy GenData Services
  hosts: srv102
  become: yes
  vars:
    project_dir: /opt/gendata-project
    
  tasks:
    - name: Install Docker
      apt:
        name: 
          - docker.io
          - docker-compose
        state: present
        update_cache: yes

    - name: Start and enable Docker
      systemd:
        name: docker
        state: started
        enabled: yes

    - name: Add user to docker group
      user:
        name: "{{ ansible_user }}"
        groups: docker
        append: yes

    - name: Create project directory
      file:
        path: "{{ project_dir }}"
        state: directory
        owner: "{{ ansible_user }}"
        group: "{{ ansible_user }}"
        mode: '0755'

    - name: Copy project files
      synchronize:
        src: ../../
        dest: "{{ project_dir }}"
        delete: yes
        rsync_opts:
          - "--exclude=deploy"
          - "--exclude=.git"

    - name: Set ownership of project files
      file:
        path: "{{ project_dir }}"
        owner: "{{ ansible_user }}"
        group: "{{ ansible_user }}"
        recurse: yes

    - name: Stop existing containers
      docker_compose:
        project_src: "{{ project_dir }}"
        state: absent
      ignore_errors: yes

    - name: Build and start services
      docker_compose:
        project_src: "{{ project_dir }}"
        build: yes
        state: present

    - name: Wait for services to be ready
      wait_for:
        host: "{{ ansible_host }}"
        port: 8080
        delay: 10
        timeout: 60

    - name: Show running containers
      command: docker ps
      register: docker_ps

    - name: Display running containers
      debug:
        var: docker_ps.stdout_lines
```

* **deploy.sh**
```bash
#!/bin/bash

# deploy.sh
# Скрипт для развертывания сервисов GenData

echo "Развертывание сервисов GenData на srv102..."

# Проверка доступности хоста
if ! ping -c 1 192.168.88.102 &> /dev/null; then
    echo "Ошибка: Хост srv102 недоступен"
    exit 1
fi

# Запуск Ansible playbook
ansible-playbook -i inventory.yml playbook.yml -v

if [ $? -eq 0 ]; then
    echo "Развертывание завершено успешно!"
    echo "Сервисы доступны по адресу: http://192.168.88.102:8080"
    echo ""
    echo "Для просмотра логов используйте:"
    echo "ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose logs -f'"
else
    echo "Ошибка при развертывании!"
    exit 1
fi
```

## 8. Makefile для удобства
```Makefile
.PHONY: build deploy logs stop clean

# Сборка локально
build:
	docker-compose build

# Развертывание на удаленном сервере
deploy:
	cd deploy/ansible && chmod +x deploy.sh && ./deploy.sh

# Просмотр логов на удаленном сервере
logs:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose logs -f'

# Остановка сервисов на удаленном сервере
stop:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose down'

# Очистка на удаленном сервере
clean:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose down --rmi all --volumes'

# Проверка статуса сервисов
status:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose ps'

# Локальный запуск для тестирования
local-run:
	docker-compose up --build

# Локальная остановка
local-stop:
	docker-compose down
```

## 9. Инструкции по запуску
### 1. Клонируйте проект и перейдите в директорию
cd gendata-project

### 2. Убедитесь, что у вас есть SSH доступ к srv102
ssh starmark@192.168.88.102 "echo 'SSH connection OK'"

### 3. Установите Ansible (если не установлен)
pip install ansible

### 4. Разверните сервисы на удаленном сервере
make deploy

### 5. Просмотр логов в реальном времени
make logs

### 6. Проверка статуса сервисов
make status

## 10. Исправление в generator-service/cmd/main.go
```GO
// main.go
package main

import (
	"encoding/binary"
	"log"
	"net"
	"time"

	"gendata-project/generator-service/internal/generator"
	"gendata-project/shared/protocol"
	"gendata-project/shared/types"
)

func main() {
	log.Println("Запуск генератора данных...")

	// Ожидание готовности получателя
	time.Sleep(10 * time.Second)

	for {
		// Подключение к получателю
		conn, err := net.Dial("tcp", "receiver-service:8080")
		if err != nil {
			log.Printf("Ошибка подключения к получателю: %v. Повторная попытка через 5 секунд...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Подключение к получателю установлено")
		
		// Генерация данных каждые 5 секунд
		ticker := time.NewTicker(5 * time.Second)
		
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
				for i, dataType := range dataTypes {
					gid := generator.GIDCreater(1, 2, 1) // Север, NPP 2, Блок 1
					
					header := generator.HeadBlockCreater(
						gid,
						dataType,
						1,                   // ListNum
						1,                   // ListVer
						int16(i),            // Index
						5,                   // NumParams
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
		time.Sleep(5 * time.Second)
	}
}


## 11. Дополнительные файлы конфигурации
* **ansible.cfg**
```cfg
[defaults]
host_key_checking = False
retry_files_enabled = False
stdout_callback = yaml
inventory = inventory.yml

[ssh_connection]
ssh_args = -o ControlMaster=auto -o ControlPersist=60s
pipelining = True
```


* **build.sh**
```bash
#!/bin/bash
# build.sh

echo "Сборка проекта GenData..."

# Проверка структуры проекта
if [ ! -f "docker-compose.yml" ]; then
    echo "Ошибка: docker-compose.yml не найден"
    exit 1
fi

# Сборка образов
echo "Сборка Docker образов..."
docker-compose build --no-cache

if [ $? -eq 0 ]; then
    echo "Сборка завершена успешно!"
    docker images | grep gendata
else
    echo "Ошибка при сборке!"
    exit 1
fi
```

* **test-local.sh**
```bash 
#!/bin/bash
# test-local.sh

echo "Локальное тестирование сервисов..."

# Запуск сервисов
docker-compose up -d

# Ожидание запуска
sleep 15

# Проверка статуса
echo "Статус контейнеров:"
docker-compose ps

# Проверка логов
echo "Логи receiver-service:"
docker-compose logs receiver-service | tail -10

echo "Логи generator-service:"
docker-compose logs generator-service | tail -10

# Проверка сетевого соединения
echo "Проверка порта 8080:"
netstat -tlnp | grep 8080

echo "Для остановки выполните: docker-compose down"
```

## 12. Документация
### GenData Project
    Проект для генерации и обработки пакетов данных в формате СПД.

### Архитектура

Проект состоит из двух микросервисов:

1. **Generator Service** - генерирует данные различных типов и отправляет их получателю
2. **Receiver Service** - принимает данные и отображает их в консоли

## Структура проекта

```
gendata-project/
├── shared/               # Общие компоненты
│   ├── types/            # Типы данных
│   ├── sensorvalue/      # Работа с сенсорными данными
│   └── protocol/         # Протокол обмена
├── generator-service/    # Сервис генерации данных
├── receiver-service/     # Сервис приема данных
├── deploy/               # Ansible для развертывания
└── scripts/              # Вспомогательные скрипты
```

## Быстрый старт

### Локальное тестирование

```bash
# Сборка и запуск
make local-run

# Просмотр логов
docker-compose logs -f

# Остановка
make local-stop
```

### Развертывание на удаленном сервере

```bash
# Развертывание на srv102
make deploy

# Просмотр логов
make logs

# Проверка статуса
make status

# Остановка сервисов
make stop
```

### Типы данных

    Поддерживаются следующие типы данных:

- **Analog (A)** - аналоговые значения (float32)
- **Float (F)** - вещественные числа (float32)
- **Fixed (X)** - числа с фиксированной точкой
- **Integer (I)** - целые числа (int32)
- **Discrete (D)** - дискретные значения

### Мониторинг

    Для мониторинга работы сервисов используйте:

```bash
# Просмотр логов в реальном времени
ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose logs -f'

# Статус контейнеров
ssh starmark@192.168.88.102 'docker ps'

# Использование ресурсов
ssh starmark@192.168.88.102 'docker stats'
```

### Порты

* **8080** - Receiver Service (проброшен на хост)
  * Доступ к сервису: http://192.168.88.102:8080


## 13. Финальная структура проекта

```bash:scripts/create-structure.sh
#!/bin/bash

echo "Создание структуры проекта GenData..."

# Создание основных директорий
mkdir -p gendata-project/{shared/{types,sensorvalue,protocol},generator-service/{cmd,internal/generator},receiver-service/{cmd,internal/{parser,display}},deploy/ansible,scripts}

# Создание файлов go.mod
cat > gendata-project/go.mod << 'EOF'
module gendata-project

go 1.21
EOF

cat > gendata-project/shared/go.mod << 'EOF'
module gendata-project/shared

go 1.21
EOF

cat > gendata-project/generator-service/go.mod << 'EOF'
module gendata-project/generator-service

go 1.21

replace gendata-project/shared => ../shared

require gendata-project/shared v0.0.0-00010101000000-000000000000
EOF

cat > gendata-project/receiver-service/go.mod << 'EOF'
module gendata-project/receiver-service

go 1.21

replace gendata-project/shared => ../shared

require gendata-project/shared v0.0.0-00010101000000-000000000000
EOF

echo "Структура проекта создана!"
echo "Теперь скопируйте все файлы кода в соответствующие директории."
```

## 14. Команды для развертывания

```bash
# 1. Создание структуры проекта
chmod +x scripts/create-structure.sh
./scripts/create-structure.sh
```

```bash
# 2. Копирование всех файлов кода в соответствующие директории
```

```bash
# 3. Сборка и тестирование локально
chmod +x scripts/test-local.sh
./scripts/test-local.sh
```

```bash
# 4. Развертывание на удаленном сервере
chmod +x deploy/ansible/deploy.sh
make deploy
```

```bash
# 5. Мониторинг
make logs
```

## README.md
### Полная инфраструктура для развертывания двух микросервисов:
* **Generator Service работает в контейнере con-gendata и генерирует данные каждые 5 секунд**
* **Receiver Service работает в контейнере con-receiver и принимает данные на порту 8080**
* **Все развертывается на удаленной VM srv102 через Ansible**
* **Логи отображаются в консоли получателя с красивым форматированием**
* **Порт 8080 проброшен на хост, так что вы можете видеть результат с вашей машины**