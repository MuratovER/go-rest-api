# Архитектура проекта

## 🏗️ Обзор архитектуры

Проект построен на принципах Clean Architecture с четким разделением ответственности между слоями.

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   HTTP Server   │  │   Middleware    │  │   Routes    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Use Cases     │  │   Interfaces    │  │   DTOs      │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │    Entities     │  │   Value Objects │  │   Services  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                       │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Repositories  │  │   Database      │  │   External  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Структура проекта

### 1. Presentation Layer (Слой представления)

#### HTTP Server (`internal/server/`)
- **server.go**: Основной HTTP сервер на Echo
- **handlers.go**: Общие обработчики

#### HTTP Handlers (`internal/inventory_asset/delivery/http/`)
- **handlers.go**: HTTP обработчики для инвентарных активов
- **routes.go**: Маршрутизация API

#### Middleware (`internal/middleware/`)
- **middlewares.go**: Общие middleware
- **csrf.go**: CSRF защита
- **debug.go**: Отладочные middleware
- **request_logger.go**: Логирование запросов
- **sanitize.go**: Санитизация данных

### 2. Application Layer (Слой приложения)

#### Use Cases (`internal/inventory_asset/usecase/`)
- **usecase.go**: Интерфейсы use cases
- **usecase.go**: Реализация бизнес-логики

#### Interfaces (`internal/inventory_asset/`)
- **delivery.go**: Интерфейсы HTTP handlers
- **usecase.go**: Интерфейсы use cases
- **pg_repository.go**: Интерфейсы репозиториев

### 3. Domain Layer (Слой домена)

#### Models (`internal/models/`)
- **inventory_asset.go**: Модели данных

### 4. Infrastructure Layer (Слой инфраструктуры)

#### Repositories (`internal/inventory_asset/repository/`)
- **pg_repository.go**: Реализация репозитория PostgreSQL

#### Database (`pkg/db/postgres/`)
- **db_conn.go**: Подключение к PostgreSQL

