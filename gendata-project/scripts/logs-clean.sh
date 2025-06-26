#!/bin/bash

echo "🧹 Очистка логов Docker..."

# Проверяем права доступа
if [ "$EUID" -ne 0 ]; then
    echo "⚠️ Для очистки системных логов Docker нужны права root"
    echo "Запустите: sudo ./scripts/clean-logs.sh"
fi

# Показать размер логов перед очисткой
echo "📊 Размер логов до очистки:"
if [ -d "/var/lib/docker/containers" ]; then
    sudo find /var/lib/docker/containers -name "*-json.log" -exec du -sh {} \; 2>/dev/null | head -10
else
    echo "Директория /var/lib/docker/containers не найдена"
fi

# Показать логи наших контейнеров
echo ""
echo "📋 Информация о наших контейнерах:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Size}}"

echo ""
echo "📊 Размер логов наших контейнеров:"
for container in con-receiver con-gendata; do
    if docker ps -a --format "{{.Names}}" | grep -q "^${container}$"; then
        log_path=$(docker inspect --format='{{.LogPath}}' $container 2>/dev/null)
        if [ -f "$log_path" ]; then
            size=$(sudo du -sh "$log_path" 2>/dev/null | cut -f1)
            echo "$container: $size ($log_path)"
        else
            echo "$container: лог-файл не найден"
        fi
    else
        echo "$container: контейнер не найден"
    fi
done

echo ""
echo "🗑️ Очистка логов контейнеров..."

# Метод 1: Очистка через Docker API (безопасно)
echo "Метод 1: Очистка через Docker API..."
for container in con-receiver con-gendata; do
    if docker ps -a --format "{{.Names}}" | grep -q "^${container}$"; then
        echo "Очистка логов $container..."
        # Получаем путь к лог-файлу
        log_path=$(docker inspect --format='{{.LogPath}}' $container 2>/dev/null)
        if [ -f "$log_path" ]; then
            # Останавливаем контейнер
            docker stop $container 2>/dev/null
            # Очищаем лог-файл
            sudo truncate -s 0 "$log_path" 2>/dev/null && echo "✅ $container: логи очищены" || echo "❌ $container: ошибка очистки"
            # Запускаем контейнер обратно
            docker start $container 2>/dev/null
        fi
    fi
done

# Метод 2: Очистка всех лог-файлов Docker (если нужно)
echo ""
echo "Метод 2: Массовая очистка всех логов Docker..."
if [ -d "/var/lib/docker/containers" ]; then
    sudo find /var/lib/docker/containers -name "*-json.log" -exec truncate -s 0 {} \; 2>/dev/null
    echo "✅ Все логи Docker очищены"
else
    echo "❌ Директория логов Docker не найдена"
fi

# Метод 3: Очистка через docker system prune
echo ""
echo "Метод 3: Очистка системы Docker..."
docker system prune -f --volumes 2>/dev/null && echo "✅ Система Docker очищена"

echo ""
echo "📊 Размер логов после очистки:"
for container in con-receiver con-gendata; do
    if docker ps -a --format "{{.Names}}" | grep -q "^${container}$"; then
        log_path=$(docker inspect --format='{{.LogPath}}' $container 2>/dev/null)
        if [ -f "$log_path" ]; then
            size=$(sudo du -sh "$log_path" 2>/dev/null | cut -f1)
            echo "$container: $size"
        fi
    fi
done

echo ""
echo "💾 Общая статистика Docker:"
docker system df

echo ""
echo "✅ Очистка логов завершена"
