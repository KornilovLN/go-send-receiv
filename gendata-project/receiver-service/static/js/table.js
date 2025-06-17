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