# qserve

**qserve** - это CLI-утилита на Go для выполнения SQL запросов через HTTP API.

## 🎯 Назначение

Утилита принимает настройки подключения к базе данных, поднимает HTTP сервер с одним эндпоинтом, который принимает SQL запросы в теле запроса и возвращает результаты в формате JSON.

## ✨ Возможности

- 🔧 **Интерактивная настройка** подключения к БД
- 🗄️ **Поддержка 4 типов БД** (PostgreSQL, MySQL, SQLite, MSSQL)
- 🔒 **Безопасное выполнение запросов** с валидацией
- ⚡ **Connection pooling** для оптимизации производительности
- 📊 **Правильная обработка типов данных** из БД
- 🛡️ **Валидация SQL запросов** на опасные операции

## 🏗️ Архитектура

Проект построен на принципах **Clean Architecture** с разделением на слои:

- **Config Layer** - интерактивная настройка подключения
- **Database Layer** - управление соединениями и выполнение запросов
- **Handler Layer** - HTTP обработчики (в разработке)
- **Service Layer** - бизнес-логика (в разработке)
- **Middleware Layer** - логирование, CORS (в разработке)

## 📁 Структура проекта

```
qserve/
├── cmd/qserve/main.go          # Точка входа
├── internal/
│   ├── config/                 # Конфигурация (интерактивный мастер)
│   ├── database/               # Подключение к БД и выполнение запросов
│   │   ├── connection.go       # Менеджер соединений
│   │   ├── query_executor.go   # Выполнение SQL запросов
│   │   └── validator.go        # Валидация и безопасность
│   ├── handler/                # HTTP обработчики (планируется)
│   ├── service/                # Бизнес-логика (планируется)
│   └── middleware/             # Middleware (планируется)
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

## 🚧 В разработке

- [ ] HTTP API с эндпоинтом `/query`
- [ ] Middleware для логирования и CORS
- [ ] Graceful shutdown
- [ ] Расширенная валидация запросов
- [ ] Поддержка транзакций
- [ ] Rate limiting
- [ ] OpenTelemetry интеграция

## 🛠️ Технологии

- **Go 1.24+** - основной язык
- **database/sql** - стандартная библиотека для работы с БД
- **net/http** - HTTP сервер (новый ServeMux Go 1.22+)
- **PostgreSQL** - драйвер `github.com/lib/pq`
- **MySQL** - драйвер `github.com/go-sql-driver/mysql`
- **SQLite** - драйвер `github.com/mattn/go-sqlite3`
- **MSSQL** - драйвер `github.com/microsoft/go-mssqldb`

## 📝 Лицензия

MIT License

## 🤝 Вклад в проект

Приветствуются любые вклады! Создавайте issues и pull requests.

---
