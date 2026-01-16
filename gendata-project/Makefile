.PHONY: build deploy logs stop clean status local-run local-stop test build-images
# Сборка локально
build:
	docker-compose build

# Развертывание на удаленном сервере
deploy:
	cd deploy/ansible && chmod +x deploy.sh && ./deploy.sh

# Просмотр логов на удаленном сервере
logs:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose logs -f'

# Остановка сервисов на удаленном сервере
stop:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose down'

# Перезапуск сервисов на удаленном сервере (без пересборки)
restart:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose restart'
	
# Очистка на удаленном сервере
clean:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose down --rmi all --volumes'

# Проверка статуса сервисов
status:
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose ps'

# Локальный запуск для тестирования
local-run:
	docker compose up --build -d

# Локальная остановка
local-stop:
	docker compose down

# Локальное тестирование
test:
	chmod +x scripts/test-local.sh && ./scripts/test-local.sh

# Сборка образов
build-images:
	chmod +x scripts/build.sh && ./scripts/build.sh

# Просмотр локальных логов
local-logs:
	docker compose logs -f

# Перезапуск локальных сервисов
local-restart:
	docker compose restart

# Показать статус локальных контейнеров
local-status:
	docker compose ps
