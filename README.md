# TEST Sender->Receiver data 
  Тест-испытание поставщика и приемника данных.
  Проект как сервисы размещен на уд. VM srv102.

## Описание проекта в:
* **[send-receiv-prj](md/send-receiv-prj.md)**
* **[Уточнения по расположению файлов](md/dockers-quest.md)**
* **[Распределение файлов](md/files-to-folders.md)** 
* **[send-receiv-prj.pdf](md/send-receiv-prj.pdf)**
* **[Добавление http в проект помимо tcp](md/add-http-to-prj.md)**

## Для запуска использовать make или команду из Makefile
### Пояснение:
* **Работа из CLI использует:**
  * Makefile
  * deploy/ansible/deploy.sh
  * прямые команды (см. Makefile) 
### Работа с проектом 
* **Соответственно для задач на локальной машине:**
  * Сборка образов
  ```bash 
  make build-images
  # or
	chmod +x scripts/build.sh && ./scripts/build.sh 
  ```
  * Запуск со сборкой: 
  ```bash
  make local-run
  # or
  docker compose up --build -d
  ```
  * Остановка: 
  ```bash
  make local-stop
  # or
	docker compose down
  ```  
  * Перезапуск локальных сервисов
  ```bash 
  make local-restart
  # or
	docker compose restart
  ```
  * Локальное тестирование
  ```bash 
  make test
  # or
	chmod +x scripts/test-local.sh && ./scripts/test-local.sh
  ```
  * Просмотр локальных логов
  ```bash 
  make local-logs
  # or
	docker compose logs -f
  ```
  * Показать статус локальных контейнеров
  ```bash 
  make local-status
  # or
	docker compose ps
  ```
* **Соответственно для задач на удаленной машине:**
  * Развертывание на удаленном сервере: 
  ```bash
  make deploy
  # or
	cd deploy/ansible && chmod +x deploy.sh && ./deploy.sh
  ```
  * Остановка на удаленном сервере: 
  ```bash
  make stop
  # or
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose down'
  ``` 
  * Перезапуск сервисов на удаленном сервере (без пересборки)
  ```bash
  make restart
  # or
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose restart'  
  ```
  * Просмотр логов на удаленном сервере
  ```bash 
  make logs
  # or
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose logs -f'
  ```
  * Очистка на удаленном сервере
  ```bash 
  make clean
  # or
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose down --rmi all --volumes'
  ```
  * Проверка статуса сервисов
  ```bash 
  make status
  # or
	ssh starmark@192.168.88.102 'cd /opt/gendata-project && docker compose ps' 
  ```
