# qserve

**qserve** - это HTTP API сервер на Go для безопасного выполнения SQL запросов к различным базам данных.

## 🎯 Назначение

qserve предоставляет RESTful API для выполнения SQL запросов к базам данных с встроенной валидацией безопасности, connection pooling и структурированным логированием. Идеально подходит для интеграции с внешними системами, администрирования БД через HTTP и создания простых SQL API.

## ✨ Возможности

- 🔧 **Интерактивная настройка** подключения к БД через CLI мастер
- 🗄️ **Поддержка 4 типов БД**: PostgreSQL, MySQL, SQLite, MSSQL
- 🔒 **Безопасное выполнение запросов** с валидацией опасных операций
- ⚡ **Connection pooling** (5 idle, 10 max connections)
- 📊 **Правильная обработка типов данных** из БД в JSON
- 🛡️ **Валидация SQL запросов** на опасные операции (DROP, TRUNCATE, ALTER, CREATE, GRANT, REVOKE)
- 📝 **Структурированное логирование** с уровнями (slog)
- 🌐 **CORS поддержка** для веб-приложений
- ❤️ **Health check** для мониторинга состояния БД
- 🏗️ **Clean Architecture** с разделением на слои

## 🏗️ Архитектура

Проект построен на принципах **Clean Architecture** с четким разделением ответственности:

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Layer (net/http)                    │
├─────────────────────────────────────────────────────────────┤
│  Handler Layer    │  Middleware Layer                       │
│  - QueryHandler   │  - CORS Middleware                      │
│  - Health Check   │  - Logging Middleware                   │
├─────────────────────────────────────────────────────────────┤
│                    Service Layer                            │
│  - Query Validation                                         │
│  - Business Logic                                           │
├─────────────────────────────────────────────────────────────┤
│                 Database Layer                              │
│  - ConnectionManager (Connection Pooling)                  │
│  - QueryExecutor (SELECT/INSERT/UPDATE/DELETE)             │
│  - QueryValidator (Security Validation)                    │
├─────────────────────────────────────────────────────────────┤
│                  Config Layer                               │
│  - Interactive Setup Wizard                                 │
│  - Configuration Validation                                 │
└─────────────────────────────────────────────────────────────┘
```

### Слои архитектуры:

- **Config Layer** - интерактивная настройка подключения с валидацией ✅
- **Database Layer** - управление соединениями, выполнение запросов, валидация безопасности ✅
- **Handler Layer** - HTTP обработчики для API endpoints ✅
- **Middleware Layer** - логирование запросов и CORS поддержка ✅

## 📁 Структура проекта

```
qserve/
├── cmd/qserve/main.go              # Точка входа приложения
├── internal/
│   ├── config/                     # Конфигурация и интерактивный мастер
│   │   └── config.go              # Настройка подключения к БД
│   ├── database/                   # Слой работы с базой данных
│   │   ├── connection.go          # Менеджер соединений с connection pooling
│   │   ├── query_executor.go      # Выполнение SQL запросов
│   │   └── validator.go           # Валидация безопасности SQL
│   ├── handler/                    # HTTP обработчики
│   │   └── query_handler.go       # Обработка API запросов
│   └── middleware/                 # HTTP middleware
│       ├── cors.go                # CORS поддержка
│       └── logger.go              # Структурированное логирование
├── qserve-linux-amd64             # Готовый бинарный файл для Linux AMD64
├── qserve-linux-arm64             # Готовый бинарный файл для Linux ARM64
├── qserve-windows.exe             # Готовый бинарный файл для Windows AMD64
├── qserve-windows-arm64.exe       # Готовый бинарный файл для Windows ARM64
├── go.mod                         # Go модули и зависимости
└── README.md                      # Документация
```

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.24+ (для сборки из исходников)
- Одна из поддерживаемых БД: PostgreSQL, MySQL, SQLite, MSSQL

### Установка

#### Вариант 1: Готовые бинарные файлы

В проекте уже собраны готовые исполняемые файлы для различных платформ:

- **Linux AMD64**: `qserve-linux-amd64`
- **Linux ARM64**: `qserve-linux-arm64` 
- **Windows AMD64**: `qserve-windows.exe`
- **Windows ARM64**: `qserve-windows-arm64.exe`

Просто скачайте нужный файл для вашей платформы и запустите:

```bash
# Linux
chmod +x qserve-linux-amd64
./qserve-linux-amd64

