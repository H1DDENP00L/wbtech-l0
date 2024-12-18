# WBTECH Задание L0

## Описание проекта

**Сервис с простым интерфейсом для отображения данных о приходящих заказах, демонстрирующий работу с брокером сообщений Kafka, базой данных PorstgreSQL, Docker и Docker Compose**

### Технологии

1.  **Go**
2.  **Docker**
3.  **Docker Compose** 
4. **База данных PostgreSQL**
5. **Kafka**
6. **Redpanda Console**

### Библиотеки
1. **gin-gonic/gin**
2. **pressly/goose/v3**
3. **IBM/sarama**
4. **lib/pq**


### Запуск
1.  **Установите зависимости:**

    ```bash
    go mod tidy
    ```
2. **Запустите Docker-контейеры**
    ```bash
    docker-compose up -d 
    ```

3. **Запустите приложение:**

    ```bash
    go run main.go
    ```

### Миграции базы данных

* Чтаблицы буду созданы автоматически если их не было при запуске
  `go run main.go`

* Для удаления таблиц используйте
  `go run migrates/dropTables.go`
