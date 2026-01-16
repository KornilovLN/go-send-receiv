# Реорганизация templates
## 1. Переместить папку templates
```bash
mv gendata-project/receiver-service/templates gendata-project/receiver-service/internal/web/
```

## 2. Создать структуру шаблонов
```html        
<!-- базовый layout и отдельные шаблоны: base.html -->

<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="/static/css/main.css">
    {{block "head" .}}{{end}}
</head>
<body>
    {{template "content" .}}
    {{block "scripts" .}}{{end}}
</body>
</html>
```

```html
<!-- index.html -->

{{define "content"}}
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
{{end}}


{{define "scripts"}}
<script src="/static/js/app.js"></script>
{{end}}
```

```html
<!-- table.html-->

{{define "head"}}
<style>
body {
    font-family: Arial, sans-serif;
    margin: 20px;
    background-color: #f5f5f5;
}
h1 {
    color: #333;
    text-align: center;
}
.controls {
    text-align: center;
    margin: 20px 0;
}
.controls button {
    padding: 10px 20px;
    margin: 0 10px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
}
.auto-scroll {
    background-color: #4CAF50;
    color: white;
}
.auto-scroll.disabled {
    background-color: #f44336;
}
.table-container {
    max-height: 1000px;
    overflow-y: auto;
    border: 1px solid #ddd;
}
table {
    width: 100%;
    border-collapse: collapse;
    background-color: white;
}
th, td {
    padding: 8px 12px;
    text-align: left;
    border-bottom: 1px solid #ddd;
}
th {
    background-color: #4CAF50;
    color: white;
    font-weight: bold;
    position: sticky;
    top: 0;
    z-index: 10;
}
tr:nth-child(even) {
    background-color: #f9f9f9;
}
tr:hover {
    background-color: #f5f5f5;
}
.loading {
    text-align: center;
    padding: 20px;
    color: #666;
}
.new-row {
    background-color: #e8f5e8 !important;
    animation: highlight 2s ease-out;
}
@keyframes highlight {
    0% { background-color: #90EE90; }
    100% { background-color: #e8f5e8; }
}
.status {
    text-align: center;
    padding: 10px;
    font-weight: bold;
}
.status.online {
    color: #4CAF50;
}
.status.offline {
    color: #f44336;
}
</style>
{{end}}

{{define "content"}}
<h1>Данные датчиков</h1>
<div class="status" id="connection-status">🟢 Подключено</div>
<div class="controls">
    <button id="auto-scroll-btn" class="auto-scroll" onclick="toggleAutoScroll()">🔄 Авто-прокрутка: ВКЛ</button>
    <button onclick="clearTable()">🗑️ Очистить</button>
    <button onclick="loadTableData()">🔄 Обновить</button>
    <span>Всего записей: <span id="total-count">0</span></span>
</div>
<div class="table-container" id="table-container">
    <table>
        <thead>
            <tr>
                <th>#</th>
                <th>Паспорт</th>
                <th>Тип</th>
                <th>Код типа</th>
                <th>Валидность</th>
                <th>Смещение</th>
                <th>Значение</th>
                <th>Байты</th>
                <th>Время</th>
            </tr>
        </thead>
        <tbody id="sensor-data">
            <tr>
                <td colspan="9" class="loading">Загрузка данных...</td>
            </tr>
        </tbody>
    </table>
</div>
{{end}}

{{define "scripts"}}
<script src="/static/js/utils.js"></script>
<script>
// Весь JavaScript код для таблицы
let autoScroll = true;
let lastMessageCount = 0;
let totalRowCount = 0;
let lastUpdateTime = 0;

function toggleAutoScroll() {
    autoScroll = !autoScroll;
    const btn = document.getElementById('auto-scroll-btn');
    if (autoScroll) {
        btn.textContent = '🔄 Авто-прокрутка: ВКЛ';
        btn.className = 'auto-scroll';
    } else {
        btn.textContent = '⏸️ Авто-прокрутка: ВЫКЛ';
        btn.className = 'auto-scroll disabled';
    }
}

function clearTable() {
    document.getElementById('sensor-data').innerHTML =
        '<tr><td colspan="9" class="loading">Таблица очищена</td></tr>';
    totalRowCount = 0;
    updateTotalCount();
}

function updateTotalCount() {
    document.getElementById('total-count').textContent = totalRowCount;
}

function scrollToBottom() {
    if (autoScroll) {
        const container = document.getElementById('table-container');
        container.scrollTop = container.scrollHeight;
    }
}

function formatTime(dateStr) {
    return new Date(dateStr).toLocaleTimeString('ru-RU');
}

function renderTableData(messages) {
    const tbody = document.getElementById('sensor-data');
    const statusEl = document.getElementById('connection-status');
    
    if (!messages || messages.length === 0) {
        tbody.innerHTML = '<tr><td colspan="9" class="loading">Нет данных</td></tr>';
        statusEl.textContent = '🔴 Нет данных';
        statusEl.className = 'status offline';
        return;
    }

    statusEl.textContent = '🟢 Подключено';
    statusEl.className = 'status online';

    let html = '';
    let paramIndex = totalRowCount + 1;
    let newRowsAdded = 0;

    const currentMessageCount = messages.length;
    const hasNewData = currentMessageCount > lastMessageCount ||
        (messages[0] && new Date(messages[0].received_at).getTime() > lastUpdateTime);

    messages.forEach(msg => {
        const messageTime = new Date(msg.received_at).getTime();
        const isNewMessage = messageTime > lastUpdateTime;
        
        msg.parameters.forEach(param => {
            const validationSymbol = param.validation.toLowerCase() === 'valid' ? '✔' : '✗';
            const typeCode = '0x' + (param.passport.replace('0x', '')).substring(0, 1) + '0';
            const rowClass = isNewMessage ? 'new-row' : '';
            
            html += '<tr class="' + rowClass + '">';
            html += '<td>' + paramIndex + '</td>';
            html += '<td>' + param.passport + '</td>';
            html += '<td>' + param.data_type + '</td>';
            html += '<td>' + typeCode + '</td>';
            html += '<td>' + validationSymbol + '</td>';
            html += '<td>' + param.shift_status + '</td>';
            html += '<td title="' + param.parsed_value + '">' + param.parsed_value + '</td>';
            html += '<td title="' + param.raw_value + '">' + param.raw_value + '</td>';
            html += '<td>' + formatTime(msg.received_at) + '</td>';
            html += '</tr>';
            
            if (isNewMessage) {
                newRowsAdded++;
            }
            paramIndex++;
        });
    });

    if (hasNewData && tbody.innerHTML.indexOf('Загрузка') === -1 && tbody.innerHTML.indexOf('Нет данных') === -1) {
        tbody.innerHTML += html;
        totalRowCount += newRowsAdded;
    } else {
        tbody.innerHTML = html;
        totalRowCount = paramIndex - 1;
    }

    updateTotalCount();
    lastMessageCount = currentMessageCount;
    lastUpdateTime = Math.max(lastUpdateTime, ...messages.map(m => new Date(m.received_at).getTime()));

    if (newRowsAdded > 0) {
        setTimeout(scrollToBottom, 100);
    }
}

async function loadTableData() {
    try {
        const response = await fetch('/api/messages');
        if (response.ok) {
            const messages = await response.json();
            renderTableData(messages);
        } else {
            throw new Error('Ошибка сети: ' + response.status);
        }
    } catch (error) {
        console.error('Ошибка загрузки данных:', error);
        document.getElementById('connection-status').textContent = '🔴 Ошибка подключения';
        document.getElementById('connection-status').className = 'status offline';
    }
}

document.addEventListener('DOMContentLoaded', loadTableData);
setInterval(loadTableData, 2000);
</script>
{{end}}
```