#### Utilities (`pkg/`)
- **converter/**: Конвертеры данных
- **csrf/**: CSRF защита
- **httpErrors/**: Обработка HTTP ошибок
- **logger/**: Логирование
- **sanitize/**: Санитизация
- **utils/**: Утилиты

## 🔄 Поток данных

### 1. HTTP Request Flow

```
HTTP Request
    ↓
Middleware (CORS, CSRF, Logging)
    ↓
Router
    ↓
HTTP Handler
    ↓
Use Case
    ↓
Repository
    ↓
Database
```

### 2. Response Flow

```
Database
    ↓
Repository
    ↓
Use Case
    ↓
HTTP Handler
    ↓
Middleware (Response Logging)
    ↓
HTTP Response
```

## 🎯 Принципы проектирования

### 1. Dependency Inversion Principle (DIP)

```go
// Интерфейс в домене
type Repository interface {
    GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error)
}

// Реализация в инфраструктуре
type inventoryAssetRepo struct {
    db     *gorm.DB
    logger logger.Logger
}
```

### 2. Single Responsibility Principle (SRP)

Каждый компонент имеет одну ответственность:
- **Handlers**: Обработка HTTP запросов
- **Use Cases**: Бизнес-логика
- **Repositories**: Доступ к данным
- **Models**: Структуры данных

### 3. Interface Segregation Principle (ISP)

```go
// Минимальные интерфейсы
type UseCase interface {
    GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, id int) (*models.InventoryAssetList, error)
    GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error)
}
```

## 🔧 Конфигурация

### Структура конфигурации

```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Logger   LoggerConfig
}

type ServerConfig struct {
    Port        string
    Mode        string
    AppVersion  string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}
```

### Загрузка конфигурации

```go
// 1. Загрузка файла конфигурации
cfgFile, err := config.LoadConfig(configPath)

// 2. Парсинг конфигурации
cfg, err := config.ParseConfig(cfgFile)

// 3. Использование в компонентах
psqlDB, err := postgres.NewPsqlDB(cfg)
```

## 🗄️ База данных

### Схема данных

```sql
CREATE TABLE inventory_assets (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    serial_code VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    price INTEGER NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### GORM модель

```go
type InventoryAsset struct {
    ID         int    `gorm:"column:id;primary_key;not_null" json:"id"`
    UserName   string `gorm:"not_null" json:"username"`
    Name       string `gorm:"not_null" json:"name"`
    SerialCode string `gorm:"not_null" json:"serial_code"`
    IsActive   bool   `gorm:"not_null" json:"is_active"`
    Price      int    `gorm:"not_null" json:"price"`
    URL        string `gorm:"not_null" json:"url"`
}
```

## 📝 Логирование

### Структура логов

```go
// Уровни логирования
logger.Debug("Debug message")
logger.Info("Info message")
logger.Warn("Warning message")
logger.Error("Error message")
logger.Fatal("Fatal message")
```

### Контекстное логирование

```go
// Логирование с контекстом
logger.Infof("Getting inventory asset by ID: %d", id)
logger.Errorf("Failed to get inventory asset by ID %d: %v", id, err)
```

## 🔒 Безопасность

### CSRF защита

```go
// Middleware для CSRF защиты
csrfMiddleware := csrf.NewCSRFMiddleware(cfg)
e.Use(csrfMiddleware.Handler)
```

### Санитизация данных

```go
// Санитизация HTML
sanitized := sanitize.HTML(input)

// Валидация данных
err := validator.Validate(input)
```

## 🧪 Тестирование

### Структура тестов

```
├── internal/
│   ├── inventory_asset/
│   │   ├── delivery/http/
│   │   │   └── handlers_test.go
│   │   ├── repository/
│   │   │   └── pg_repository_test.go
│   │   └── usecase/
│   │       └── usecase_test.go
│   └── models/
│       └── inventory_asset_test.go
└── pkg/
    └── utils/
        └── pagination_test.go
```

### Моки

```go
//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
type UseCase interface {
    GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error)
}
```

## 🚀 Развертывание

### Docker

```dockerfile
# Многоэтапная сборка
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Docker Compose

```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      
  postgres:
    image: postgres:13
    environment:
      - POSTGRES_DB=inventory_assets
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
```

## 📊 Мониторинг

### Health Checks

```go
// Endpoint для проверки здоровья
func (h *Handler) HealthCheck(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{
        "status": "ok",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}
```

### Метрики

```go
// Базовые метрики
type Metrics struct {
    RequestCount   int64
    ErrorCount     int64
    ResponseTime   time.Duration
}
```

## 🔄 Миграции

### GORM Auto Migration

```go
// Автоматическая миграция
db.AutoMigrate(&models.InventoryAsset{})
```

### Ручные миграции

```sql
-- Создание индексов
CREATE INDEX idx_inventory_assets_username ON inventory_assets(username);
CREATE INDEX idx_inventory_assets_serial_code ON inventory_assets(serial_code);
```

## 📈 Масштабирование

### Горизонтальное масштабирование

1. **Load Balancer**: Nginx или HAProxy
2. **Multiple Instances**: Docker Swarm или Kubernetes
3. **Database**: Read replicas для чтения

### Вертикальное масштабирование

1. **Resources**: Увеличение CPU/RAM
2. **Database**: Улучшение производительности БД
3. **Caching**: Redis для кэширования

## 🔧 Отладка

### Профилирование

```go
import _ "net/http/pprof"

// Включение pprof
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

### Трассировка

```go
// OpenTracing
span, ctx := opentracing.StartSpanFromContext(context.Background(), "operation")
defer span.Finish()
```

## 📚 Дополнительные ресурсы

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [Echo Framework](https://echo.labstack.com/)
- [GORM Documentation](https://gorm.io/docs/) 