# Windows
qserve-windows.exe
```

#### Вариант 2: Сборка из исходников

```bash
git clone <repository-url>
cd qserve
go mod tidy
go run cmd/qserve/main.go
```

#### Вариант 3: Сборка бинарного файла

```bash
# Для текущей платформы
go build -o qserve cmd/qserve/main.go

# Для Linux AMD64
GOOS=linux GOARCH=amd64 go build -o qserve-linux-amd64 cmd/qserve/main.go

# Для Windows AMD64
GOOS=windows GOARCH=amd64 go build -o qserve-windows.exe cmd/qserve/main.go
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

После настройки сервер будет доступен на указанном порту.

## 🔒 Безопасность

### Валидация SQL запросов

- ✅ **Блокировка опасных операций**: DROP, TRUNCATE, ALTER, CREATE, GRANT, REVOKE
- ✅ **Проверка целых слов** (не подстрок) - "DROPPED" не будет заблокирован
- ✅ **Удаление комментариев** перед валидацией
- ✅ **Блокировка системных процедур** SQL Server (sp_, xp_, fn_)
- ✅ **Определение типа запроса** (SELECT, INSERT, UPDATE, DELETE)

### Connection Pooling

- Максимум 5 idle соединений
- Максимум 10 активных соединений
- Автоматическое управление жизненным циклом соединений

## 📡 API Endpoints

### POST /query
Выполняет SQL запросы с валидацией безопасности

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

**Ответ (Ошибка):**
```json
{
  "status": "error",
  "error": "dangerous keyword found: DROP"
}
```

### GET /health
Проверка состояния сервиса и подключения к БД

**Ответ (OK):**
```json
{
  "status": "ok"
}
```

**Ответ (Ошибка БД):**
```json
{
  "status": "error",
  "error": "database connection failed"
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

### Проверка валидации безопасности:
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "DROP TABLE users"}'
# Ответ: {"status": "error", "error": "dangerous keyword found: DROP"}
```

## 🎯 Целевое использование

Идеально для:
- **Быстрого выполнения SQL запросов** через REST API
- **Интеграции с внешними системами** и микросервисами
- **Администрирования БД** через HTTP интерфейс
- **Прототипирования и тестирования** SQL запросов
- **Создания простых SQL API** для веб-приложений
- **Мониторинга состояния БД** через health check
- **Быстрого развертывания** - готовые бинарные файлы для всех популярных платформ
- **Контейнеризации** - статические бинарные файлы без зависимостей

## 🛠️ Технологии

- **Go 1.24+** - основной язык программирования
- **database/sql** - стандартная библиотека для работы с БД
- **net/http** - HTTP сервер с новым ServeMux (Go 1.22+)
- **log/slog** - структурированное логирование
- **PostgreSQL** - драйвер `github.com/lib/pq`
- **MySQL** - драйвер `github.com/go-sql-driver/mysql`
- **SQLite** - драйвер `github.com/mattn/go-sqlite3`
- **MSSQL** - драйвер `github.com/microsoft/go-mssqldb`

## 🚧 Планы развития

- [ ] **Graceful shutdown** с корректным закрытием соединений
- [ ] **Поддержка транзакций** для сложных операций
- [ ] **Rate limiting** для защиты от злоупотреблений
- [ ] **OpenTelemetry интеграция** для observability
- [ ] **Аутентификация и авторизация** (JWT, API keys)
- [ ] **Конфигурация через файлы** (YAML, JSON, ENV)
- [ ] **Метрики Prometheus** для мониторинга
- [ ] **Поддержка prepared statements** для безопасности
- [ ] **Batch операции** для множественных запросов
- [ ] **WebSocket поддержка** для real-time уведомлений

## 📝 Лицензия

MIT License

## 🤝 Вклад в проект

Приветствуются любые вклады! Создавайте issues и pull requests.

### Разработка

```bash
# Запуск тестов
go test ./...

# Проверка кода
go vet ./...
golangci-lint run

# Форматирование
go fmt ./...
```

---

**qserve** - простой, безопасный и эффективный способ выполнения SQL запросов через HTTP API.