## 3. Вынести JavaScript в отдельные файлы
```js
// app.js

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

    const typesContainer = document.getElementById('message-types');
    if (stats.messages_per_type && Object.keys(stats.messages_per_type).length > 0) {
        let typesHtml = '';
        for (const [type, count] of Object.entries(stats.messages_per_type)) {
            typesHtml += `
                <div class="message-type-item">
                    <span class="message-type-name">${type}</span>
                    <span class="message-type-count">${count}</span>
                </div>
            `;
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
        html += `
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
                    <table class="parameters-table">
                        <thead>
                            <tr>
                                <th>#</th>
                                <th>Тип</th>
                                <th>Статус</th>
                                <th>Значение</th>
                                <th>Raw данные</th>
                                <th>Паспорт</th>
                            </tr>
                        </thead>
                        <tbody>
        `;
        
        msg.parameters.forEach(param => {
            const validationClass = param.validation.toLowerCase() === 'valid' ? 'valid' : 'invalid';
            html += `
                <tr>
                    <td class="parameter-index">[${param.index}]</td>
                    <td><span class="parameter-type">${param.data_type}</span></td>
                    <td><span class="parameter-validation ${validationClass}">${param.validation}</span></td>
                    <td class="parameter-value"><strong>${param.parsed_value}</strong></td>
                    <td class="parameter-raw">${param.raw_value}</td>
                    <td class="parameter-passport">${param.passport}</td>
                </tr>
            `;
        });
        
        html += `
                        </tbody>
                    </table>
                </div>
            </div>
        `;
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
```

