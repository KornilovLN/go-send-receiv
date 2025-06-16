# Выбор лицензии для проекта

## 1. MIT License (Рекомендуется)
* **Плюсы:**
  * Очень простая и понятная
  * Разрешает коммерческое использование
  * Минимальные ограничения
  * Широко используется в Go-сообществе

    MIT License

    Copyright (c) 2024 [Ваше имя/организация]

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.

## 2. Apache License 2.0
    Если нужна защита от патентных претензий:

    Apache License
    Version 2.0, January 2004
    http://www.apache.org/licenses/

    [Полный текст лицензии Apache 2.0]

## 3. GPL v3
    Если надо, чтобы все производные работы тоже были открытыми:

    GNU GENERAL PUBLIC LICENSE
    Version 3, 29 June 2007

    [Полный текст GPL v3]



# Создание новой ветки с улучшениями
## Шаг 1: Создание и переключение на новую ветку
* **Убедитесь, что находитесь в главной ветке**
```bash
git checkout main
```
* **Получите последние изменения**
```bash 
git pull origin main
```
* **Создайте новую ветку для веб-интерфейса**
```bash 
git checkout -b feature/web-interface
```
* **Или создайте ветку для общих улучшений**
```bash 
git checkout -b feature/improvements-v2
```

## Шаг 2: Внесение изменений
* **Добавьте лицензию**
```bash
echo "MIT License..." > LICENSE
```
* **Обновите README.md**
```
# GenData Project

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-blue.svg)](https://docker.com)

Система генерации и обработки пакетов данных в формате СПД с веб-интерфейсом мониторинга.

## 🚀 Возможности

- **Генератор данных**: Создает пакеты различных типов данных (Analog, Float, Fixed, Integer, Discrete)
- **Приемник данных**: Обрабатывает и отображает полученные данные
- **Веб-интерфейс**: Мониторинг в реальном времени через браузер
- **Docker**: Контейнеризация для легкого развертывания
- **Ansible**: Автоматизированное развертывание на удаленных серверах

## 📋 Требования

- Go 1.21+
- Docker & Docker Compose
- Ansible (для развертывания)

## 🏗️ Архитектура
```

┌─────────────────┐ TCP:8080 ┌─────────────────┐ │ Generator │ ──────────────► │ Receiver │ │ Service │ │ Service │ └─────────────────┘ └─────────────────┘ │ HTTP:8081 ▼ ┌─────────────────┐ │ Web Interface │ │ (Monitoring) │ └─────────────────┘

```
## 🚀 Быстрый старт

### Локальное тестирование

```bash
# Клонирование репозитория
git clone <repository-url>
cd gendata-project

# Запуск сервисов
make local-run

# Открыть веб-интерфейс
open http://localhost:8081
```

## Развертывание на сервере
```
# Настройка inventory.yml с вашими серверами
# Развертывание
make deploy

# Мониторинг
make logs
```

```
📊 Веб-интерфейс
После запуска веб-интерфейс доступен по адресу:

Локально: http://localhost:8081
На сервере: http://your-server:8081
Интерфейс показывает:

📈 Статистику сообщений в реальном времени
📝 Последние полученные сообщения
🔄 Автообновление каждые 2 секунды
🛠️ Команды Make
make local-run      # Локальный запуск
make deploy         # Развертывание на сервере
make logs          # Просмотр логов
make status        # Статус сервисов
make stop          # Остановка сервисов
make clean         # Очистка
```

```
📁 Структура проекта
gendata-project/
├── shared/                 # Общие компоненты
│   ├── types/             # Типы данных СПД
│   ├── sensorvalue/       # Обработка сенсорных данных
│   └── protocol/          # Протокол обмена
├── generator-service/     # Сервис генерации
├── receiver-service/      # Сервис приема с веб-интерфейсом
├── deploy/               # Ansible playbooks
└── scripts/              # Вспомогательные скрипты
```

```
🔧 Типы данных
Тип	Символ	Описание
Analog	A	Аналоговые значения (float32)
Float	F	Вещественные числа
Fixed	X	Числа с фиксированной точкой
Integer	I	Целые числа (int32)
Discrete	D	Дискретные значения
📝 Лицензия
Этот проект лицензирован под MIT License - см. файл LICENSE.
```
```
🤝 Участие в разработке
Форкните проект
Создайте ветку для новой функции (git checkout -b feature/amazing-feature)
Зафиксируйте изменения (git commit -m 'Add amazing feature')
Отправьте в ветку (git push origin feature/amazing-feature)
Откройте Pull Request
📞 Поддержка
Если у вас есть вопросы или предложения, создайте Issue в репозитории. EOF
```

### Шаг 3: Добавление улучшений
```bash
# Добавьте веб-интерфейс (код из предыдущего ответа)
# Обновите docker-compose.yml
# Добавьте новые функции
```

### Шаг 4: Коммит изменений
```
# Добавьте все изменения
git add .

# Создайте коммит с описанием
git commit -m "feat: добавлен веб-интерфейс мониторинга

- Добавлен HTTP сервер на порту 8081
- Веб-интерфейс с реальным временем мониторинга
- API endpoints для статистики
- Обновлен docker-compose.yml для нового порта
- Добавлена MIT лицензия
- Обновлена документация"
```

### Шаг 5: Отправка в репозиторий
```
# Отправьте ветку в удаленный репозиторий
git push origin feature/web-interface

# Или если это первая отправка ветки
git push -u origin feature/web-interface
```

### Шаг 6: Создание Pull Request
```
После отправки ветки создайте Pull Request через веб-интерфейс Git-платформы (GitHub/GitLab):

Перейдите в репозиторий
Нажмите "New Pull Request" или "Create Merge Request"
Выберите ветку feature/web-interface → main
Добавьте описание:
## 🚀 Добавлен веб-интерфейс мониторинга

### Что добавлено:
- ✅ HTTP сервер для веб-интерфейса (порт 8081)
- ✅ Реальное время мониторинга данных
- ✅ API endpoints для статистики
- ✅ Автообновление интерфейса каждые 2 секунды
- ✅ MIT лицензия
- ✅ Обновленная документация

### Как тестировать:
1. `make local-run`
2. Открыть http://localhost:8081
3. Проверить отображение статистики

### Скриншоты:
[Добавьте скриншоты веб-интерфейса]

Closes #[номер issue, если есть]
```

### Шаг 7: Слияние и очистка

После одобрения и слияния PR:

```bash
# Переключитесь на main
git checkout main

# Получите обновления
git pull origin main

# Удалите локальную ветку
git branch -d feature/web-interface

# Удалите удаленную ветку (опционально)
git push origin --delete feature/web-interface
```

## Дополнительные файлы для проекта

### .gitignore

```gitignore:.gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
gendata-project/generator-service/generator
gendata-project/receiver-service/receiver

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log