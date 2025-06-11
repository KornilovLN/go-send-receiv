package web

// GetIndexTemplate возвращает HTML шаблон главной страницы
func GetIndexTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>GenData Monitor</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body>
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

    <script>
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

            // Обновляем типы сообщений
            const typesContainer = document.getElementById('message-types');
            if (stats.messages_per_type && Object.keys(stats.messages_per_type).length > 0) {
                let typesHtml = '';
                for (const [type, count] of Object.entries(stats.messages_per_type)) {
                    typesHtml += ` + "`" + `<div class="message-type-item">
                        <span class="message-type-name">${type}</span>
                        <span class="message-type-count">${count}</span>
                    </div>` + "`" + `;
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
                html += ` + "`" + `<div class="message">
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
                        <div class="parameters-header">🔧 Параметры (${msg.parameters.length})</div>` + "`" + `;

                msg.parameters.forEach(param => {
                    const validationClass = param.validation.toLowerCase() === 'valid' ? 'valid' : 'invalid';
                    html += ` + "`" + `<div class="parameter">
                        <div class="parameter-grid">
                            <div class="parameter-index">[${param.index}]</div>
                            <div class="parameter-type">${param.data_type}</div>
                            <div class="parameter-validation ${validationClass}">${param.validation}</div>
                            <div class="parameter-value">
                                <strong>Значение:</strong> ${param.parsed_value}<br>
                                <strong>Raw:</strong> ${param.raw_value}<br>
                                <strong>Passport:</strong> ${param.passport}
                            </div>
                        </div>
                    </div>` + "`" + `;
                });

                html += ` + "`" + `</div>
                </div>` + "`" + `;
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
    </script>
</body>
</html>`
}
