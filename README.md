# Inventory Assets API

REST API для управления инвентарными активами, построенная на Go с использованием Clean Architecture.

## 🏗️ Архитектура

Проект следует принципам Clean Architecture с четким разделением слоев:

```
├── cmd/api/                    # Точка входа в приложение
├── internal/                   # Внутренняя логика приложения
│   ├── inventory_asset/        # Домен инвентарных активов
│   │   ├── delivery/http/      # HTTP handlers
│   │   ├── repository/         # Слой доступа к данным
│   │   └── usecase/           # Бизнес-логика
│   ├── middleware/            # HTTP middleware
│   ├── models/                # Модели данных
│   └── server/                # HTTP сервер
├── pkg/                       # Переиспользуемые пакеты
│   ├── converter/             # Конвертеры данных
│   ├── csrf/                  # CSRF защита
│   ├── db/postgres/           # Подключение к PostgreSQL
│   ├── httpErrors/            # Обработка HTTP ошибок
│   ├── logger/                # Логирование
│   ├── sanitize/              # Санитизация данных
│   └── utils/                 # Утилиты
└── config/                    # Конфигурация
```

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.24+
- PostgreSQL 12+
- Docker (опционально)

### Установка

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd go-rest-api
```

2. Установите зависимости:
```bash
go mod download
```

3. Настройте переменные окружения:
```bash
cp .env.example .env
# Отредактируйте .env файл
```

4. Запустите PostgreSQL:
```bash
# Используя Docker
docker-compose up -d postgres

# Или локально
# Убедитесь, что PostgreSQL запущен и доступен
```

5. Запустите приложение:
```bash
go run cmd/api/main.go
```

### Использование Docker

```bash
# Сборка и запуск
docker-compose up --build

# Только для разработки
docker-compose -f docker-compose.dev.yml up --build
```

## 📚 API Документация

### Swagger UI

После запуска приложения, документация доступна по адресу:
- http://localhost:8080/swagger/index.html

### Основные эндпоинты

#### Получить все инвентарные активы
```http
GET /api/v1/inventory-assets?page=1&size=10
```

**Параметры:**
- `page` (int, optional): Номер страницы (по умолчанию: 1)
- `size` (int, optional): Количество элементов на странице (по умолчанию: 10)
- `orderBy` (string, optional): Поле для сортировки

**Ответ:**
```json
{
  "items": [
    {
      "id": 1,
      "username": "john_doe",
      "name": "Laptop Dell XPS 13",
      "serial_code": "DLXPS13001",
      "is_active": true,
      "price": 120000,
      "url": "https://example.com/laptop1"
    }
  ],
  "total": 100,
  "page": 1,
  "size": 10,
  "pages": 10
}
```

#### Получить инвентарный актив по ID
```http
GET /api/v1/inventory-assets/{id}
```

**Параметры:**
- `id` (int, required): ID инвентарного актива

**Ответ:**
```json
{
  "items": [
    {
      "id": 1,
      "username": "john_doe",
      "name": "Laptop Dell XPS 13",
      "serial_code": "DLXPS13001",
      "is_active": true,
      "price": 120000,
      "url": "https://example.com/laptop1"
    }
  ],
  "total": 1,
  "page": 1,
  "size": 10,
  "pages": 1
}
```

## 🛠️ Разработка

### Структура проекта

#### Модели данных (`internal/models/`)

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

#### Repository (`internal/inventory_asset/repository/`)

```go
type Repository interface {
    GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, id int) (*models.InventoryAssetList, error)
    GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error)
}
```

#### Use Case (`internal/inventory_asset/usecase/`)

```go
type UseCase interface {
    GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, id int) (*models.InventoryAssetList, error)
    GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error)
}
```

#### HTTP Handlers (`internal/inventory_asset/delivery/http/`)

```go
type Handlers interface {
    GetInventoryAssetsById() echo.HandlerFunc
    GetInventoryAssets() echo.HandlerFunc
}
```

### Тестирование

```bash
# Запуск всех тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...

# Запуск тестов с подробным выводом
go test -v ./...
```

### Линтинг

```bash
# Запуск golangci-lint
golangci-lint run

# Форматирование кода
go fmt ./...

# Проверка кода
go vet ./...
```

## 🔧 Конфигурация

### Переменные окружения

```env
# Сервер
SERVER_PORT=8080
SERVER_MODE=debug
SERVER_APP_VERSION=1.0.0

# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=inventory_assets
DB_SSL_MODE=disable

# Логирование
LOGGER_LEVEL=debug
LOGGER_DEV_MODE=true
```

### Конфигурационные файлы

Поддерживаются форматы:
- `.env`
- `.yaml`
- `.json`
- `.toml`

## 📦 Зависимости

### Основные зависимости

- **Echo v4** - HTTP фреймворк
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - Основная база данных
- **Zap** - Логирование
- **Viper** - Управление конфигурацией
- **Validator** - Валидация данных

### Разработка

- **Swagger** - Документация API
- **OpenTracing** - Трассировка
- **Bluemonday** - Санитизация HTML

## 🚀 Развертывание

### Production

```bash
# Сборка для production
go build -o bin/api cmd/api/main.go

# Запуск
./bin/api
```

### Docker

```bash
# Сборка образа
docker build -t inventory-assets-api .

# Запуск контейнера
docker run -p 8080:8080 inventory-assets-api
```

## 📝 Логирование

Приложение использует структурированное логирование с помощью Zap:

```go
logger.Infof("Successfully retrieved inventory asset by ID: %d", id)
logger.Errorf("Failed to get inventory assets: %v", err)
```

## 🔒 Безопасность

- CSRF защита
- Санитизация входных данных
- Валидация запросов
- Логирование безопасности

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## 📄 Лицензия

MIT License - см. файл [LICENSE](LICENSE) для подробностей.

## 🆘 Поддержка

Если у вас есть вопросы или проблемы:

1. Проверьте [Issues](https://github.com/your-repo/issues)
2. Создайте новое Issue с подробным описанием
3. Обратитесь к команде разработки

---

**Версия:** 1.0.0  
**Последнее обновление:** 2024
