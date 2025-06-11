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
