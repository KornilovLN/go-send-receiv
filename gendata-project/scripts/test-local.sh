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