## 4. Обновить templates.go
```GO
// templates.go

package web

import (
    "bytes"
    "fmt"
    "html/template"
    "os"
    "path/filepath"
)

var templates *template.Template

// LoadTemplates загружает все HTML шаблоны из папки templates
func LoadTemplates() error {
    templatesDir := "internal/web/templates"
    
    tmpl := template.New("")
    
    // Загружаем все HTML файлы рекурсивно
    err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if filepath.Ext(path) == ".html" {
            // Получаем относительный путь для имени шаблона
            relPath, _ := filepath.Rel(templatesDir, path)
            
            // Читаем содержимое файла
            content, err := os.ReadFile(path)
            if err != nil {
                return err
            }
            
            // Парсим шаблон
            _, err = tmpl.New(relPath).Parse(string(content))
            if err != nil {
                return err
            }
        }
        return nil
    })
    
    if err != nil {
        return err
    }
    
    templates = tmpl
    return nil
}

// GetTemplates возвращает загруженные шаблоны
func GetTemplates() *template.Template {
    return templates
}

// RenderTemplate рендерит указанный шаблон с данными
func RenderTemplate(templateName string, data interface{}) (string, error) {
    if templates == nil {
        return "", fmt.Errorf("templates not loaded")
    }
    
    var buf bytes.Buffer
    err := templates.ExecuteTemplate(&buf, templateName, data)
    if err != nil {
        return "", err
    }
    
    return buf.String(), nil
}
```

## 5. Обновить handlers.go
```GO
// handlers.go

package web

import (
    "net/http"
)

// IndexHandler обрабатывает главную страницу
func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
    data := struct {
        Title string
    }{
        Title: "GenData Monitor",
    }
    
    // Используем базовый layout с шаблоном index
    err := templates.ExecuteTemplate(w, "layouts/base.html", data)
    if err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
        return
    }
}

// TableHandler обрабатывает страницу с таблицей
func (s *Server) TableHandler(w http.ResponseWriter, r *http.Request) {
    data := struct {
        Title string
    }{
        Title: "Данные датчиков",
    }
    
    // Используем базовый layout с шаблоном table
    err := templates.ExecuteTemplate(w, "layouts/base.html", data)
    if err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
        return
    }
}
```

## 6. Обновить main.go для инициализации шаблонов
    Основные изменения:

    ✅ Добавлен вызов web.LoadTemplates() в начале main()
    ✅ Добавлено логирование успешной загрузки шаблонов
    ✅ Добавлены недостающие импорты в templates.go
    Теперь при запуске приложения шаблоны будут загружены до создания веб-сервера!
```GO
// main.go

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
    
    // Загружаем шаблоны при запуске
    err := web.LoadTemplates()
    if err != nil {
        log.Fatal("Failed to load templates:", err)
    }
    log.Println("Шаблоны успешно загружены")
    
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
```

## 7. Создать отдельные шаблоны для каждой страницы
* **index.html**
```html
{{template "layouts/base.html" .}}
{{define "content"}}
<!-- Содержимое главной страницы -->
{{template "index.html" .}}
{{end}}
```
* **table.html**
```html
{{template "layouts/base.html" .}}
{{define "content"}}
<!-- Содержимое страницы таблицы -->
{{template "table.html" .}}
{{end}}
```

### Архитектура с разделением HTML, CSS и JavaScript по отдельным файлам.