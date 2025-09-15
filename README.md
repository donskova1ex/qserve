# qserve

**qserve** - это CLI-утилита на Go для выполнения SQL запросов через HTTP API.

## 🎯 Назначение

Утилита принимает настройки подключения к базе данных, поднимает HTTP сервер с эндпоинтами для выполнения SQL запросов и возвращает результаты в формате JSON.

## ✨ Возможности

- 🔧 **Интерактивная настройка** подключения к БД
- 🗄️ **Поддержка 4 типов БД** (PostgreSQL, MySQL, SQLite, MSSQL)
- 🔒 **Безопасное выполнение запросов** с валидацией
- ⚡ **Connection pooling** для оптимизации производительности
- 📊 **Правильная обработка типов данных** из БД
- 🛡️ **Валидация SQL запросов** на опасные операции
- 📝 **Структурированное логирование** с уровнями
- 🌐 **CORS поддержка** для веб-приложений
- ❤️ **Health check** для мониторинга

## 🏗️ Архитектура

Проект построен на принципах **Clean Architecture** с разделением на слои:

- **Config Layer** - интерактивная настройка подключения ✅
- **Database Layer** - управление соединениями и выполнение запросов ✅
- **Handler Layer** - HTTP обработчики ✅
- **Middleware Layer** - логирование, CORS ✅

## 📁 Структура проекта

```
qserve/
├── cmd/qserve/main.go          # Точка входа
├── internal/
│   ├── config/                 # Конфигурация (интерактивный мастер)
│   │   └── config.go
│   ├── database/               # Подключение к БД и выполнение запросов
│   │   ├── connection.go       # Менеджер соединений
│   │   ├── query_executor.go   # Выполнение SQL запросов
│   │   └── validator.go        # Валидация и безопасность
│   ├── handler/                # HTTP обработчики
│   │   └── query_handler.go    # Обработка запросов
│   └── middleware/             # Middleware
│       ├── cors.go             # CORS поддержка
│       └── logger.go           # Логирование
└── go.mod
```

## 🚀 Быстрый старт

### Установка

```bash
git clone <repository-url>
cd qserve
go mod tidy
```

### Запуск

```bash
go run cmd/qserve/main.go
```

Утилита запустит интерактивный мастер настройки:

```
🔧 Welcome to qserve setup!
============================
Select DB type (1: PostgreSQL, 2: MySQL, 3: SQLite, 4: MSSQL): 1
Select DB host [localhost]: localhost
Select DB port: 5432
Select DB user: postgres
Select DB password: 
Select DB name: mydb
Select service port [8080]: 8080
```

## 🔒 Безопасность

- ✅ Валидация SQL запросов на опасные операции (DROP, TRUNCATE, ALTER, CREATE, GRANT, REVOKE)
- ✅ Проверка целых слов (не подстрок) - "DROPPED" не будет заблокирован
- ✅ Удаление комментариев перед валидацией
- ✅ Блокировка системных процедур SQL Server (sp_, xp_, fn_)
- ✅ Определение типа запроса (SELECT, INSERT, UPDATE, DELETE)

## 🎯 Целевое использование

Идеально для:
- Быстрого выполнения SQL запросов через API
- Интеграции с внешними системами
- Администрирования БД через HTTP
- Прототипирования и тестирования
- Создания простых SQL API для веб-приложений

## 📡 API Endpoints

### POST /query
Выполняет SQL запросы

**Запрос:**
```json
{
  "query": "SELECT * FROM users LIMIT 5"
}
```

**Ответ (SELECT):**
```json
{
  "status": "success",
  "data": [
    {"id": 1, "name": "John", "email": "john@example.com"}
  ]
}
```

**Ответ (INSERT/UPDATE/DELETE):**
```json
{
  "status": "success",
  "rows_affected": 1
}
```

### GET /health
Проверка состояния сервиса

**Ответ:**
```json
{
  "status": "success"
}
```

## 🧪 Примеры использования

### SELECT запрос:
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "SELECT * FROM users LIMIT 5"}'
```

### INSERT запрос:
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "INSERT INTO users (name, email) VALUES (\"Alice\", \"alice@example.com\")"}'
```

### Health check:
```bash
curl -X GET http://localhost:8080/health
```

## 🚧 Планы развития

- [ ] Graceful shutdown
- [ ] Поддержка транзакций
- [ ] Rate limiting
- [ ] OpenTelemetry интеграция
- [ ] Аутентификация и авторизация
- [ ] Конфигурация через файлы

## 🛠️ Технологии

- **Go 1.24+** - основной язык
- **database/sql** - стандартная библиотека для работы с БД
- **net/http** - HTTP сервер (новый ServeMux Go 1.22+)
- **log/slog** - структурированное логирование
- **PostgreSQL** - драйвер `github.com/lib/pq`
- **MySQL** - драйвер `github.com/go-sql-driver/mysql`
- **SQLite** - драйвер `github.com/mattn/go-sqlite3`
- **MSSQL** - драйвер `github.com/microsoft/go-mssqldb`

## 📝 Лицензия

MIT License

## 🤝 Вклад в проект

Приветствуются любые вклады! Создавайте issues и pull requests.

---
