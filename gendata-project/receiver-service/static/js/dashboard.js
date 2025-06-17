
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
					typesHtml += `<div class="message-type-item">
						<span class="message-type-name">${type}</span>
						<span class="message-type-count">${count}</span>
					</div>`;
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
				html += `<div class="message">
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
							<tbody>`;

				msg.parameters.forEach(param => {
					const validationClass = param.validation.toLowerCase() === 'valid' ? 'valid' : 'invalid';
					html += `<tr>
						<td class="parameter-index">[${param.index}]</td>
						<td><span class="parameter-type">${param.data_type}</span></td>
						<td><span class="parameter-validation ${validationClass}">${param.validation}</span></td>
						<td class="parameter-value"><strong>${param.parsed_value}</strong></td>
						<td class="parameter-raw">${param.raw_value}</td>
						<td class="parameter-passport">${param.passport}</td>
					</tr>`;
				});

				html += `			</tbody>
						</table>
					</div>
				</div>`;
			});

			container.innerHTML = html;
		}

		async function loadData() {
			if (isLoading) return;
			isLoading = true;

			try {
				const statsResponse = await fetch('/api/stats');
				if (statsResponse.ok) {
					const stats = await statsResponse.json();
					updateStats(stats);
				}

				const messagesResponse = await fetch('/api/messages');
				if (messagesResponse.ok) {
					const messages = await messagesResponse.json();
					renderMessages(messages);
				}

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

		document.addEventListener('DOMContentLoaded', loadData);
		setInterval(loadData, 3000);
