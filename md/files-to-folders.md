# Структура создана правильно.
    Теперь нужно переместить файлы в правильные места. Вот что нужно сделать:

## 1. Перемещение файлов в правильные директории
* **Переместить файлы Ansible**
mv inventory.yml deploy/ansible/
mv playbook.yml deploy/ansible/
mv deploy.sh deploy/ansible/

## 2. Создать недостающие файлы
```yml
[defaults]
host_key_checking = False
retry_files_enabled = False
stdout_callback = yaml
inventory = inventory.yml

[ssh_connection]
ssh_args = -o ControlMaster=auto -o ControlPersist=60s
pipelining = True
```

* **ansible.cfg**
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
netstat -tlnp | grep 8080 || ss -tlnp | grep 8080

echo "Для остановки выполните: docker-compose down"
```

```bash
#!/bin/bash
# create-structure.sh

echo "Создание структуры проекта GenData..."

# Создание основных директорий
mkdir -p gendata-project/{shared/{types,sensorvalue,protocol},generator-service/{cmd,internal/generator},receiver-service/{cmd,internal/{parser,display}},deploy/ansible,scripts}

echo "Структура проекта создана!"
```

## 3. Исправить содержимое shared/sensorvalue/sensor.go
* **Файл должен называться sensorvalue.go, а не sensor.go:**
```bash
mv shared/sensorvalue/sensor.go shared/sensorvalue/sensorvalue.go
```

## 4. Обновить Makefile с правильными путями
```Makefile
.PHONY: build deploy logs stop clean status local-run local-stop

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

# Локальное тестирование
test:
	chmod +x scripts/test-local.sh && ./scripts/test-local.sh

# Сборка образов
build-images:
	chmod +x scripts/build.sh && ./scripts/build.sh
```

## 5. Обновить deploy/ansible/deploy.sh с правильными путями
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

# Проверка наличия файлов
if [ ! -f "inventory.yml" ]; then
    echo "Ошибка: inventory.yml не найден"
    exit 1
fi

if [ ! -f "playbook.yml" ]; then
    echo "Ошибка: playbook.yml не найден"
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
    echo ""
    echo "Или используйте: make logs"
else
    echo "Ошибка при развертывании!"
    exit 1
fi
```

## 6. Создать README.md
# GenData Project

    Проект для генерации и обработки пакетов данных в формате СПД.

## Архитектура

Проект состоит из двух микросервисов:

1. **Generator Service** - генерирует данные различных типов и отправляет их получателю
2. **Receiver Service** - принимает данные и отображает их в консоли

## Структура проекта
```
gendata-project/
├── shared/                 # Общие компоненты
│   ├── types/             # Типы данных
│   ├── sensorvalue/       # Работа с сенсорными данными
│   └── protocol/          # Протокол обмена
├── generator-service/     # Сервис генерации данных
├── receiver-service/      # Сервис приема данных
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

# Автоматическое тестирование
make test
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

## Типы данных

Поддерживаются следующие типы данных:

- **Analog (A)** - аналоговые значения (float32)
- **Float (F)** - вещественные числа (float32)
- **Fixed (X)** - числа с фиксированной точкой
- **Integer (I)** - целые числа (int32)
- **Discrete (D)** - дискретные значения

## Команды Make

- `make build` - сборка Docker образов
- `make deploy` - развертывание на удаленном сервере
- `make logs` - просмотр логов удаленного сервера
- `make status` - статус сервисов на удаленном сервере
- `make stop` - остановка сервисов на удаленном сервере
- `make clean` - полная очистка на удаленном сервере
- `make local-run` - локальный запуск
- `make local-stop` - локальная остановка
- `make test` - локальное тестирование

## Мониторинг

Для мониторинга работы сервисов используйте:

```bash
# Просмотр логов в реальном времени
make logs

# Статус контейнеров
ssh starmark@192.168.88.102 'docker ps'

# Использование ресурсов
ssh starmark@192.168.88.102 'docker stats'
```

## Порты

- **8080** - Receiver Service (проброшен на хост)

Доступ к сервису: http://192.168.88.102:8080

## Устранение неполадок

### Проблемы с подключением

```bash
# Проверка доступности хоста
ping 192.168.88.102

# Проверка SSH подключения
ssh starmark@192.168.88.102 "echo 'SSH OK'"

# Проверка Docker на удаленном хосте
ssh starmark@192.168.88.102 'docker --version'
```

### Проблемы с сервисами

```bash
# Перезапуск сервисов
ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose restart'

# Пересборка образов
ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker-compose build --no-cache'

# Полная переустановка
make clean
make deploy
```
```

## 7. Команды для финальной настройки

```bash
# Сделать скрипты исполняемыми
chmod +x scripts/*.sh
chmod +x deploy/ansible/deploy.sh

# Переместить файлы в правильные места
mv inventory.yml deploy/ansible/
mv playbook.yml deploy/ansible/
mv deploy.sh deploy/ansible/

# Переименовать файл сенсора
mv shared/sensorvalue/sensor.go shared/sensorvalue/sensorvalue.go

# Проверить структуру
tree
```

## 8. Тестирование

```bash
# Локальное тестирование
make test

# Если все работает локально, разворачиваем на удаленном сервере
make deploy

# Просмотр логов
make logs

Copy

Apply

README.md
Теперь структура проекта правильная и все файлы находятся в нужных местах. Вы можете запускать локальное тестирование и развертывание на удаленном сервере.