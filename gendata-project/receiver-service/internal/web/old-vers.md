```GO
// old-vesion-templates.go
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
                    <p></p>
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
							<tbody>` + "`" + `;

				msg.parameters.forEach(param => {
					const validationClass = param.validation.toLowerCase() === 'valid' ? 'valid' : 'invalid';
					html += ` + "`" + `<tr>
						<td class="parameter-index">[${param.index}]</td>
						<td><span class="parameter-type">${param.data_type}</span></td>
						<td><span class="parameter-validation ${validationClass}">${param.validation}</span></td>
						<td class="parameter-value"><strong>${param.parsed_value}</strong></td>
						<td class="parameter-raw">${param.raw_value}</td>
						<td class="parameter-passport">${param.passport}</td>
					</tr>` + "`" + `;
				});

				html += ` + "`" + `			</tbody>
						</table>
					</div>
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

// GetTableTemplate возвращает HTML шаблон для табличного отображения данных
func GetTableTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>Данные датчиков</title>
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
            max-height: 1000px; /* Ограничиваем высоту таблицы */
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
</head>
<body>
    <h1>Данные датчиков</h1>
    
    <div class="status" id="connection-status">🟢 Подключено</div>
    
    <div class="controls">
        <button id="auto-scroll-btn" class="auto-scroll" onclick="toggleAutoScroll()">
            🔄 Авто-прокрутка: ВКЛ
        </button>
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
                <tr><td colspan="9" class="loading">Загрузка данных...</td></tr>
            </tbody>
        </table>
    </div>

    <script>
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

            // Обновляем статус подключения
            statusEl.textContent = '🟢 Подключено';
            statusEl.className = 'status online';

            let html = '';
            let paramIndex = totalRowCount + 1;
            let newRowsAdded = 0;
            
            // Проверяем, есть ли новые сообщения
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

            // Если режим добавления новых данных
            if (hasNewData && tbody.innerHTML.indexOf('Загрузка') === -1 && tbody.innerHTML.indexOf('Нет данных') === -1) {
                tbody.innerHTML += html; // Добавляем к существующим данным
                totalRowCount += newRowsAdded;
            } else {
                tbody.innerHTML = html; // Полная замена
                totalRowCount = paramIndex - 1;
            }

            updateTotalCount();
            lastMessageCount = currentMessageCount;
            lastUpdateTime = Math.max(lastUpdateTime, ...messages.map(m => new Date(m.received_at).getTime()));
            
            // Автопрокрутка к новым данным
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

        // Загружаем данные при загрузке страницы
        document.addEventListener('DOMContentLoaded', loadTableData);

        // Автоматическое обновление каждые 2 секунды
        setInterval(loadTableData, 2000);
    </script>
</body>
</html>`
}
```