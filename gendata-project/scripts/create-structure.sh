#!/bin/bash
# create-structure.sh

echo "Создание структуры проекта GenData..."

# Создание основных директорий
mkdir -p gendata-project/{shared/{types,sensorvalue,protocol},generator-service/{cmd,internal/generator},receiver-service/{cmd,internal/{parser,display}},deploy/ansible,scripts}

echo "Структура проекта создана!"
