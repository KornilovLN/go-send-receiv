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
